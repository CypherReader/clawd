#!/bin/bash
# Proton Pass CLI Setup Script
# This script sets up Proton Pass CLI
# Store this in your workspace for persistence between sessions

set -e

echo "=== Proton Pass CLI Setup ==="

# Create necessary directories
mkdir -p ~/.local/bin

# Check if jq is installed (required for installation)
if ! command -v jq &> /dev/null; then
    echo "Installing jq (required dependency)..."
    apt-get update && apt-get install -y jq
fi

# Check if pass-cli is already installed
if command -v pass-cli &> /dev/null; then
    echo "✓ Proton Pass CLI already installed: $(pass-cli --version)"
else
    echo "Installing Proton Pass CLI..."
    curl -fsSL https://proton.me/download/pass-cli/install.sh | bash
    
    # Add to PATH
    if [[ ":$PATH:" != *":$HOME/.local/bin:"* ]]; then
        echo "export PATH=\$PATH:\$HOME/.local/bin" >> ~/.bashrc
        export PATH=$PATH:$HOME/.local/bin
        echo "Added ~/.local/bin to PATH"
    fi
    
    echo "✓ Proton Pass CLI installed: $(pass-cli --version)"
fi

echo ""
echo "=== Next Steps ==="
echo "1. Log in to Proton Pass:"
echo "   pass-cli login"
echo ""
echo "2. List your vaults:"
echo "   pass-cli vault list"
echo ""
echo "3. Store a secret (e.g., Git token):"
echo "   pass-cli item create --vault 'Personal' --title 'GitHub Token' \\"
echo "     --field 'token=YOUR_TOKEN' --field 'url=https://github.com'"
echo ""
echo "4. Retrieve a secret:"
echo "   pass-cli item get 'GitHub Token' --field token"