package georgia

import (
	"bd_bot/internal/scraper"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// GPRScraper handles the Georgia Procurement Registry
type GPRScraper struct {
	BaseURL string
}

// NewGPRScraper creates a new instance
func NewGPRScraper() *GPRScraper {
	return &GPRScraper{
		BaseURL: "https://ssl.doas.state.ga.us/gpr/index",
	}
}

func (s *GPRScraper) Name() string {
	return "Georgia Procurement Registry (GPR)"
}

func (s *GPRScraper) Scrape(ctx context.Context) ([]scraper.Solicitation, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", s.BaseURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch GPR: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GPR returned status: %s", resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	title := doc.Find("title").Text()
	slog.Debug("Connected to GPR", "title", title)

	return []scraper.Solicitation{}, nil
}