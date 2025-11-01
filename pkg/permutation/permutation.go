package permutation

import (
	"fmt"
	"strings"
)

// Generator handles subdomain permutation and mutation
type Generator struct {
	commonWords []string
	numbers     []string
	separators  []string
}

// NewGenerator creates a new permutation generator
func NewGenerator() *Generator {
	return &Generator{
		commonWords: []string{
			"admin", "api", "app", "auth", "backup", "beta", "blog", "cdn", "chat", "cms",
			"dashboard", "db", "demo", "dev", "docs", "email", "ftp", "git", "help", "img",
			"internal", "lab", "mail", "mobile", "new", "old", "portal", "prod", "secure",
			"shop", "stage", "static", "support", "test", "vpn", "web", "wiki", "www",
			"assets", "cache", "cloud", "data", "files", "forum", "home", "media", "news",
			"office", "panel", "proxy", "search", "server", "store", "upload", "video",
		},
		numbers: []string{
			"1", "2", "3", "01", "02", "03", "2020", "2021", "2022", "2023", "2024", "2025",
		},
		separators: []string{"-", "_", ""},
	}
}

// GeneratePermutations creates permutations based on found subdomains
func (g *Generator) GeneratePermutations(foundSubdomains []string, baseDomain string) []string {
	permutations := make(map[string]struct{})
	
	for _, subdomain := range foundSubdomains {
		// Extract the subdomain part (remove base domain)
		subPart := strings.TrimSuffix(subdomain, "."+baseDomain)
		if subPart == subdomain {
			continue // Skip if it doesn't end with base domain
		}
		
		// Generate various permutations
		g.addWordPermutations(subPart, baseDomain, permutations)
		g.addNumberPermutations(subPart, baseDomain, permutations)
		g.addEnvironmentPermutations(subPart, baseDomain, permutations)
		g.addTLDPermutations(subPart, baseDomain, permutations)
	}
	
	// Convert map to slice
	result := make([]string, 0, len(permutations))
	for perm := range permutations {
		result = append(result, perm)
	}
	
	return result
}

// addWordPermutations adds common word-based permutations
func (g *Generator) addWordPermutations(subPart, baseDomain string, permutations map[string]struct{}) {
	for _, word := range g.commonWords {
		for _, sep := range g.separators {
			// Prepend word
			perm1 := fmt.Sprintf("%s%s%s.%s", word, sep, subPart, baseDomain)
			permutations[perm1] = struct{}{}
			
			// Append word
			perm2 := fmt.Sprintf("%s%s%s.%s", subPart, sep, word, baseDomain)
			permutations[perm2] = struct{}{}
		}
	}
}

// addNumberPermutations adds number-based permutations
func (g *Generator) addNumberPermutations(subPart, baseDomain string, permutations map[string]struct{}) {
	for _, num := range g.numbers {
		for _, sep := range g.separators {
			// Prepend number
			perm1 := fmt.Sprintf("%s%s%s.%s", num, sep, subPart, baseDomain)
			permutations[perm1] = struct{}{}
			
			// Append number
			perm2 := fmt.Sprintf("%s%s%s.%s", subPart, sep, num, baseDomain)
			permutations[perm2] = struct{}{}
		}
	}
	
	// Sequential numbers
	for i := 1; i <= 10; i++ {
		perm := fmt.Sprintf("%s%d.%s", subPart, i, baseDomain)
		permutations[perm] = struct{}{}
		
		perm2 := fmt.Sprintf("%s-%d.%s", subPart, i, baseDomain)
		permutations[perm2] = struct{}{}
	}
}

// addEnvironmentPermutations adds environment-based permutations
func (g *Generator) addEnvironmentPermutations(subPart, baseDomain string, permutations map[string]struct{}) {
	environments := []string{"dev", "test", "stage", "staging", "prod", "production", "beta", "alpha", "demo", "uat"}
	
	for _, env := range environments {
		for _, sep := range g.separators {
			// Environment prefix
			perm1 := fmt.Sprintf("%s%s%s.%s", env, sep, subPart, baseDomain)
			permutations[perm1] = struct{}{}
			
			// Environment suffix
			perm2 := fmt.Sprintf("%s%s%s.%s", subPart, sep, env, baseDomain)
			permutations[perm2] = struct{}{}
		}
	}
}

// addTLDPermutations adds TLD-based permutations (for subdomain discovery)
func (g *Generator) addTLDPermutations(subPart, baseDomain string, permutations map[string]struct{}) {
	regions := []string{"us", "eu", "asia", "uk", "ca", "au", "de", "fr", "jp", "cn"}
	
	for _, region := range regions {
		for _, sep := range g.separators {
			perm1 := fmt.Sprintf("%s%s%s.%s", region, sep, subPart, baseDomain)
			permutations[perm1] = struct{}{}
			
			perm2 := fmt.Sprintf("%s%s%s.%s", subPart, sep, region, baseDomain)
			permutations[perm2] = struct{}{}
		}
	}
}

// GenerateTypoSquatting generates typosquatting variations
func (g *Generator) GenerateTypoSquatting(subdomain, baseDomain string) []string {
	variations := make(map[string]struct{})
	subPart := strings.TrimSuffix(subdomain, "."+baseDomain)
	
	if subPart == subdomain {
		return nil
	}
	
	// Character substitution
	g.addCharacterSubstitutions(subPart, baseDomain, variations)
	
	// Character insertion
	g.addCharacterInsertions(subPart, baseDomain, variations)
	
	// Character deletion
	g.addCharacterDeletions(subPart, baseDomain, variations)
	
	// Character transposition
	g.addCharacterTranspositions(subPart, baseDomain, variations)
	
	result := make([]string, 0, len(variations))
	for variation := range variations {
		result = append(result, variation)
	}
	
	return result
}

func (g *Generator) addCharacterSubstitutions(subPart, baseDomain string, variations map[string]struct{}) {
	substitutions := map[rune][]rune{
		'a': {'@', '4'},
		'e': {'3'},
		'i': {'1', '!'},
		'o': {'0'},
		's': {'5', '$'},
		't': {'7'},
		'l': {'1'},
	}
	
	for i, char := range subPart {
		if subs, exists := substitutions[char]; exists {
			for _, sub := range subs {
				newSub := subPart[:i] + string(sub) + subPart[i+1:]
				variations[fmt.Sprintf("%s.%s", newSub, baseDomain)] = struct{}{}
			}
		}
	}
}

func (g *Generator) addCharacterInsertions(subPart, baseDomain string, variations map[string]struct{}) {
	chars := "abcdefghijklmnopqrstuvwxyz0123456789"
	
	for i := 0; i <= len(subPart); i++ {
		for _, char := range chars {
			newSub := subPart[:i] + string(char) + subPart[i:]
			if len(newSub) <= 20 { // Limit length
				variations[fmt.Sprintf("%s.%s", newSub, baseDomain)] = struct{}{}
			}
		}
	}
}

func (g *Generator) addCharacterDeletions(subPart, baseDomain string, variations map[string]struct{}) {
	for i := 0; i < len(subPart); i++ {
		newSub := subPart[:i] + subPart[i+1:]
		if len(newSub) > 0 {
			variations[fmt.Sprintf("%s.%s", newSub, baseDomain)] = struct{}{}
		}
	}
}

func (g *Generator) addCharacterTranspositions(subPart, baseDomain string, variations map[string]struct{}) {
	for i := 0; i < len(subPart)-1; i++ {
		chars := []rune(subPart)
		chars[i], chars[i+1] = chars[i+1], chars[i]
		newSub := string(chars)
		variations[fmt.Sprintf("%s.%s", newSub, baseDomain)] = struct{}{}
	}
}

// GenerateFromWordlist creates subdomains from a custom wordlist
func (g *Generator) GenerateFromWordlist(wordlist []string, baseDomain string) []string {
	subdomains := make([]string, 0, len(wordlist))
	
	for _, word := range wordlist {
		subdomain := fmt.Sprintf("%s.%s", word, baseDomain)
		subdomains = append(subdomains, subdomain)
		
		// Add variations with separators and numbers
		for _, sep := range g.separators {
			for _, num := range g.numbers {
				if sep != "" {
					variation1 := fmt.Sprintf("%s%s%s.%s", word, sep, num, baseDomain)
					variation2 := fmt.Sprintf("%s%s%s.%s", num, sep, word, baseDomain)
					subdomains = append(subdomains, variation1, variation2)
				}
			}
		}
	}
	
	return subdomains
}
