package cache

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// Cache provides caching functionality for subdomain results
type Cache struct {
	cacheDir string
	ttl      time.Duration
}

// CacheEntry represents a cached result
type CacheEntry struct {
	Domain    string    `json:"domain"`
	Results   []string  `json:"results"`
	Timestamp time.Time `json:"timestamp"`
	Source    string    `json:"source"`
}

// NewCache creates a new cache instance
func NewCache(cacheDir string, ttl time.Duration) *Cache {
	if cacheDir == "" {
		homeDir, _ := os.UserHomeDir()
		cacheDir = filepath.Join(homeDir, ".stormfinder", "cache")
	}
	
	// Create cache directory if it doesn't exist
	os.MkdirAll(cacheDir, 0755)
	
	return &Cache{
		cacheDir: cacheDir,
		ttl:      ttl,
	}
}

// Get retrieves cached results for a domain and source
func (c *Cache) Get(domain, source string) ([]string, bool) {
	cacheKey := c.generateCacheKey(domain, source)
	cachePath := filepath.Join(c.cacheDir, cacheKey+".json")
	
	// Check if cache file exists
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		return nil, false
	}
	
	// Read cache file
	data, err := ioutil.ReadFile(cachePath)
	if err != nil {
		return nil, false
	}
	
	var entry CacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, false
	}
	
	// Check if cache is expired
	if time.Since(entry.Timestamp) > c.ttl {
		// Remove expired cache
		os.Remove(cachePath)
		return nil, false
	}
	
	return entry.Results, true
}

// Set stores results in cache
func (c *Cache) Set(domain, source string, results []string) error {
	cacheKey := c.generateCacheKey(domain, source)
	cachePath := filepath.Join(c.cacheDir, cacheKey+".json")
	
	entry := CacheEntry{
		Domain:    domain,
		Results:   results,
		Timestamp: time.Now(),
		Source:    source,
	}
	
	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return err
	}
	
	return ioutil.WriteFile(cachePath, data, 0644)
}

// Clear removes all cached entries
func (c *Cache) Clear() error {
	return os.RemoveAll(c.cacheDir)
}

// ClearExpired removes expired cache entries
func (c *Cache) ClearExpired() error {
	files, err := ioutil.ReadDir(c.cacheDir)
	if err != nil {
		return err
	}
	
	for _, file := range files {
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}
		
		cachePath := filepath.Join(c.cacheDir, file.Name())
		data, err := ioutil.ReadFile(cachePath)
		if err != nil {
			continue
		}
		
		var entry CacheEntry
		if err := json.Unmarshal(data, &entry); err != nil {
			continue
		}
		
		if time.Since(entry.Timestamp) > c.ttl {
			os.Remove(cachePath)
		}
	}
	
	return nil
}

// GetStats returns cache statistics
func (c *Cache) GetStats() (map[string]interface{}, error) {
	files, err := ioutil.ReadDir(c.cacheDir)
	if err != nil {
		return nil, err
	}
	
	stats := map[string]interface{}{
		"total_entries": 0,
		"expired_entries": 0,
		"cache_size_mb": 0.0,
		"cache_dir": c.cacheDir,
		"ttl_hours": c.ttl.Hours(),
	}
	
	var totalSize int64
	expiredCount := 0
	
	for _, file := range files {
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}
		
		totalSize += file.Size()
		
		cachePath := filepath.Join(c.cacheDir, file.Name())
		data, err := ioutil.ReadFile(cachePath)
		if err != nil {
			continue
		}
		
		var entry CacheEntry
		if err := json.Unmarshal(data, &entry); err != nil {
			continue
		}
		
		if time.Since(entry.Timestamp) > c.ttl {
			expiredCount++
		}
	}
	
	stats["total_entries"] = len(files)
	stats["expired_entries"] = expiredCount
	stats["cache_size_mb"] = float64(totalSize) / (1024 * 1024)
	
	return stats, nil
}

// generateCacheKey creates a unique cache key for domain and source
func (c *Cache) generateCacheKey(domain, source string) string {
	key := fmt.Sprintf("%s:%s", domain, source)
	hash := md5.Sum([]byte(key))
	return hex.EncodeToString(hash[:])
}
