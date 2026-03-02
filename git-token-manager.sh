#!/bin/bash
# Git Token Manager for Proton Pass
# Use this to store and retrieve Git tokens securely

set -e

# Load Proton Pass environment
export PATH="/root/.local/bin:$PATH"
export PROTON_PASS_KEY_PROVIDER=fs

VAULT_NAME="Personal"
ITEM_NAME="GitHub Personal Access Token"

case "$1" in
    store)
        if [ -z "$2" ]; then
            echo "Usage: $0 store <token>"
            echo "Example: $0 store ghp_xxxxxxxxxxxxxxxxxxxx"
            exit 1
        fi
        
        TOKEN="$2"
        
        echo "Storing Git token in Proton Pass..."
        
        # Check if item already exists
        if pass-cli item get "$ITEM_NAME" --field token 2>/dev/null; then
            echo "Token already exists. Updating..."
            # Note: Update command might be different - check documentation
            echo "For now, please delete the existing item first if needed"
        else
            # Create new item
            pass-cli item create \
                --vault "$VAULT_NAME" \
                --title "$ITEM_NAME" \
                --field "token=$TOKEN" \
                --field "purpose=GitHub repository access" \
                --field "repository=https://github.com/CypherReader/Oracle" \
                --field "type=Personal Access Token" \
                --note "Git token for CypherReader/Oracle repository"
            
            echo "✓ Token stored successfully in Proton Pass"
        fi
        ;;
    
    get)
        echo "Retrieving Git token from Proton Pass..."
        TOKEN=$(pass-cli item get "$ITEM_NAME" --field token 2>/dev/null || echo "")
        
        if [ -n "$TOKEN" ]; then
            echo "✓ Token retrieved"
            echo "$TOKEN"
        else
            echo "❌ Token not found or error retrieving"
            exit 1
        fi
        ;;
    
    test)
        echo "Testing Git token..."
        TOKEN=$(pass-cli item get "$ITEM_NAME" --field token 2>/dev/null || echo "")
        
        if [ -z "$TOKEN" ]; then
            echo "❌ No token found"
            exit 1
        fi
        
        echo "Testing access to GitHub..."
        RESPONSE=$(curl -s -H "Authorization: token $TOKEN" \
            -H "Accept: application/vnd.github.v3+json" \
            https://api.github.com/user 2>/dev/null | grep -o '"login":"[^"]*"' | head -1 || echo "")
        
        if [ -n "$RESPONSE" ]; then
            echo "✓ GitHub authentication successful: $RESPONSE"
        else
            echo "⚠️ Could not verify token (might still work for Git)"
        fi
        ;;
    
    setup-git)
        echo "Setting up Git with token..."
        TOKEN=$(pass-cli item get "$ITEM_NAME" --field token 2>/dev/null || echo "")
        
        if [ -z "$TOKEN" ]; then
            echo "❌ No token found in Proton Pass"
            exit 1
        fi
        
        # Configure Git to use token
        git config --global credential.helper 'store --file ~/.git-credentials'
        echo "https://${TOKEN}:x-oauth-basic@github.com" > ~/.git-credentials
        chmod 600 ~/.git-credentials
        
        echo "✓ Git credential helper configured"
        echo "You can now clone: git clone https://github.com/CypherReader/Oracle.git"
        ;;
    
    clone)
        echo "Cloning repository..."
        TOKEN=$(pass-cli item get "$ITEM_NAME" --field token 2>/dev/null || echo "")
        
        if [ -z "$TOKEN" ]; then
            echo "❌ No token found in Proton Pass"
            exit 1
        fi
        
        REPO="https://github.com/CypherReader/Oracle.git"
        CLONE_URL="https://${TOKEN}@github.com/CypherReader/Oracle.git"
        
        echo "Cloning: $REPO"
        git clone "$CLONE_URL" 2>&1 | tail -20
        
        if [ $? -eq 0 ]; then
            echo "✓ Repository cloned successfully"
        else
            echo "❌ Clone failed"
        fi
        ;;
    
    help|*)
        echo "Git Token Manager for Proton Pass"
        echo "Usage: $0 <command>"
        echo ""
        echo "Commands:"
        echo "  store <token>      Store a Git token in Proton Pass"
        echo "  get               Retrieve the Git token"
        echo "  test              Test the Git token with GitHub API"
        echo "  setup-git         Configure Git to use the token"
        echo "  clone             Clone the Oracle repository"
        echo "  help              Show this help"
        echo ""
        echo "Example workflow:"
        echo "  1. $0 store ghp_xxxxxxxxxxxxxxxxxxxx"
        echo "  2. $0 test"
        echo "  3. $0 clone"
        ;;
esac