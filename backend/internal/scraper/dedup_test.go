package scraper_test

import (
	"encoding/hex"
	"testing"

	"github.com/deadlium/eventbeep/backend/internal/models"
	"github.com/deadlium/eventbeep/backend/internal/scraper"
)

func TestComputeHashKey_Deterministic(t *testing.T) {
	company := "Razorpay"
	title := "Software Engineer - Fresher"

	hash1 := scraper.ComputeHashKey(company, title)
	hash2 := scraper.ComputeHashKey(company, title)

	if hash1 != hash2 {
		t.Fatalf("expected identical hash for same inputs, got %s vs %s", hash1, hash2)
	}

	if len(hash1) != 64 {
		t.Errorf("expected 64-character hex string, got length %d (%s)", len(hash1), hash1)
	}

	// Verify valid hex decoding
	decoded, err := hex.DecodeString(hash1)
	if err != nil || len(decoded) != 32 {
		t.Fatalf("expected 32-byte sha256 output, got error: %v", err)
	}
}

func TestComputeHashKey_CaseInsensitive(t *testing.T) {
	hashLower := scraper.ComputeHashKey("razorpay", "software engineer - fresher")
	hashUpper := scraper.ComputeHashKey("RAZORPAY", "SOFTWARE ENGINEER - FRESHER")
	hashMixed := scraper.ComputeHashKey("  RazorPay  ", "  Software Engineer - Fresher  ")

	if hashLower != hashUpper {
		t.Errorf("expected case-insensitive hash match: %s != %s", hashLower, hashUpper)
	}
	if hashLower != hashMixed {
		t.Errorf("expected whitespace-trimmed hash match: %s != %s", hashLower, hashMixed)
	}
}

func TestComputeHashKey_Distinct(t *testing.T) {
	h1 := scraper.ComputeHashKey("Swiggy", "Frontend Engineer")
	h2 := scraper.ComputeHashKey("Zomato", "Frontend Engineer")
	h3 := scraper.ComputeHashKey("Swiggy", "Backend Engineer")

	if h1 == h2 || h1 == h3 || h2 == h3 {
		t.Errorf("expected distinct hashes for distinct job listings, got h1=%s, h2=%s, h3=%s", h1, h2, h3)
	}
}

func TestComputeHashKeyWithURL(t *testing.T) {
	h1 := scraper.ComputeHashKeyWithURL("Cred", "Go Developer", "https://example.com/job/1")
	h2 := scraper.ComputeHashKeyWithURL("cred", "go developer", "https://example.com/job/1")
	h3 := scraper.ComputeHashKeyWithURL("Cred", "Go Developer", "https://example.com/job/2")

	if h1 != h2 {
		t.Errorf("expected case-insensitive match for hash with URL")
	}
	if h1 == h3 {
		t.Errorf("expected different hash when source URL differs")
	}
	if len(h1) != 64 {
		t.Errorf("expected 64-char hex string, got %d", len(h1))
	}
}

func TestGenerateJobHash(t *testing.T) {
	job := &models.Job{
		Company: "Postman",
		Title:   "Junior Backend Developer",
	}

	hash := scraper.GenerateJobHash(job)
	if hash == "" {
		t.Fatal("expected non-empty hash")
	}
	if job.HashKey != hash {
		t.Errorf("expected job.HashKey to be populated, got %s vs %s", job.HashKey, hash)
	}
	if scraper.GenerateJobHash(nil) != "" {
		t.Error("expected empty string for nil job")
	}
}
