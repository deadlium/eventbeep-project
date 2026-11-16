package scraper

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/deadlium/eventbeep/backend/internal/models"
)

// ComputeHashKey computes a deterministic 64-character hex SHA-256 hash
// of lowercase(company_name + job_title). This satisfies the EventBeep deduplication
// contract for idempotent PostgreSQL upserts on jobs.hash_key.
func ComputeHashKey(company, title string) string {
	normalized := strings.ToLower(strings.TrimSpace(company) + strings.TrimSpace(title))
	hash := sha256.Sum256([]byte(normalized))
	return hex.EncodeToString(hash[:])
}

// ComputeHashKeyWithURL computes a deterministic 64-character hex SHA-256 hash
// incorporating source_url: lowercase(company + title + source_url).
func ComputeHashKeyWithURL(company, title, sourceURL string) string {
	normalized := strings.ToLower(strings.TrimSpace(company) + strings.TrimSpace(title) + strings.TrimSpace(sourceURL))
	hash := sha256.Sum256([]byte(normalized))
	return hex.EncodeToString(hash[:])
}

// GenerateJobHash populates or returns the HashKey for a models.Job instance.
func GenerateJobHash(job *models.Job) string {
	if job == nil {
		return ""
	}
	hash := ComputeHashKey(job.Company, job.Title)
	job.HashKey = hash
	return hash
}
