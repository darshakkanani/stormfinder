package ai

import (
	"context"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/projectdiscovery/gologger"
)

// SubdomainPredictor uses AI/ML techniques to predict likely subdomains
type SubdomainPredictor struct {
	patterns       map[string]float64 // Pattern frequency weights
	ngramModel     map[string]float64 // N-gram language model
	contextModel   map[string][]string // Context-based predictions
	industryModel  map[string][]string // Industry-specific patterns
	technologyMap  map[string][]string // Technology-based predictions
}

// PredictionResult represents an AI prediction
type PredictionResult struct {
	Subdomain   string  `json:"subdomain"`
	Confidence  float64 `json:"confidence"`
	Reasoning   string  `json:"reasoning"`
	Category    string  `json:"category"`
	Priority    int     `json:"priority"`
}

// NewSubdomainPredictor creates a new AI predictor
func NewSubdomainPredictor() *SubdomainPredictor {
	return &SubdomainPredictor{
		patterns:      make(map[string]float64),
		ngramModel:    make(map[string]float64),
		contextModel:  initializeContextModel(),
		industryModel: initializeIndustryModel(),
		technologyMap: initializeTechnologyMap(),
	}
}

// PredictSubdomains uses AI to predict likely subdomains
func (sp *SubdomainPredictor) PredictSubdomains(ctx context.Context, domain string, knownSubdomains []string, maxPredictions int) []PredictionResult {
	gologger.Info().Msgf("🤖 AI Engine: Analyzing patterns for %s", domain)
	
	// Train the model on known subdomains
	sp.trainOnKnownSubdomains(knownSubdomains, domain)
	
	var predictions []PredictionResult
	
	// 1. Pattern-based predictions
	patternPredictions := sp.generatePatternPredictions(domain, knownSubdomains)
	predictions = append(predictions, patternPredictions...)
	
	// 2. N-gram based predictions
	ngramPredictions := sp.generateNgramPredictions(domain, knownSubdomains)
	predictions = append(predictions, ngramPredictions...)
	
	// 3. Context-aware predictions
	contextPredictions := sp.generateContextPredictions(domain, knownSubdomains)
	predictions = append(predictions, contextPredictions...)
	
	// 4. Industry-specific predictions
	industryPredictions := sp.generateIndustryPredictions(domain, knownSubdomains)
	predictions = append(predictions, industryPredictions...)
	
	// 5. Technology stack predictions
	techPredictions := sp.generateTechnologyPredictions(domain, knownSubdomains)
	predictions = append(predictions, techPredictions...)
	
	// 6. Semantic similarity predictions
	semanticPredictions := sp.generateSemanticPredictions(domain, knownSubdomains)
	predictions = append(predictions, semanticPredictions...)
	
	// Deduplicate and rank predictions
	uniquePredictions := sp.deduplicateAndRank(predictions, knownSubdomains)
	
	// Limit results
	if len(uniquePredictions) > maxPredictions {
		uniquePredictions = uniquePredictions[:maxPredictions]
	}
	
	gologger.Info().Msgf("🎯 AI Engine: Generated %d high-confidence predictions", len(uniquePredictions))
	return uniquePredictions
}

// trainOnKnownSubdomains trains the AI model on existing subdomains
func (sp *SubdomainPredictor) trainOnKnownSubdomains(subdomains []string, baseDomain string) {
	for _, subdomain := range subdomains {
		subPart := strings.TrimSuffix(subdomain, "."+baseDomain)
		if subPart == subdomain {
			continue
		}
		
		// Extract patterns
		sp.extractPatterns(subPart)
		
		// Build n-gram model
		sp.buildNgramModel(subPart)
	}
}

// generatePatternPredictions creates predictions based on learned patterns
func (sp *SubdomainPredictor) generatePatternPredictions(domain string, known []string) []PredictionResult {
	var predictions []PredictionResult
	
	// Analyze existing patterns
	patterns := sp.analyzeExistingPatterns(known, domain)
	
	// Generate variations based on patterns
	for pattern, confidence := range patterns {
		if confidence > 0.3 { // Confidence threshold
			variations := sp.generatePatternVariations(pattern, domain)
			for _, variation := range variations {
				predictions = append(predictions, PredictionResult{
					Subdomain:  variation,
					Confidence: confidence * 0.8, // Slightly lower confidence for predictions
					Reasoning:  fmt.Sprintf("Pattern-based prediction from '%s'", pattern),
					Category:   "pattern",
					Priority:   1,
				})
			}
		}
	}
	
	return predictions
}

// generateNgramPredictions uses n-gram language modeling
func (sp *SubdomainPredictor) generateNgramPredictions(domain string, known []string) []PredictionResult {
	var predictions []PredictionResult
	
	// Generate character-level predictions
	for _, subdomain := range known {
		subPart := strings.TrimSuffix(subdomain, "."+domain)
		if len(subPart) < 3 {
			continue
		}
		
		// Generate variations using n-gram model
		variations := sp.generateNgramVariations(subPart)
		for _, variation := range variations {
			confidence := sp.calculateNgramConfidence(variation)
			if confidence > 0.4 {
				predictions = append(predictions, PredictionResult{
					Subdomain:  variation + "." + domain,
					Confidence: confidence,
					Reasoning:  "N-gram language model prediction",
					Category:   "ngram",
					Priority:   2,
				})
			}
		}
	}
	
	return predictions
}

// generateContextPredictions uses context-aware predictions
func (sp *SubdomainPredictor) generateContextPredictions(domain string, known []string) []PredictionResult {
	var predictions []PredictionResult
	
	// Detect context from existing subdomains
	contexts := sp.detectContexts(known, domain)
	
	for context, confidence := range contexts {
		if contextSuggestions, exists := sp.contextModel[context]; exists {
			for _, suggestion := range contextSuggestions {
				predictions = append(predictions, PredictionResult{
					Subdomain:  suggestion + "." + domain,
					Confidence: confidence * 0.9,
					Reasoning:  fmt.Sprintf("Context-aware prediction based on '%s' pattern", context),
					Category:   "context",
					Priority:   1,
				})
			}
		}
	}
	
	return predictions
}

// generateIndustryPredictions creates industry-specific predictions
func (sp *SubdomainPredictor) generateIndustryPredictions(domain string, known []string) []PredictionResult {
	var predictions []PredictionResult
	
	// Detect industry based on existing subdomains and domain name
	industry := sp.detectIndustry(domain, known)
	
	if industryPatterns, exists := sp.industryModel[industry]; exists {
		for _, pattern := range industryPatterns {
			predictions = append(predictions, PredictionResult{
				Subdomain:  pattern + "." + domain,
				Confidence: 0.7,
				Reasoning:  fmt.Sprintf("Industry-specific prediction for %s sector", industry),
				Category:   "industry",
				Priority:   2,
			})
		}
	}
	
	return predictions
}

// generateTechnologyPredictions predicts based on technology stack
func (sp *SubdomainPredictor) generateTechnologyPredictions(domain string, known []string) []PredictionResult {
	var predictions []PredictionResult
	
	// Detect technology stack
	technologies := sp.detectTechnologies(known, domain)
	
	for tech, confidence := range technologies {
		if techPatterns, exists := sp.technologyMap[tech]; exists {
			for _, pattern := range techPatterns {
				predictions = append(predictions, PredictionResult{
					Subdomain:  pattern + "." + domain,
					Confidence: confidence * 0.8,
					Reasoning:  fmt.Sprintf("Technology-based prediction for %s stack", tech),
					Category:   "technology",
					Priority:   2,
				})
			}
		}
	}
	
	return predictions
}

// generateSemanticPredictions uses semantic similarity
func (sp *SubdomainPredictor) generateSemanticPredictions(domain string, known []string) []PredictionResult {
	var predictions []PredictionResult
	
	// Semantic word groups
	semanticGroups := map[string][]string{
		"api":     {"rest", "graphql", "v1", "v2", "gateway", "service"},
		"admin":   {"control", "panel", "dashboard", "manage", "console"},
		"dev":     {"test", "staging", "beta", "alpha", "sandbox"},
		"mail":    {"smtp", "pop", "imap", "webmail", "exchange"},
		"secure":  {"ssl", "tls", "vpn", "auth", "login"},
		"mobile":  {"app", "ios", "android", "m", "mobile"},
		"cdn":     {"static", "assets", "media", "files", "img"},
		"db":      {"database", "mysql", "postgres", "mongo", "redis"},
	}
	
	// Find semantic matches in known subdomains
	for _, subdomain := range known {
		subPart := strings.TrimSuffix(subdomain, "."+domain)
		for keyword, related := range semanticGroups {
			if strings.Contains(strings.ToLower(subPart), keyword) {
				for _, relatedWord := range related {
					predictions = append(predictions, PredictionResult{
						Subdomain:  relatedWord + "." + domain,
						Confidence: 0.6,
						Reasoning:  fmt.Sprintf("Semantic similarity to '%s'", subPart),
						Category:   "semantic",
						Priority:   3,
					})
				}
			}
		}
	}
	
	return predictions
}

// Helper functions for AI model initialization
func initializeContextModel() map[string][]string {
	return map[string][]string{
		"api":        {"api-v1", "api-v2", "rest-api", "graphql", "api-gateway", "microservice"},
		"admin":      {"admin-panel", "control-panel", "dashboard", "management", "console"},
		"dev":        {"development", "testing", "staging", "sandbox", "preview"},
		"prod":       {"production", "live", "release", "stable"},
		"auth":       {"authentication", "authorization", "sso", "oauth", "login", "signin"},
		"mail":       {"email", "smtp", "pop3", "imap", "webmail", "mailserver"},
		"cdn":        {"static", "assets", "media", "images", "files", "resources"},
		"mobile":     {"m", "mobile-api", "app", "ios", "android"},
		"monitoring": {"metrics", "logs", "health", "status", "alerts"},
		"security":   {"secure", "ssl", "vpn", "firewall", "scanner"},
	}
}

func initializeIndustryModel() map[string][]string {
	return map[string][]string{
		"ecommerce": {"shop", "cart", "checkout", "payment", "orders", "products", "catalog"},
		"fintech":   {"banking", "payments", "wallet", "trading", "crypto", "finance"},
		"healthcare": {"patient", "medical", "health", "clinic", "doctor", "pharmacy"},
		"education": {"student", "course", "learning", "classroom", "library", "academic"},
		"gaming":    {"game", "player", "match", "tournament", "leaderboard", "guild"},
		"media":     {"content", "video", "streaming", "news", "blog", "podcast"},
		"saas":      {"tenant", "workspace", "organization", "billing", "subscription"},
		"iot":       {"device", "sensor", "gateway", "telemetry", "monitoring"},
	}
}

func initializeTechnologyMap() map[string][]string {
	return map[string][]string{
		"kubernetes": {"k8s", "cluster", "pod", "service", "ingress", "namespace"},
		"aws":        {"s3", "ec2", "lambda", "rds", "cloudfront", "elb"},
		"docker":     {"container", "registry", "hub", "swarm"},
		"microservices": {"service", "gateway", "discovery", "config", "circuit-breaker"},
		"react":      {"spa", "frontend", "ui", "components"},
		"nodejs":     {"node", "npm", "express", "backend"},
		"python":     {"django", "flask", "fastapi", "celery"},
		"java":       {"spring", "tomcat", "maven", "gradle"},
	}
}

// Helper methods for pattern analysis and prediction
func (sp *SubdomainPredictor) extractPatterns(subdomain string) {
	// Extract various patterns
	patterns := []string{
		// Length patterns
		fmt.Sprintf("length_%d", len(subdomain)),
		// Character patterns
		sp.getCharacterPattern(subdomain),
		// Word patterns
		sp.getWordPattern(subdomain),
	}
	
	for _, pattern := range patterns {
		sp.patterns[pattern]++
	}
}

func (sp *SubdomainPredictor) getCharacterPattern(s string) string {
	pattern := ""
	for _, char := range s {
		if char >= 'a' && char <= 'z' {
			pattern += "a"
		} else if char >= 'A' && char <= 'Z' {
			pattern += "A"
		} else if char >= '0' && char <= '9' {
			pattern += "0"
		} else {
			pattern += string(char)
		}
	}
	return pattern
}

func (sp *SubdomainPredictor) getWordPattern(s string) string {
	// Detect common word patterns
	re := regexp.MustCompile(`[a-zA-Z]+`)
	words := re.FindAllString(s, -1)
	if len(words) > 0 {
		return fmt.Sprintf("words_%d", len(words))
	}
	return "no_words"
}

func (sp *SubdomainPredictor) buildNgramModel(subdomain string) {
	// Build character-level trigrams
	for i := 0; i <= len(subdomain)-3; i++ {
		trigram := subdomain[i : i+3]
		sp.ngramModel[trigram]++
	}
}

func (sp *SubdomainPredictor) analyzeExistingPatterns(subdomains []string, baseDomain string) map[string]float64 {
	patterns := make(map[string]float64)
	
	for _, subdomain := range subdomains {
		subPart := strings.TrimSuffix(subdomain, "."+baseDomain)
		if subPart == subdomain {
			continue
		}
		
		// Extract base words
		words := regexp.MustCompile(`[a-zA-Z]+`).FindAllString(subPart, -1)
		for _, word := range words {
			patterns[word] += 1.0
		}
	}
	
	// Normalize confidence scores
	total := float64(len(subdomains))
	for pattern := range patterns {
		patterns[pattern] = patterns[pattern] / total
	}
	
	return patterns
}

func (sp *SubdomainPredictor) generatePatternVariations(pattern, domain string) []string {
	var variations []string
	
	// Add common prefixes and suffixes
	prefixes := []string{"", "new-", "old-", "beta-", "test-", "dev-", "prod-"}
	suffixes := []string{"", "-api", "-v1", "-v2", "-new", "-old", "-test"}
	
	for _, prefix := range prefixes {
		for _, suffix := range suffixes {
			variation := prefix + pattern + suffix
			if variation != pattern { // Don't include the original
				variations = append(variations, variation+"."+domain)
			}
		}
	}
	
	return variations
}

func (sp *SubdomainPredictor) generateNgramVariations(subdomain string) []string {
	// This is a simplified n-gram variation generator
	// In a real implementation, you'd use more sophisticated NLP techniques
	var variations []string
	
	// Generate variations by character substitution based on n-gram probabilities
	common_substitutions := map[string][]string{
		"1": {"01", "2", "3"},
		"2": {"02", "1", "3"},
		"api": {"rest", "v1", "v2"},
		"dev": {"test", "stage", "beta"},
		"prod": {"live", "www", "main"},
	}
	
	for original, substitutes := range common_substitutions {
		if strings.Contains(subdomain, original) {
			for _, substitute := range substitutes {
				variation := strings.Replace(subdomain, original, substitute, 1)
				variations = append(variations, variation)
			}
		}
	}
	
	return variations
}

func (sp *SubdomainPredictor) calculateNgramConfidence(subdomain string) float64 {
	if len(subdomain) < 3 {
		return 0.1
	}
	
	score := 0.0
	count := 0
	
	for i := 0; i <= len(subdomain)-3; i++ {
		trigram := subdomain[i : i+3]
		if freq, exists := sp.ngramModel[trigram]; exists {
			score += freq
			count++
		}
	}
	
	if count == 0 {
		return 0.1
	}
	
	return math.Min(score/float64(count)/10.0, 1.0) // Normalize to 0-1
}

func (sp *SubdomainPredictor) detectContexts(subdomains []string, baseDomain string) map[string]float64 {
	contexts := make(map[string]float64)
	
	for _, subdomain := range subdomains {
		subPart := strings.TrimSuffix(subdomain, "."+baseDomain)
		subLower := strings.ToLower(subPart)
		
		// Detect various contexts
		if strings.Contains(subLower, "api") {
			contexts["api"] += 1.0
		}
		if strings.Contains(subLower, "admin") {
			contexts["admin"] += 1.0
		}
		if strings.Contains(subLower, "dev") || strings.Contains(subLower, "test") {
			contexts["dev"] += 1.0
		}
		if strings.Contains(subLower, "prod") || strings.Contains(subLower, "www") {
			contexts["prod"] += 1.0
		}
		if strings.Contains(subLower, "mail") || strings.Contains(subLower, "smtp") {
			contexts["mail"] += 1.0
		}
	}
	
	// Normalize
	total := float64(len(subdomains))
	for context := range contexts {
		contexts[context] = contexts[context] / total
	}
	
	return contexts
}

func (sp *SubdomainPredictor) detectIndustry(domain string, subdomains []string) string {
	// Simple industry detection based on keywords
	domainLower := strings.ToLower(domain)
	allText := domainLower + " " + strings.Join(subdomains, " ")
	
	industryKeywords := map[string][]string{
		"ecommerce":  {"shop", "cart", "store", "product", "order", "payment"},
		"fintech":    {"bank", "pay", "finance", "money", "crypto", "wallet"},
		"healthcare": {"health", "medical", "patient", "doctor", "clinic"},
		"education":  {"edu", "school", "student", "course", "learn"},
		"gaming":     {"game", "play", "player", "match", "guild"},
		"media":      {"news", "blog", "content", "video", "stream"},
		"saas":       {"app", "service", "platform", "tool", "dashboard"},
	}
	
	maxScore := 0
	detectedIndustry := "general"
	
	for industry, keywords := range industryKeywords {
		score := 0
		for _, keyword := range keywords {
			if strings.Contains(allText, keyword) {
				score++
			}
		}
		if score > maxScore {
			maxScore = score
			detectedIndustry = industry
		}
	}
	
	return detectedIndustry
}

func (sp *SubdomainPredictor) detectTechnologies(subdomains []string, domain string) map[string]float64 {
	technologies := make(map[string]float64)
	
	allText := strings.ToLower(domain + " " + strings.Join(subdomains, " "))
	
	techKeywords := map[string][]string{
		"kubernetes": {"k8s", "kube", "cluster", "pod"},
		"aws":        {"aws", "s3", "ec2", "lambda", "cloudfront"},
		"docker":     {"docker", "container", "registry"},
		"microservices": {"service", "micro", "gateway", "discovery"},
		"react":      {"react", "spa", "frontend", "ui"},
		"nodejs":     {"node", "npm", "express", "js"},
		"python":     {"python", "django", "flask", "py"},
		"java":       {"java", "spring", "tomcat", "jar"},
	}
	
	for tech, keywords := range techKeywords {
		score := 0.0
		for _, keyword := range keywords {
			if strings.Contains(allText, keyword) {
				score += 1.0
			}
		}
		if score > 0 {
			technologies[tech] = score / float64(len(keywords))
		}
	}
	
	return technologies
}

func (sp *SubdomainPredictor) deduplicateAndRank(predictions []PredictionResult, knownSubdomains []string) []PredictionResult {
	// Create a map to deduplicate
	uniquePredictions := make(map[string]PredictionResult)
	knownSet := make(map[string]bool)
	
	// Create set of known subdomains
	for _, known := range knownSubdomains {
		knownSet[known] = true
	}
	
	// Deduplicate and filter out known subdomains
	for _, pred := range predictions {
		if knownSet[pred.Subdomain] {
			continue // Skip already known subdomains
		}
		
		if existing, exists := uniquePredictions[pred.Subdomain]; exists {
			// Keep the prediction with higher confidence
			if pred.Confidence > existing.Confidence {
				uniquePredictions[pred.Subdomain] = pred
			}
		} else {
			uniquePredictions[pred.Subdomain] = pred
		}
	}
	
	// Convert back to slice and sort by confidence
	var result []PredictionResult
	for _, pred := range uniquePredictions {
		result = append(result, pred)
	}
	
	sort.Slice(result, func(i, j int) bool {
		if result[i].Confidence != result[j].Confidence {
			return result[i].Confidence > result[j].Confidence
		}
		return result[i].Priority < result[j].Priority
	})
	
	return result
}

// GetModelStats returns statistics about the AI model
func (sp *SubdomainPredictor) GetModelStats() map[string]interface{} {
	return map[string]interface{}{
		"patterns_learned":    len(sp.patterns),
		"ngram_entries":      len(sp.ngramModel),
		"context_models":     len(sp.contextModel),
		"industry_models":    len(sp.industryModel),
		"technology_models":  len(sp.technologyMap),
		"model_version":      "1.0",
		"last_trained":       time.Now().Format(time.RFC3339),
	}
}
