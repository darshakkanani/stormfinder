package runner

import (
	"github.com/projectdiscovery/gologger"
	updateutils "github.com/projectdiscovery/utils/update"
)

const banner = `
╔═══════════════════════════════════════════════════════════════════════════════╗
║                                                                               ║
║   ███████╗████████╗ ██████╗ ██████╗ ███╗   ███╗███████╗██╗███╗   ██╗██████╗  ║
║   ██╔════╝╚══██╔══╝██╔═══██╗██╔══██╗████╗ ████║██╔════╝██║████╗  ██║██╔══██╗ ║
║   ███████╗   ██║   ██║   ██║██████╔╝██╔████╔██║█████╗  ██║██╔██╗ ██║██║  ██║ ║
║   ╚════██║   ██║   ██║   ██║██╔══██╗██║╚██╔╝██║██╔══╝  ██║██║╚██╗██║██║  ██║ ║
║   ███████║   ██║   ╚██████╔╝██║  ██║██║ ╚═╝ ██║██║     ██║██║ ╚████║██████╔╝ ║
║   ╚══════╝   ╚═╝    ╚═════╝ ╚═╝  ╚═╝╚═╝     ╚═╝╚═╝     ╚═╝╚═╝  ╚═══╝╚═════╝  ║
║                                                                               ║
║           🚀 Next-Gen AI-Powered Subdomain Discovery Platform 🚀              ║
║                                                                               ║
║  💡 Features: AI Prediction | Advanced CT Mining | Social Mining | Mapping   ║
║  ⚡ Enhanced: Real-time Monitoring | Relationship Analysis | ML Classification║
║                                                                               ║
╚═══════════════════════════════════════════════════════════════════════════════╝
`

// Name
const ToolName = `stormfinder`

// Version is the current version of stormfinder
const version = `v2.9.0`

// showBanner is used to show the banner to the user
func showBanner() {
	gologger.Print().Msgf("%s\n", banner)
	gologger.Print().Msgf("                    🌟 Enhanced by Advanced AI & ML Technologies 🌟\n")
	gologger.Print().Msgf("                         github.com/darshakkanani/stormfinder\n\n")
}

// GetUpdateCallback returns a callback function that updates stormfinder
func GetUpdateCallback() func() {
	return func() {
		showBanner()
		updateutils.GetUpdateToolCallback("stormfinder", version)()
	}
}
