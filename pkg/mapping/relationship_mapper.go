package mapping

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/projectdiscovery/gologger"
)

// RelationshipMapper creates visual maps of subdomain relationships
type RelationshipMapper struct {
	subdomains map[string]*SubdomainNode
	edges      []RelationshipEdge
	mutex      sync.RWMutex
}

// SubdomainNode represents a subdomain in the relationship graph
type SubdomainNode struct {
	Name         string            `json:"name"`
	Level        int               `json:"level"`
	IPAddresses  []string          `json:"ip_addresses"`
	Technologies []string          `json:"technologies"`
	Services     []ServiceInfo     `json:"services"`
	Certificates []string          `json:"certificates"`
	Metadata     map[string]string `json:"metadata"`
	Confidence   float64           `json:"confidence"`
	LastSeen     time.Time         `json:"last_seen"`
	Sources      []string          `json:"sources"`
}

// RelationshipEdge represents a connection between subdomains
type RelationshipEdge struct {
	From         string  `json:"from"`
	To           string  `json:"to"`
	Relationship string  `json:"relationship"`
	Strength     float64 `json:"strength"`
	Evidence     string  `json:"evidence"`
}

// ServiceInfo contains service information for a subdomain
type ServiceInfo struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	Service  string `json:"service"`
	Version  string `json:"version"`
	Banner   string `json:"banner"`
}

// SubdomainMap represents the complete subdomain relationship map
type SubdomainMap struct {
	Domain       string             `json:"domain"`
	Nodes        []*SubdomainNode   `json:"nodes"`
	Edges        []RelationshipEdge `json:"edges"`
	Statistics   MapStatistics      `json:"statistics"`
	Visualization string            `json:"visualization"`
	GeneratedAt  time.Time          `json:"generated_at"`
}

// MapStatistics contains mapping statistics
type MapStatistics struct {
	TotalNodes      int                    `json:"total_nodes"`
	TotalEdges      int                    `json:"total_edges"`
	MaxDepth        int                    `json:"max_depth"`
	Clusters        int                    `json:"clusters"`
	Technologies    map[string]int         `json:"technologies"`
	ServicePorts    map[int]int           `json:"service_ports"`
	IPRanges        []string              `json:"ip_ranges"`
	ProcessingTime  time.Duration         `json:"processing_time"`
}

// NewRelationshipMapper creates a new relationship mapper
func NewRelationshipMapper() *RelationshipMapper {
	return &RelationshipMapper{
		subdomains: make(map[string]*SubdomainNode),
		edges:      make([]RelationshipEdge, 0),
	}
}

// BuildRelationshipMap creates a comprehensive subdomain relationship map
func (rm *RelationshipMapper) BuildRelationshipMap(ctx context.Context, domain string, subdomains []string) (*SubdomainMap, error) {
	startTime := time.Now()
	gologger.Info().Msgf("🗺️  Relationship Mapping: Building comprehensive map for %s", domain)
	
	// Initialize nodes
	for _, subdomain := range subdomains {
		rm.addSubdomainNode(subdomain, domain)
	}
	
	// Analyze relationships
	rm.analyzeIPRelationships(ctx)
	rm.analyzeTechnologyRelationships(ctx)
	rm.analyzeNamingPatterns(ctx, domain)
	rm.analyzeCertificateRelationships(ctx)
	rm.analyzeServiceRelationships(ctx)
	
	// Generate visualization
	visualization := rm.generateVisualization(domain)
	
	// Create final map
	nodes := make([]*SubdomainNode, 0, len(rm.subdomains))
	rm.mutex.RLock()
	for _, node := range rm.subdomains {
		nodes = append(nodes, node)
	}
	edges := rm.edges
	rm.mutex.RUnlock()
	
	subdomainMap := &SubdomainMap{
		Domain:        domain,
		Nodes:         nodes,
		Edges:         edges,
		Visualization: visualization,
		GeneratedAt:   time.Now(),
	}
	
	// Calculate statistics
	subdomainMap.Statistics = rm.calculateMapStatistics(startTime)
	
	gologger.Info().Msgf("✅ Relationship Mapping: Generated map with %d nodes, %d relationships", 
		len(nodes), len(edges))
	
	return subdomainMap, nil
}

// addSubdomainNode adds a subdomain node to the map
func (rm *RelationshipMapper) addSubdomainNode(subdomain, baseDomain string) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()
	
	if _, exists := rm.subdomains[subdomain]; exists {
		return
	}
	
	level := rm.calculateSubdomainLevel(subdomain, baseDomain)
	
	node := &SubdomainNode{
		Name:         subdomain,
		Level:        level,
		IPAddresses:  make([]string, 0),
		Technologies: make([]string, 0),
		Services:     make([]ServiceInfo, 0),
		Certificates: make([]string, 0),
		Metadata:     make(map[string]string),
		Confidence:   1.0,
		LastSeen:     time.Now(),
		Sources:      []string{"enumeration"},
	}
	
	// Resolve IP addresses
	if ips, err := net.LookupIP(subdomain); err == nil {
		for _, ip := range ips {
			node.IPAddresses = append(node.IPAddresses, ip.String())
		}
	}
	
	rm.subdomains[subdomain] = node
}

// analyzeIPRelationships finds relationships based on IP addresses
func (rm *RelationshipMapper) analyzeIPRelationships(ctx context.Context) {
	rm.mutex.RLock()
	nodes := make([]*SubdomainNode, 0, len(rm.subdomains))
	for _, node := range rm.subdomains {
		nodes = append(nodes, node)
	}
	rm.mutex.RUnlock()
	
	// Group subdomains by IP addresses
	ipGroups := make(map[string][]string)
	for _, node := range nodes {
		for _, ip := range node.IPAddresses {
			ipGroups[ip] = append(ipGroups[ip], node.Name)
		}
	}
	
	// Create relationships for subdomains sharing IPs
	for ip, subdomains := range ipGroups {
		if len(subdomains) > 1 {
			for i := 0; i < len(subdomains); i++ {
				for j := i + 1; j < len(subdomains); j++ {
					edge := RelationshipEdge{
						From:         subdomains[i],
						To:           subdomains[j],
						Relationship: "shared_ip",
						Strength:     0.8,
						Evidence:     fmt.Sprintf("Both resolve to IP %s", ip),
					}
					rm.addEdge(edge)
				}
			}
		}
	}
}

// analyzeTechnologyRelationships finds relationships based on technology stacks
func (rm *RelationshipMapper) analyzeTechnologyRelationships(ctx context.Context) {
	// This would analyze HTTP headers, server responses, etc.
	// Simplified implementation for demonstration
	
	rm.mutex.RLock()
	defer rm.mutex.RUnlock()
	
	for _, node := range rm.subdomains {
		// Detect technologies based on subdomain name patterns
		if strings.Contains(node.Name, "api") {
			node.Technologies = append(node.Technologies, "API")
		}
		if strings.Contains(node.Name, "cdn") || strings.Contains(node.Name, "static") {
			node.Technologies = append(node.Technologies, "CDN")
		}
		if strings.Contains(node.Name, "mail") || strings.Contains(node.Name, "smtp") {
			node.Technologies = append(node.Technologies, "Email")
		}
		if strings.Contains(node.Name, "admin") || strings.Contains(node.Name, "panel") {
			node.Technologies = append(node.Technologies, "Admin")
		}
	}
}

// analyzeNamingPatterns finds relationships based on naming patterns
func (rm *RelationshipMapper) analyzeNamingPatterns(ctx context.Context, baseDomain string) {
	rm.mutex.RLock()
	nodes := make([]*SubdomainNode, 0, len(rm.subdomains))
	for _, node := range rm.subdomains {
		nodes = append(nodes, node)
	}
	rm.mutex.RUnlock()
	
	// Analyze naming patterns
	patterns := make(map[string][]string)
	
	for _, node := range nodes {
		subPart := strings.TrimSuffix(node.Name, "."+baseDomain)
		
		// Extract base patterns
		if strings.Contains(subPart, "-") {
			parts := strings.Split(subPart, "-")
			if len(parts) > 1 {
				basePattern := parts[0]
				patterns[basePattern] = append(patterns[basePattern], node.Name)
			}
		}
		
		// Number patterns
		if len(subPart) > 1 && subPart[len(subPart)-1] >= '0' && subPart[len(subPart)-1] <= '9' {
			basePattern := subPart[:len(subPart)-1]
			patterns[basePattern] = append(patterns[basePattern], node.Name)
		}
	}
	
	// Create pattern-based relationships
	for pattern, subdomains := range patterns {
		if len(subdomains) > 1 {
			for i := 0; i < len(subdomains); i++ {
				for j := i + 1; j < len(subdomains); j++ {
					edge := RelationshipEdge{
						From:         subdomains[i],
						To:           subdomains[j],
						Relationship: "naming_pattern",
						Strength:     0.6,
						Evidence:     fmt.Sprintf("Both follow pattern: %s*", pattern),
					}
					rm.addEdge(edge)
				}
			}
		}
	}
}

// analyzeCertificateRelationships finds relationships based on SSL certificates
func (rm *RelationshipMapper) analyzeCertificateRelationships(ctx context.Context) {
	// This would analyze SSL certificates for shared issuers, SANs, etc.
	// Placeholder implementation
	gologger.Debug().Msg("Analyzing certificate relationships...")
}

// analyzeServiceRelationships finds relationships based on running services
func (rm *RelationshipMapper) analyzeServiceRelationships(ctx context.Context) {
	// This would perform service detection and create relationships
	// Placeholder implementation
	gologger.Debug().Msg("Analyzing service relationships...")
}

// generateVisualization creates a visual representation of the subdomain map
func (rm *RelationshipMapper) generateVisualization(domain string) string {
	// Generate a simple text-based visualization
	var viz strings.Builder
	
	viz.WriteString(fmt.Sprintf("Subdomain Relationship Map for %s\n", domain))
	viz.WriteString(strings.Repeat("=", 50) + "\n\n")
	
	// Group by levels
	levels := make(map[int][]string)
	rm.mutex.RLock()
	for name, node := range rm.subdomains {
		levels[node.Level] = append(levels[node.Level], name)
	}
	rm.mutex.RUnlock()
	
	// Display by levels
	for level := 0; level <= 5; level++ {
		if subdomains, exists := levels[level]; exists {
			viz.WriteString(fmt.Sprintf("Level %d:\n", level))
			for _, subdomain := range subdomains {
				viz.WriteString(fmt.Sprintf("  - %s\n", subdomain))
			}
			viz.WriteString("\n")
		}
	}
	
	// Add relationships
	viz.WriteString("Relationships:\n")
	for _, edge := range rm.edges {
		viz.WriteString(fmt.Sprintf("  %s --[%s]--> %s (%.2f)\n", 
			edge.From, edge.Relationship, edge.To, edge.Strength))
	}
	
	return viz.String()
}

// Helper methods

func (rm *RelationshipMapper) calculateSubdomainLevel(subdomain, baseDomain string) int {
	subPart := strings.TrimSuffix(subdomain, "."+baseDomain)
	if subPart == subdomain {
		return 0 // Not a subdomain of baseDomain
	}
	
	return strings.Count(subPart, ".") + 1
}

func (rm *RelationshipMapper) addEdge(edge RelationshipEdge) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()
	
	// Check if edge already exists
	for _, existing := range rm.edges {
		if (existing.From == edge.From && existing.To == edge.To) ||
		   (existing.From == edge.To && existing.To == edge.From) {
			return // Edge already exists
		}
	}
	
	rm.edges = append(rm.edges, edge)
}

func (rm *RelationshipMapper) calculateMapStatistics(startTime time.Time) MapStatistics {
	rm.mutex.RLock()
	defer rm.mutex.RUnlock()
	
	stats := MapStatistics{
		TotalNodes:     len(rm.subdomains),
		TotalEdges:     len(rm.edges),
		Technologies:   make(map[string]int),
		ServicePorts:   make(map[int]int),
		ProcessingTime: time.Since(startTime),
	}
	
	// Calculate max depth
	maxLevel := 0
	for _, node := range rm.subdomains {
		if node.Level > maxLevel {
			maxLevel = node.Level
		}
		
		// Count technologies
		for _, tech := range node.Technologies {
			stats.Technologies[tech]++
		}
		
		// Count service ports
		for _, service := range node.Services {
			stats.ServicePorts[service.Port]++
		}
	}
	stats.MaxDepth = maxLevel
	
	return stats
}

// ExportToJSON exports the relationship map to JSON
func (rm *RelationshipMapper) ExportToJSON(subdomainMap *SubdomainMap) ([]byte, error) {
	return json.MarshalIndent(subdomainMap, "", "  ")
}

// ExportToGraphviz exports the relationship map to Graphviz DOT format
func (rm *RelationshipMapper) ExportToGraphviz(subdomainMap *SubdomainMap) string {
	var dot strings.Builder
	
	dot.WriteString("digraph SubdomainMap {\n")
	dot.WriteString("  rankdir=TB;\n")
	dot.WriteString("  node [shape=box, style=rounded];\n\n")
	
	// Add nodes
	for _, node := range subdomainMap.Nodes {
		color := rm.getNodeColor(node)
		dot.WriteString(fmt.Sprintf("  \"%s\" [label=\"%s\\nIPs: %d\\nTech: %v\" fillcolor=\"%s\" style=\"filled,rounded\"];\n",
			node.Name, node.Name, len(node.IPAddresses), node.Technologies, color))
	}
	
	dot.WriteString("\n")
	
	// Add edges
	for _, edge := range subdomainMap.Edges {
		style := rm.getEdgeStyle(edge.Relationship)
		dot.WriteString(fmt.Sprintf("  \"%s\" -> \"%s\" [label=\"%s\" %s];\n",
			edge.From, edge.To, edge.Relationship, style))
	}
	
	dot.WriteString("}\n")
	
	return dot.String()
}

func (rm *RelationshipMapper) getNodeColor(node *SubdomainNode) string {
	if len(node.Technologies) > 0 {
		switch node.Technologies[0] {
		case "API":
			return "lightblue"
		case "CDN":
			return "lightgreen"
		case "Email":
			return "lightyellow"
		case "Admin":
			return "lightcoral"
		default:
			return "lightgray"
		}
	}
	return "white"
}

func (rm *RelationshipMapper) getEdgeStyle(relationship string) string {
	switch relationship {
	case "shared_ip":
		return "color=red, style=bold"
	case "naming_pattern":
		return "color=blue, style=dashed"
	case "certificate":
		return "color=green, style=dotted"
	default:
		return "color=gray"
	}
}
