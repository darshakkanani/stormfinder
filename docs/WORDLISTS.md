# 📝 Stormfinder Wordlist Management

## 🚀 **Enhanced Wordlist Capabilities**

Stormfinder now supports multiple advanced wordlist options for maximum subdomain discovery effectiveness.

## 📋 **Wordlist Options**

### **📄 Single Wordlist File**
```bash
# Use a custom wordlist file
stormfinder -d target.com -b -w /path/to/custom-wordlist.txt

# Use with verbose output to see loading progress
stormfinder -d target.com -b -w custom.txt -v
```

### **📁 Multiple Wordlists from Directory**
```bash
# Load all wordlist files from a directory
stormfinder -d target.com -b --wordlist-dir /path/to/wordlists/

# Use the built-in wordlists directory
stormfinder -d target.com -b --wordlist-dir wordlists/
```

### **🌐 Download Wordlists from URLs**
```bash
# Download wordlists from URLs
stormfinder -d target.com -b --wordlist-urls "https://example.com/wordlist1.txt,https://example.com/wordlist2.txt"

# Popular wordlist URLs
stormfinder -d target.com -b --wordlist-urls "https://raw.githubusercontent.com/danielmiessler/SecLists/master/Discovery/DNS/subdomains-top1million-5000.txt"
```

### **🔄 Combine Multiple Sources**
```bash
# Use all wordlist sources together
stormfinder -d target.com -b \
  -w custom.txt \
  --wordlist-dir wordlists/ \
  --wordlist-urls "https://example.com/extra.txt" \
  -v
```

## 📚 **Built-in Wordlists**

Stormfinder includes several specialized wordlists:

### **🌐 `wordlists/common.txt`**
- **500+ common subdomains**
- General-purpose subdomain patterns
- Web services, infrastructure, and development terms

### **🔧 `wordlists/tech-stack.txt`**
- **Technology-specific subdomains**
- Web frameworks, databases, cloud services
- DevOps tools, monitoring systems, containers

### **🏢 `wordlists/industry.txt`**
- **Industry-specific subdomains**
- E-commerce, finance, healthcare, education
- Government, media, real estate, and more

## 🎯 **Wordlist File Formats**

### **Supported File Extensions**
- `.txt` - Plain text wordlists
- `.list` - List files
- `.wordlist` - Wordlist files

### **File Format Requirements**
```
# Comments start with #
subdomain1
subdomain2
subdomain3

# Empty lines are ignored

# More subdomains
api
admin
dev
```

### **Best Practices**
- **One subdomain per line**
- **No protocols** (http://, https://)
- **No domains** (just the subdomain part)
- **Comments with #** for organization
- **UTF-8 encoding** recommended

## 🚀 **Advanced Usage Examples**

### **Technology-Focused Scan**
```bash
# Target technology companies with tech-specific wordlists
stormfinder -d techcompany.com -b \
  -w wordlists/tech-stack.txt \
  --brute-threads 50 \
  -v
```

### **E-commerce Platform Scan**
```bash
# Target e-commerce sites with industry-specific terms
stormfinder -d shop.com -b \
  -w wordlists/industry.txt \
  -p \
  --recursive-enum \
  -v
```

### **Comprehensive Discovery**
```bash
# Maximum discovery with all wordlist sources
stormfinder -d target.com -b -p \
  --wordlist-dir wordlists/ \
  --wordlist-urls "https://raw.githubusercontent.com/danielmiessler/SecLists/master/Discovery/DNS/subdomains-top1million-20000.txt" \
  --brute-threads 100 \
  --cache \
  -v
```

### **Custom Wordlist Creation**
```bash
# Create a custom wordlist for specific target
echo -e "portal\ndashboard\nclient\npartner\nvendor" > custom-target.txt
stormfinder -d target.com -b -w custom-target.txt -v
```

## 📊 **Performance Optimization**

### **Wordlist Size vs Speed**
- **Small wordlists (< 1,000)**: Fast, targeted discovery
- **Medium wordlists (1,000-10,000)**: Balanced approach
- **Large wordlists (> 10,000)**: Comprehensive but slower

### **Threading Optimization**
```bash
# Optimize for speed with more threads
stormfinder -d target.com -b \
  --wordlist-dir wordlists/ \
  --brute-threads 100 \
  --optimize-speed

# Optimize for memory with fewer threads
stormfinder -d target.com -b \
  -w wordlists/common.txt \
  --brute-threads 25 \
  --optimize-memory
```

## 🔍 **Wordlist Statistics**

### **Verbose Output Shows**
- Number of words loaded from each source
- Total unique words after deduplication
- Loading progress for each wordlist
- Download progress for URL-based wordlists

### **Example Verbose Output**
```
⚡ Loading wordlist from file: custom.txt
✅ Loaded 500 words from custom.txt
⚡ Loading wordlists from directory: wordlists/
✅ Loaded 300 words from directory wordlists/
⚡ Downloading wordlists from 1 URLs
✅ Downloaded 1000 words from https://example.com/wordlist.txt
✅ Loaded 1650 words from custom wordlists
```

## 🌐 **Popular Wordlist URLs**

### **SecLists (Recommended)**
```bash
# Top 1 million subdomains (5,000 most common)
https://raw.githubusercontent.com/danielmiessler/SecLists/master/Discovery/DNS/subdomains-top1million-5000.txt

# Top 1 million subdomains (20,000 most common)
https://raw.githubusercontent.com/danielmiessler/SecLists/master/Discovery/DNS/subdomains-top1million-20000.txt

# Comprehensive subdomain list
https://raw.githubusercontent.com/danielmiessler/SecLists/master/Discovery/DNS/bitquark-subdomains-top100000.txt
```

### **Other Popular Sources**
```bash
# Assetnote wordlists
https://wordlists-cdn.assetnote.io/data/manual/best-dns-wordlist.txt

# Jhaddix All.txt
https://gist.githubusercontent.com/jhaddix/86a06c5dc309d08580a018c66354a056/raw/f58e82c9abfa46a932eb92edbe6b18214141439b/all.txt
```

## 🛠️ **Custom Wordlist Creation Tips**

### **Target-Specific Wordlists**
1. **Research the target organization**
2. **Identify their technology stack**
3. **Look for naming conventions**
4. **Include industry-specific terms**
5. **Add common variations**

### **Example Custom Wordlist**
```
# Company-specific
companyname-api
companyname-portal
companyname-admin

# Technology stack
django-app
react-frontend
nodejs-api

# Environment variations
dev-portal
staging-api
prod-dashboard

# Department-specific
hr-portal
finance-system
marketing-tools
```

## 🎯 **Best Practices**

### **Wordlist Selection**
- **Start with built-in wordlists** for general discovery
- **Add technology-specific lists** based on target analysis
- **Include industry-specific terms** for targeted organizations
- **Use large public wordlists** for comprehensive scans

### **Performance Tuning**
- **Use caching** for repeated scans with large wordlists
- **Adjust thread count** based on target responsiveness
- **Monitor memory usage** with very large wordlists
- **Use rate limiting** to avoid being blocked

### **Operational Security**
- **Respect rate limits** to avoid detection
- **Use proxy rotation** for large wordlist scans
- **Monitor for blocking** and adjust accordingly
- **Cache results** to avoid repeated requests

## 🚀 **Integration Examples**

### **With Other Tools**
```bash
# Generate custom wordlist from discovered subdomains
stormfinder -d target.com -silent | cut -d. -f1 | sort -u > discovered-patterns.txt

# Use discovered patterns for other targets
stormfinder -d newtarget.com -b -w discovered-patterns.txt
```

### **Automation Scripts**
```bash
#!/bin/bash
# Automated wordlist-based discovery

DOMAIN=$1
WORDLIST_DIR="wordlists"
OUTPUT_DIR="results"

# Create output directory
mkdir -p $OUTPUT_DIR

# Run with different wordlist combinations
stormfinder -d $DOMAIN -b -w $WORDLIST_DIR/common.txt -o $OUTPUT_DIR/common.txt
stormfinder -d $DOMAIN -b -w $WORDLIST_DIR/tech-stack.txt -o $OUTPUT_DIR/tech.txt
stormfinder -d $DOMAIN -b -w $WORDLIST_DIR/industry.txt -o $OUTPUT_DIR/industry.txt

# Combine results
cat $OUTPUT_DIR/*.txt | sort -u > $OUTPUT_DIR/all-subdomains.txt
```

Your wordlist management in Stormfinder is now incredibly powerful and flexible! 🌪️📝
