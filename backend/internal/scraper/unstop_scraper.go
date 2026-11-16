package scraper

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const DefaultUnstopURL = "https://unstop.com/jobs?experience=0"

// UnstopScraper scrapes fresher/entry-level job opportunities from Unstop.
type UnstopScraper struct {
	client    HTTPClient
	sourceURL string
}

// NewUnstopScraper creates an Unstop scraper instance.
func NewUnstopScraper(client HTTPClient, customURL ...string) *UnstopScraper {
	target := DefaultUnstopURL
	if len(customURL) > 0 && customURL[0] != "" {
		target = customURL[0]
	}
	return &UnstopScraper{
		client:    client,
		sourceURL: target,
	}
}

func (s *UnstopScraper) Name() string {
	return "unstop"
}

func (s *UnstopScraper) SourceURL() string {
	return s.sourceURL
}

func (s *UnstopScraper) Scrape(ctx context.Context) ([]ScrapedJob, error) {
	body, err := FetchHTML(ctx, s.client, s.sourceURL)
	if err != nil {
		return nil, fmt.Errorf("unstop scraper fetch error: %w", err)
	}
	defer body.Close()

	return s.ParseHTML(body, s.sourceURL)
}

// ParseHTML extracts jobs from an HTML stream.
func (s *UnstopScraper) ParseHTML(r io.Reader, baseURL string) ([]ScrapedJob, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("failed to parse unstop HTML document: %w", err)
	}

	var jobs []ScrapedJob

	// Query opportunity / job cards
	cardSelectors := ".opportunity-card, .job-card, .single-item, article.job-item, div[class*='opportunity_card']"
	doc.Find(cardSelectors).Each(func(i int, card *goquery.Selection) {
		title := cleanText(card.Find("h3, h2, .job-title, .title, a[class*='title']").First().Text())
		company := cleanText(card.Find(".company-name, .organisation-name, .org-name, div[class*='organisation']").First().Text())
		city := cleanText(card.Find(".job-location, .location, span[class*='location']").First().Text())
		role := cleanText(card.Find(".badge, .job-role, span[class*='role'], .category").First().Text())
		rawDesc := cleanText(card.Find(".description, .job-desc, p").First().Text())

		// Extract link
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
