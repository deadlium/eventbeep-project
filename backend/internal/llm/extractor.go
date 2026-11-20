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
	"sync"
	"time"
)

const (
	DefaultGeminiModel   = "gemini-1.5-flash"
	DefaultGeminiBaseURL = "https://generativelanguage.googleapis.com"
	SystemPrompt         = "Extract 'role' and 'city' into JSON: {\"role\": \"string\", \"city\": \"string\"}"
	DefaultTimeout       = 15 * time.Second
	MaxRetries           = 3
	CircuitCooldown      = 60 * time.Second
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
// with exponential backoff retries on 503/429 spikes, circuit breaker protection, and fallback heuristics.
type GeminiExtractor struct {
	apiKey            string
	model             string
	baseURL           string
	client            HTTPClient
	fallback          Extractor
	consecutiveErrors int
	circuitOpenUntil  time.Time
	mu                sync.RWMutex
}

// NewGeminiExtractor constructs a GeminiExtractor with configurable options.
func NewGeminiExtractor(apiKey string, opts ...Option) *GeminiExtractor {
	e := &GeminiExtractor{
		apiKey:   strings.TrimSpace(apiKey),
		model:    DefaultGeminiModel,
		baseURL:  DefaultGeminiBaseURL,
		client:   &http.Client{Timeout: DefaultTimeout},
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

func WithTimeout(timeout time.Duration) Option {
	return func(e *GeminiExtractor) {
		if timeout > 0 {
			if c, ok := e.client.(*http.Client); ok {
				c.Timeout = timeout
			} else {
				e.client = &http.Client{Timeout: timeout}
			}
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
	Contents          []geminiContent          `json:"contents"`
	SystemInstruction *geminiSystemInstruction `json:"system_instruction,omitempty"`
	GenerationConfig  *geminiGenerationConfig  `json:"generationConfig,omitempty"`
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

// Gemini Interactions API request schema (POST /v1beta/interactions)
type interactionsRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

// Gemini Interactions API response schema
type interactionsResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Steps  []struct {
		Type    string `json:"type"`
		Content []struct {
			Text string `json:"text"`
			Type string `json:"type"`
		} `json:"content"`
	} `json:"steps"`
	Error *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
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

// Extract calls the Gemini REST API to extract role and city, delegating to fallback on error or circuit trip.
func (e *GeminiExtractor) Extract(ctx context.Context, title, description string) (*ExtractedEntity, error) {
	// If no API key configured, use rule-based fallback directly
	if e.apiKey == "" {
		return e.fallback.Extract(ctx, title, description)
	}

	// Check circuit breaker status
	e.mu.RLock()
	isOpen := time.Now().Before(e.circuitOpenUntil)
	e.mu.RUnlock()

	if isOpen {
		return e.fallback.Extract(ctx, title, description)
	}

	entity, err := e.callGeminiAPI(ctx, title, description)
	if err != nil {
		e.mu.Lock()
		e.consecutiveErrors++
		if e.consecutiveErrors >= 3 {
			e.circuitOpenUntil = time.Now().Add(CircuitCooldown)
			log.Printf("[llm] Gemini API temporary cooldown active (tripped by consecutive errors, resuming in %v)", CircuitCooldown)
		}
		e.mu.Unlock()

		if strings.Contains(err.Error(), "503") || strings.Contains(err.Error(), "service_unavailable") || strings.Contains(err.Error(), "high demand") {
			log.Printf("[llm] Gemini API temporarily busy (503 high demand spike), seamlessly activating fallback extractor")
		} else {
			log.Printf("[llm] Gemini API extraction failed (%v), activating fallback extractor", err)
		}
		return e.fallback.Extract(ctx, title, description)
	}

	// Successful call resets consecutive error counter
	e.mu.Lock()
	e.consecutiveErrors = 0
	e.circuitOpenUntil = time.Time{}
	e.mu.Unlock()

	return entity, nil
}

// ValidateKey directly tests the Gemini API credentials and model connection without activating the fallback.
func (e *GeminiExtractor) ValidateKey(ctx context.Context) (*ExtractedEntity, error) {
	if e.apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is empty or not set")
	}
	return e.callGeminiAPI(ctx, "Junior Golang Developer", "Fresher software engineer position based in Bengaluru, India.")
}

func (e *GeminiExtractor) callGeminiAPI(ctx context.Context, title, description string) (*ExtractedEntity, error) {
	// Keys with AQ. prefix or models containing 3.8 target the Gemini Interactions API
	if strings.HasPrefix(e.apiKey, "AQ.") || strings.Contains(e.model, "3.8") {
		return e.callInteractionsAPI(ctx, title, description)
	}

	entity, err := e.callGenerateContentAPI(ctx, title, description)
	if err != nil && (strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "unsupported")) {
		// Attempt interactions API as fallback
		if intEntity, intErr := e.callInteractionsAPI(ctx, title, description); intErr == nil {
			return intEntity, nil
		}
	}
	return entity, err
}

func isTransientHTTPStatus(code int) bool {
	return code == http.StatusServiceUnavailable || // 503
		code == http.StatusTooManyRequests || // 429
		code == http.StatusBadGateway || // 502
		code == http.StatusGatewayTimeout || // 504
		code == http.StatusInternalServerError // 500
}

func (e *GeminiExtractor) callInteractionsAPI(ctx context.Context, title, description string) (*ExtractedEntity, error) {
	model := e.model
	if model == DefaultGeminiModel || model == "" {
		model = "gemini-3.8-flash"
	}

	prompt := fmt.Sprintf("%s\n\nTitle: %s\nDescription: %s\nReturn ONLY the JSON object without code fences or additional commentary.", SystemPrompt, title, description)
	reqBody := interactionsRequest{
		Model: model,
		Input: prompt,
	}

	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal interactions request: %w", err)
	}

	endpoint := fmt.Sprintf("%s/v1beta/interactions", e.baseURL)

	var lastErr error
	for attempt := 1; attempt <= MaxRetries; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(jsonBytes))
		if err != nil {
			return nil, fmt.Errorf("failed to create interactions request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-goog-api-key", e.apiKey)

		resp, err := e.client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("http request failed: %w", err)
			if attempt < MaxRetries {
				delay := time.Duration(attempt*500) * time.Millisecond
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				case <-time.After(delay):
					continue
				}
			}
			return nil, lastErr
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			lastErr = fmt.Errorf("gemini interactions API error (status %d): %s", resp.StatusCode, string(bodyBytes))
			if isTransientHTTPStatus(resp.StatusCode) && attempt < MaxRetries {
				delay := time.Duration(attempt*500) * time.Millisecond
				log.Printf("[llm] Gemini %s high demand/busy (%d), retrying in %v (attempt %d/%d)...", model, resp.StatusCode, delay, attempt, MaxRetries)
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				case <-time.After(delay):
					continue
				}
			}
			return nil, lastErr
		}

		var intResp interactionsResponse
		if err := json.Unmarshal(bodyBytes, &intResp); err != nil {
			return nil, fmt.Errorf("failed to unmarshal interactions response: %w", err)
		}

		if intResp.Error != nil {
			return nil, fmt.Errorf("gemini interactions error %d: %s", intResp.Error.Code, intResp.Error.Message)
		}

		var rawText string
		for _, step := range intResp.Steps {
			if step.Type == "model_output" {
				for _, part := range step.Content {
					if part.Text != "" {
						rawText += part.Text
					}
				}
			}
		}
		if rawText == "" {
			for _, step := range intResp.Steps {
				for _, part := range step.Content {
					if part.Text != "" {
						rawText += part.Text
					}
				}
			}
		}

		if strings.TrimSpace(rawText) == "" {
			return nil, fmt.Errorf("empty output in gemini interactions response: %s", string(bodyBytes))
		}

		cleanedJSON := cleanJSONText(rawText)

		var entity ExtractedEntity
		if err := json.Unmarshal([]byte(cleanedJSON), &entity); err != nil {
			if start := strings.Index(cleanedJSON, "{"); start != -1 {
				if end := strings.LastIndex(cleanedJSON, "}"); end != -1 && end > start {
					if err2 := json.Unmarshal([]byte(cleanedJSON[start:end+1]), &entity); err2 == nil {
						entity.Role = strings.TrimSpace(entity.Role)
						entity.City = strings.TrimSpace(entity.City)
						return &entity, nil
					}
				}
			}
			return nil, fmt.Errorf("failed to parse extracted entity JSON (%s): %w", rawText, err)
		}

		entity.Role = strings.TrimSpace(entity.Role)
		entity.City = strings.TrimSpace(entity.City)

		return &entity, nil
	}

	return nil, lastErr
}

func (e *GeminiExtractor) callGenerateContentAPI(ctx context.Context, title, description string) (*ExtractedEntity, error) {
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

	var lastErr error
	for attempt := 1; attempt <= MaxRetries; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(jsonBytes))
		if err != nil {
			return nil, fmt.Errorf("failed to create http request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-goog-api-key", e.apiKey)

		resp, err := e.client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("http request failed: %w", err)
			if attempt < MaxRetries {
				delay := time.Duration(attempt*500) * time.Millisecond
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				case <-time.After(delay):
					continue
				}
			}
			return nil, lastErr
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			lastErr = fmt.Errorf("gemini API error (status %d): %s", resp.StatusCode, string(bodyBytes))
			if isTransientHTTPStatus(resp.StatusCode) && attempt < MaxRetries {
				delay := time.Duration(attempt*500) * time.Millisecond
				log.Printf("[llm] Gemini %s transient error (%d), retrying in %v (attempt %d/%d)...", e.model, resp.StatusCode, delay, attempt, MaxRetries)
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				case <-time.After(delay):
					continue
				}
			}
			return nil, lastErr
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

	return nil, lastErr
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
