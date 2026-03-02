#!/usr/bin/env python3
import requests
import json
import sys

def fetch_tweet_via_api(tweet_id):
    """Try to fetch tweet using Twitter API v2"""
    # Note: This requires a Twitter API bearer token
    # For now, we'll try public endpoints
    headers = {
        'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36',
        'Accept': 'application/json',
        'Accept-Language': 'en-US,en;q=0.9',
    }
    
    # Try public API endpoint (may not work without authentication)
    url = f"https://api.twitter.com/2/tweets/{tweet_id}?tweet.fields=created_at,author_id,text,public_metrics"
    
    print(f"Attempting to fetch tweet {tweet_id}...")
    print("Note: This may fail without proper Twitter API authentication")
    
    try:
        response = requests.get(url, headers=headers, timeout=10)
        if response.status_code == 200:
            data = response.json()
            print(f"Success via API: {json.dumps(data, indent=2)}")
            return data
        else:
            print(f"API returned status {response.status_code}")
            print(f"Response: {response.text[:200]}")
    except Exception as e:
        print(f"API request failed: {e}")
    
    return None

def fetch_tweet_via_embed(tweet_id):
    """Try to fetch tweet via oEmbed endpoint"""
    url = f"https://publish.twitter.com/oembed?url=https://twitter.com/i/status/{tweet_id}"
    
    try:
        response = requests.get(url, timeout=10)
        if response.status_code == 200:
            data = response.json()
            print(f"\n=== Tweet {tweet_id} (via oEmbed) ===")
            print(f"Author: {data.get('author_name', 'Unknown')}")
            print(f"HTML content (first 500 chars): {data.get('html', '')[:500]}...")
            
            # Extract text from HTML
            html = data.get('html', '')
            # Simple extraction - look for tweet text in HTML
            import re
            text_match = re.search(r'<p[^>]*>(.*?)</p>', html, re.DOTALL)
            if text_match:
                text = text_match.group(1)
                # Clean HTML tags
                text = re.sub(r'<[^>]+>', '', text)
                print(f"Extracted text: {text}")
            
            return data
        else:
            print(f"oEmbed returned status {response.status_code}")
    except Exception as e:
        print(f"oEmbed request failed: {e}")
    
    return None

def main():
    tweet_ids = [
        "2025178949417046130",
        "2026396408694386984"
    ]
    
    print("Attempting to fetch tweets...")
    print("=" * 50)
    
    for tweet_id in tweet_ids:
        print(f"\n📱 Processing tweet ID: {tweet_id}")
        print("-" * 30)
        
        # Try oEmbed first (more likely to work)
        result = fetch_tweet_via_embed(tweet_id)
        
        if not result:
            print(f"\nCould not fetch tweet {tweet_id} via available methods.")
            print(f"Direct URL: https://x.com/i/status/{tweet_id}")
        
        print("\n")

if __name__ == "__main__":
    main()