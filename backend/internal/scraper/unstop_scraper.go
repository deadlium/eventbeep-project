package scraper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

const (
	DefaultUnstopURL = "https://unstop.com/api/public/opportunity/search-result?opportunity=jobs&per_page=20"
	DefaultUnstopWeb = "https://unstop.com/jobs"
)

// Default search keywords to rotate across cycles to constantly discover new listings
var defaultUnstopKeywords = []string{
	"developer",
	"software",
	"fresher",
	"intern",
	"backend",
	"frontend",
	"python",
	"data",
	"react",
	"fullstack",
}

// UnstopScraper scrapes fresher/entry-level job opportunities from Unstop.
type UnstopScraper struct {
	client       HTTPClient
	sourceURL    string
	isCustomURL  bool
	keywordIndex int
	pageIndex    int
	mu           sync.Mutex
}

// NewUnstopScraper creates an Unstop scraper instance.
func NewUnstopScraper(client HTTPClient, customURL ...string) *UnstopScraper {
	target := DefaultUnstopURL
	isCustom := false
	if len(customURL) > 0 && customURL[0] != "" {
		target = customURL[0]
		isCustom = true
	}
	return &UnstopScraper{
		client:       client,
		sourceURL:    target,
		isCustomURL:  isCustom,
		keywordIndex: 0,
		pageIndex:    1,
	}
}

func (s *UnstopScraper) Name() string {
	return "unstop"
}

func (s *UnstopScraper) SourceURL() string {
	return s.sourceURL
}

func (s *UnstopScraper) nextTargetURL() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isCustomURL {
		return s.sourceURL
	}

	keyword := defaultUnstopKeywords[s.keywordIndex%len(defaultUnstopKeywords)]
	page := s.pageIndex

	// Increment page, wrap to next keyword every 3 pages
	s.pageIndex++
	if s.pageIndex > 3 {
		s.pageIndex = 1
		s.keywordIndex++
	}

	return fmt.Sprintf("https://unstop.com/api/public/opportunity/search-result?opportunity=jobs&searchTerm=%s&page=%d&per_page=20",
		url.QueryEscape(keyword), page)
}

func (s *UnstopScraper) Scrape(ctx context.Context) ([]ScrapedJob, error) {
	targetURL := s.nextTargetURL()

	body, err := FetchHTML(ctx, s.client, targetURL)
	if err != nil {
		return nil, fmt.Errorf("unstop scraper fetch error: %w", err)
	}
	defer body.Close()

	// Read content into buffer to inspect if JSON or HTML
	data, err := io.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("failed to read unstop response: %w", err)
	}

	trimmed := bytes.TrimSpace(data)
	if len(trimmed) > 0 && trimmed[0] == '{' {
		jobs, jsonErr := s.ParseJSON(bytes.NewReader(trimmed), DefaultUnstopWeb)
		if jsonErr == nil && len(jobs) > 0 {
			log.Printf("[unstop] Scraped %d fresh opportunities from live API: %s", len(jobs), targetURL)
			return jobs, nil
		}
	}

	// Fallback to HTML parser
	return s.ParseHTML(bytes.NewReader(trimmed), targetURL)
}

// UnstopAPIResponse maps the Unstop search-result JSON structure.
type unstopAPIResponse struct {
	Data struct {
		CurrentPage int `json:"current_page"`
		Data        []struct {
			ID           int    `json:"id"`
			Title        string `json:"title"`
			PublicURL    string `json:"public_url"`
			SEOURL       string `json:"seo_url"`
			Details      string `json:"details"`
			Organisation struct {
				Name string `json:"name"`
			} `json:"organisation"`
			JobDetail struct {
				Locations []string `json:"locations"`
			} `json:"jobDetail"`
			Locations []struct {
				City  string `json:"city"`
				State string `json:"state"`
			} `json:"locations"`
			Workfunction []struct {
				Name string `json:"name"`
			} `json:"workfunction"`
		} `json:"data"`
	} `json:"data"`
}

// ParseJSON extracts jobs from the Unstop JSON API response.
func (s *UnstopScraper) ParseJSON(r io.Reader, baseURL string) ([]ScrapedJob, error) {
	var resp unstopAPIResponse
	if err := json.NewDecoder(r).Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode unstop API JSON: %w", err)
	}

	var jobs []ScrapedJob
	for _, item := range resp.Data.Data {
		title := cleanText(item.Title)
		company := cleanText(item.Organisation.Name)
		if title == "" || company == "" {
			continue
		}

		city := ""
		if len(item.JobDetail.Locations) > 0 {
			city = cleanText(item.JobDetail.Locations[0])
		} else if len(item.Locations) > 0 {
			city = cleanText(item.Locations[0].City)
		}

		role := ""
		if len(item.Workfunction) > 0 {
			role = cleanText(item.Workfunction[0].Name)
		}

		link := item.SEOURL
		if link == "" && item.PublicURL != "" {
			link = resolveURL("https://unstop.com/", item.PublicURL)
		}
		if link == "" {
			link = baseURL
		}

		rawDesc := cleanHTMLText(item.Details)

		job := ScrapedJob{
			Title:          title,
			Company:        company,
			RawDescription: rawDesc,
			Role:           role,
			City:           city,
			SourceURL:      link,
			HashKey:        ComputeHashKey(company, title),
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}

// ParseHTML extracts jobs from an HTML stream.
func (s *UnstopScraper) ParseHTML(r io.Reader, baseURL string) ([]ScrapedJob, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("failed to parse unstop HTML document: %w", err)
	}

	var jobs []ScrapedJob

	cardSelectors := ".opportunity-card, .job-card, .single-item, article.job-item, div[class*='opportunity_card']"
	doc.Find(cardSelectors).Each(func(i int, card *goquery.Selection) {
		title := cleanText(card.Find("h3, h2, .job-title, .title, a[class*='title']").First().Text())
		company := cleanText(card.Find(".company-name, .organisation-name, .org-name, div[class*='organisation']").First().Text())
		city := cleanText(card.Find(".job-location, .location, span[class*='location']").First().Text())
		role := cleanText(card.Find(".badge, .job-role, span[class*='role'], .category").First().Text())
		rawDesc := cleanText(card.Find(".description, .job-desc, p").First().Text())

		link := ""
		if href, exists := card.Find("a[href]").First().Attr("href"); exists {
			link = resolveURL(baseURL, href)
		} else if href, exists := card.Attr("href"); exists {
			link = resolveURL(baseURL, href)
		}

		if title == "" || company == "" {
			return
		}

		if link == "" {
			link = baseURL
		}

		job := ScrapedJob{
			Title:          title,
			Company:        company,
			RawDescription: rawDesc,
			Role:           role,
			City:           city,
			SourceURL:      link,
			HashKey:        ComputeHashKey(company, title),
		}
		jobs = append(jobs, job)
	})

	return jobs, nil
}

func cleanText(s string) string {
	parts := strings.Fields(strings.TrimSpace(s))
	return strings.Join(parts, " ")
}

func cleanHTMLText(htmlStr string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	if err != nil {
		return cleanText(htmlStr)
	}
	return cleanText(doc.Text())
}

func resolveURL(baseStr, refStr string) string {
	refStr = strings.TrimSpace(refStr)
	if refStr == "" {
		return baseStr
	}

	base, err := url.Parse(baseStr)
	if err != nil {
		return refStr
	}
	ref, err := url.Parse(refStr)
	if err != nil {
		return refStr
	}
	return base.ResolveReference(ref).String()
}
