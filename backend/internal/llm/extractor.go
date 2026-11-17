package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	DefaultGeminiModel   = "gemini-1.5-flash"
	DefaultGeminiBaseURL = "https://generativelanguage.googleapis.com"
	SystemPrompt         = "Extract 'role' and 'city' into JSON: {\"role\": \"string\", \"city\": \"string\"}"
)

// ExtractedEntity contains the standardized role and city parsed by the LLM or fallback.
type ExtractedEntity struct {
	Role string `json:"role"`
	City string `json:"city"`
}

// Extractor defines the interface for extracting role and city from raw job text.
type Extractor interface {
	Extract(ctx context.Context, title, description string) (*ExtractedEntity, error)
}

// HTTPClient abstracts the HTTP transport for Gemini API calls.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// GeminiExtractor calls the Google Gemini API with system instructions and JSON structured output,
// and falls back gracefully to rule-based heuristics on missing keys or network errors.
type GeminiExtractor struct {
	apiKey   string
	model    string
	baseURL  string
	client   HTTPClient
	fallback Extractor
}

// NewGeminiExtractor constructs a GeminiExtractor with configurable options.
func NewGeminiExtractor(apiKey string, opts ...Option) *GeminiExtractor {
	e := &GeminiExtractor{
		apiKey:   strings.TrimSpace(apiKey),
		model:    DefaultGeminiModel,
		baseURL:  DefaultGeminiBaseURL,
		client:   &http.Client{Timeout: 15 * time.Second},
		fallback: NewFallbackExtractor(),
	}

	for _, opt := range opts {
		opt(e)
	}

	return e
}

// Option allows configuring the GeminiExtractor.
type Option func(*GeminiExtractor)

func WithModel(model string) Option {
	return func(e *GeminiExtractor) {
		if strings.TrimSpace(model) != "" {
			e.model = strings.TrimSpace(model)
		}
	}
}

func WithBaseURL(url string) Option {
	return func(e *GeminiExtractor) {
		if strings.TrimSpace(url) != "" {
			e.baseURL = strings.TrimRight(strings.TrimSpace(url), "/")
		}
	}
}

func WithHTTPClient(client HTTPClient) Option {
	return func(e *GeminiExtractor) {
		if client != nil {
			e.client = client
		}
	}
}

func WithFallback(fallback Extractor) Option {
	return func(e *GeminiExtractor) {
		if fallback != nil {
			e.fallback = fallback
		}
	}
}

// Gemini API request schemas
type geminiRequest struct {
	Contents          []geminiContent         `json:"contents"`
	SystemInstruction *geminiSystemInstruction `json:"system_instruction,omitempty"`
	GenerationConfig  *geminiGenerationConfig `json:"generationConfig,omitempty"`
}

type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiSystemInstruction struct {
	Parts []geminiPart `json:"parts"`
}

type geminiGenerationConfig struct {
	ResponseMimeType string `json:"response_mime_type"`
}

// Gemini API response schemas
type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
	Error *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// Extract calls the Gemini REST API to extract role and city, delegating to fallback on error.
func (e *GeminiExtractor) Extract(ctx context.Context, title, description string) (*ExtractedEntity, error) {
	// If no API key configured, use rule-based fallback directly
	if e.apiKey == "" {
		return e.fallback.Extract(ctx, title, description)
	}

	entity, err := e.callGeminiAPI(ctx, title, description)
	if err != nil {
		log.Printf("[llm] Gemini API extraction failed (%v), activating fallback extractor", err)
		return e.fallback.Extract(ctx, title, description)
	}

	return entity, nil
}

func (e *GeminiExtractor) callGeminiAPI(ctx context.Context, title, description string) (*ExtractedEntity, error) {
	reqBody := geminiRequest{
		SystemInstruction: &geminiSystemInstruction{
			Parts: []geminiPart{
				{Text: SystemPrompt},
			},
		},
		Contents: []geminiContent{
			{
				Parts: []geminiPart{
					{Text: fmt.Sprintf("Title: %s\nDescription: %s", title, description)},
				},
			},
		},
		GenerationConfig: &geminiGenerationConfig{
			ResponseMimeType: "application/json",
		},
	}

	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	endpoint := fmt.Sprintf("%s/v1beta/models/%s:generateContent?key=%s", e.baseURL, e.model, e.apiKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("gemini API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	var geminiResp geminiResponse
	if err := json.Unmarshal(bodyBytes, &geminiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal gemini response: %w", err)
	}

	if geminiResp.Error != nil {
		return nil, fmt.Errorf("gemini error %d: %s", geminiResp.Error.Code, geminiResp.Error.Message)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("empty candidates in gemini response")
	}

	rawText := strings.TrimSpace(geminiResp.Candidates[0].Content.Parts[0].Text)
	cleanedJSON := cleanJSONText(rawText)

	var entity ExtractedEntity
	if err := json.Unmarshal([]byte(cleanedJSON), &entity); err != nil {
		return nil, fmt.Errorf("failed to parse extracted entity JSON (%s): %w", rawText, err)
	}

	entity.Role = strings.TrimSpace(entity.Role)
	entity.City = strings.TrimSpace(entity.City)

	return &entity, nil
}

// cleanJSONText strips markdown code fences (```json ... ```) if returned by the LLM.
func cleanJSONText(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "```") {
		// Strip starting ``` or ```json
		if idx := strings.Index(s, "\n"); idx != -1 {
			s = s[idx+1:]
		}
		// Strip trailing ```
		if idx := strings.LastIndex(s, "```"); idx != -1 {
			s = s[:idx]
		}
	}
	return strings.TrimSpace(s)
}
