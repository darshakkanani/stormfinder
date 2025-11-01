package passive

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/exp/maps"

	"github.com/projectdiscovery/gologger"
	"github.com/darshakkanani/stormfinder/v2/pkg/subscraping"
	"github.com/darshakkanani/stormfinder/v2/pkg/subscraping/sources/alienvault"
	"github.com/darshakkanani/stormfinder/v2/pkg/subscraping/sources/anubis"
	"github.com/darshakkanani/stormfinder/v2/pkg/subscraping/sources/bevigil"
	"github.com/darshakkanani/stormfinder/v2/pkg/subscraping/sources/bufferover"
	"github.com/darshakkanani/stormfinder/v2/pkg/subscraping/sources/builtwith"
	"github.com/darshakkanani/stormfinder/v2/pkg/subscraping/sources/c99"
	"github.com/darshakkanani/stormfinder/v2/pkg/subscraping/sources/censys"
	"github.com/darshakkanani/stormfinder/v2/pkg/subscraping/sources/certspotter"
	"github.com/darshakkanani/stormfinder/v2/pkg/subscraping/sources/chaos"
	"github.com/darshakkanani/stormfinder/v2/pkg/subscraping/sources/chinaz"
	"github.com/darshakkanani/stormfinder/v2/pkg/subscraping/sources/commoncrawl"
	"github.com/darshakkanani/stormfinder/v2/pkg/subscraping/sources/crtsh"
	"github.com/darshakkanani/stormfinder/v2/pkg/subscraping/sources/digitalyama"
	"github.com/darshakkanani/stormfinder/v2/pkg/subscraping/sources/digitorus"
	"github.com/darshakkanani/stormfinder/v2/pkg/subscraping/sources/dnsdb"
	"github.com/darshakkanani/stormfinder/v2/pkg/subscraping/sources/dnsdumpster"
	"github.com/darshakkanani/stormfinder/v2/pkg/subscraping/sources/dnsrepo"
	"github.com/darshakkanani/stormfinder/v2/pkg/subscraping/sources/driftnet"
	"github.com/darshakkanani/stormfinder/v2/pkg/subscraping/sources/facebook"
	"github.com/darshakkanani/stormfinder/v2/pkg/subscraping/sources/fofa"
	"github.com/darshakkanani/stormfinder/v2/pkg/subscraping/sources/fullhunt"
	"github.com/darshakkanani/stormfinder/v2/pkg/subscraping/sources/github"
	"github.com/darshakkanani/stormfinder/v2/pkg/subscraping/sources/hackertarget"
	"github.com/darshakkanani/stormfinder/v2/pkg/subscraping/sources/hudsonrock"
	"github.com/darshakkanani/stormfinder/v2/pkg/subscraping/sources/intelx"
	"github.com/darshakkanani/stormfinder/v2/pkg/subscraping/sources/leakix"
	"github.com/darshakkanani/stormfinder/v2/pkg/subscraping/sources/netlas"
	"github.com/darshakkanani/stormfinder/v2/pkg/subscraping/sources/onyphe"
	"github.com/darshakkanani/stormfinder/v2/pkg/subscraping/sources/pugrecon"
	"github.com/darshakkanani/stormfinder/v2/pkg/subscraping/sources/quake"
	"github.com/darshakkanani/stormfinder/v2/pkg/subscraping/sources/rapiddns"
	"github.com/darshakkanani/stormfinder/v2/pkg/subscraping/sources/redhuntlabs"
	"github.com/darshakkanani/stormfinder/v2/pkg/subscraping/sources/robtex"
	"github.com/darshakkanani/stormfinder/v2/pkg/subscraping/sources/rsecloud"
	"github.com/darshakkanani/stormfinder/v2/pkg/subscraping/sources/securitytrails"
	"github.com/darshakkanani/stormfinder/v2/pkg/subscraping/sources/shodan"
	"github.com/darshakkanani/stormfinder/v2/pkg/subscraping/sources/sitedossier"
	"github.com/darshakkanani/stormfinder/v2/pkg/subscraping/sources/threatbook"
	"github.com/darshakkanani/stormfinder/v2/pkg/subscraping/sources/threatcrowd"
	"github.com/darshakkanani/stormfinder/v2/pkg/subscraping/sources/urlscan"
	"github.com/darshakkanani/stormfinder/v2/pkg/subscraping/sources/virustotal"
	"github.com/darshakkanani/stormfinder/v2/pkg/subscraping/sources/wayback"
	"github.com/darshakkanani/stormfinder/v2/pkg/subscraping/sources/waybackarchive"
	"github.com/darshakkanani/stormfinder/v2/pkg/subscraping/sources/whoisxmlapi"
	"github.com/darshakkanani/stormfinder/v2/pkg/subscraping/sources/windvane"
	"github.com/darshakkanani/stormfinder/v2/pkg/subscraping/sources/zoomeyeapi"
	mapsutil "github.com/projectdiscovery/utils/maps"
)

var AllSources = [...]subscraping.Source{
	&alienvault.Source{},
	&anubis.Source{},
	&bevigil.Source{},
	&bufferover.Source{},
	&c99.Source{},
	&censys.Source{},
	&certspotter.Source{},
	&chaos.Source{},
	&chinaz.Source{},
	&commoncrawl.Source{},
	&crtsh.Source{},
	&digitorus.Source{},
	&dnsdb.Source{},
	&dnsdumpster.Source{},
	&dnsrepo.Source{},
	&driftnet.Source{},
	&fofa.Source{},
	&fullhunt.Source{},
	&github.Source{},
	&hackertarget.Source{},
	&intelx.Source{},
	&netlas.Source{},
	&onyphe.Source{},
	&leakix.Source{},
	&quake.Source{},
	&pugrecon.Source{},
	&rapiddns.Source{},
	&redhuntlabs.Source{},
	// &riddler.Source{}, // failing due to cloudfront protection
	&robtex.Source{},
	&rsecloud.Source{},
	&securitytrails.Source{},
	&shodan.Source{},
	&sitedossier.Source{},
	&threatbook.Source{},
	&threatcrowd.Source{},
	&virustotal.Source{},
	&waybackarchive.Source{},
	&whoisxmlapi.Source{},
	&windvane.Source{},
	&zoomeyeapi.Source{},
	&facebook.Source{},
	// &threatminer.Source{}, // failing  api
	// &reconcloud.Source{}, // failing due to cloudflare bot protection
	&builtwith.Source{},
	&hudsonrock.Source{},
	&digitalyama.Source{},
	&urlscan.Source{},
	&wayback.Source{},
}

var sourceWarnings = mapsutil.NewSyncLockMap[string, string](
	mapsutil.WithMap(mapsutil.Map[string, string]{}))

var NameSourceMap = make(map[string]subscraping.Source, len(AllSources))

func init() {
	for _, currentSource := range AllSources {
		NameSourceMap[strings.ToLower(currentSource.Name())] = currentSource
	}
}

// Agent is a struct for running passive subdomain enumeration
// against a given host. It wraps subscraping package and provides
// a layer to build upon.
type Agent struct {
	sources []subscraping.Source
}

// New creates a new agent for passive subdomain discovery
func New(sourceNames, excludedSourceNames []string, useAllSources, useSourcesSupportingRecurse bool) *Agent {
	sources := make(map[string]subscraping.Source, len(AllSources))

	if useAllSources {
		maps.Copy(sources, NameSourceMap)
	} else {
		if len(sourceNames) > 0 {
			for _, source := range sourceNames {
				if NameSourceMap[source] == nil {
					gologger.Warning().Msgf("There is no source with the name: %s", source)
				} else {
					sources[source] = NameSourceMap[source]
				}
			}
		} else {
			for _, currentSource := range AllSources {
				if currentSource.IsDefault() {
					sources[currentSource.Name()] = currentSource
				}
			}
		}
	}

	if len(excludedSourceNames) > 0 {
		for _, sourceName := range excludedSourceNames {
			delete(sources, sourceName)
		}
	}

	if useSourcesSupportingRecurse {
		for sourceName, source := range sources {
			if !source.HasRecursiveSupport() {
				delete(sources, sourceName)
			}
		}
	}

	if len(sources) == 0 {
		gologger.Fatal().Msg("No sources selected for this search")
	}

	gologger.Debug().Msgf("Selected source(s) for this search: %s", strings.Join(maps.Keys(sources), ", "))

	for _, currentSource := range sources {
		if warning, ok := sourceWarnings.Get(strings.ToLower(currentSource.Name())); ok {
			gologger.Warning().Msg(warning)
		}
	}

	// TODO: Consider refactoring this to avoid potential duplication issues
	for _, source := range sources {
		if source.NeedsKey() {
			if apiKey := os.Getenv(fmt.Sprintf("%s_API_KEY", strings.ToUpper(source.Name()))); apiKey != "" {
				source.AddApiKeys([]string{apiKey})
			}
		}
	}

	// Create the agent, insert the sources and remove the excluded sources
	agent := &Agent{sources: maps.Values(sources)}

	return agent
}
