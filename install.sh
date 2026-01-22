#!/bin/bash
# BakaSub Installation Script
# Automatically detects OS/architecture and installs the latest release

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# ASCII Art
print_banner() {
    echo -e "${CYAN}"
    cat << "EOF"
    ____        _          ____        _     
   | __ )  __ _| | ____ _ / ___| _   _| |__  
   |  _ \ / _` | |/ / _` |\___ \| | | | '_ \ 
   | |_) | (_| |   < (_| | ___) | |_| | |_) |
   |____/ \__,_|_|\_\__,_||____/ \__,_|_.__/ 
                                             
EOF
    echo -e "${NC}"
    echo -e "${BLUE}BakaSub Installer${NC}"
    echo -e "${YELLOW}It's not like we want you to use this tool or anything... B-Baka!${NC}"
    echo ""
}

# Detect OS and Architecture
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    
    case "$OS" in
        linux*)
            OS="linux"
            ;;
        darwin*)
            OS="darwin"
            ;;
        msys*|mingw*|cygwin*)
            OS="windows"
            ;;
        *)
            echo -e "${RED}❌ Unsupported OS: $OS${NC}"
            exit 1
            ;;
    esac
    
    case "$ARCH" in
        x86_64|amd64)
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
    
    echo -e "${GREEN}✓ Detected platform: ${OS}-${ARCH}${NC}"
}

# Get latest release version
get_latest_version() {
    echo -e "${BLUE}🔍 Fetching latest release...${NC}"
    
    LATEST_VERSION=$(curl -s https://api.github.com/repos/lsilvatti/bakasub/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    
    if [ -z "$LATEST_VERSION" ]; then
        echo -e "${RED}❌ Failed to fetch latest version${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}✓ Latest version: ${LATEST_VERSION}${NC}"
}

# Download binary
download_binary() {
    BINARY_NAME="bakasub-${OS}-${ARCH}"
    if [ "$OS" = "windows" ]; then
        BINARY_NAME="${BINARY_NAME}.exe"
    fi
    
    DOWNLOAD_URL="https://github.com/lsilvatti/bakasub/releases/download/${LATEST_VERSION}/${BINARY_NAME}"
    
    echo -e "${BLUE}📦 Downloading BakaSub...${NC}"
    echo -e "${CYAN}   URL: ${DOWNLOAD_URL}${NC}"
    
    TMP_DIR=$(mktemp -d)
    TMP_FILE="${TMP_DIR}/${BINARY_NAME}"
    
    if ! curl -L -o "${TMP_FILE}" "${DOWNLOAD_URL}"; then
        echo -e "${RED}❌ Download failed${NC}"
        rm -rf "${TMP_DIR}"
        exit 1
    fi
    
    echo -e "${GREEN}✓ Download complete${NC}"
}

# Verify checksum (if available)
verify_checksum() {
    CHECKSUM_URL="${DOWNLOAD_URL}.sha256"
    CHECKSUM_FILE="${TMP_DIR}/${BINARY_NAME}.sha256"
    
    echo -e "${BLUE}🔐 Verifying checksum...${NC}"
    
    if curl -s -L -o "${CHECKSUM_FILE}" "${CHECKSUM_URL}" 2>/dev/null; then
        cd "${TMP_DIR}"
        if sha256sum -c "${CHECKSUM_FILE}" --quiet 2>/dev/null; then
            echo -e "${GREEN}✓ Checksum verified${NC}"
        else
            echo -e "${YELLOW}⚠ Checksum verification failed (continuing anyway)${NC}"
        fi
        cd - > /dev/null
    else
        echo -e "${YELLOW}⚠ Checksum file not available (skipping verification)${NC}"
    fi
}

# Install binary
install_binary() {
    echo -e "${BLUE}📥 Installing BakaSub...${NC}"
    
    # Determine installation directory
    if [ -w "/usr/local/bin" ]; then
        INSTALL_DIR="/usr/local/bin"
        NEEDS_SUDO=false
    elif [ -w "$HOME/.local/bin" ]; then
        INSTALL_DIR="$HOME/.local/bin"
        NEEDS_SUDO=false
    else
        INSTALL_DIR="/usr/local/bin"
        NEEDS_SUDO=true
    fi
    
    # Create directory if needed
    if [ ! -d "$INSTALL_DIR" ]; then
        mkdir -p "$INSTALL_DIR" 2>/dev/null || {
            echo -e "${YELLOW}⚠ Cannot create $INSTALL_DIR, trying with sudo...${NC}"
            sudo mkdir -p "$INSTALL_DIR"
            NEEDS_SUDO=true
        }
    fi
    
    # Install
    TARGET="${INSTALL_DIR}/bakasub"
    if [ "$OS" = "windows" ]; then
        TARGET="${TARGET}.exe"
    fi
    
    if [ "$NEEDS_SUDO" = true ]; then
        echo -e "${YELLOW}⚠ Requesting sudo permissions for installation...${NC}"
        sudo cp "${TMP_FILE}" "${TARGET}"
        sudo chmod +x "${TARGET}"
    else
        cp "${TMP_FILE}" "${TARGET}"
        chmod +x "${TARGET}"
    fi
    
    # Cleanup
    rm -rf "${TMP_DIR}"
    
    echo -e "${GREEN}✓ Installed to: ${TARGET}${NC}"
}

# Check if binary is in PATH
check_path() {
    if ! command -v bakasub &> /dev/null; then
        echo -e "${YELLOW}⚠ Warning: ${INSTALL_DIR} is not in your PATH${NC}"
        echo -e "${YELLOW}  Add this line to your ~/.bashrc or ~/.zshrc:${NC}"
        echo -e "${CYAN}  export PATH=\"${INSTALL_DIR}:\$PATH\"${NC}"
        return 1
    fi
    return 0
}

# Test installation
test_installation() {
    echo -e "${BLUE}🧪 Testing installation...${NC}"
    
    if bakasub --version &> /dev/null; then
        VERSION_OUTPUT=$(bakasub --version)
        echo -e "${GREEN}✓ BakaSub is working!${NC}"
        echo -e "${CYAN}  ${VERSION_OUTPUT}${NC}"
    else
        echo -e "${YELLOW}⚠ Installation complete but cannot run bakasub${NC}"
        return 1
    fi
}

# Print success message
print_success() {
    echo ""
    echo -e "${GREEN}╔════════════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║  🎉 BakaSub installed successfully!           ║${NC}"
    echo -e "${GREEN}╚════════════════════════════════════════════════╝${NC}"
    echo ""
    echo -e "${CYAN}📚 Getting Started:${NC}"
    echo -e "   1. Run: ${YELLOW}bakasub${NC}"
    echo -e "   2. Follow the setup wizard"
    echo -e "   3. Start translating!"
    echo ""
    echo -e "${CYAN}📖 Documentation:${NC}"
    echo -e "   https://github.com/lsilvatti/bakasub#readme"
    echo ""
    echo -e "${CYAN}💖 Like BakaSub?${NC}"
    echo -e "   ⭐ Star us: https://github.com/lsilvatti/bakasub"
    echo -e "   ☕ Support: https://ko-fi.com/bakasub"
    echo ""
}

# Main installation flow
main() {
    print_banner
    detect_platform
    get_latest_version
    download_binary
    verify_checksum
    install_binary
    
    if check_path; then
        test_installation
    fi
    
    print_success
}

main
