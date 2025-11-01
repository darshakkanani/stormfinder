#!/bin/bash

# 🔨 Stormfinder Build Script
# Automated build process for all platforms

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Build information
VERSION=$(grep 'const version' pkg/runner/banners.go | cut -d'"' -f2)
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

echo -e "${BLUE}🌪️  Building Stormfinder v${VERSION}${NC}"
echo -e "${BLUE}========================================${NC}"

# Build flags
LDFLAGS="-s -w -X main.version=${VERSION} -X main.buildTime=${BUILD_TIME} -X main.gitCommit=${GIT_COMMIT}"

# Clean previous builds
echo -e "${YELLOW}🧹 Cleaning previous builds...${NC}"
rm -rf dist/
mkdir -p dist/

# Build for current platform
echo -e "${YELLOW}🔨 Building for current platform...${NC}"
go build -ldflags="${LDFLAGS}" -o dist/stormfinder ./cmd/stormfinder

# Cross-compilation for major platforms
echo -e "${YELLOW}🌍 Cross-compiling for multiple platforms...${NC}"

# Linux
echo -e "${BLUE}  📦 Building for Linux...${NC}"
GOOS=linux GOARCH=amd64 go build -ldflags="${LDFLAGS}" -o dist/stormfinder-linux-amd64 ./cmd/stormfinder
GOOS=linux GOARCH=arm64 go build -ldflags="${LDFLAGS}" -o dist/stormfinder-linux-arm64 ./cmd/stormfinder

# macOS
echo -e "${BLUE}  🍎 Building for macOS...${NC}"
GOOS=darwin GOARCH=amd64 go build -ldflags="${LDFLAGS}" -o dist/stormfinder-darwin-amd64 ./cmd/stormfinder
GOOS=darwin GOARCH=arm64 go build -ldflags="${LDFLAGS}" -o dist/stormfinder-darwin-arm64 ./cmd/stormfinder

# Windows
echo -e "${BLUE}  🪟 Building for Windows...${NC}"
GOOS=windows GOARCH=amd64 go build -ldflags="${LDFLAGS}" -o dist/stormfinder-windows-amd64.exe ./cmd/stormfinder
GOOS=windows GOARCH=arm64 go build -ldflags="${LDFLAGS}" -o dist/stormfinder-windows-arm64.exe ./cmd/stormfinder

# Create archives
echo -e "${YELLOW}📦 Creating release archives...${NC}"
cd dist/

# Linux archives
tar -czf stormfinder-linux-amd64.tar.gz stormfinder-linux-amd64
tar -czf stormfinder-linux-arm64.tar.gz stormfinder-linux-arm64

# macOS archives
tar -czf stormfinder-darwin-amd64.tar.gz stormfinder-darwin-amd64
tar -czf stormfinder-darwin-arm64.tar.gz stormfinder-darwin-arm64

# Windows archives
zip stormfinder-windows-amd64.zip stormfinder-windows-amd64.exe
zip stormfinder-windows-arm64.zip stormfinder-windows-arm64.exe

cd ..

# Generate checksums
echo -e "${YELLOW}🔐 Generating checksums...${NC}"
cd dist/
sha256sum * > checksums.txt
cd ..

# Build summary
echo -e "\n${GREEN}✅ Build completed successfully!${NC}"
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}📊 Build Summary:${NC}"
echo -e "${GREEN}  Version: ${VERSION}${NC}"
echo -e "${GREEN}  Build Time: ${BUILD_TIME}${NC}"
echo -e "${GREEN}  Git Commit: ${GIT_COMMIT}${NC}"
echo -e "${GREEN}  Artifacts: $(ls dist/ | wc -l) files${NC}"

echo -e "\n${BLUE}📁 Generated files:${NC}"
ls -la dist/

echo -e "\n${GREEN}🚀 Ready for release!${NC}"
