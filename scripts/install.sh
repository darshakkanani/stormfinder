#!/bin/bash

# 🚀 Stormfinder Installation Script
# One-command installation for Stormfinder

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
REPO_URL="https://github.com/darshakkanani/stormfinder"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="$HOME/.config/stormfinder"

echo -e "${BLUE}🌪️  Stormfinder Installation Script${NC}"
echo -e "${BLUE}====================================${NC}"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    arm64|aarch64)
        ARCH="arm64"
        ;;
    *)
        echo -e "${RED}❌ Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
esac

echo -e "${BLUE}🔍 Detected platform: ${OS}-${ARCH}${NC}"

# Check if Go is installed
if command -v go &> /dev/null; then
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    echo -e "${GREEN}✅ Go ${GO_VERSION} found${NC}"
    INSTALL_METHOD="source"
else
    echo -e "${YELLOW}⚠️  Go not found. Will install from binary.${NC}"
    INSTALL_METHOD="binary"
fi

# Function to install from source
install_from_source() {
    echo -e "${YELLOW}📦 Installing from source...${NC}"
    
    # Create temporary directory
    TEMP_DIR=$(mktemp -d)
    cd "$TEMP_DIR"
    
    # Clone repository
    echo -e "${BLUE}📥 Cloning repository...${NC}"
    git clone "$REPO_URL.git"
    cd stormfinder
    
    # Build
    echo -e "${BLUE}🔨 Building Stormfinder...${NC}"
    go build ./cmd/stormfinder
    
    # Install
    echo -e "${BLUE}📦 Installing to ${INSTALL_DIR}...${NC}"
    sudo mv stormfinder "$INSTALL_DIR/"
    sudo chmod +x "$INSTALL_DIR/stormfinder"
    
    # Cleanup
    cd /
    rm -rf "$TEMP_DIR"
}

# Function to install from binary
install_from_binary() {
    echo -e "${YELLOW}📦 Installing from binary...${NC}"
    
    # Determine download URL
    BINARY_NAME="stormfinder-${OS}-${ARCH}"
    if [ "$OS" = "windows" ]; then
        BINARY_NAME="${BINARY_NAME}.exe"
    fi
    
    DOWNLOAD_URL="${REPO_URL}/releases/latest/download/${BINARY_NAME}.tar.gz"
    
    # Create temporary directory
    TEMP_DIR=$(mktemp -d)
    cd "$TEMP_DIR"
    
    # Download
    echo -e "${BLUE}📥 Downloading ${BINARY_NAME}...${NC}"
    if command -v curl &> /dev/null; then
        curl -L -o stormfinder.tar.gz "$DOWNLOAD_URL"
    elif command -v wget &> /dev/null; then
        wget -O stormfinder.tar.gz "$DOWNLOAD_URL"
    else
        echo -e "${RED}❌ Neither curl nor wget found. Please install one of them.${NC}"
        exit 1
    fi
    
    # Extract
    echo -e "${BLUE}📦 Extracting...${NC}"
    tar -xzf stormfinder.tar.gz
    
    # Install
    echo -e "${BLUE}📦 Installing to ${INSTALL_DIR}...${NC}"
    sudo mv "$BINARY_NAME" "$INSTALL_DIR/stormfinder"
    sudo chmod +x "$INSTALL_DIR/stormfinder"
    
    # Cleanup
    cd /
    rm -rf "$TEMP_DIR"
}

# Create configuration directory
echo -e "${BLUE}⚙️  Setting up configuration...${NC}"
mkdir -p "$CONFIG_DIR"

# Download configuration templates
echo -e "${BLUE}📥 Downloading configuration templates...${NC}"
if command -v curl &> /dev/null; then
    curl -L -o "$CONFIG_DIR/config.yaml.example" "${REPO_URL}/raw/main/configs/config.yaml.example"
    curl -L -o "$CONFIG_DIR/providers.yaml.example" "${REPO_URL}/raw/main/configs/providers.yaml.example"
elif command -v wget &> /dev/null; then
    wget -O "$CONFIG_DIR/config.yaml.example" "${REPO_URL}/raw/main/configs/config.yaml.example"
    wget -O "$CONFIG_DIR/providers.yaml.example" "${REPO_URL}/raw/main/configs/providers.yaml.example"
fi

# Install based on method
if [ "$INSTALL_METHOD" = "source" ]; then
    install_from_source
else
    install_from_binary
fi

# Verify installation
echo -e "${BLUE}🔍 Verifying installation...${NC}"
if command -v stormfinder &> /dev/null; then
    VERSION=$(stormfinder -version 2>&1 | grep "Current Version" | awk '{print $3}')
    echo -e "${GREEN}✅ Stormfinder ${VERSION} installed successfully!${NC}"
else
    echo -e "${RED}❌ Installation failed. Stormfinder not found in PATH.${NC}"
    exit 1
fi

# Show next steps
echo -e "\n${GREEN}🎉 Installation completed!${NC}"
echo -e "${GREEN}========================${NC}"
echo -e "${BLUE}📖 Next steps:${NC}"
echo -e "${BLUE}  1. Test installation: ${YELLOW}stormfinder -h${NC}"
echo -e "${BLUE}  2. Basic scan: ${YELLOW}stormfinder -d example.com${NC}"
echo -e "${BLUE}  3. Configure API keys: ${YELLOW}cp $CONFIG_DIR/providers.yaml.example $CONFIG_DIR/provider-config.yaml${NC}"
echo -e "${BLUE}  4. Read documentation: ${YELLOW}${REPO_URL}${NC}"

echo -e "\n${GREEN}🌪️  Happy subdomain hunting!${NC}"
