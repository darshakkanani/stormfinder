# 🚀 Stormfinder Installation Guide

## 📋 **Prerequisites**

- **Go**: Version 1.24 or higher
- **Git**: For cloning the repository
- **Internet Connection**: For passive enumeration sources

## 🔧 **Installation Methods**

### **Method 1: From Source (Recommended)**

```bash
# Clone the repository
git clone https://github.com/darshakkanani/stormfinder.git
cd stormfinder

# Build the binary
go build ./cmd/stormfinder

# Make it executable
chmod +x stormfinder

# Optional: Move to PATH
sudo mv stormfinder /usr/local/bin/
```

### **Method 2: Go Install**

```bash
go install github.com/darshakkanani/stormfinder/v2/cmd/stormfinder@latest
```

### **Method 3: Download Binary** *(Coming Soon)*

Pre-compiled binaries will be available in the releases section.

## ⚙️ **Configuration**

### **Basic Setup**

Stormfinder works out of the box with default settings. No configuration required for basic usage!

```bash
# Test installation
stormfinder -h

# Basic enumeration
stormfinder -d example.com
```

### **Advanced Configuration**

#### **Provider Configuration**

Create a provider configuration file for API keys:

```bash
# Create config directory
mkdir -p ~/.config/stormfinder

# Create provider config
cat > ~/.config/stormfinder/provider-config.yaml << EOF
shodan:
  - YOUR_SHODAN_API_KEY
virustotal:
  - YOUR_VIRUSTOTAL_API_KEY
github:
  - YOUR_GITHUB_TOKEN
securitytrails:
  - YOUR_SECURITYTRAILS_API_KEY
chaos:
  - YOUR_CHAOS_API_KEY
EOF
```

#### **Main Configuration**

```bash
# Create main config
cat > ~/.config/stormfinder/config.yaml << EOF
# Default settings
threads: 25
timeout: 30
max-time: 10
cache: true
cache-ttl: 24
optimize-speed: true
EOF
```

## 🧪 **Verification**

Run the comprehensive test suite:

```bash
# Run all feature tests
./test_all_features.sh
```

Expected output:
```
🎉 ALL TESTS PASSED! Stormfinder is ready for GitHub release! 🚀
```

## 🚀 **Quick Start Examples**

### **Basic Discovery**
```bash
# Simple enumeration
stormfinder -d target.com

# Multiple domains
stormfinder -d target1.com,target2.com

# From file
echo "target.com" > domains.txt
stormfinder -dL domains.txt
```

### **Enhanced Discovery**
```bash
# Brute force + permutations
stormfinder -d target.com -b -p

# AI-powered discovery
stormfinder -d target.com --ai

# Full power scan
stormfinder -d target.com -b -p --ai --advanced-ct --social --map
```

### **Output Options**
```bash
# Save to file
stormfinder -d target.com -o results.txt

# JSON output with sources
stormfinder -d target.com -oJ -cs -o results.json

# Silent mode for piping
stormfinder -d target.com -silent | grep "api\."
```

## 🔍 **API Keys Setup** *(Optional but Recommended)*

### **Free API Keys**

1. **GitHub**: https://github.com/settings/tokens
   - Scope: `public_repo` (for public repositories)

2. **VirusTotal**: https://www.virustotal.com/gui/join-us
   - Free tier: 4 requests/minute

3. **SecurityTrails**: https://securitytrails.com/corp/api
   - Free tier: 50 queries/month

### **Paid API Keys** *(For Enhanced Results)*

1. **Shodan**: https://account.shodan.io/
   - Paid plans for unlimited queries

2. **Chaos**: https://chaos.projectdiscovery.io/
   - ProjectDiscovery's subdomain dataset

## 🐛 **Troubleshooting**

### **Common Issues**

#### **Build Errors**
```bash
# Update Go modules
go mod tidy

# Clean build
go clean -cache
go build ./cmd/stormfinder
```

#### **Permission Denied**
```bash
# Make executable
chmod +x stormfinder

# Or run with go
go run ./cmd/stormfinder -d example.com
```

#### **No Results Found**
```bash
# Check internet connection
ping google.com

# Test with verbose mode
stormfinder -d example.com -v

# List available sources
stormfinder -ls
```

#### **Rate Limiting**
```bash
# Reduce rate limit
stormfinder -d target.com -rl 5

# Use fewer threads
stormfinder -d target.com -t 5
```

## 📊 **Performance Tuning**

### **For Speed**
```bash
stormfinder -d target.com --optimize-speed --cache --brute-threads 100
```

### **For Memory**
```bash
stormfinder -d target.com --optimize-memory --max-memory 256
```

### **For Large Scans**
```bash
stormfinder -d target.com -b -p --cache --max-time 60 --brute-threads 50
```

## 🔄 **Updates**

### **Update from Source**
```bash
cd stormfinder
git pull origin main
go build ./cmd/stormfinder
```

### **Check for Updates**
```bash
stormfinder -version
stormfinder -up
```

## 🆘 **Support**

- **Issues**: https://github.com/darshakkanani/stormfinder/issues
- **Discussions**: https://github.com/darshakkanani/stormfinder/discussions
- **Documentation**: See `FEATURES.md` for complete feature list

## 🎯 **Next Steps**

1. **Read the Features**: Check `FEATURES.md` for complete capabilities
2. **Join the Community**: Star the repository and follow for updates
3. **Contribute**: Submit issues, feature requests, or pull requests
4. **Share**: Help others discover this powerful tool

**Happy hunting! 🌪️**
