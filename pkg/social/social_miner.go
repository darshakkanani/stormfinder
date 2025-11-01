package social

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/projectdiscovery/gologger"
)

// SocialMiner mines subdomains from social media and code repositories
type SocialMiner struct {
	client      *http.Client
	rateLimiter map[string]*time.Ticker
	apiKeys     map[string]string
	cache       map[string][]SocialResult
	cacheMutex  sync.RWMutex
}

// SocialResult represents a result from social media mining
type SocialResult struct {
	Subdomain   string    `json:"subdomain"`
	Source      string    `json:"source"`
	Platform    string    `json:"platform"`
	URL         string    `json:"url"`
	Context     string    `json:"context"`
	Confidence  float64   `json:"confidence"`
	Timestamp   time.Time `json:"timestamp"`
	Author      string    `json:"author"`
	Engagement  int       `json:"engagement"`
}

// SocialMiningResult contains comprehensive social mining results
type SocialMiningResult struct {
	Subdomains      []string                 `json:"subdomains"`
	Results         []SocialResult           `json:"results"`
	PlatformStats   map[string]int           `json:"platform_stats"`
	ConfidenceStats map[string]int           `json:"confidence_stats"`
	Timeline        []TimelineEntry          `json:"timeline"`
	TopAuthors      []AuthorStats            `json:"top_authors"`
	Statistics      SocialMiningStatistics   `json:"statistics"`
}

// TimelineEntry represents social media timeline
type TimelineEntry struct {
	Date        time.Time `json:"date"`
	Platform    string    `json:"platform"`
	Subdomains  []string  `json:"subdomains"`
	Activity    string    `json:"activity"`
}

// AuthorStats contains author statistics
type AuthorStats struct {
	Author     string `json:"author"`
	Platform   string `json:"platform"`
	Mentions   int    `json:"mentions"`
	Subdomains int    `json:"subdomains"`
}

// SocialMiningStatistics contains mining statistics
type SocialMiningStatistics struct {
	TotalResults      int           `json:"total_results"`
	UniqueSubdomains  int           `json:"unique_subdomains"`
	PlatformsQueried  int           `json:"platforms_queried"`
	ProcessingTime    time.Duration `json:"processing_time"`
	HighConfidence    int           `json:"high_confidence"`
	MediumConfidence  int           `json:"medium_confidence"`
	LowConfidence     int           `json:"low_confidence"`
}

// NewSocialMiner creates a new social media miner
func NewSocialMiner() *SocialMiner {
	return &SocialMiner{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		rateLimiter: map[string]*time.Ticker{
			"github":    time.NewTicker(1 * time.Second),   // 1 request per second
			"twitter":   time.NewTicker(2 * time.Second),   // 0.5 requests per second
			"reddit":    time.NewTicker(500 * time.Millisecond), // 2 requests per second
			"pastebin":  time.NewTicker(200 * time.Millisecond), // 5 requests per second
			"gitlab":    time.NewTicker(1 * time.Second),   // 1 request per second
		},
		apiKeys: make(map[string]string),
		cache:   make(map[string][]SocialResult),
	}
}

// SetAPIKey sets an API key for a specific platform
func (sm *SocialMiner) SetAPIKey(platform, key string) {
	sm.apiKeys[platform] = key
}

// MineFromSocialMedia performs comprehensive social media mining
func (sm *SocialMiner) MineFromSocialMedia(ctx context.Context, domain string) (*SocialMiningResult, error) {
	startTime := time.Now()
	gologger.Info().Msgf("📱 Social Mining: Starting comprehensive analysis for %s", domain)
	
	result := &SocialMiningResult{
		Subdomains:      make([]string, 0),
		Results:         make([]SocialResult, 0),
		PlatformStats:   make(map[string]int),
		ConfidenceStats: make(map[string]int),
		Timeline:        make([]TimelineEntry, 0),
		TopAuthors:      make([]AuthorStats, 0),
		Statistics:      SocialMiningStatistics{},
	}
	
	var wg sync.WaitGroup
	var mutex sync.Mutex
	
	// Mine from different platforms in parallel
	platforms := []string{"github", "gitlab", "pastebin", "reddit", "twitter", "stackoverflow"}
	
	for _, platform := range platforms {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			
			gologger.Debug().Msgf("Mining from platform: %s", p)
			platformResults, err := sm.mineFromPlatform(ctx, p, domain)
			if err != nil {
				gologger.Warning().Msgf("Failed to mine from %s: %v", p, err)
				return
			}
			
			mutex.Lock()
			for _, res := range platformResults {
				result.Results = append(result.Results, res)
				result.PlatformStats[p]++
				
				// Add unique subdomains
				if !contains(result.Subdomains, res.Subdomain) {
					result.Subdomains = append(result.Subdomains, res.Subdomain)
				}
				
				// Categorize by confidence
				if res.Confidence >= 0.8 {
					result.ConfidenceStats["high"]++
				} else if res.Confidence >= 0.5 {
					result.ConfidenceStats["medium"]++
				} else {
					result.ConfidenceStats["low"]++
				}
			}
			mutex.Unlock()
			
			gologger.Debug().Msgf("Found %d results from %s", len(platformResults), p)
		}(platform)
	}
	
	wg.Wait()
	
	// Post-process results
	sm.generateTimeline(result)
	sm.calculateAuthorStats(result)
	sm.calculateStatistics(result, startTime)
	
	gologger.Info().Msgf("✅ Social Mining: Found %d results, %d unique subdomains from %d platforms", 
		len(result.Results), len(result.Subdomains), len(platforms))
	
	return result, nil
}

// mineFromPlatform mines subdomains from a specific platform
func (sm *SocialMiner) mineFromPlatform(ctx context.Context, platform, domain string) ([]SocialResult, error) {
	switch platform {
	case "github":
		return sm.mineFromGitHub(ctx, domain)
	case "gitlab":
		return sm.mineFromGitLab(ctx, domain)
	case "pastebin":
		return sm.mineFromPastebin(ctx, domain)
	case "reddit":
		return sm.mineFromReddit(ctx, domain)
	case "twitter":
		return sm.mineFromTwitter(ctx, domain)
	case "stackoverflow":
		return sm.mineFromStackOverflow(ctx, domain)
	default:
		return nil, fmt.Errorf("unsupported platform: %s", platform)
	}
}

// mineFromGitHub mines subdomains from GitHub
func (sm *SocialMiner) mineFromGitHub(ctx context.Context, domain string) ([]SocialResult, error) {
	// Check cache first
	cacheKey := fmt.Sprintf("github:%s", domain)
	if cached := sm.getFromCache(cacheKey); cached != nil {
		return cached, nil
	}
	
	// Rate limiting
	select {
	case <-sm.rateLimiter["github"].C:
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	
	var results []SocialResult
	
	// Search GitHub for domain mentions
	searchQueries := []string{
		fmt.Sprintf("\"%s\" language:yaml", domain),
		fmt.Sprintf("\"%s\" language:json", domain),
		fmt.Sprintf("\"%s\" language:dockerfile", domain),
		fmt.Sprintf("\"%s\" filename:config", domain),
		fmt.Sprintf("\"%s\" filename:.env", domain),
	}
	
	for _, query := range searchQueries {
		url := fmt.Sprintf("https://api.github.com/search/code?q=%s&per_page=100", query)
		
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			continue
		}
		
		// Add GitHub API key if available
		if apiKey, exists := sm.apiKeys["github"]; exists {
			req.Header.Set("Authorization", "token "+apiKey)
		}
		req.Header.Set("User-Agent", "Stormfinder-SocialMiner/1.0")
		
		resp, err := sm.client.Do(req)
		if err != nil {
			continue
		}
		
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			continue
		}
		
		var searchResult struct {
			Items []struct {
				Name        string `json:"name"`
				Path        string `json:"path"`
				HTMLURL     string `json:"html_url"`
				Repository  struct {
					FullName string `json:"full_name"`
					Owner    struct {
						Login string `json:"login"`
					} `json:"owner"`
				} `json:"repository"`
			} `json:"items"`
		}
		
		if err := json.NewDecoder(resp.Body).Decode(&searchResult); err != nil {
			resp.Body.Close()
			continue
		}
		resp.Body.Close()
		
		// Extract subdomains from search results
		for _, item := range searchResult.Items {
			// Get file content to extract subdomains
			subdomains := sm.extractSubdomainsFromGitHubFile(ctx, item.HTMLURL, domain)
			
			for _, subdomain := range subdomains {
				result := SocialResult{
					Subdomain:  subdomain,
					Source:     "github",
					Platform:   "GitHub",
					URL:        item.HTMLURL,
					Context:    fmt.Sprintf("Found in %s/%s", item.Repository.FullName, item.Path),
					Confidence: sm.calculateGitHubConfidence(item.Path, item.Name),
					Timestamp:  time.Now(),
					Author:     item.Repository.Owner.Login,
					Engagement: 0, // Would need additional API calls to get stars/forks
				}
				results = append(results, result)
			}
		}
	}
	
	// Cache results
	sm.setCache(cacheKey, results)
	
	return results, nil
}

// mineFromGitLab mines subdomains from GitLab
func (sm *SocialMiner) mineFromGitLab(ctx context.Context, domain string) ([]SocialResult, error) {
	// Similar implementation to GitHub but for GitLab API
	var results []SocialResult
	
	// GitLab search API
	url := fmt.Sprintf("https://gitlab.com/api/v4/search?scope=blobs&search=%s", domain)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return results, err
	}
	
	if apiKey, exists := sm.apiKeys["gitlab"]; exists {
		req.Header.Set("PRIVATE-TOKEN", apiKey)
	}
	
	// Rate limiting
	select {
	case <-sm.rateLimiter["gitlab"].C:
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	
	resp, err := sm.client.Do(req)
	if err != nil {
		return results, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return results, fmt.Errorf("GitLab API returned %d", resp.StatusCode)
	}
	
	var searchResults []struct {
		Basename   string `json:"basename"`
		Data       string `json:"data"`
		Path       string `json:"path"`
		Filename   string `json:"filename"`
		ProjectID  int    `json:"project_id"`
		Ref        string `json:"ref"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&searchResults); err != nil {
		return results, err
	}
	
	// Extract subdomains from GitLab results
	for _, item := range searchResults {
		subdomains := sm.extractSubdomainsFromText(item.Data, domain)
		
		for _, subdomain := range subdomains {
			result := SocialResult{
				Subdomain:  subdomain,
				Source:     "gitlab",
				Platform:   "GitLab",
				URL:        fmt.Sprintf("https://gitlab.com/project/%d/-/blob/%s/%s", item.ProjectID, item.Ref, item.Path),
				Context:    fmt.Sprintf("Found in %s", item.Path),
				Confidence: sm.calculateGitLabConfidence(item.Path, item.Filename),
				Timestamp:  time.Now(),
				Author:     "Unknown",
				Engagement: 0,
			}
			results = append(results, result)
		}
	}
	
	return results, nil
}

// mineFromPastebin mines subdomains from Pastebin
func (sm *SocialMiner) mineFromPastebin(ctx context.Context, domain string) ([]SocialResult, error) {
	var results []SocialResult
	
	// Pastebin scraping API (if available) or web scraping
	// This is a simplified implementation
	
	searchURL := fmt.Sprintf("https://psbdmp.ws/api/search/%s", domain)
	
	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return results, err
	}
	
	// Rate limiting
	select {
	case <-sm.rateLimiter["pastebin"].C:
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	
	resp, err := sm.client.Do(req)
	if err != nil {
		return results, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return results, fmt.Errorf("Pastebin search returned %d", resp.StatusCode)
	}
	
	var pastebinResults []struct {
		ID   string `json:"id"`
		Text string `json:"text"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&pastebinResults); err != nil {
		return results, err
	}
	
	// Extract subdomains from Pastebin results
	for _, item := range pastebinResults {
		subdomains := sm.extractSubdomainsFromText(item.Text, domain)
		
		for _, subdomain := range subdomains {
			result := SocialResult{
				Subdomain:  subdomain,
				Source:     "pastebin",
				Platform:   "Pastebin",
				URL:        fmt.Sprintf("https://pastebin.com/%s", item.ID),
				Context:    "Found in paste",
				Confidence: 0.6, // Medium confidence for pastebin
				Timestamp:  time.Now(),
				Author:     "Anonymous",
				Engagement: 0,
			}
			results = append(results, result)
		}
	}
	
	return results, nil
}

// mineFromReddit mines subdomains from Reddit
func (sm *SocialMiner) mineFromReddit(ctx context.Context, domain string) ([]SocialResult, error) {
	var results []SocialResult
	
	// Reddit API search
	searchURL := fmt.Sprintf("https://www.reddit.com/search.json?q=%s&limit=100", domain)
	
	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return results, err
	}
	
	req.Header.Set("User-Agent", "Stormfinder-SocialMiner/1.0")
	
	// Rate limiting
	select {
	case <-sm.rateLimiter["reddit"].C:
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	
	resp, err := sm.client.Do(req)
	if err != nil {
		return results, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return results, fmt.Errorf("Reddit API returned %d", resp.StatusCode)
	}
	
	var redditResponse struct {
		Data struct {
			Children []struct {
				Data struct {
					Title     string  `json:"title"`
					Selftext  string  `json:"selftext"`
					URL       string  `json:"url"`
					Author    string  `json:"author"`
					Score     int     `json:"score"`
					Subreddit string  `json:"subreddit"`
					Permalink string  `json:"permalink"`
				} `json:"data"`
			} `json:"children"`
		} `json:"data"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&redditResponse); err != nil {
		return results, err
	}
	
	// Extract subdomains from Reddit posts
	for _, child := range redditResponse.Data.Children {
		post := child.Data
		text := post.Title + " " + post.Selftext + " " + post.URL
		subdomains := sm.extractSubdomainsFromText(text, domain)
		
		for _, subdomain := range subdomains {
			result := SocialResult{
				Subdomain:  subdomain,
				Source:     "reddit",
				Platform:   "Reddit",
				URL:        "https://reddit.com" + post.Permalink,
				Context:    fmt.Sprintf("Found in r/%s post", post.Subreddit),
				Confidence: sm.calculateRedditConfidence(post.Score, post.Subreddit),
				Timestamp:  time.Now(),
				Author:     post.Author,
				Engagement: post.Score,
			}
			results = append(results, result)
		}
	}
	
	return results, nil
}

// mineFromTwitter mines subdomains from Twitter
func (sm *SocialMiner) mineFromTwitter(ctx context.Context, domain string) ([]SocialResult, error) {
	// Note: Twitter API v2 requires authentication
	// This is a placeholder implementation
	var results []SocialResult
	
	if _, exists := sm.apiKeys["twitter"]; !exists {
		return results, fmt.Errorf("Twitter API key required")
	}
	
	// Twitter API v2 search
	searchURL := fmt.Sprintf("https://api.twitter.com/2/tweets/search/recent?query=%s&max_results=100", domain)
	
	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return results, err
	}
	
	req.Header.Set("Authorization", "Bearer "+sm.apiKeys["twitter"])
	
	// Rate limiting
	select {
	case <-sm.rateLimiter["twitter"].C:
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	
	resp, err := sm.client.Do(req)
	if err != nil {
		return results, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return results, fmt.Errorf("Twitter API returned %d", resp.StatusCode)
	}
	
	var twitterResponse struct {
		Data []struct {
			ID     string `json:"id"`
			Text   string `json:"text"`
			Author string `json:"author_id"`
		} `json:"data"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&twitterResponse); err != nil {
		return results, err
	}
	
	// Extract subdomains from tweets
	for _, tweet := range twitterResponse.Data {
		subdomains := sm.extractSubdomainsFromText(tweet.Text, domain)
		
		for _, subdomain := range subdomains {
			result := SocialResult{
				Subdomain:  subdomain,
				Source:     "twitter",
				Platform:   "Twitter",
				URL:        fmt.Sprintf("https://twitter.com/i/status/%s", tweet.ID),
				Context:    "Found in tweet",
				Confidence: 0.7, // Medium-high confidence for Twitter
				Timestamp:  time.Now(),
				Author:     tweet.Author,
				Engagement: 0, // Would need additional API calls for engagement metrics
			}
			results = append(results, result)
		}
	}
	
	return results, nil
}

// mineFromStackOverflow mines subdomains from Stack Overflow
func (sm *SocialMiner) mineFromStackOverflow(ctx context.Context, domain string) ([]SocialResult, error) {
	var results []SocialResult
	
	// Stack Exchange API
	searchURL := fmt.Sprintf("https://api.stackexchange.com/2.3/search?order=desc&sort=activity&intitle=%s&site=stackoverflow", domain)
	
	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return results, err
	}
	
	resp, err := sm.client.Do(req)
	if err != nil {
		return results, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return results, fmt.Errorf("Stack Overflow API returned %d", resp.StatusCode)
	}
	
	var soResponse struct {
		Items []struct {
			QuestionID int    `json:"question_id"`
			Title      string `json:"title"`
			Body       string `json:"body"`
			Owner      struct {
				DisplayName string `json:"display_name"`
			} `json:"owner"`
			Score int `json:"score"`
		} `json:"items"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&soResponse); err != nil {
		return results, err
	}
	
	// Extract subdomains from Stack Overflow questions
	for _, item := range soResponse.Items {
		text := item.Title + " " + item.Body
		subdomains := sm.extractSubdomainsFromText(text, domain)
		
		for _, subdomain := range subdomains {
			result := SocialResult{
				Subdomain:  subdomain,
				Source:     "stackoverflow",
				Platform:   "Stack Overflow",
				URL:        fmt.Sprintf("https://stackoverflow.com/questions/%d", item.QuestionID),
				Context:    "Found in question/answer",
				Confidence: sm.calculateStackOverflowConfidence(item.Score),
				Timestamp:  time.Now(),
				Author:     item.Owner.DisplayName,
				Engagement: item.Score,
			}
			results = append(results, result)
		}
	}
	
	return results, nil
}

// Helper methods for subdomain extraction and confidence calculation

func (sm *SocialMiner) extractSubdomainsFromGitHubFile(ctx context.Context, fileURL, domain string) []string {
	// This would fetch the actual file content and extract subdomains
	// For now, return empty slice as placeholder
	return []string{}
}

func (sm *SocialMiner) extractSubdomainsFromText(text, domain string) []string {
	var subdomains []string
	
	// Regex to find subdomains
	domainRegex := regexp.MustCompile(`([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+` + regexp.QuoteMeta(domain))
	matches := domainRegex.FindAllString(text, -1)
	
	for _, match := range matches {
		if !contains(subdomains, match) {
			subdomains = append(subdomains, match)
		}
	}
	
	return subdomains
}

func (sm *SocialMiner) calculateGitHubConfidence(path, filename string) float64 {
	confidence := 0.5 // Base confidence
	
	// Higher confidence for configuration files
	if strings.Contains(path, "config") || strings.Contains(filename, "config") {
		confidence += 0.2
	}
	
	// Higher confidence for environment files
	if strings.Contains(filename, ".env") || strings.Contains(filename, "environment") {
		confidence += 0.3
	}
	
	// Higher confidence for Docker files
	if strings.Contains(filename, "docker") || strings.Contains(filename, "Dockerfile") {
		confidence += 0.2
	}
	
	if confidence > 1.0 {
		confidence = 1.0
	}
	
	return confidence
}

func (sm *SocialMiner) calculateGitLabConfidence(path, filename string) float64 {
	return sm.calculateGitHubConfidence(path, filename) // Same logic
}

func (sm *SocialMiner) calculateRedditConfidence(score int, subreddit string) float64 {
	confidence := 0.4 // Base confidence for Reddit
	
	// Higher confidence for technical subreddits
	techSubreddits := []string{"programming", "webdev", "sysadmin", "netsec", "devops"}
	for _, tech := range techSubreddits {
		if strings.Contains(strings.ToLower(subreddit), tech) {
			confidence += 0.2
			break
		}
	}
	
	// Higher confidence for higher scored posts
	if score > 100 {
		confidence += 0.2
	} else if score > 10 {
		confidence += 0.1
	}
	
	if confidence > 1.0 {
		confidence = 1.0
	}
	
	return confidence
}

func (sm *SocialMiner) calculateStackOverflowConfidence(score int) float64 {
	confidence := 0.6 // Base confidence for Stack Overflow
	
	// Higher confidence for higher scored questions/answers
	if score > 50 {
		confidence += 0.3
	} else if score > 10 {
		confidence += 0.2
	} else if score > 0 {
		confidence += 0.1
	}
	
	if confidence > 1.0 {
		confidence = 1.0
	}
	
	return confidence
}

// Cache management methods
func (sm *SocialMiner) getFromCache(key string) []SocialResult {
	sm.cacheMutex.RLock()
	defer sm.cacheMutex.RUnlock()
	
	if results, exists := sm.cache[key]; exists {
		return results
	}
	return nil
}

func (sm *SocialMiner) setCache(key string, results []SocialResult) {
	sm.cacheMutex.Lock()
	defer sm.cacheMutex.Unlock()
	
	sm.cache[key] = results
}

// Post-processing methods
func (sm *SocialMiner) generateTimeline(result *SocialMiningResult) {
	// Group results by date and platform
	timelineMap := make(map[string]map[string][]string)
	
	for _, res := range result.Results {
		dateKey := res.Timestamp.Format("2006-01-02")
		if timelineMap[dateKey] == nil {
			timelineMap[dateKey] = make(map[string][]string)
		}
		timelineMap[dateKey][res.Platform] = append(timelineMap[dateKey][res.Platform], res.Subdomain)
	}
	
	// Convert to timeline entries
	for dateStr, platforms := range timelineMap {
		date, _ := time.Parse("2006-01-02", dateStr)
		for platform, subdomains := range platforms {
			entry := TimelineEntry{
				Date:       date,
				Platform:   platform,
				Subdomains: subdomains,
				Activity:   fmt.Sprintf("Found %d subdomains", len(subdomains)),
			}
			result.Timeline = append(result.Timeline, entry)
		}
	}
}

func (sm *SocialMiner) calculateAuthorStats(result *SocialMiningResult) {
	authorMap := make(map[string]map[string]*AuthorStats)
	
	for _, res := range result.Results {
		if authorMap[res.Author] == nil {
			authorMap[res.Author] = make(map[string]*AuthorStats)
		}
		
		if authorMap[res.Author][res.Platform] == nil {
			authorMap[res.Author][res.Platform] = &AuthorStats{
				Author:   res.Author,
				Platform: res.Platform,
			}
		}
		
		stats := authorMap[res.Author][res.Platform]
		stats.Mentions++
		// Count unique subdomains per author
		// This is simplified - would need to track unique subdomains properly
		stats.Subdomains++
	}
	
	// Convert to slice and sort by mentions
	for _, platforms := range authorMap {
		for _, stats := range platforms {
			result.TopAuthors = append(result.TopAuthors, *stats)
		}
	}
	
	// Sort by mentions (simplified - would use sort.Slice in real implementation)
}

func (sm *SocialMiner) calculateStatistics(result *SocialMiningResult, startTime time.Time) {
	result.Statistics.TotalResults = len(result.Results)
	result.Statistics.UniqueSubdomains = len(result.Subdomains)
	result.Statistics.PlatformsQueried = len(result.PlatformStats)
	result.Statistics.ProcessingTime = time.Since(startTime)
	
	// Count confidence levels
	for _, res := range result.Results {
		if res.Confidence >= 0.8 {
			result.Statistics.HighConfidence++
		} else if res.Confidence >= 0.5 {
			result.Statistics.MediumConfidence++
		} else {
			result.Statistics.LowConfidence++
		}
	}
}

// Helper function
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
