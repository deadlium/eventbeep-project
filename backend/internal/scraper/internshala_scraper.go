package scraper

import (
	"context"
	"fmt"
	"io"

	"github.com/PuerkitoBio/goquery"
)

const DefaultInternshalaURL = "https://internshala.com/jobs/fresher-jobs/"

// InternshalaScraper scrapes fresher jobs and internships from Internshala.
type InternshalaScraper struct {
	client    HTTPClient
	sourceURL string
}

// NewInternshalaScraper creates an Internshala scraper instance.
func NewInternshalaScraper(client HTTPClient, customURL ...string) *InternshalaScraper {
	target := DefaultInternshalaURL
	if len(customURL) > 0 && customURL[0] != "" {
		target = customURL[0]
	}
	return &InternshalaScraper{
		client:    client,
		sourceURL: target,
	}
}

func (s *InternshalaScraper) Name() string {
	return "internshala"
}

func (s *InternshalaScraper) SourceURL() string {
	return s.sourceURL
}

func (s *InternshalaScraper) Scrape(ctx context.Context) ([]ScrapedJob, error) {
	body, err := FetchHTML(ctx, s.client, s.sourceURL)
	if err != nil {
		return nil, fmt.Errorf("internshala scraper fetch error: %w", err)
	}
	defer body.Close()

	return s.ParseHTML(body, s.sourceURL)
}

// ParseHTML extracts jobs from Internshala HTML content.
func (s *InternshalaScraper) ParseHTML(r io.Reader, baseURL string) ([]ScrapedJob, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("failed to parse internshala HTML document: %w", err)
	}

	var jobs []ScrapedJob

	cardSelectors := ".individual_internship, .internship_meta, .job-container, article.job"
	doc.Find(cardSelectors).Each(func(i int, card *goquery.Selection) {
		title := cleanText(card.Find(".job-title-href, .profile, h3.heading_4_5, .job_title").First().Text())
		company := cleanText(card.Find(".company_name, .link_display_like_text, .company-name").First().Text())
		city := cleanText(card.Find(".location_link, #location_names, .locations, .row-1-item.locations").First().Text())
		role := cleanText(card.Find(".profile, .status-info, .role").First().Text())
		rawDesc := cleanText(card.Find(".internship_other_details, .other_detail_item_row, p").First().Text())

		link := ""
		if href, exists := card.Find("a.view_detail_button, a.job-title-href, a[href*='/job/detail/'], a[href*='/internship/detail/'], a[href]").First().Attr("href"); exists {
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
