#!/usr/bin/env python3
import asyncio
from playwright.async_api import async_playwright
import sys

async def fetch_tweet(tweet_id):
    """Fetch tweet content using Playwright"""
    async with async_playwright() as p:
        # Launch browser
        browser = await p.chromium.launch(headless=True)
        context = await browser.new_context(
            user_agent='Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36'
        )
        page = await context.new_page()
        
        tweet_url = f"https://x.com/i/status/{tweet_id}"
        
        try:
            print(f"Fetching tweet: {tweet_url}")
            await page.goto(tweet_url, wait_until='networkidle', timeout=30000)
            
            # Wait for tweet content to load
            await page.wait_for_selector('article[data-testid="tweet"]', timeout=10000)
            
            # Extract tweet content
            tweet_data = await page.evaluate('''() => {
                const tweet = document.querySelector('article[data-testid="tweet"]');
                if (!tweet) return null;
                
                // Get tweet text
                const textElement = tweet.querySelector('div[data-testid="tweetText"]');
                const text = textElement ? textElement.innerText : '';
                
                // Get author
                const authorElement = tweet.querySelector('div[data-testid="User-Name"]');
                const author = authorElement ? authorElement.innerText.split('·')[0].trim() : '';
                
                // Get timestamp
                const timeElement = tweet.querySelector('time');
                const timestamp = timeElement ? timeElement.getAttribute('datetime') : '';
                
                // Get engagement stats
                const stats = {};
                const statElements = tweet.querySelectorAll('[data-testid="reply"] span, [data-testid="retweet"] span, [data-testid="like"] span, [data-testid="view"] span');
                statElements.forEach(el => {
                    const parent = el.closest('[data-testid]');
                    if (parent) {
                        const statType = parent.getAttribute('data-testid');
                        stats[statType] = el.innerText;
                    }
                });
                
                return {
                    text: text,
                    author: author,
                    timestamp: timestamp,
                    stats: stats,
                    url: window.location.href
                };
            }''')
            
            if tweet_data:
                print(f"\n=== Tweet {tweet_id} ===")
                print(f"Author: {tweet_data['author']}")
                print(f"Timestamp: {tweet_data['timestamp']}")
                print(f"Content: {tweet_data['text']}")
                print(f"Stats: {tweet_data['stats']}")
                print(f"URL: {tweet_data['url']}")
                return tweet_data
            else:
                print(f"Could not extract tweet content for {tweet_id}")
                return None
                
        except Exception as e:
            print(f"Error fetching tweet {tweet_id}: {e}")
            # Try to get page content for debugging
            try:
                content = await page.content()
                if "suspended" in content.lower() or "rate limit" in content.lower():
                    print("Twitter is rate limiting or blocking access")
                elif "login" in content.lower():
                    print("Twitter requires login to view this content")
            except:
                pass
            return None
        finally:
            await browser.close()

async def main():
    tweet_ids = [
        "2025178949417046130",
        "2026396408694386984"
    ]
    
    results = []
    for tweet_id in tweet_ids:
        result = await fetch_tweet(tweet_id)
        results.append(result)
        await asyncio.sleep(2)  # Be polite with delays
    
    return results

if __name__ == "__main__":
    results = asyncio.run(main())
    print(f"\nFetched {len([r for r in results if r])} out of {len(results)} tweets")