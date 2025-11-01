package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/projectdiscovery/gologger"
)

// RealtimeMonitor provides continuous subdomain monitoring
type RealtimeMonitor struct {
	domains     map[string]*DomainWatch
	alerts      chan Alert
	webhooks    []WebhookConfig
	running     bool
	mutex       sync.RWMutex
	stopChan    chan struct{}
}

// DomainWatch represents a domain being monitored
type DomainWatch struct {
	Domain          string            `json:"domain"`
	KnownSubdomains map[string]bool   `json:"known_subdomains"`
	LastCheck       time.Time         `json:"last_check"`
	CheckInterval   time.Duration     `json:"check_interval"`
	AlertThreshold  int               `json:"alert_threshold"`
	Config          MonitoringConfig  `json:"config"`
	Statistics      MonitoringStats   `json:"statistics"`
}

// Alert represents a monitoring alert
type Alert struct {
	ID          string                 `json:"id"`
	Domain      string                 `json:"domain"`
	Type        string                 `json:"type"`
	Severity    string                 `json:"severity"`
	Message     string                 `json:"message"`
	NewSubdomains []string             `json:"new_subdomains,omitempty"`
	ChangedSubdomains []SubdomainChange `json:"changed_subdomains,omitempty"`
	Timestamp   time.Time              `json:"timestamp"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// SubdomainChange represents a change in subdomain status
type SubdomainChange struct {
	Subdomain   string    `json:"subdomain"`
	ChangeType  string    `json:"change_type"`
	OldValue    string    `json:"old_value,omitempty"`
	NewValue    string    `json:"new_value,omitempty"`
	DetectedAt  time.Time `json:"detected_at"`
}

// WebhookConfig represents webhook configuration
type WebhookConfig struct {
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Enabled bool              `json:"enabled"`
}

// MonitoringConfig contains monitoring configuration
type MonitoringConfig struct {
	EnableNewSubdomainAlerts bool          `json:"enable_new_subdomain_alerts"`
	EnableIPChangeAlerts     bool          `json:"enable_ip_change_alerts"`
	EnableCertChangeAlerts   bool          `json:"enable_cert_change_alerts"`
	EnableServiceAlerts      bool          `json:"enable_service_alerts"`
	CheckInterval           time.Duration `json:"check_interval"`
	AlertCooldown           time.Duration `json:"alert_cooldown"`
	MaxAlertsPerHour        int           `json:"max_alerts_per_hour"`
}

// MonitoringStats contains monitoring statistics
type MonitoringStats struct {
	TotalChecks       int       `json:"total_checks"`
	NewSubdomains     int       `json:"new_subdomains"`
	ChangedSubdomains int       `json:"changed_subdomains"`
	AlertsSent        int       `json:"alerts_sent"`
	LastAlert         time.Time `json:"last_alert"`
	Uptime            time.Duration `json:"uptime"`
	StartTime         time.Time `json:"start_time"`
}

// NewRealtimeMonitor creates a new realtime monitor
func NewRealtimeMonitor() *RealtimeMonitor {
	return &RealtimeMonitor{
		domains:  make(map[string]*DomainWatch),
		alerts:   make(chan Alert, 1000),
		webhooks: make([]WebhookConfig, 0),
		stopChan: make(chan struct{}),
	}
}

// AddDomain adds a domain to monitor
func (rm *RealtimeMonitor) AddDomain(domain string, config MonitoringConfig) error {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	if config.CheckInterval == 0 {
		config.CheckInterval = 5 * time.Minute // Default 5 minutes
	}

	watch := &DomainWatch{
		Domain:          domain,
		KnownSubdomains: make(map[string]bool),
		LastCheck:       time.Time{},
		CheckInterval:   config.CheckInterval,
		AlertThreshold:  5, // Alert after 5 new subdomains
		Config:          config,
		Statistics: MonitoringStats{
			StartTime: time.Now(),
		},
	}

	rm.domains[domain] = watch
	gologger.Info().Msgf("📡 Monitor: Added domain %s for realtime monitoring", domain)
	return nil
}

// RemoveDomain removes a domain from monitoring
func (rm *RealtimeMonitor) RemoveDomain(domain string) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	delete(rm.domains, domain)
	gologger.Info().Msgf("📡 Monitor: Removed domain %s from monitoring", domain)
}

// AddWebhook adds a webhook for alerts
func (rm *RealtimeMonitor) AddWebhook(webhook WebhookConfig) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	rm.webhooks = append(rm.webhooks, webhook)
	gologger.Info().Msgf("📡 Monitor: Added webhook %s", webhook.URL)
}

// Start begins realtime monitoring
func (rm *RealtimeMonitor) Start(ctx context.Context) error {
	rm.mutex.Lock()
	if rm.running {
		rm.mutex.Unlock()
		return fmt.Errorf("monitor is already running")
	}
	rm.running = true
	rm.mutex.Unlock()

	gologger.Info().Msg("🚀 Realtime Monitor: Starting continuous monitoring...")

	// Start alert processor
	go rm.processAlerts(ctx)

	// Start monitoring each domain
	var wg sync.WaitGroup
	for domain := range rm.domains {
		wg.Add(1)
		go func(d string) {
			defer wg.Done()
			rm.monitorDomain(ctx, d)
		}(domain)
	}

	// Wait for context cancellation or stop signal
	select {
	case <-ctx.Done():
		gologger.Info().Msg("📡 Monitor: Context cancelled, stopping...")
	case <-rm.stopChan:
		gologger.Info().Msg("📡 Monitor: Stop signal received, stopping...")
	}

	rm.mutex.Lock()
	rm.running = false
	rm.mutex.Unlock()

	return nil
}

// Stop stops the realtime monitoring
func (rm *RealtimeMonitor) Stop() {
	close(rm.stopChan)
}

// monitorDomain monitors a specific domain
func (rm *RealtimeMonitor) monitorDomain(ctx context.Context, domain string) {
	ticker := time.NewTicker(rm.getDomainCheckInterval(domain))
	defer ticker.Stop()

	gologger.Debug().Msgf("📡 Monitor: Starting monitoring for %s", domain)

	for {
		select {
		case <-ctx.Done():
			return
		case <-rm.stopChan:
			return
		case <-ticker.C:
			rm.checkDomain(ctx, domain)
		}
	}
}

// checkDomain performs a check on a domain
func (rm *RealtimeMonitor) checkDomain(ctx context.Context, domain string) {
	rm.mutex.RLock()
	watch, exists := rm.domains[domain]
	if !exists {
		rm.mutex.RUnlock()
		return
	}
	rm.mutex.RUnlock()

	gologger.Debug().Msgf("📡 Monitor: Checking domain %s", domain)

	// Perform subdomain enumeration (simplified)
	currentSubdomains := rm.enumerateSubdomains(ctx, domain)

	// Compare with known subdomains
	newSubdomains := rm.findNewSubdomains(watch.KnownSubdomains, currentSubdomains)
	changedSubdomains := rm.findChangedSubdomains(ctx, domain, currentSubdomains)

	// Update statistics
	rm.mutex.Lock()
	watch.LastCheck = time.Now()
	watch.Statistics.TotalChecks++
	watch.Statistics.Uptime = time.Since(watch.Statistics.StartTime)

	// Update known subdomains
	for _, subdomain := range currentSubdomains {
		watch.KnownSubdomains[subdomain] = true
	}
	rm.mutex.Unlock()

	// Generate alerts if needed
	if len(newSubdomains) > 0 && watch.Config.EnableNewSubdomainAlerts {
		alert := Alert{
			ID:            fmt.Sprintf("new-subdomains-%d", time.Now().Unix()),
			Domain:        domain,
			Type:          "new_subdomains",
			Severity:      rm.calculateSeverity(len(newSubdomains)),
			Message:       fmt.Sprintf("Discovered %d new subdomains for %s", len(newSubdomains), domain),
			NewSubdomains: newSubdomains,
			Timestamp:     time.Now(),
			Metadata: map[string]interface{}{
				"count": len(newSubdomains),
				"check_id": watch.Statistics.TotalChecks,
			},
		}

		select {
		case rm.alerts <- alert:
			watch.Statistics.AlertsSent++
			watch.Statistics.NewSubdomains += len(newSubdomains)
		default:
			gologger.Warning().Msg("Alert queue full, dropping alert")
		}
	}

	if len(changedSubdomains) > 0 {
		alert := Alert{
			ID:                fmt.Sprintf("changed-subdomains-%d", time.Now().Unix()),
			Domain:            domain,
			Type:              "changed_subdomains",
			Severity:          "medium",
			Message:           fmt.Sprintf("Detected %d subdomain changes for %s", len(changedSubdomains), domain),
			ChangedSubdomains: changedSubdomains,
			Timestamp:         time.Now(),
			Metadata: map[string]interface{}{
				"count": len(changedSubdomains),
			},
		}

		select {
		case rm.alerts <- alert:
			watch.Statistics.AlertsSent++
			watch.Statistics.ChangedSubdomains += len(changedSubdomains)
		default:
			gologger.Warning().Msg("Alert queue full, dropping alert")
		}
	}

	gologger.Debug().Msgf("📡 Monitor: Check completed for %s - New: %d, Changed: %d", 
		domain, len(newSubdomains), len(changedSubdomains))
}

// processAlerts processes alerts and sends notifications
func (rm *RealtimeMonitor) processAlerts(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case alert := <-rm.alerts:
			rm.handleAlert(alert)
		}
	}
}

// handleAlert handles a single alert
func (rm *RealtimeMonitor) handleAlert(alert Alert) {
	gologger.Info().Msgf("🚨 ALERT [%s]: %s", alert.Severity, alert.Message)

	// Log alert details
	if len(alert.NewSubdomains) > 0 {
		gologger.Info().Msgf("   New subdomains: %v", alert.NewSubdomains)
	}

	if len(alert.ChangedSubdomains) > 0 {
		gologger.Info().Msgf("   Changed subdomains: %d", len(alert.ChangedSubdomains))
	}

	// Send to webhooks
	for _, webhook := range rm.webhooks {
		if webhook.Enabled {
			go rm.sendWebhook(webhook, alert)
		}
	}

	// Update domain statistics
	rm.mutex.Lock()
	if watch, exists := rm.domains[alert.Domain]; exists {
		watch.Statistics.LastAlert = alert.Timestamp
	}
	rm.mutex.Unlock()
}

// sendWebhook sends an alert to a webhook
func (rm *RealtimeMonitor) sendWebhook(webhook WebhookConfig, alert Alert) {
	// This would send HTTP POST to webhook URL
	// Simplified implementation
	gologger.Debug().Msgf("📡 Monitor: Sending alert to webhook %s", webhook.URL)
}

// Helper methods

func (rm *RealtimeMonitor) enumerateSubdomains(ctx context.Context, domain string) []string {
	// This would perform actual subdomain enumeration
	// For demonstration, return some sample subdomains
	return []string{
		"www." + domain,
		"api." + domain,
		"mail." + domain,
	}
}

func (rm *RealtimeMonitor) findNewSubdomains(known map[string]bool, current []string) []string {
	var newSubdomains []string
	for _, subdomain := range current {
		if !known[subdomain] {
			newSubdomains = append(newSubdomains, subdomain)
		}
	}
	return newSubdomains
}

func (rm *RealtimeMonitor) findChangedSubdomains(ctx context.Context, domain string, subdomains []string) []SubdomainChange {
	var changes []SubdomainChange
	
	// This would check for IP changes, certificate changes, etc.
	// Simplified implementation
	
	return changes
}

func (rm *RealtimeMonitor) calculateSeverity(count int) string {
	if count >= 10 {
		return "high"
	} else if count >= 5 {
		return "medium"
	}
	return "low"
}

func (rm *RealtimeMonitor) getDomainCheckInterval(domain string) time.Duration {
	rm.mutex.RLock()
	defer rm.mutex.RUnlock()
	
	if watch, exists := rm.domains[domain]; exists {
		return watch.CheckInterval
	}
	return 5 * time.Minute // Default
}

// GetStatus returns the current monitoring status
func (rm *RealtimeMonitor) GetStatus() map[string]interface{} {
	rm.mutex.RLock()
	defer rm.mutex.RUnlock()

	status := map[string]interface{}{
		"running":        rm.running,
		"domains_count":  len(rm.domains),
		"webhooks_count": len(rm.webhooks),
		"alert_queue":    len(rm.alerts),
	}

	domains := make(map[string]interface{})
	for name, watch := range rm.domains {
		domains[name] = map[string]interface{}{
			"last_check":       watch.LastCheck,
			"known_subdomains": len(watch.KnownSubdomains),
			"total_checks":     watch.Statistics.TotalChecks,
			"alerts_sent":      watch.Statistics.AlertsSent,
			"uptime":          watch.Statistics.Uptime.String(),
		}
	}
	status["domains"] = domains

	return status
}

// GetAlertHistory returns recent alerts
func (rm *RealtimeMonitor) GetAlertHistory(domain string, limit int) []Alert {
	// This would return stored alert history
	// Placeholder implementation
	return []Alert{}
}

// ExportConfiguration exports monitoring configuration
func (rm *RealtimeMonitor) ExportConfiguration() ([]byte, error) {
	rm.mutex.RLock()
	defer rm.mutex.RUnlock()

	config := map[string]interface{}{
		"domains":  rm.domains,
		"webhooks": rm.webhooks,
	}

	return json.MarshalIndent(config, "", "  ")
}

// ImportConfiguration imports monitoring configuration
func (rm *RealtimeMonitor) ImportConfiguration(data []byte) error {
	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		return err
	}

	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	// Import domains and webhooks
	// Simplified implementation

	return nil
}
