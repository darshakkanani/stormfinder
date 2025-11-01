# 🎯 Basic Usage Examples

## 🚀 **Quick Start**

### **Simple Enumeration**
```bash
# Basic subdomain discovery
stormfinder -d target.com

# Multiple domains
stormfinder -d target1.com,target2.com,target3.com

# From file
echo "target.com" > domains.txt
stormfinder -dL domains.txt
```

### **Enhanced Discovery**
```bash
# Brute force enumeration
stormfinder -d target.com -b

# Permutation generation
stormfinder -d target.com -p

# Combined enhanced techniques
stormfinder -d target.com -b -p --recursive-enum
```

## 📊 **Output Options**

### **File Output**
```bash
# Save to text file
stormfinder -d target.com -o results.txt

# JSON output with source attribution
stormfinder -d target.com -oJ -cs -o results.json

# Directory output for multiple domains
stormfinder -dL domains.txt -oD ./results/
```

### **Display Options**
```bash
# Silent mode (only subdomains)
stormfinder -d target.com -silent

# Verbose mode (detailed progress)
stormfinder -d target.com -v

# No colors (for scripts)
stormfinder -d target.com -nc
```

## ⚡ **Performance Tuning**

### **Speed Optimization**
```bash
# High-speed scan
stormfinder -d target.com --optimize-speed --brute-threads 100

# With caching for repeated scans
stormfinder -d target.com --cache --cache-ttl 48
```

### **Resource Management**
```bash
# Memory-optimized scan
stormfinder -d target.com --optimize-memory --max-memory 256

# Rate limiting for respectful scanning
stormfinder -d target.com -rl 10 -t 5
```

## 🔍 **Source Management**

### **Source Selection**
```bash
# List all available sources
stormfinder -ls

# Use specific sources only
stormfinder -d target.com -s crtsh,hackertarget,anubis

# Exclude specific sources
stormfinder -d target.com -es alienvault,zoomeyeapi

# Use all sources (slow but comprehensive)
stormfinder -d target.com -all
```

### **Source Statistics**
```bash
# Show source effectiveness
stormfinder -d target.com -stats

# Verbose with source attribution
stormfinder -d target.com -v -cs
```

## 🎯 **Filtering & Matching**

### **Pattern Matching**
```bash
# Match specific patterns
stormfinder -d target.com -m "api.*,admin.*"

# Filter out patterns
stormfinder -d target.com -f "test.*,dev.*"

# Match from file
echo "api" > patterns.txt
echo "admin" >> patterns.txt
stormfinder -d target.com -m patterns.txt
```

### **IP Resolution**
```bash
# Include IP addresses
stormfinder -d target.com -nW -oI

# Exclude IP addresses from output
stormfinder -d target.com -ei
```

## 🌐 **Network Configuration**

### **Proxy Usage**
```bash
# HTTP proxy
stormfinder -d target.com -proxy http://proxy:8080

# SOCKS proxy
stormfinder -d target.com -proxy socks5://proxy:1080
```

### **Custom DNS**
```bash
# Custom resolvers
stormfinder -d target.com -r 8.8.8.8,1.1.1.1

# Resolver list from file
echo "8.8.8.8" > resolvers.txt
echo "1.1.1.1" >> resolvers.txt
stormfinder -d target.com -rL resolvers.txt
```

## 📋 **Common Workflows**

### **Bug Bounty Hunting**
```bash
# Comprehensive discovery
stormfinder -d target.com -b -p --recursive-enum --cache -v -o bounty-results.txt

# Quick reconnaissance
stormfinder -d target.com --optimize-speed -silent | grep -E "(api|admin|dev|staging)"
```

### **Security Assessment**
```bash
# Detailed enumeration with sources
stormfinder -d target.com -b -p -v -cs -oJ -o assessment.json

# Active subdomain verification
stormfinder -d target.com -nW -oI -o active-subdomains.txt
```

### **Monitoring Setup**
```bash
# Regular monitoring scan
stormfinder -d target.com --cache -silent -o current-scan.txt

# Compare with previous scan
diff previous-scan.txt current-scan.txt | grep "^>" | sed 's/^> //' > new-subdomains.txt
```

## 🔧 **Configuration Examples**

### **Using Config Files**
```bash
# Custom configuration
stormfinder -config /path/to/config.yaml -d target.com

# Custom provider config
stormfinder -pc /path/to/providers.yaml -d target.com
```

### **Environment Variables**
```bash
# Set API keys via environment
export STORMFINDER_GITHUB_TOKEN="your_token_here"
export STORMFINDER_SHODAN_KEY="your_key_here"
stormfinder -d target.com
```

## 📊 **Output Analysis**

### **Processing Results**
```bash
# Count subdomains
stormfinder -d target.com -silent | wc -l

# Extract unique second-level domains
stormfinder -d target.com -silent | cut -d. -f2- | sort -u

# Find API endpoints
stormfinder -d target.com -silent | grep -i api

# Find development/staging environments
stormfinder -d target.com -silent | grep -E "(dev|test|staging|beta)"
```

### **Integration with Other Tools**
```bash
# Pipe to httpx for HTTP probing
stormfinder -d target.com -silent | httpx -silent

# Pipe to nmap for port scanning
stormfinder -d target.com -silent | nmap -iL - -p 80,443

# Pipe to nuclei for vulnerability scanning
stormfinder -d target.com -silent | nuclei -t vulnerabilities/
```

## 🚨 **Troubleshooting**

### **Common Issues**
```bash
# No results found
stormfinder -d target.com -v  # Check verbose output

# Rate limiting issues
stormfinder -d target.com -rl 5 -t 3  # Reduce rate and threads

# Memory issues
stormfinder -d target.com --optimize-memory --max-memory 256

# Network issues
stormfinder -d target.com -proxy http://proxy:8080 -timeout 60
```

### **Debugging**
```bash
# Maximum verbosity
stormfinder -d target.com -v -stats

# Test specific sources
stormfinder -d target.com -s crtsh -v

# Check configuration
stormfinder -version
stormfinder -ls
```
