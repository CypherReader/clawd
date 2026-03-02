#!/bin/bash
# Secure Git Setup for CypherReader/Oracle repository
# Uses file-based token storage with proper permissions

set -e

TOKEN_FILE=".git-token"
CREDENTIALS_FILE="$HOME/.git-credentials"

case "$1" in
    setup)
        echo "=== Secure Git Setup ==="
        
        if [ -z "$2" ]; then
            echo "Please provide your GitHub Personal Access Token:"
            echo "Usage: $0 setup <token>"
            echo "Example: $0 setup ghp_xxxxxxxxxxxxxxxxxxxx"
            exit 1
        fi
        
        TOKEN="$2"
        
        # Store token in secure file
        echo "$TOKEN" > "$TOKEN_FILE"
        chmod 600 "$TOKEN_FILE"  # Only owner can read/write
        echo "✓ Token stored in $TOKEN_FILE (permissions: 600)"
        
        # Set up Git credentials
        git config --global credential.helper 'store --file '"$CREDENTIALS_FILE"
        echo "https://${TOKEN}:x-oauth-basic@github.com" > "$CREDENTIALS_FILE"
        chmod 600 "$CREDENTIALS_FILE"
        echo "✓ Git credentials configured"
        
        # Test the token
        echo "Testing token with GitHub API..."
        RESPONSE=$(curl -s -H "Authorization: token $TOKEN" \
            -H "Accept: application/vnd.github.v3+json" \
            https://api.github.com/user 2>/dev/null | grep -o '"login":"[^"]*"' | head -1 || echo "")
        
        if [ -n "$RESPONSE" ]; then
            echo "✓ GitHub authentication successful: $RESPONSE"
        else
            echo "⚠️ Could not verify token via API (might still work for Git)"
        fi
        
        echo ""
        echo "Setup complete! You can now clone the repository."
        ;;
    
    clone)
        echo "=== Cloning Repository ==="
        
        if [ ! -f "$TOKEN_FILE" ]; then
            echo "❌ Token file not found. Run '$0 setup <token>' first."
            exit 1
        fi
        
        TOKEN=$(cat "$TOKEN_FILE")
        REPO="https://github.com/CypherReader/Oracle.git"
        CLONE_URL="https://${TOKEN}@github.com/CypherReader/Oracle.git"
        
        echo "Repository: $REPO"
        echo "Cloning..."
        
        git clone "$CLONE_URL" 2>&1 | tail -20
        
        if [ $? -eq 0 ]; then
            echo "✓ Repository cloned successfully"
            
            # Navigate to repository
            cd Oracle 2>/dev/null && echo "Changed to Oracle directory" || echo "Could not change to Oracle directory"
        else
            echo "❌ Clone failed"
            echo ""
            echo "Possible issues:"
            echo "1. Repository might not exist: $REPO"
            echo "2. Token might not have access"
            echo "3. Repository might be private"
            echo ""
            echo "Check: curl -I https://github.com/CypherReader/Oracle"
        fi
        ;;
    
    test)
        echo "=== Testing Setup ==="
        
        if [ ! -f "$TOKEN_FILE" ]; then
            echo "❌ Token file not found"
            exit 1
        fi
        
        TOKEN=$(cat "$TOKEN_FILE")
        
        echo "1. Testing token file permissions..."
        ls -la "$TOKEN_FILE"
        
        echo ""
        echo "2. Testing GitHub API access..."
        curl -s -H "Authorization: token $TOKEN" \
            -H "Accept: application/vnd.github.v3+json" \
            https://api.github.com/user | grep -o '"login":"[^"]*"\|"name":"[^"]*"' | head -5
        
        echo ""
        echo "3. Testing repository access..."
        curl -s -I "https://api.github.com/repos/CypherReader/Oracle" | head -5
        ;;
    
    clean)
        echo "=== Cleaning Up ==="
        
        read -p "Are you sure you want to remove token files? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            rm -f "$TOKEN_FILE" "$CREDENTIALS_FILE"
            echo "✓ Token files removed"
        else
            echo "Cleanup cancelled"
        fi
        ;;
    
    help|*)
        echo "Secure Git Setup for CypherReader/Oracle"
        echo "Usage: $0 <command>"
        echo ""
        echo "Commands:"
        echo "  setup <token>   Store GitHub token and configure Git"
        echo "  clone           Clone the Oracle repository"
        echo "  test            Test the current setup"
        echo "  clean           Remove token files (secure cleanup)"
        echo "  help            Show this help"
        echo ""
        echo "Example:"
        echo "  $0 setup ghp_xxxxxxxxxxxxxxxxxxxx"
        echo "  $0 clone"
        echo ""
        echo "Note: Token is stored in $TOKEN_FILE with 600 permissions"
        ;;
esac