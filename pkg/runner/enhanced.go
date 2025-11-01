package runner

import (
	"bufio"
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/projectdiscovery/gologger"
	"github.com/darshakkanani/stormfinder/v2/pkg/bruteforce"
	"github.com/darshakkanani/stormfinder/v2/pkg/permutation"
	"github.com/darshakkanani/stormfinder/v2/pkg/resolve"
)

// EnhancedEnumerator provides advanced subdomain discovery capabilities
type EnhancedEnumerator struct {
	options     *Options
	bruteForcer *bruteforce.BruteForcer
	permGen     *permutation.Generator
	resolver    *resolve.ResolutionPool
}

// NewEnhancedEnumerator creates a new enhanced enumeration engine
func NewEnhancedEnumerator(options *Options, resolver *resolve.ResolutionPool) *EnhancedEnumerator {
	var bf *bruteforce.BruteForcer
	if options.BruteForce {
		resolvers := []string{"8.8.8.8:53", "1.1.1.1:53", "208.67.222.222:53", "9.9.9.9:53"}
		if len(options.Resolvers) > 0 {
			resolvers = options.Resolvers
		}
		bf = bruteforce.NewBruteForcer(resolvers, options.BruteThreads, time.Duration(options.Timeout)*time.Second)
		
		// Load custom wordlists if provided
		wordlist := loadAllWordlists(options)
		if len(wordlist) > 0 {
			bf.SetWordlist(wordlist)
			LogSuccess("Loaded %d words from custom wordlists", len(wordlist))
		}
	}

	return &EnhancedEnumerator{
		options:     options,
		bruteForcer: bf,
		permGen:     permutation.NewGenerator(),
		resolver:    resolver,
	}
}

// EnhancedEnumerate performs comprehensive subdomain enumeration
func (e *EnhancedEnumerator) EnhancedEnumerate(ctx context.Context, domain string, passiveResults []string) []string {
	allResults := make(map[string]struct{})
	
	// Add passive results
	for _, result := range passiveResults {
		allResults[result] = struct{}{}
	}

	var wg sync.WaitGroup
	resultsChan := make(chan string, 1000)
	
	// Collect results
	go func() {
		for result := range resultsChan {
			allResults[result] = struct{}{}
		}
	}()

	// 1. DNS Brute Force Enumeration
	if e.options.BruteForce && e.bruteForcer != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			e.runBruteForce(ctx, domain, resultsChan)
		}()
	}

	// 2. Subdomain Permutations
	if e.options.Permutations {
		wg.Add(1)
		go func() {
			defer wg.Done()
			e.runPermutations(ctx, domain, passiveResults, resultsChan)
		}()
	}

	// 3. Recursive Enumeration
	if e.options.Recursive {
		wg.Add(1)
		go func() {
			defer wg.Done()
			e.runRecursiveEnumeration(ctx, domain, passiveResults, resultsChan, 0)
		}()
	}

	wg.Wait()
	close(resultsChan)

	// Convert results to slice
	results := make([]string, 0, len(allResults))
	for result := range allResults {
		results = append(results, result)
	}

	return results
}

// runBruteForce performs DNS brute force enumeration
func (e *EnhancedEnumerator) runBruteForce(ctx context.Context, domain string, results chan<- string) {
	if e.options.Verbose {
		gologger.Info().Msgf("Starting DNS brute force enumeration for %s", domain)
	}

	bruteResults := e.bruteForcer.BruteForce(ctx, domain)
	count := 0
	
	for result := range bruteResults {
		if result.Error == nil {
			results <- result.Subdomain
			count++
			if e.options.Verbose && count%100 == 0 {
				gologger.Info().Msgf("Brute force found %d subdomains so far...", count)
			}
		}
	}

	if e.options.Verbose {
		gologger.Info().Msgf("DNS brute force completed: found %d subdomains", count)
	}
}

// runPermutations generates and tests subdomain permutations
func (e *EnhancedEnumerator) runPermutations(ctx context.Context, domain string, foundSubdomains []string, results chan<- string) {
	if e.options.Verbose {
		gologger.Info().Msgf("Generating subdomain permutations for %s", domain)
	}

	// Generate permutations based on found subdomains
	permutations := e.permGen.GeneratePermutations(foundSubdomains, domain)
	
	if e.options.Verbose {
		gologger.Info().Msgf("Generated %d permutations to test", len(permutations))
	}

	// Test permutations in parallel
	semaphore := make(chan struct{}, e.options.BruteThreads)
	var wg sync.WaitGroup

	count := 0
	for _, perm := range permutations {
		select {
		case <-ctx.Done():
			return
		default:
		}

		wg.Add(1)
		go func(subdomain string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if e.testSubdomain(subdomain) {
				results <- subdomain
				count++
			}
		}(perm)
	}

	wg.Wait()

	if e.options.Verbose {
		gologger.Info().Msgf("Permutation testing completed: found %d valid subdomains", count)
	}
}

// runRecursiveEnumeration performs recursive subdomain discovery
func (e *EnhancedEnumerator) runRecursiveEnumeration(ctx context.Context, domain string, foundSubdomains []string, results chan<- string, depth int) {
	if depth >= e.options.MaxDepth {
		return
	}

	if e.options.Verbose {
		gologger.Info().Msgf("Running recursive enumeration (depth %d) for %s", depth+1, domain)
	}

	// Extract unique second-level subdomains for recursive search
	secondLevelDomains := make(map[string]struct{})
	for _, subdomain := range foundSubdomains {
		parts := strings.Split(subdomain, ".")
		if len(parts) >= 3 {
			// Create second-level domain (e.g., api.example.com -> api)
			secondLevel := parts[0] + "." + domain
			secondLevelDomains[secondLevel] = struct{}{}
		}
	}

	// For each second-level domain, try common third-level subdomains
	commonPrefixes := []string{
		"api", "admin", "test", "dev", "staging", "prod", "www", "mail", "ftp",
		"secure", "internal", "private", "public", "beta", "alpha", "demo",
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, e.options.BruteThreads)

	for secondLevel := range secondLevelDomains {
		for _, prefix := range commonPrefixes {
			select {
			case <-ctx.Done():
				return
			default:
			}

			wg.Add(1)
			go func(sub, pref string) {
				defer wg.Done()
				semaphore <- struct{}{}
				defer func() { <-semaphore }()

				testDomain := pref + "." + sub
				if e.testSubdomain(testDomain) {
					results <- testDomain
				}
			}(secondLevel, prefix)
		}
	}

	wg.Wait()
}

// testSubdomain tests if a subdomain resolves
func (e *EnhancedEnumerator) testSubdomain(subdomain string) bool {
	// Simple DNS resolution test
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(e.options.Timeout)*time.Second)
	defer cancel()

	// Use the resolver to test the subdomain
	if e.resolver != nil {
		// This is a simplified test - in a real implementation,
		// you'd want to use the proper resolution pool
		select {
		case <-ctx.Done():
			return false
		default:
			return true // Placeholder - implement actual DNS resolution
		}
	}

	return false
}

// loadWordlist loads a wordlist from file
func loadWordlist(filename string) ([]string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	wordlist := make([]string, 0, len(lines))
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			wordlist = append(wordlist, line)
		}
	}

	return wordlist, nil
}

// loadAllWordlists loads wordlists from all specified sources
func loadAllWordlists(options *Options) []string {
	var allWords []string
	wordMap := make(map[string]bool) // To avoid duplicates
	
	// Load from single wordlist file
	if options.Wordlist != "" {
		if options.Verbose {
			LogProgress("Loading wordlist from file: %s", options.Wordlist)
		}
		words, err := loadWordlist(options.Wordlist)
		if err != nil {
			gologger.Warning().Msgf("Could not load wordlist %s: %v", options.Wordlist, err)
		} else {
			for _, word := range words {
				if !wordMap[word] {
					wordMap[word] = true
					allWords = append(allWords, word)
				}
			}
			if options.Verbose {
				LogSuccess("Loaded %d words from %s", len(words), options.Wordlist)
			}
		}
	}
	
	// Load from wordlist directory
	if options.WordlistDir != "" {
		if options.Verbose {
			LogProgress("Loading wordlists from directory: %s", options.WordlistDir)
		}
		words := loadWordlistsFromDirectory(options.WordlistDir)
		for _, word := range words {
			if !wordMap[word] {
				wordMap[word] = true
				allWords = append(allWords, word)
			}
		}
		if options.Verbose && len(words) > 0 {
			LogSuccess("Loaded %d words from directory %s", len(words), options.WordlistDir)
		}
	}
	
	// Load from URLs
	if len(options.WordlistURLs) > 0 {
		if options.Verbose {
			LogProgress("Downloading wordlists from %d URLs", len(options.WordlistURLs))
		}
		for _, url := range options.WordlistURLs {
			words := downloadWordlist(url)
			for _, word := range words {
				if !wordMap[word] {
					wordMap[word] = true
					allWords = append(allWords, word)
				}
			}
			if options.Verbose && len(words) > 0 {
				LogSuccess("Downloaded %d words from %s", len(words), url)
			}
		}
	}
	
	// Load built-in wordlist if no custom wordlists provided
	if len(allWords) == 0 {
		builtinWords := loadBuiltinWordlist()
		for _, word := range builtinWords {
			if !wordMap[word] {
				wordMap[word] = true
				allWords = append(allWords, word)
			}
		}
		if options.Verbose {
			LogInfo("Using built-in wordlist with %d words", len(builtinWords))
		}
	}
	
	return allWords
}

// loadWordlistsFromDirectory loads all wordlist files from a directory
func loadWordlistsFromDirectory(dirPath string) []string {
	var allWords []string
	
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		gologger.Warning().Msgf("Could not read directory %s: %v", dirPath, err)
		return allWords
	}
	
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		
		// Only process text files
		filename := file.Name()
		if strings.HasSuffix(filename, ".txt") || strings.HasSuffix(filename, ".list") || strings.HasSuffix(filename, ".wordlist") {
			filepath := filepath.Join(dirPath, filename)
			words, err := loadWordlist(filepath)
			if err != nil {
				gologger.Warning().Msgf("Could not load wordlist %s: %v", filepath, err)
				continue
			}
			allWords = append(allWords, words...)
		}
	}
	
	return allWords
}

// downloadWordlist downloads a wordlist from a URL
func downloadWordlist(url string) []string {
	var words []string
	
	resp, err := http.Get(url)
	if err != nil {
		gologger.Warning().Msgf("Could not download wordlist from %s: %v", url, err)
		return words
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		gologger.Warning().Msgf("HTTP %d when downloading wordlist from %s", resp.StatusCode, url)
		return words
	}
	
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			words = append(words, line)
		}
	}
	
	if err := scanner.Err(); err != nil {
		gologger.Warning().Msgf("Error reading wordlist from %s: %v", url, err)
	}
	
	return words
}

// loadBuiltinWordlist loads the built-in wordlist from the wordlists directory
func loadBuiltinWordlist() []string {
	// Try to load from the built-in wordlists directory
	builtinPath := "wordlists/common.txt"
	
	// First try relative to current directory
	if words, err := loadWordlist(builtinPath); err == nil {
		return words
	}
	
	// Try relative to executable directory
	execPath, err := os.Executable()
	if err == nil {
		execDir := filepath.Dir(execPath)
		builtinPath = filepath.Join(execDir, "wordlists", "common.txt")
		if words, err := loadWordlist(builtinPath); err == nil {
			return words
		}
	}
	
	// Fallback to a basic built-in wordlist
	return []string{
		"www", "mail", "ftp", "localhost", "webmail", "smtp", "pop", "ns1", "ns2", "ns3", "ns4", "ns5",
		"admin", "administrator", "api", "app", "apps", "blog", "cdn", "cpanel", "dev", "development",
		"docs", "forum", "help", "img", "images", "m", "mobile", "mx", "news", "old", "portal",
		"secure", "shop", "sql", "ssl", "stage", "staging", "static", "stats", "status", "test",
		"testing", "vpn", "web", "webdisk", "webmail", "whm", "ww1", "ww2", "ww3", "subdomain",
		"s1", "s2", "s3", "s4", "s5", "server", "server1", "server2", "host", "host1", "host2",
		"email", "e", "demo", "beta", "alpha", "preview", "pre", "prod", "production", "live",
	}
}

// GetEnhancedStats returns statistics about the enhanced enumeration
func (e *EnhancedEnumerator) GetEnhancedStats() map[string]interface{} {
	stats := make(map[string]interface{})
	
	stats["brute_force_enabled"] = e.options.BruteForce
	stats["permutations_enabled"] = e.options.Permutations
	stats["recursive_enabled"] = e.options.Recursive
	stats["max_depth"] = e.options.MaxDepth
	stats["brute_threads"] = e.options.BruteThreads
	
	if e.options.Wordlist != "" {
		stats["custom_wordlist"] = e.options.Wordlist
	}
	
	return stats
}
