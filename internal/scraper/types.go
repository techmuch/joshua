package scraper

import (
	"context"
	"time"
)

// Solicitation represents a scraped business opportunity
type Solicitation struct {
	SourceID    string                 `json:"source_id"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Agency      string                 `json:"agency"`
	DueDate     time.Time              `json:"due_date"`
	URL         string                 `json:"url"`
	RawData     map[string]interface{} `json:"raw_data"`
}

// Scraper defines the interface that all site-specific scrapers must implement
type Scraper interface {
	// Name returns the unique identifier for this scraper (e.g., "georgia-gpr")
	Name() string
	
	// Scrape executes the scraping logic and returns a list of found solicitations
	Scrape(ctx context.Context) ([]Solicitation, error)
}
