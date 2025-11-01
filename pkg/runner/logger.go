package runner

import (
	"fmt"
	"strings"
	"time"

	"github.com/projectdiscovery/gologger"
)

// Custom logging functions with cooler prefixes
func LogInfo(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	gologger.Print().Msgf("🔍 %s", message)
}

func LogSuccess(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	gologger.Print().Msgf("✅ %s", message)
}

func LogProgress(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	gologger.Print().Msgf("⚡ %s", message)
}

func LogDiscovery(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	gologger.Print().Msgf("🎯 %s", message)
}

func LogEnhanced(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	gologger.Print().Msgf("🚀 %s", message)
}

func LogAI(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	gologger.Print().Msgf("🤖 %s", message)
}

func LogSocial(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	gologger.Print().Msgf("📱 %s", message)
}

func LogCT(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	gologger.Print().Msgf("🔐 %s", message)
}

func LogMapping(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	gologger.Print().Msgf("🗺️  %s", message)
}

func LogMonitor(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	gologger.Print().Msgf("📡 %s", message)
}

func LogStats(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	gologger.Print().Msgf("📊 %s", message)
}

// LogResults shows final results with enhanced formatting
func LogResults(domain string, count int, duration time.Duration) {
	// Create a fancy results display
	border := strings.Repeat("═", 80)
	
	gologger.Print().Msgf("╔%s╗", border)
	gologger.Print().Msgf("║                           🎉 ENUMERATION COMPLETE 🎉                        ║")
	gologger.Print().Msgf("╠%s╣", border)
	gologger.Print().Msgf("║  🎯 Target Domain: %-58s ║", domain)
	gologger.Print().Msgf("║  📊 Subdomains Found: %-51d ║", count)
	gologger.Print().Msgf("║  ⏱️  Execution Time: %-53s ║", duration.String())
	gologger.Print().Msgf("║  🚀 Status: SUCCESS - All enumeration techniques completed successfully     ║")
	gologger.Print().Msgf("╚%s╝", border)
}

// LogStartup shows startup information with style
func LogStartup(domain string) {
	gologger.Print().Msgf("🌟 Initializing advanced subdomain enumeration for: %s", domain)
	gologger.Print().Msgf("⚡ Loading passive sources and enhanced discovery engines...")
}

// LogConfig shows configuration loading
func LogConfig(configPath string) {
	gologger.Print().Msgf("⚙️  Loading configuration from: %s", configPath)
}

// LogSources shows source statistics
func LogSources(count int) {
	gologger.Print().Msgf("🔧 Available enumeration sources: %d", count)
	gologger.Print().Msgf("💡 Sources marked with (*) require API keys for enhanced functionality")
}
