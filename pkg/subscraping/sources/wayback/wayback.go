package wayback

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/darshakkanani/stormfinder/v2/pkg/subscraping"
)

// Source is the passive scraping agent
type Source struct {
	apiKeys   []string
	timeTaken time.Duration
	errors    int
	results   int
}

type waybackResponse [][]string

// Run function returns all subdomains found with the service
func (s *Source) Run(ctx context.Context, domain string, session *subscraping.Session) <-chan subscraping.Result {
	results := make(chan subscraping.Result)
	s.errors = 0
	s.results = 0

	go func() {
		defer func(startTime time.Time) {
			s.timeTaken = time.Since(startTime)
			close(results)
		}(time.Now())

		// Query Wayback Machine CDX API
		searchURL := fmt.Sprintf("http://web.archive.org/cdx/search/cdx?url=*.%s/*&output=json&collapse=urlkey&fl=original", domain)
		
		resp, err := session.SimpleGet(ctx, searchURL)
		if err != nil {
			s.errors++
			results <- subscraping.Result{Source: s.Name(), Type: subscraping.Error, Error: err}
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			s.errors++
			results <- subscraping.Result{Source: s.Name(), Type: subscraping.Error, Error: fmt.Errorf("unexpected status code: %d", resp.StatusCode)}
			return
		}

		var response waybackResponse
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			s.errors++
			results <- subscraping.Result{Source: s.Name(), Type: subscraping.Error, Error: err}
			return
		}

		subdomains := make(map[string]struct{})
		
		// Regex to extract subdomains from URLs
		subdomainRegex := regexp.MustCompile(`https?://([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+` + regexp.QuoteMeta(domain))
		
		for _, entry := range response {
			if len(entry) > 0 {
				url := entry[0]
				matches := subdomainRegex.FindStringSubmatch(url)
				if len(matches) > 0 {
					// Extract the full subdomain
					fullMatch := matches[0]
					// Remove protocol
					subdomain := strings.TrimPrefix(fullMatch, "https://")
					subdomain = strings.TrimPrefix(subdomain, "http://")
					
					// Extract just the domain part (remove path)
					if idx := strings.Index(subdomain, "/"); idx != -1 {
						subdomain = subdomain[:idx]
					}
					
					if strings.HasSuffix(subdomain, "."+domain) || subdomain == domain {
						subdomains[subdomain] = struct{}{}
					}
				}
			}
		}

		for subdomain := range subdomains {
			s.results++
			results <- subscraping.Result{Source: s.Name(), Type: subscraping.Subdomain, Value: subdomain}
		}
	}()

	return results
}

// Name returns the name of the source
func (s *Source) Name() string {
	return "wayback"
}

func (s *Source) IsDefault() bool {
	return true
}

func (s *Source) HasRecursiveSupport() bool {
	return false
}

func (s *Source) NeedsKey() bool {
	return false
}

func (s *Source) AddApiKeys(keys []string) {
	s.apiKeys = keys
}

func (s *Source) Statistics() subscraping.Statistics {
	return subscraping.Statistics{
		Errors:    s.errors,
		Results:   s.results,
		TimeTaken: s.timeTaken,
	}
}
