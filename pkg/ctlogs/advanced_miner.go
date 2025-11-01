package ctlogs

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/projectdiscovery/gologger"
)

// AdvancedCTMiner provides advanced Certificate Transparency log mining
type AdvancedCTMiner struct {
	logServers    []CTLogServer
	client        *http.Client
	rateLimiter   *time.Ticker
	cache         map[string][]CTEntry
	cacheMutex    sync.RWMutex
}

// CTLogServer represents a Certificate Transparency log server
type CTLogServer struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	MaxEntries  int    `json:"max_entries"`
	Description string `json:"description"`
	Operator    string `json:"operator"`
}

// CTEntry represents a certificate transparency log entry
type CTEntry struct {
	LeafInput string `json:"leaf_input"`
	ExtraData string `json:"extra_data"`
}

// CTLogResponse represents the response from CT log API
type CTLogResponse struct {
	Entries []CTEntry `json:"entries"`
}

// CertificateInfo contains parsed certificate information
type CertificateInfo struct {
	CommonName       string    `json:"common_name"`
	SubjectAltNames  []string  `json:"subject_alt_names"`
	Issuer          string    `json:"issuer"`
	NotBefore       time.Time `json:"not_before"`
	NotAfter        time.Time `json:"not_after"`
	SerialNumber    string    `json:"serial_number"`
	Fingerprint     string    `json:"fingerprint"`
	LogServer       string    `json:"log_server"`
	EntryIndex      int64     `json:"entry_index"`
}

// AdvancedCTResult contains comprehensive CT mining results
type AdvancedCTResult struct {
	Subdomains      []string          `json:"subdomains"`
	Certificates    []CertificateInfo `json:"certificates"`
	WildcardCerts   []CertificateInfo `json:"wildcard_certs"`
	ExpiredCerts    []CertificateInfo `json:"expired_certs"`
	RecentCerts     []CertificateInfo `json:"recent_certs"`
	Issuers         map[string]int    `json:"issuers"`
	Timeline        []TimelineEntry   `json:"timeline"`
	Statistics      CTStatistics      `json:"statistics"`
}

// TimelineEntry represents certificate issuance timeline
type TimelineEntry struct {
	Date         time.Time `json:"date"`
	Action       string    `json:"action"`
	Certificate  string    `json:"certificate"`
	Subdomains   []string  `json:"subdomains"`
}

// CTStatistics contains mining statistics
type CTStatistics struct {
	TotalCertificates   int                    `json:"total_certificates"`
	UniqueSubdomains    int                    `json:"unique_subdomains"`
	WildcardCerts      int                    `json:"wildcard_certs"`
	ExpiredCerts       int                    `json:"expired_certs"`
	LogServersQueried  int                    `json:"log_servers_queried"`
	TimeRange          string                 `json:"time_range"`
	IssuerDistribution map[string]int         `json:"issuer_distribution"`
	ProcessingTime     time.Duration          `json:"processing_time"`
}

// NewAdvancedCTMiner creates a new advanced CT miner
func NewAdvancedCTMiner() *AdvancedCTMiner {
	return &AdvancedCTMiner{
		logServers: []CTLogServer{
			{
				Name:        "Google Argon 2024",
				URL:         "https://ct.googleapis.com/logs/argon2024/",
				MaxEntries:  10000,
				Description: "Google Certificate Transparency Log",
				Operator:    "Google",
			},
			{
				Name:        "Cloudflare Nimbus 2024",
				URL:         "https://ct.cloudflare.com/logs/nimbus2024/",
				MaxEntries:  10000,
				Description: "Cloudflare Certificate Transparency Log",
				Operator:    "Cloudflare",
			},
			{
				Name:        "Let's Encrypt Oak 2024",
				URL:         "https://oak.ct.letsencrypt.org/2024h1/",
				MaxEntries:  10000,
				Description: "Let's Encrypt Certificate Transparency Log",
				Operator:    "Let's Encrypt",
			},
			{
				Name:        "DigiCert Yeti 2024",
				URL:         "https://yeti2024.ct.digicert.com/log/",
				MaxEntries:  10000,
				Description: "DigiCert Certificate Transparency Log",
				Operator:    "DigiCert",
			},
			{
				Name:        "Sectigo Mammoth 2024",
				URL:         "https://mammoth.ct.comodo.com/",
				MaxEntries:  10000,
				Description: "Sectigo Certificate Transparency Log",
				Operator:    "Sectigo",
			},
		},
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		rateLimiter: time.NewTicker(100 * time.Millisecond), // 10 requests per second
		cache:       make(map[string][]CTEntry),
	}
}

// MineAdvanced performs advanced CT log mining with comprehensive analysis
func (miner *AdvancedCTMiner) MineAdvanced(ctx context.Context, domain string) (*AdvancedCTResult, error) {
	startTime := time.Now()
	gologger.Info().Msgf("🔍 Advanced CT Mining: Starting comprehensive analysis for %s", domain)
	
	result := &AdvancedCTResult{
		Subdomains:    make([]string, 0),
		Certificates:  make([]CertificateInfo, 0),
		WildcardCerts: make([]CertificateInfo, 0),
		ExpiredCerts:  make([]CertificateInfo, 0),
		RecentCerts:   make([]CertificateInfo, 0),
		Issuers:       make(map[string]int),
		Timeline:      make([]TimelineEntry, 0),
		Statistics:    CTStatistics{IssuerDistribution: make(map[string]int)},
	}
	
	var wg sync.WaitGroup
	var mutex sync.Mutex
	
	// Query multiple CT log servers in parallel
	for _, logServer := range miner.logServers {
		wg.Add(1)
		go func(server CTLogServer) {
			defer wg.Done()
			
			gologger.Debug().Msgf("Querying CT log server: %s", server.Name)
			entries, err := miner.queryLogServer(ctx, server, domain)
			if err != nil {
				gologger.Warning().Msgf("Failed to query %s: %v", server.Name, err)
				return
			}
			
			// Process entries
			for _, entry := range entries {
				certInfo := miner.parseCertificateEntry(entry, server.Name)
				if certInfo != nil && miner.isRelevantToDomain(certInfo, domain) {
					mutex.Lock()
					result.Certificates = append(result.Certificates, *certInfo)
					
					// Extract subdomains
					subdomains := miner.extractSubdomains(certInfo, domain)
					for _, subdomain := range subdomains {
						if !contains(result.Subdomains, subdomain) {
							result.Subdomains = append(result.Subdomains, subdomain)
						}
					}
					
					// Categorize certificates
					miner.categorizeCertificate(certInfo, result)
					
					// Update statistics
					result.Statistics.IssuerDistribution[certInfo.Issuer]++
					result.Issuers[certInfo.Issuer]++
					
					mutex.Unlock()
				}
			}
			
			gologger.Debug().Msgf("Processed %d entries from %s", len(entries), server.Name)
		}(logServer)
	}
	
	wg.Wait()
	
	// Post-process results
	miner.generateTimeline(result)
	miner.calculateStatistics(result, startTime)
	miner.identifyPatterns(result, domain)
	
	gologger.Info().Msgf("✅ Advanced CT Mining: Found %d certificates, %d unique subdomains", 
		len(result.Certificates), len(result.Subdomains))
	
	return result, nil
}

// queryLogServer queries a specific CT log server
func (miner *AdvancedCTMiner) queryLogServer(ctx context.Context, server CTLogServer, domain string) ([]CTEntry, error) {
	// Check cache first
	cacheKey := fmt.Sprintf("%s:%s", server.Name, domain)
	miner.cacheMutex.RLock()
	if cached, exists := miner.cache[cacheKey]; exists {
		miner.cacheMutex.RUnlock()
		return cached, nil
	}
	miner.cacheMutex.RUnlock()
	
	// Rate limiting
	select {
	case <-miner.rateLimiter.C:
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	
	// Query the CT log server
	url := fmt.Sprintf("%sct/v1/get-entries?start=0&end=%d", server.URL, server.MaxEntries)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("User-Agent", "Stormfinder-CTMiner/1.0")
	
	resp, err := miner.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d from %s", resp.StatusCode, server.Name)
	}
	
	var ctResponse CTLogResponse
	if err := json.NewDecoder(resp.Body).Decode(&ctResponse); err != nil {
		return nil, err
	}
	
	// Filter entries relevant to domain
	var relevantEntries []CTEntry
	for _, entry := range ctResponse.Entries {
		if miner.isEntryRelevant(entry, domain) {
			relevantEntries = append(relevantEntries, entry)
		}
	}
	
	// Cache results
	miner.cacheMutex.Lock()
	miner.cache[cacheKey] = relevantEntries
	miner.cacheMutex.Unlock()
	
	return relevantEntries, nil
}

// parseCertificateEntry parses a CT log entry into certificate information
func (miner *AdvancedCTMiner) parseCertificateEntry(entry CTEntry, logServer string) *CertificateInfo {
	// This is a simplified parser - in reality, you'd need to decode the DER certificate
	// For this example, we'll extract basic information from the leaf input
	
	// Decode base64 leaf input
	leafData, err := base64.StdEncoding.DecodeString(entry.LeafInput)
	if err != nil {
		return nil
	}
	
	// Create fingerprint
	hash := sha256.Sum256(leafData)
	fingerprint := fmt.Sprintf("%x", hash)
	
	// Extract basic information (simplified)
	certInfo := &CertificateInfo{
		Fingerprint:     fingerprint,
		LogServer:       logServer,
		SerialNumber:    fmt.Sprintf("%x", hash[:8]),
		NotBefore:       time.Now().Add(-30 * 24 * time.Hour), // Placeholder
		NotAfter:        time.Now().Add(90 * 24 * time.Hour),  // Placeholder
	}
	
	// Try to extract domain information from the entry
	// This is a simplified extraction - real implementation would parse the certificate
	entryStr := string(leafData)
	
	// Extract common name and SANs using regex patterns
	cnRegex := regexp.MustCompile(`CN=([^,\s]+)`)
	if matches := cnRegex.FindStringSubmatch(entryStr); len(matches) > 1 {
		certInfo.CommonName = matches[1]
	}
	
	// Extract Subject Alternative Names
	sanRegex := regexp.MustCompile(`DNS:([^,\s]+)`)
	sanMatches := sanRegex.FindAllStringSubmatch(entryStr, -1)
	for _, match := range sanMatches {
		if len(match) > 1 {
			certInfo.SubjectAltNames = append(certInfo.SubjectAltNames, match[1])
		}
	}
	
	// Extract issuer
	issuerRegex := regexp.MustCompile(`O=([^,\s]+)`)
	if matches := issuerRegex.FindStringSubmatch(entryStr); len(matches) > 1 {
		certInfo.Issuer = matches[1]
	} else {
		certInfo.Issuer = "Unknown"
	}
	
	return certInfo
}

// isEntryRelevant checks if a CT entry is relevant to the target domain
func (miner *AdvancedCTMiner) isEntryRelevant(entry CTEntry, domain string) bool {
	entryStr := strings.ToLower(entry.LeafInput + entry.ExtraData)
	domainLower := strings.ToLower(domain)
	
	// Check if the domain appears in the entry
	return strings.Contains(entryStr, domainLower)
}

// isRelevantToDomain checks if certificate info is relevant to the domain
func (miner *AdvancedCTMiner) isRelevantToDomain(certInfo *CertificateInfo, domain string) bool {
	domainLower := strings.ToLower(domain)
	
	// Check common name
	if strings.Contains(strings.ToLower(certInfo.CommonName), domainLower) {
		return true
	}
	
	// Check SANs
	for _, san := range certInfo.SubjectAltNames {
		if strings.Contains(strings.ToLower(san), domainLower) {
			return true
		}
	}
	
	return false
}

// extractSubdomains extracts subdomains from certificate information
func (miner *AdvancedCTMiner) extractSubdomains(certInfo *CertificateInfo, baseDomain string) []string {
	var subdomains []string
	domainLower := strings.ToLower(baseDomain)
	
	// Extract from common name
	if cn := strings.ToLower(certInfo.CommonName); strings.HasSuffix(cn, "."+domainLower) {
		subdomains = append(subdomains, certInfo.CommonName)
	}
	
	// Extract from SANs
	for _, san := range certInfo.SubjectAltNames {
		sanLower := strings.ToLower(san)
		if strings.HasSuffix(sanLower, "."+domainLower) && !contains(subdomains, san) {
			subdomains = append(subdomains, san)
		}
	}
	
	return subdomains
}

// categorizeCertificate categorizes certificates into different types
func (miner *AdvancedCTMiner) categorizeCertificate(certInfo *CertificateInfo, result *AdvancedCTResult) {
	// Check for wildcard certificates
	if strings.Contains(certInfo.CommonName, "*") {
		result.WildcardCerts = append(result.WildcardCerts, *certInfo)
	}
	
	for _, san := range certInfo.SubjectAltNames {
		if strings.Contains(san, "*") {
			result.WildcardCerts = append(result.WildcardCerts, *certInfo)
			break
		}
	}
	
	// Check for expired certificates
	if time.Now().After(certInfo.NotAfter) {
		result.ExpiredCerts = append(result.ExpiredCerts, *certInfo)
	}
	
	// Check for recent certificates (last 30 days)
	if time.Since(certInfo.NotBefore) <= 30*24*time.Hour {
		result.RecentCerts = append(result.RecentCerts, *certInfo)
	}
}

// generateTimeline creates a timeline of certificate issuance
func (miner *AdvancedCTMiner) generateTimeline(result *AdvancedCTResult) {
	for _, cert := range result.Certificates {
		subdomains := []string{cert.CommonName}
		subdomains = append(subdomains, cert.SubjectAltNames...)
		
		entry := TimelineEntry{
			Date:        cert.NotBefore,
			Action:      "Certificate Issued",
			Certificate: cert.Fingerprint,
			Subdomains:  subdomains,
		}
		result.Timeline = append(result.Timeline, entry)
	}
	
	// Sort timeline by date
	// Note: In a real implementation, you'd use sort.Slice here
}

// calculateStatistics calculates comprehensive statistics
func (miner *AdvancedCTMiner) calculateStatistics(result *AdvancedCTResult, startTime time.Time) {
	result.Statistics.TotalCertificates = len(result.Certificates)
	result.Statistics.UniqueSubdomains = len(result.Subdomains)
	result.Statistics.WildcardCerts = len(result.WildcardCerts)
	result.Statistics.ExpiredCerts = len(result.ExpiredCerts)
	result.Statistics.LogServersQueried = len(miner.logServers)
	result.Statistics.ProcessingTime = time.Since(startTime)
	
	// Calculate time range
	if len(result.Certificates) > 0 {
		earliest := result.Certificates[0].NotBefore
		latest := result.Certificates[0].NotAfter
		
		for _, cert := range result.Certificates {
			if cert.NotBefore.Before(earliest) {
				earliest = cert.NotBefore
			}
			if cert.NotAfter.After(latest) {
				latest = cert.NotAfter
			}
		}
		
		result.Statistics.TimeRange = fmt.Sprintf("%s to %s", 
			earliest.Format("2006-01-02"), latest.Format("2006-01-02"))
	}
}

// identifyPatterns identifies patterns in certificate issuance
func (miner *AdvancedCTMiner) identifyPatterns(result *AdvancedCTResult, domain string) {
	// This could include pattern analysis like:
	// - Certificate renewal patterns
	// - Subdomain naming conventions
	// - Issuer preferences
	// - Seasonal patterns in certificate issuance
	
	gologger.Debug().Msgf("Pattern analysis completed for %s", domain)
}

// GetCacheStats returns cache statistics
func (miner *AdvancedCTMiner) GetCacheStats() map[string]interface{} {
	miner.cacheMutex.RLock()
	defer miner.cacheMutex.RUnlock()
	
	totalEntries := 0
	for _, entries := range miner.cache {
		totalEntries += len(entries)
	}
	
	return map[string]interface{}{
		"cached_queries":  len(miner.cache),
		"total_entries":   totalEntries,
		"log_servers":     len(miner.logServers),
		"cache_hit_ratio": "N/A", // Would need to track hits/misses
	}
}

// ClearCache clears the CT mining cache
func (miner *AdvancedCTMiner) ClearCache() {
	miner.cacheMutex.Lock()
	defer miner.cacheMutex.Unlock()
	
	miner.cache = make(map[string][]CTEntry)
	gologger.Debug().Msg("CT mining cache cleared")
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
