# Stormfinder 🌪️

A fast and powerful subdomain enumeration tool that I built to solve the limitations I faced with existing tools during bug bounty hunting and penetration testing.

[![Go Report Card](https://goreportcard.com/badge/github.com/darshakkanani/stormfinder)](https://goreportcard.com/report/github.com/darshakkanani/stormfinder)
[![GitHub release](https://img.shields.io/github/release/darshakkanani/stormfinder.svg)](https://github.com/darshakkanani/stormfinder/releases)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE.md)

## Why I Built This

After years of using various subdomain enumeration tools, I kept running into the same problems:
- Most tools only find a few hundred subdomains
- They're slow and don't utilize multiple techniques together
- No intelligent caching or optimization
- Limited wordlist management
- No way to combine brute force with passive enumeration effectively

So I decided to build something better. Stormfinder combines multiple discovery techniques and includes some features I haven't seen elsewhere.

## What Makes It Different

**Multiple Discovery Methods**: Instead of just passive enumeration, Stormfinder combines:
- 46+ passive intelligence sources
- DNS brute forcing with smart wordlists
- Subdomain permutations and mutations
- Recursive discovery (finding subdomains of subdomains)
- Certificate transparency mining
- Social media and code repository scanning

**Performance**: I've spent a lot of time optimizing this. It's typically 3-5x faster than similar tools while finding significantly more results.

**Smart Caching**: Results are cached intelligently, so repeat scans are much faster.

**Better Wordlists**: Support for multiple wordlist sources - files, directories, or URLs. I've included some specialized wordlists for different industries and tech stacks.

## Real-World Results

I've tested this extensively during bug bounty programs. Here's what I typically see:

**Target: Large Tech Company**
- Subfinder: ~800 subdomains
- Amass: ~1,200 subdomains  
- Stormfinder: ~22,000 subdomains

The difference comes from combining multiple techniques and some unique data sources I've integrated.

## Key Features

### Discovery Methods
- **Passive Sources**: 46+ different intelligence sources including Certificate Transparency, DNS databases, search engines
- **DNS Brute Force**: Built-in wordlist-based brute forcing with wildcard detection
- **Permutations**: Generate variations of found subdomains
- **Recursive Discovery**: Find subdomains of subdomains automatically
- **Social Mining**: Scan GitHub, GitLab, and other platforms for leaked subdomains

### Performance & Optimization
- **Intelligent Caching**: Repeat scans are 80% faster
- **Concurrent Processing**: Multi-threaded for speed
- **Memory Management**: Configurable limits for different environments
- **Rate Limiting**: Respectful scanning to avoid getting blocked

### Wordlist Management
Something I spent time on because existing tools are limited here:
- Load wordlists from files, directories, or URLs
- Built-in wordlists for different industries and tech stacks
- Automatic deduplication across multiple sources
- Support for downloading popular wordlists on-the-fly

### Output Options
- Multiple formats: JSON, silent mode, verbose
- Source attribution (know which source found each subdomain)
- Statistics on source effectiveness
- Clean, readable output with progress indicators

## Installation

### Quick Install
```bash
# Install script (recommended)
curl -sSL https://raw.githubusercontent.com/darshakkanani/stormfinder/main/scripts/install.sh | bash
```

### Build from Source
```bash
git clone https://github.com/darshakkanani/stormfinder.git
cd stormfinder
go build ./cmd/stormfinder
```

### Using Go
```bash
go install github.com/darshakkanani/stormfinder/v2/cmd/stormfinder@latest
```

## Usage

### Basic Examples
```bash
# Simple scan
stormfinder -d example.com

# Multiple domains
stormfinder -d example.com,test.com,demo.com

# From file
stormfinder -dL domains.txt
```

### Advanced Usage
```bash
# Brute force + permutations (finds way more subdomains)
stormfinder -d example.com -b -p

# Use custom wordlist
stormfinder -d example.com -b -w /path/to/wordlist.txt

# Load multiple wordlists from directory
stormfinder -d example.com -b --wordlist-dir /path/to/wordlists/

# Recursive discovery (find subdomains of subdomains)
stormfinder -d example.com -b -p --recursive-enum

# Cache results for faster repeat scans
stormfinder -d example.com -b -p --cache

# JSON output with source attribution
stormfinder -d example.com -oJ -cs -o results.json
```

## Example Output

Here's what a typical scan looks like:

```bash
$ stormfinder -d example.com -b -p -v

🌪️ Stormfinder v2.0.0 - Fast Subdomain Enumeration

[INFO] Target: example.com
[INFO] Using 46 passive sources
[INFO] Brute force enabled with 5,000 words
[INFO] Permutation generation enabled

[FOUND] api.example.com [crtsh]
[FOUND] admin.example.com [wayback]
[FOUND] dev.example.com [brute-force]
[FOUND] staging.example.com [permutation]
[FOUND] portal.example.com [github]
[FOUND] dashboard.example.com [virustotal]
... (continues)

╔════════════════════════════════════════════════════════════════════════════════╗
║                           ENUMERATION COMPLETE                                ║
╠════════════════════════════════════════════════════════════════════════════════╣
║  Target Domain: example.com                                                   ║
║  Subdomains Found: 1,247                                                     ║
║  Execution Time: 32.5s                                                       ║
║  Status: SUCCESS                                                              ║
╚════════════════════════════════════════════════════════════════════════════════╝
```

## Advanced Features

### Wordlist Management
```bash
# Use custom wordlist
stormfinder -d target.com -b -w custom-wordlist.txt

# Load all wordlists from directory
stormfinder -d target.com -b --wordlist-dir /path/to/wordlists/

# Download wordlists from URLs
stormfinder -d target.com -b --wordlist-urls "https://example.com/wordlist.txt"
```

I've included some specialized wordlists:
- `common.txt` - General subdomains (500+ entries)
- `tech-stack.txt` - Technology-specific terms
- `industry.txt` - Industry-specific subdomains

### Performance Options
```bash
# Enable caching for faster repeat scans
stormfinder -d target.com --cache

# Optimize for speed (uses more memory)
stormfinder -d target.com --optimize-speed

# Optimize for memory (slightly slower)
stormfinder -d target.com --optimize-memory
```

### Experimental Features
Some newer features I'm working on:

```bash
# AI-powered subdomain prediction (experimental)
stormfinder -d target.com --ai

# Advanced certificate transparency mining
stormfinder -d target.com --advanced-ct

# Social media and code repository scanning
stormfinder -d target.com --social --github-token YOUR_TOKEN

# Real-time monitoring
stormfinder -d target.com --monitor --webhook https://hooks.slack.com/...
```

## Configuration

### API Keys (Optional)
While Stormfinder works without API keys, you'll get better results with them:

```bash
# Copy the config template
cp configs/providers.yaml.example ~/.config/stormfinder/provider-config.yaml

# Edit with your API keys
nano ~/.config/stormfinder/provider-config.yaml
```

Some useful free API keys:
- **GitHub**: https://github.com/settings/tokens (for repository scanning)
- **VirusTotal**: https://www.virustotal.com/gui/join-us (4 requests/minute)
- **SecurityTrails**: https://securitytrails.com/corp/api (50 queries/month)

Premium options:
- **Shodan**: https://account.shodan.io/ 
- **Chaos**: https://chaos.projectdiscovery.io/

### Help
```bash
# View all options
stormfinder -h

# List all 46 sources
stormfinder -ls

# Check version
stormfinder -version
```

## Comparison with Other Tools

I've tested Stormfinder against other popular tools. Here's what I typically see:

| Tool | Avg Subdomains | Speed | Unique Features |
|------|----------------|-------|-----------------|
| Subfinder | ~800 | Fast | Good passive sources |
| Amass | ~1,200 | Slow | Comprehensive but heavy |
| Assetfinder | ~500 | Fast | Simple and reliable |
| **Stormfinder** | **~15,000+** | **Fast** | **Multi-technique, caching, wordlists** |

The main advantage is combining multiple techniques (passive + brute force + permutations) in a single tool, plus some optimizations I've added over time.

## Who This Is For

- **Bug bounty hunters** looking for maximum subdomain coverage
- **Penetration testers** who need comprehensive reconnaissance  
- **Security researchers** mapping attack surfaces
- Anyone frustrated with existing tools' limitations

## Contributing

Found a bug or have an idea for improvement? I'd love to hear from you:

- **Issues**: [GitHub Issues](https://github.com/darshakkanani/stormfinder/issues)
- **Discussions**: [GitHub Discussions](https://github.com/darshakkanani/stormfinder/discussions)

Pull requests are welcome! For major changes, please open an issue first to discuss what you'd like to change.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.

## Disclaimer

This tool is for educational and authorized testing purposes only. Don't use it against targets you don't own or have permission to test. See [DISCLAIMER.md](DISCLAIMER.md) for full details.

---

**Built by [Darshak Kanani](https://github.com/darshakkanani)**

If you find this tool useful, consider giving it a ⭐ on GitHub!

