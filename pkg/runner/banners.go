package runner

import (
	"github.com/projectdiscovery/gologger"
	updateutils "github.com/projectdiscovery/utils/update"
)

const banner = `
🌪️ Stormfinder v2.9.0 - Fast Subdomain Enumeration

A powerful subdomain discovery tool that combines multiple techniques:
• 46+ passive intelligence sources
• DNS brute forcing with smart wordlists  
• Subdomain permutations and mutations
• Recursive discovery and caching
• Social media and code repository scanning
`

// Name
const ToolName = `stormfinder`

// Version is the current version of stormfinder
const version = `v2.9.0`

// showBanner is used to show the banner to the user
func showBanner() {
	gologger.Print().Msgf("%s", banner)
	gologger.Print().Msgf("                         github.com/darshakkanani/stormfinder\n\n")
}

// GetUpdateCallback returns a callback function that updates stormfinder
func GetUpdateCallback() func() {
	return func() {
		showBanner()
		updateutils.GetUpdateToolCallback("stormfinder", version)()
	}
}
