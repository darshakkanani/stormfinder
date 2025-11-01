package urlscan

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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

type urlscanResponse struct {
	Results []struct {
		Page struct {
			Domain string `json:"domain"`
		} `json:"page"`
		Task struct {
			Domain string `json:"domain"`
		} `json:"task"`
	} `json:"results"`
}

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

		// Search for subdomains in URLScan.io
		searchURL := fmt.Sprintf("https://urlscan.io/api/v1/search/?q=domain:%s&size=10000", domain)
		
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

		var response urlscanResponse
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			s.errors++
			results <- subscraping.Result{Source: s.Name(), Type: subscraping.Error, Error: err}
			return
		}

		subdomains := make(map[string]struct{})
		for _, result := range response.Results {
			if result.Page.Domain != "" && strings.HasSuffix(result.Page.Domain, "."+domain) {
				subdomains[result.Page.Domain] = struct{}{}
			}
			if result.Task.Domain != "" && strings.HasSuffix(result.Task.Domain, "."+domain) {
				subdomains[result.Task.Domain] = struct{}{}
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
	return "urlscan"
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
