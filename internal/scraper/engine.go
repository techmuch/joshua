package scraper

import (
	"context"
	"log/slog"
	"sync"
)

// Engine manages the execution of multiple scrapers
type Engine struct {
	scrapers []Scraper
}

// NewEngine creates a new scraper engine
func NewEngine() *Engine {
	return &Engine{
		scrapers: make([]Scraper, 0),
	}
}

// Register adds a scraper to the engine
func (e *Engine) Register(s Scraper) {
	e.scrapers = append(e.scrapers, s)
}

// Run executes all registered scrapers
func (e *Engine) Run(ctx context.Context) ([]Solicitation, error) {
	var allSolicitations []Solicitation
	var mu sync.Mutex
	var wg sync.WaitGroup

	slog.Info("Starting Scraper Engine", "source_count", len(e.scrapers))

	for _, s := range e.scrapers {
		wg.Add(1)
		go func(scraper Scraper) {
			defer wg.Done()
			slog.Info("Running scraper", "scraper", scraper.Name())
			
			results, err := scraper.Scrape(ctx)
			if err != nil {
				slog.Error("Scraper failed", "scraper", scraper.Name(), "error", err)
				return
			}

			mu.Lock()
			allSolicitations = append(allSolicitations, results...)
			mu.Unlock()
			
			slog.Info("Scraper finished", "scraper", scraper.Name(), "items_found", len(results))
		}(s)
	}

	wg.Wait()
	
	slog.Info("Scraper Engine run complete", "total_items", len(allSolicitations))

	return allSolicitations, nil
}