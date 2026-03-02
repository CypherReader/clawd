package main

import (
	"context"
	"fmt"
	"log"
	"sort"
	"time"

	"kalshi"
)

func main() {
	client := kalshi.NewClient()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	
	fmt.Println("Fetching active Kalshi markets with trading volume...")
	fmt.Println()
	
	// Get more markets to find active ones
	opts := &kalshi.GetMarketsOptions{
		Limit:  50,
		Status: "open", // Try "open" first
	}
	
	marketsResp, err := client.GetMarkets(ctx, opts)
	if err != nil {
		// Try without status filter
		opts.Status = ""
		marketsResp, err = client.GetMarkets(ctx, opts)
		if err != nil {
			log.Fatalf("Failed to get markets: %v", err)
		}
	}
	
	// Filter for markets with volume
	var activeMarkets []kalshi.Market
	for _, market := range marketsResp.Markets {
		if market.Volume > 0 || market.Volume24h > 0 {
			activeMarkets = append(activeMarkets, market)
		}
	}
	
	if len(activeMarkets) == 0 {
		fmt.Println("No markets with trading volume found. Showing all markets:")
		activeMarkets = marketsResp.Markets
	}
	
	fmt.Printf("Found %d total markets, %d with volume\n", 
		len(marketsResp.Markets), len(activeMarkets))
	fmt.Println()
	
	// Sort by volume (descending)
	sort.Slice(activeMarkets, func(i, j int) bool {
		return activeMarkets[i].Volume24h > activeMarkets[j].Volume24h
	})
	
	// Display top active markets
	fmt.Println("Top 5 markets by 24h volume:")
	for i := 0; i < len(activeMarkets) && i < 5; i++ {
		market := activeMarkets[i]
		fmt.Printf("%d. %s\n", i+1, market.Ticker)
		fmt.Printf("   Title: %s\n", market.Title)
		fmt.Printf("   Status: %s\n", market.Status)
		fmt.Printf("   YES Bid/Ask: %d/%d ($%s/$%s)\n",
			market.YesBid, market.YesAsk, market.YesBidDollars, market.YesAskDollars)
		fmt.Printf("   NO Bid/Ask: %d/%d ($%s/$%s)\n",
			market.NoBid, market.NoAsk, market.NoBidDollars, market.NoAskDollars)
		fmt.Printf("   Volume: %d (24h: %d)\n", market.Volume, market.Volume24h)
		fmt.Printf("   Expires: %s\n", market.ExpirationTime.Format("2006-01-02 15:04"))
		
		// Check arbitrage
		arbPercent, exists := kalshi.CalculateArbitrageOpportunity(
			market.YesBid,
			market.YesAsk,
			market.NoBid,
			market.NoAsk,
		)
		
		if exists {
			fmt.Printf("   ⚠️  ARBITRAGE: %.2f%%\n", arbPercent)
		} else {
			fmt.Printf("   ✓ No arbitrage\n")
		}
		fmt.Println()
	}
	
	// Analyze market statistics
	fmt.Println("Market Analysis:")
	fmt.Printf("Total markets analyzed: %d\n", len(activeMarkets))
	
	var totalVolume24h int
	var marketsWithBidAsk int
	var totalBidAskSpread float64
	
	for _, market := range activeMarkets {
		totalVolume24h += market.Volume24h
		
		// Check if market has bid/ask prices
		if market.YesBid > 0 && market.YesAsk > 0 {
			marketsWithBidAsk++
			spread := float64(market.YesAsk-market.YesBid) / 100.0
			totalBidAskSpread += spread
		}
	}
	
	fmt.Printf("Total 24h volume: %d\n", totalVolume24h)
	fmt.Printf("Markets with bid/ask prices: %d\n", marketsWithBidAsk)
	
	if marketsWithBidAsk > 0 {
		avgSpread := totalBidAskSpread / float64(marketsWithBidAsk)
		fmt.Printf("Average bid-ask spread: $%.4f (%.2f%%)\n", 
			avgSpread, avgSpread*100)
	}
	
	// Check for any arbitrage opportunities across all markets
	fmt.Println("\nScanning for arbitrage opportunities...")
	var arbitrageCount int
	for _, market := range activeMarkets {
		arbPercent, exists := kalshi.CalculateArbitrageOpportunity(
			market.YesBid,
			market.YesAsk,
			market.NoBid,
			market.NoAsk,
		)
		
		if exists {
			arbitrageCount++
			fmt.Printf("Found arbitrage in %s: %.2f%%\n", market.Ticker, arbPercent)
		}
	}
	
	if arbitrageCount == 0 {
		fmt.Println("No arbitrage opportunities found in active markets.")
	} else {
		fmt.Printf("Found %d arbitrage opportunity(ies)\n", arbitrageCount)
	}
}