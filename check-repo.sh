#!/bin/bash
echo "Checking GitHub repository: https://github.com/CypherReader/Oracle"
echo ""

# Check if repository exists (public)
echo "1. Checking public access..."
STATUS=$(curl -s -o /dev/null -w "%{http_code}" "https://github.com/CypherReader/Oracle")
echo "   HTTP Status: $STATUS"

if [ "$STATUS" = "200" ]; then
    echo "   ✓ Repository exists and is publicly accessible"
elif [ "$STATUS" = "404" ]; then
    echo "   ✗ Repository not found (404)"
    echo "   Possible reasons:"
    echo "   - Repository is private"
    echo "   - Repository doesn't exist"
    echo "   - Repository name is different (case-sensitive)"
    echo "   - URL has a typo"
else
    echo "   ? Unexpected status: $STATUS"
fi

echo ""
echo "2. Testing with different cases..."
for repo in Oracle oracle ORACLE; do
    STATUS=$(curl -s -o /dev/null -w "%{http_code}" "https://github.com/CypherReader/$repo")
    echo "   CypherReader/$repo: $STATUS"
done

echo ""
echo "3. Suggested next steps:"
echo "   - Double-check the exact repository name"
echo "   - Ensure you have a GitHub Personal Access Token"
echo "   - Use: ./secure-git-setup.sh setup <token>"
echo "   - Then try: ./secure-git-setup.sh clone"