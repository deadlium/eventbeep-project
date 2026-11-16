package llm

import (
	"context"
	"regexp"
	"strings"
)

// FallbackExtractor provides heuristic, rule-based role and city extraction
// when the LLM API is unavailable, unconfigured, or encounters transient failures.
type FallbackExtractor struct{}

// NewFallbackExtractor creates a new rule-based fallback extractor.
func NewFallbackExtractor() *FallbackExtractor {
	return &FallbackExtractor{}
}

// Pre-compiled regex patterns for roles
var (
	reFrontend   = regexp.MustCompile(`(?i)\b(frontend|front-end|front\s+end|react|vue|angular|ui\s+developer)\b`)
	reBackend    = regexp.MustCompile(`(?i)\b(backend|back-end|back\s+end|golang|go\s+developer|node\.?js|python\s+developer|java\s+developer)\b`)
	reFullstack  = regexp.MustCompile(`(?i)\b(fullstack|full-stack|full\s+stack)\b`)
	reDevOps     = regexp.MustCompile(`(?i)\b(devops|dev-ops|sre|site\s+reliability|cloud\s+engineer|infrastructure)\b`)
	reDataSci    = regexp.MustCompile(`(?i)\b(data\s+scientist|machine\s+learning|ml\s+engineer|ai\s+engineer|deep\s+learning)\b`)
	reDataAnal   = regexp.MustCompile(`(?i)\b(data\s+analyst|business\s+analyst|data\s+analytics|bi\s+analyst)\b`)
	reQA         = regexp.MustCompile(`(?i)\b(qa|quality\s+assurance|sdet|test\s+engineer|software\s+tester)\b`)
	reMobile     = regexp.MustCompile(`(?i)\b(mobile|android|ios|flutter|react\s+native)\b`)
	reProduct    = regexp.MustCompile(`(?i)\b(product\s+manager|associate\s+product\s+manager|apm)\b`)
	reGrowth     = regexp.MustCompile(`(?i)\b(growth|marketing|operations)\b`)
	reSoftwareEng = regexp.MustCompile(`(?i)\b(software\s+engineer|software\s+developer|sde|sde[- ]?1|programmer|developer)\b`)
)

// Pre-compiled regex patterns for Indian tech hubs & Remote
var (
	reBangalore = regexp.MustCompile(`(?i)\b(bangalore|bengaluru)\b`)
	reMumbai    = regexp.MustCompile(`(?i)\b(mumbai|bombay|navi\s+mumbai|thane)\b`)
	reDelhiNCR  = regexp.MustCompile(`(?i)\b(delhi|new\s+delhi|ncr|gurgaon|gurugram|noida)\b`)
	reHyderabad = regexp.MustCompile(`(?i)\b(hyderabad|secunderabad)\b`)
	rePune      = regexp.MustCompile(`(?i)\b(pune)\b`)
	reChennai   = regexp.MustCompile(`(?i)\b(chennai|madras)\b`)
	reKolkata   = regexp.MustCompile(`(?i)\b(kolkata|calcutta)\b`)
	reAhmedabad = regexp.MustCompile(`(?i)\b(ahmedabad)\b`)
	reRemote    = regexp.MustCompile(`(?i)\b(remote|work\s+from\s+home|wfh|anywhere)\b`)
)

func (f *FallbackExtractor) Extract(ctx context.Context, title, description string) (*ExtractedEntity, error) {
	combined := title + " " + description

	role := extractRole(title, combined)
	city := extractCity(combined)

	return &ExtractedEntity{
		Role: role,
		City: city,
	}, nil
}

func extractRole(title, combined string) string {
	// First check title priority
	switch {
	case reFullstack.MatchString(title):
		return "Full Stack Engineer"
	case reFrontend.MatchString(title):
		return "Frontend Engineer"
	case reBackend.MatchString(title):
		return "Backend Engineer"
	case reDevOps.MatchString(title):
		return "DevOps Engineer"
	case reDataSci.MatchString(title):
		return "Data Scientist"
	case reDataAnal.MatchString(title):
		return "Data Analyst"
	case reQA.MatchString(title):
		return "QA Engineer"
	case reMobile.MatchString(title):
		return "Mobile Developer"
	case reProduct.MatchString(title):
		return "Product Manager"
	case reGrowth.MatchString(title):
		return "Growth Associate"
	case reSoftwareEng.MatchString(title):
		return "Software Engineer"
	}

	// Fallback to combined text
	switch {
	case reFullstack.MatchString(combined):
		return "Full Stack Engineer"
	case reFrontend.MatchString(combined):
		return "Frontend Engineer"
	case reBackend.MatchString(combined):
		return "Backend Engineer"
	case reDevOps.MatchString(combined):
		return "DevOps Engineer"
	case reDataSci.MatchString(combined):
		return "Data Scientist"
	case reDataAnal.MatchString(combined):
		return "Data Analyst"
	case reQA.MatchString(combined):
		return "QA Engineer"
	case reMobile.MatchString(combined):
		return "Mobile Developer"
	case reSoftwareEng.MatchString(combined):
		return "Software Engineer"
	}

	// Default sanitized title if reasonable length, or generic fallback
	clean := strings.TrimSpace(title)
	if len(clean) > 0 && len(clean) <= 60 {
		return clean
	}
	return "Software Engineer"
}

func extractCity(combined string) string {
	switch {
	case reRemote.MatchString(combined):
		return "Remote"
	case reBangalore.MatchString(combined):
		return "Bangalore"
	case reMumbai.MatchString(combined):
		return "Mumbai"
	case reDelhiNCR.MatchString(combined):
		return "Delhi NCR"
	case reHyderabad.MatchString(combined):
		return "Hyderabad"
	case rePune.MatchString(combined):
		return "Pune"
	case reChennai.MatchString(combined):
		return "Chennai"
	case reKolkata.MatchString(combined):
		return "Kolkata"
	case reAhmedabad.MatchString(combined):
		return "Ahmedabad"
	}

	return "Remote"
}
