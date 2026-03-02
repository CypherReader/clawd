#!/bin/bash
# Proton Pass CLI Setup Script for Server/Headless Environments
# This script sets up Proton Pass CLI with filesystem key storage
# Store this in your workspace for persistence between sessions

set -e

echo "=== Proton Pass CLI Server Setup ==="

# Create necessary directories
mkdir -p ~/.local/bin
mkdir -p ~/.config/proton-pass
mkdir -p ~/.local/share/proton-pass-cli

# Add to PATH if not already there
if [[ ":$PATH:" != *":$HOME/.local/bin:"* ]]; then
    echo "export PATH=\$PATH:\$HOME/.local/bin" >> ~/.bashrc
    export PATH=$PATH:$HOME/.local/bin
    echo "Added ~/.local/bin to PATH"
fi

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
    echo "✓ Proton Pass CLI installed: $(pass-cli --version)"
fi

# Set up environment for server/headless use
echo ""
echo "=== Server Environment Configuration ==="

# Create environment setup script
cat > ~/.protonpass-env << 'EOF'
# Proton Pass CLI Server Environment
# Source this file before using pass-cli in server environments

# Use filesystem key storage (required for headless servers)
export PROTON_PASS_KEY_PROVIDER=fs

# Optional: Set custom session directory
# export PROTON_PASS_SESSION_DIR="$HOME/.local/share/proton-pass-cli/.session"

# Optional: Enable debug logging if needed
# export PASS_LOG_LEVEL=debug

echo "Proton Pass CLI server environment loaded"
echo "Key provider: $PROTON_PASS_KEY_PROVIDER"
EOF

# Make it executable
chmod +x ~/.protonpass-env

echo "Created environment setup: ~/.protonpass-env"
echo ""
echo "=== Usage ==="
echo ""
echo "1. Load server environment:"
echo "   source ~/.protonpass-env"
echo ""
echo "2. Force logout if previously logged in with keyring:"
echo "   pass-cli logout --force"
echo ""
echo "3. Login (will use filesystem key storage):"
echo "   pass-cli login"
echo ""
echo "Note: In server environments, you may need to use"
echo "      device code or alternative login methods."