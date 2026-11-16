package scraper

import (
	"context"
	"fmt"
	"io"

	"github.com/PuerkitoBio/goquery"
)

const DefaultWellfoundURL = "https://wellfound.com/jobs"

// WellfoundScraper scrapes startup fresher/junior roles from Wellfound.
type WellfoundScraper struct {
	client    HTTPClient
	sourceURL string
}

// NewWellfoundScraper creates a Wellfound scraper instance.
func NewWellfoundScraper(client HTTPClient, customURL ...string) *WellfoundScraper {
	target := DefaultWellfoundURL
	if len(customURL) > 0 && customURL[0] != "" {
		target = customURL[0]
	}
	return &WellfoundScraper{
		client:    client,
		sourceURL: target,
	}
}

func (s *WellfoundScraper) Name() string {
	return "wellfound"
}

func (s *WellfoundScraper) SourceURL() string {
	return s.sourceURL
}

func (s *WellfoundScraper) Scrape(ctx context.Context) ([]ScrapedJob, error) {
	body, err := FetchHTML(ctx, s.client, s.sourceURL)
	if err != nil {
		return nil, fmt.Errorf("wellfound scraper fetch error: %w", err)
	}
	defer body.Close()

	return s.ParseHTML(body, s.sourceURL)
}

// ParseHTML extracts jobs from Wellfound HTML content.
func (s *WellfoundScraper) ParseHTML(r io.Reader, baseURL string) ([]ScrapedJob, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("failed to parse wellfound HTML document: %w", err)
	}

	var jobs []ScrapedJob

	cardSelectors := "[data-test='JobListingCard'], .job-listing, article.job-card, .listing-card, div[class*='JobListing']"
	doc.Find(cardSelectors).Each(func(i int, card *goquery.Selection) {
		title := cleanText(card.Find("[data-test='job-title'], .job-title, h2, h3, a[class*='jobTitle']").First().Text())
		company := cleanText(card.Find("[data-test='company-name'], .company-name, h4, .startup-name, span[class*='company']").First().Text())
		city := cleanText(card.Find("[data-test='job-location'], .location, .locations, span[class*='location']").First().Text())
		role := cleanText(card.Find("[data-test='job-role'], .role, .role-tag, span[class*='role']").First().Text())
		rawDesc := cleanText(card.Find(".description, [data-test='job-description'], p").First().Text())

		link := ""
		if href, exists := card.Find("a[href*='/jobs/'], a[data-test='job-link'], a[href]").First().Attr("href"); exists {
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
