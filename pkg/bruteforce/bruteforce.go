package bruteforce

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/projectdiscovery/gologger"
)

// BruteForcer handles DNS brute force enumeration
type BruteForcer struct {
	resolvers    []string
	threads      int
	timeout      time.Duration
	wordlist     []string
	wildcardIPs  map[string]struct{}
}

// Result represents a brute force result
type Result struct {
	Subdomain string
	IPs       []string
	Error     error
}

// NewBruteForcer creates a new brute force instance
func NewBruteForcer(resolvers []string, threads int, timeout time.Duration) *BruteForcer {
	if len(resolvers) == 0 {
		resolvers = []string{"8.8.8.8:53", "1.1.1.1:53", "208.67.222.222:53"}
	}
	
	return &BruteForcer{
		resolvers: resolvers,
		threads:   threads,
		timeout:   timeout,
		wordlist:  getDefaultWordlist(),
	}
}

// SetWordlist sets a custom wordlist for brute forcing
func (b *BruteForcer) SetWordlist(wordlist []string) {
	b.wordlist = wordlist
}

// DetectWildcards detects wildcard DNS responses for a domain
func (b *BruteForcer) DetectWildcards(domain string) error {
	b.wildcardIPs = make(map[string]struct{})
	
	// Test with random subdomains to detect wildcards
	testSubdomains := []string{
		"nonexistent-" + generateRandomString(10),
		"random-" + generateRandomString(8),
		"test-" + generateRandomString(12),
	}
	
	for _, testSub := range testSubdomains {
		testDomain := testSub + "." + domain
		ips, err := b.resolve(testDomain)
		if err == nil && len(ips) > 0 {
			for _, ip := range ips {
				b.wildcardIPs[ip] = struct{}{}
			}
		}
	}
	
	if len(b.wildcardIPs) > 0 {
		gologger.Info().Msgf("Detected wildcard DNS for %s with IPs: %v", domain, getKeys(b.wildcardIPs))
	}
	
	return nil
}

// BruteForce performs DNS brute force enumeration
func (b *BruteForcer) BruteForce(ctx context.Context, domain string) <-chan Result {
	results := make(chan Result, 100)
	
	go func() {
		defer close(results)
		
		// Detect wildcards first
		if err := b.DetectWildcards(domain); err != nil {
			gologger.Warning().Msgf("Could not detect wildcards for %s: %v", domain, err)
		}
		
		// Create worker pool
		jobs := make(chan string, len(b.wordlist))
		var wg sync.WaitGroup
		
		// Start workers
		for i := 0; i < b.threads; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for subdomain := range jobs {
					select {
					case <-ctx.Done():
						return
					default:
						fullDomain := subdomain + "." + domain
						ips, err := b.resolve(fullDomain)
						
						if err != nil {
							continue // Skip DNS errors
						}
						
						// Filter out wildcard IPs
						filteredIPs := []string{}
						for _, ip := range ips {
							if _, isWildcard := b.wildcardIPs[ip]; !isWildcard {
								filteredIPs = append(filteredIPs, ip)
							}
						}
						
						if len(filteredIPs) > 0 {
							results <- Result{
								Subdomain: fullDomain,
								IPs:       filteredIPs,
							}
						}
					}
				}
			}()
		}
		
		// Send jobs
		for _, word := range b.wordlist {
			select {
			case <-ctx.Done():
				close(jobs)
				wg.Wait()
				return
			case jobs <- word:
			}
		}
		close(jobs)
		wg.Wait()
	}()
	
	return results
}

// resolve performs DNS resolution with timeout and multiple resolvers
func (b *BruteForcer) resolve(domain string) ([]string, error) {
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{Timeout: b.timeout}
			return d.DialContext(ctx, network, b.resolvers[0]) // Use first resolver for now
		},
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), b.timeout)
	defer cancel()
	
	ips, err := resolver.LookupIPAddr(ctx, domain)
	if err != nil {
		return nil, err
	}
	
	var result []string
	for _, ip := range ips {
		result = append(result, ip.IP.String())
	}
	
	return result, nil
}

// getDefaultWordlist returns a comprehensive subdomain wordlist
func getDefaultWordlist() []string {
	return []string{
		// Common subdomains
		"www", "mail", "ftp", "localhost", "webmail", "smtp", "pop", "ns1", "webdisk", "ns2",
		"cpanel", "whm", "autodiscover", "autoconfig", "m", "imap", "test", "ns", "blog",
		"pop3", "dev", "www2", "admin", "forum", "news", "vpn", "ns3", "mail2", "new",
		"mysql", "old", "lists", "support", "mobile", "mx", "static", "docs", "beta", "shop",
		"sql", "secure", "demo", "cp", "calendar", "wiki", "web", "media", "email", "images",
		"img", "www1", "intranet", "portal", "video", "sip", "dns2", "api", "cdn", "stats",
		"dns1", "ns4", "www3", "dns", "search", "staging", "server", "mx1", "chat", "wap",
		"my", "svn", "mail1", "sites", "proxy", "ads", "host", "crm", "cms", "backup",
		"mx2", "lyncdiscover", "info", "apps", "download", "remote", "db", "forums", "store",
		"relay", "files", "newsletter", "app", "live", "owa", "en", "start", "sms", "office",
		"exchange", "ipv4", "mail3", "help", "blogs", "helpdesk", "web1", "home", "library",
		"ftp2", "ntp", "monitor", "login", "service", "correo", "www4", "moodle", "it",
		"gateway", "gw", "i", "stat", "stage", "ldap", "tv", "ssl", "web2", "ns5", "upload",
		"nagios", "smtp2", "online", "ad", "survey", "data", "radio", "extranet", "test2",
		"mssql", "dns3", "jobs", "services", "panel", "irc", "hosting", "cloud", "de", "gmail",
		"s", "bbs", "cs", "ww", "mrtg", "review", "ddns", "lab", "r", "analytics", "sandbox",
		"ja", "www5", "postgres", "www6", "rs", "mail4", "travel", "spanish", "secure2", "tv2",
		"ping", "direct", "survey2", "trace", "www7", "ftp1", "files2", "c", "b", "mobile2",
		"facebook", "s2", "s1", "www-dev", "twitter", "devtest", "f", "ecommerce", "social",
		"backup2", "oracle", "sun", "msoid", "share", "v2", "magento", "photos", "redmine",
		"node", "pma", "mt", "zendesk", "sub", "s3", "movie", "secure3", "ps", "training",
		"labs", "linux", "sc", "love", "fax", "php", "lp", "tracking", "thumbs", "up", "tw",
		"campus", "reg", "digital", "demo2", "da", "tr", "otrs", "web3", "home2", "uat", "v",
		"tmall", "union", "noc", "netmail", "beta2", "archive", "s4", "photo", "eb", "video2",
		"web-dev", "v1", "mail5", "ham", "ops", "lab2", "dev2", "img2", "vps", "driver",
		
		// Technical subdomains
		"api", "cdn", "assets", "static", "media", "content", "files", "images", "js", "css",
		"fonts", "uploads", "downloads", "resources", "data", "cache", "tmp", "temp",
		
		// Environment-based
		"prod", "production", "staging", "stage", "dev", "development", "test", "testing",
		"qa", "uat", "demo", "sandbox", "preview", "beta", "alpha", "rc", "pre", "preprod",
		
		// Geographic/Language
		"us", "eu", "asia", "uk", "ca", "au", "de", "fr", "es", "it", "jp", "cn", "br",
		"mx", "in", "ru", "nl", "se", "no", "dk", "fi", "pl", "cz", "hu", "ro", "bg",
		
		// Services
		"auth", "sso", "oauth", "login", "signin", "signup", "register", "account", "profile",
		"dashboard", "admin", "panel", "control", "manage", "console", "cp", "cpanel",
		"plesk", "whm", "webmin", "phpmyadmin", "pma", "adminer",
		
		// Applications
		"app", "apps", "application", "service", "services", "microservice", "ms", "ws",
		"webservice", "rest", "graphql", "grpc", "soap",
		
		// Infrastructure
		"lb", "loadbalancer", "proxy", "reverse-proxy", "gateway", "firewall", "router",
		"switch", "hub", "bridge", "tunnel", "vpn", "bastion", "jump", "relay",
		
		// Monitoring & Logging
		"monitor", "monitoring", "metrics", "stats", "analytics", "logs", "logging",
		"kibana", "grafana", "prometheus", "nagios", "zabbix", "cacti", "munin",
		
		// Databases
		"db", "database", "mysql", "postgres", "postgresql", "mongo", "mongodb", "redis",
		"elastic", "elasticsearch", "solr", "cassandra", "neo4j", "influx", "influxdb",
		
		// CI/CD & DevOps
		"ci", "cd", "jenkins", "gitlab", "github", "bitbucket", "bamboo", "teamcity",
		"travis", "circleci", "drone", "concourse", "spinnaker", "argo", "tekton",
		
		// Cloud & Containers
		"k8s", "kubernetes", "docker", "registry", "harbor", "quay", "gcr", "ecr",
		"acs", "aks", "eks", "gke", "openshift", "rancher", "nomad", "consul",
		
		// Security
		"vault", "secrets", "keystore", "cert", "certificate", "ca", "pki", "acme",
		"security", "sec", "scanner", "scan", "pentest", "audit",
		
		// Backup & Storage
		"backup", "backups", "archive", "storage", "s3", "blob", "object", "file",
		"nfs", "smb", "ftp", "sftp", "rsync", "sync",
		
		// Communication
		"mail", "email", "smtp", "pop", "pop3", "imap", "webmail", "exchange", "outlook",
		"chat", "slack", "teams", "discord", "irc", "xmpp", "sip", "voip", "pbx",
		
		// Content Management
		"cms", "wordpress", "wp", "drupal", "joomla", "ghost", "hugo", "jekyll",
		"contentful", "strapi", "directus", "craft", "concrete5",
		
		// E-commerce
		"shop", "store", "ecommerce", "cart", "checkout", "payment", "pay", "billing",
		"invoice", "magento", "shopify", "woocommerce", "prestashop", "opencart",
		
		// Single letters and numbers (for short subdomains)
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p",
		"q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
		"01", "02", "03", "04", "05", "06", "07", "08", "09", "10",
	}
}

// Helper functions
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

func getKeys(m map[string]struct{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
