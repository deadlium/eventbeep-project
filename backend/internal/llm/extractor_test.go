package llm_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/deadlium/eventbeep/backend/internal/llm"
)

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
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("key") != "secret-test-key" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
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

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer mockServer.Close()

	extractor := llm.NewGeminiExtractor(
		"secret-test-key",
		llm.WithBaseURL(mockServer.URL),
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
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer mockServer.Close()

	extractor := llm.NewGeminiExtractor(
		"secret-test-key",
		llm.WithBaseURL(mockServer.URL),
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
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "server error", http.StatusInternalServerError)
	}))
	defer mockServer.Close()

	extractor := llm.NewGeminiExtractor(
		"secret-test-key",
		llm.WithBaseURL(mockServer.URL),
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

	// Test server returning invalid JSON body (corrupt candidates)
	corruptServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"candidates": [{"content": {"parts": [{"text": "not-valid-json"}]}}]}`))
	}))
	defer corruptServer.Close()

	eCorrupt := llm.NewGeminiExtractor("test-key", llm.WithBaseURL(corruptServer.URL))
	res, err := eCorrupt.Extract(context.Background(), "Backend Developer", "Bangalore")
	if err != nil {
		t.Fatalf("expected fallback on corrupt JSON, got error: %v", err)
	}
	if res.Role != "Backend Engineer" {
		t.Errorf("expected fallback role 'Backend Engineer', got %s", res.Role)
	}

	// Test server returning API error in payload
	errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"error": {"code": 403, "message": "API key expired"}}`))
	}))
	defer errorServer.Close()

	eErr := llm.NewGeminiExtractor("test-key", llm.WithBaseURL(errorServer.URL))
	resErr, err := eErr.Extract(context.Background(), "QA Tester", "Hyderabad")
	if err != nil {
		t.Fatalf("expected fallback on API error payload, got error: %v", err)
	}
	if resErr.Role != "QA Engineer" {
		t.Errorf("expected fallback QA Engineer, got %s", resErr.Role)
	}
}
