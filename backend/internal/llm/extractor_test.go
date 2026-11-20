package llm_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/deadlium/eventbeep/backend/internal/llm"
)

type mockHTTPClient struct {
	doFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.doFunc(req)
}

func TestFallbackExtractor_RolesAndCities(t *testing.T) {
	fb := llm.NewFallbackExtractor()

	testCases := []struct {
		title       string
		desc        string
		expectedRole string
		expectedCity string
	}{
		{
			title:       "Junior Golang Backend Developer",
			desc:        "Looking for freshers in Bangalore with Go and PostgreSQL experience.",
			expectedRole: "Backend Engineer",
			expectedCity: "Bangalore",
		},
		{
			title:       "React UI Intern",
			desc:        "Work on frontend user experience in Mumbai.",
			expectedRole: "Frontend Engineer",
			expectedCity: "Mumbai",
		},
		{
			title:       "Full Stack Developer Fresher",
			desc:        "Opportunity in Gurgaon / Delhi NCR for React and Node.",
			expectedRole: "Full Stack Engineer",
			expectedCity: "Delhi NCR",
		},
		{
			title:       "Data Analyst Trainee",
			desc:        "Remote position for Python and SQL enthusiasts.",
			expectedRole: "Data Analyst",
			expectedCity: "Remote",
		},
		{
			title:       "DevOps Cloud Engineer",
			desc:        "Join our cloud infrastructure team in Pune.",
			expectedRole: "DevOps Engineer",
			expectedCity: "Pune",
		},
		{
			title:       "QA Software Tester",
			desc:        "Automation testing role based in Hyderabad.",
			expectedRole: "QA Engineer",
			expectedCity: "Hyderabad",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			entity, err := fb.Extract(context.Background(), tc.title, tc.desc)
			if err != nil {
				t.Fatalf("unexpected extract error: %v", err)
			}
			if entity.Role != tc.expectedRole {
				t.Errorf("expected role '%s', got '%s'", tc.expectedRole, entity.Role)
			}
			if entity.City != tc.expectedCity {
				t.Errorf("expected city '%s', got '%s'", tc.expectedCity, entity.City)
			}
		})
	}
}

func TestGeminiExtractor_EmptyAPIKey_UsesFallback(t *testing.T) {
	extractor := llm.NewGeminiExtractor("")

	entity, err := extractor.Extract(context.Background(), "Junior Go Backend Developer", "Bangalore office")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if entity.Role != "Backend Engineer" {
		t.Errorf("expected fallback role 'Backend Engineer', got '%s'", entity.Role)
	}
	if entity.City != "Bangalore" {
		t.Errorf("expected fallback city 'Bangalore', got '%s'", entity.City)
	}
}

func TestGeminiExtractor_SuccessWithMockServer(t *testing.T) {
	client := &mockHTTPClient{
		doFunc: func(r *http.Request) (*http.Response, error) {
			if r.URL.Query().Get("key") != "secret-test-key" && r.Header.Get("x-goog-api-key") != "secret-test-key" {
				return &http.Response{
					StatusCode: http.StatusUnauthorized,
					Body:       io.NopCloser(bytes.NewBufferString("unauthorized")),
				}, nil
			}

			resp := map[string]interface{}{
				"candidates": []map[string]interface{}{
					{
						"content": map[string]interface{}{
							"parts": []map[string]interface{}{
								{
									"text": `{"role": "Backend Engineer", "city": "Bengaluru"}`,
								},
							},
						},
					},
				},
			}
			b, _ := json.Marshal(resp)
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(bytes.NewReader(b)),
			}, nil
		},
	}

	extractor := llm.NewGeminiExtractor(
		"secret-test-key",
		llm.WithHTTPClient(client),
		llm.WithModel("gemini-1.5-flash"),
	)

	entity, err := extractor.Extract(context.Background(), "Backend Developer", "Bengaluru office")
	if err != nil {
		t.Fatalf("unexpected extraction error: %v", err)
	}

	if entity.Role != "Backend Engineer" {
		t.Errorf("expected role 'Backend Engineer', got '%s'", entity.Role)
	}
	if entity.City != "Bengaluru" {
		t.Errorf("expected city 'Bengaluru', got '%s'", entity.City)
	}
}

func TestGeminiExtractor_MarkdownCodeFenceHandling(t *testing.T) {
	client := &mockHTTPClient{
		doFunc: func(r *http.Request) (*http.Response, error) {
			resp := map[string]interface{}{
				"candidates": []map[string]interface{}{
					{
						"content": map[string]interface{}{
							"parts": []map[string]interface{}{
								{
									"text": "```json\n{\"role\": \"Frontend Engineer\", \"city\": \"Remote\"}\n```",
								},
							},
						},
					},
				},
			}
			b, _ := json.Marshal(resp)
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(bytes.NewReader(b)),
			}, nil
		},
	}

	extractor := llm.NewGeminiExtractor(
		"secret-test-key",
		llm.WithHTTPClient(client),
	)

	entity, err := extractor.Extract(context.Background(), "React Dev", "Anywhere")
	if err != nil {
		t.Fatalf("unexpected extraction error: %v", err)
	}

	if entity.Role != "Frontend Engineer" {
		t.Errorf("expected role 'Frontend Engineer', got '%s'", entity.Role)
	}
	if entity.City != "Remote" {
		t.Errorf("expected city 'Remote', got '%s'", entity.City)
	}
}

func TestGeminiExtractor_ErrorTriggersFallback(t *testing.T) {
	client := &mockHTTPClient{
		doFunc: func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(bytes.NewBufferString("server error")),
			}, nil
		},
	}

	extractor := llm.NewGeminiExtractor(
		"secret-test-key",
		llm.WithHTTPClient(client),
	)

	entity, err := extractor.Extract(context.Background(), "Junior Golang Developer", "Bangalore location")
	if err != nil {
		t.Fatalf("expected fallback to succeed, got error: %v", err)
	}

	if entity.Role != "Backend Engineer" {
		t.Errorf("expected fallback role 'Backend Engineer', got '%s'", entity.Role)
	}
	if entity.City != "Bangalore" {
		t.Errorf("expected fallback city 'Bangalore', got '%s'", entity.City)
	}
}

func TestFallbackExtractor_MoreRolesAndCities(t *testing.T) {
	fb := llm.NewFallbackExtractor()

	testCases := []struct {
		title string
		desc  string
		role  string
		city  string
	}{
		{"AI Engineer", "Deep learning NLP in Chennai", "Data Scientist", "Chennai"},
		{"Mobile App Developer", "Flutter developer in Kolkata", "Mobile Developer", "Kolkata"},
		{"Associate Product Manager", "Roadmaps and telemetry in Ahmedabad", "Product Manager", "Ahmedabad"},
		{"Growth Hacker", "Marketing and conversions in Mumbai", "Growth Associate", "Mumbai"},
		{"Random Title", "Nothing special", "Random Title", "Remote"},
	}

	for _, tc := range testCases {
		e, err := fb.Extract(context.Background(), tc.title, tc.desc)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if e.Role != tc.role {
			t.Errorf("expected role %s, got %s", tc.role, e.Role)
		}
		if e.City != tc.city {
			t.Errorf("expected city %s, got %s", tc.city, e.City)
		}
	}
}

func TestGeminiExtractor_OptionsAndErrorBranches(t *testing.T) {
	fb := llm.NewFallbackExtractor()
	customClient := &http.Client{}

	e := llm.NewGeminiExtractor("key",
		llm.WithHTTPClient(customClient),
		llm.WithFallback(fb),
		llm.WithModel("gemini-2.0-flash"),
	)
	if e == nil {
		t.Fatal("expected non-nil extractor")
	}

	// Test client returning invalid JSON body (corrupt candidates)
	corruptClient := &mockHTTPClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(bytes.NewBufferString(`{"candidates": [{"content": {"parts": [{"text": "not-valid-json"}]}}]}`)),
			}, nil
		},
	}

	eCorrupt := llm.NewGeminiExtractor("test-key", llm.WithHTTPClient(corruptClient))
	res, err := eCorrupt.Extract(context.Background(), "Backend Developer", "Bangalore")
	if err != nil {
		t.Fatalf("expected fallback on corrupt JSON, got error: %v", err)
	}
	if res.Role != "Backend Engineer" {
		t.Errorf("expected fallback role 'Backend Engineer', got %s", res.Role)
	}

	// Test client returning API error in payload
	errorClient := &mockHTTPClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(bytes.NewBufferString(`{"error": {"code": 403, "message": "API key expired"}}`)),
			}, nil
		},
	}

	eErr := llm.NewGeminiExtractor("test-key", llm.WithHTTPClient(errorClient))
	resErr, err := eErr.Extract(context.Background(), "QA Tester", "Hyderabad")
	if err != nil {
		t.Fatalf("expected fallback on API error payload, got error: %v", err)
	}
	if resErr.Role != "QA Engineer" {
		t.Errorf("expected fallback QA Engineer, got %s", resErr.Role)
	}
}

func TestGeminiExtractor_InteractionsAPI(t *testing.T) {
	client := &mockHTTPClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/v1beta/interactions" {
				t.Fatalf("expected path /v1beta/interactions, got %s", req.URL.Path)
			}
			if req.Header.Get("x-goog-api-key") != "AQ.TestKey123" {
				t.Fatalf("expected x-goog-api-key header, got %s", req.Header.Get("x-goog-api-key"))
			}

			resp := map[string]interface{}{
				"id":     "interaction-123",
				"status": "completed",
				"steps": []map[string]interface{}{
					{
						"type": "thought",
					},
					{
						"type": "model_output",
						"content": []map[string]interface{}{
							{
								"type": "text",
								"text": `{"role": "Full Stack Engineer", "city": "Bengaluru"}`,
							},
						},
					},
				},
			}
			b, _ := json.Marshal(resp)
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(bytes.NewReader(b)),
			}, nil
		},
	}

	extractor := llm.NewGeminiExtractor(
		"AQ.TestKey123",
		llm.WithHTTPClient(client),
		llm.WithModel("gemini-3.8-flash"),
	)

	entity, err := extractor.Extract(context.Background(), "Full Stack Dev", "Bengaluru office")
	if err != nil {
		t.Fatalf("unexpected extraction error: %v", err)
	}

	if entity.Role != "Full Stack Engineer" {
		t.Errorf("expected role 'Full Stack Engineer', got '%s'", entity.Role)
	}
	if entity.City != "Bengaluru" {
		t.Errorf("expected city 'Bengaluru', got '%s'", entity.City)
	}
}

func TestGeminiExtractor_503RetrySuccess(t *testing.T) {
	attempts := 0
	client := &mockHTTPClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			attempts++
			if attempts == 1 {
				// Return 503 high demand spike on first attempt
				return &http.Response{
					StatusCode: http.StatusServiceUnavailable,
					Header:     http.Header{"Content-Type": []string{"application/json"}},
					Body:       io.NopCloser(bytes.NewBufferString(`{"error":{"message":"gemini-3.8-flash is currently experiencing high demand","code":"service_unavailable"}}`)),
				}, nil
			}

			// Return success on second attempt
			resp := map[string]interface{}{
				"id":     "interaction-retry",
				"status": "completed",
				"steps": []map[string]interface{}{
					{
						"type": "model_output",
						"content": []map[string]interface{}{
							{
								"type": "text",
								"text": `{"role": "Backend Engineer", "city": "Pune"}`,
							},
						},
					},
				},
			}
			b, _ := json.Marshal(resp)
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(bytes.NewReader(b)),
			}, nil
		},
	}

	extractor := llm.NewGeminiExtractor(
		"AQ.TestRetryKey",
		llm.WithHTTPClient(client),
		llm.WithModel("gemini-3.8-flash"),
	)

	entity, err := extractor.Extract(context.Background(), "Golang Developer", "Pune location")
	if err != nil {
		t.Fatalf("expected retry to succeed, got error: %v", err)
	}

	if attempts != 2 {
		t.Errorf("expected 2 attempts (1 failure + 1 retry), got %d", attempts)
	}
	if entity.Role != "Backend Engineer" {
		t.Errorf("expected role 'Backend Engineer', got '%s'", entity.Role)
	}
	if entity.City != "Pune" {
		t.Errorf("expected city 'Pune', got '%s'", entity.City)
	}
}

func TestGeminiExtractor_503PersistentFallback(t *testing.T) {
	client := &mockHTTPClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusServiceUnavailable,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(bytes.NewBufferString(`{"error":{"message":"high demand spike","code":"service_unavailable"}}`)),
			}, nil
		},
	}

	extractor := llm.NewGeminiExtractor(
		"AQ.TestPersistent503",
		llm.WithHTTPClient(client),
		llm.WithModel("gemini-3.8-flash"),
	)

	// Deliberate persistent 503 should gracefully activate rule fallback without returning an error
	entity, err := extractor.Extract(context.Background(), "React UI Intern", "Mumbai office")
	if err != nil {
		t.Fatalf("expected fallback extractor to succeed, got error: %v", err)
	}

	if entity.Role != "Frontend Engineer" {
		t.Errorf("expected fallback role 'Frontend Engineer', got '%s'", entity.Role)
	}
	if entity.City != "Mumbai" {
		t.Errorf("expected fallback city 'Mumbai', got '%s'", entity.City)
	}
}
