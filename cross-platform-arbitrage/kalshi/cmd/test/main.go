package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"kalshi"
)

func main() {
	// Create Kalshi client
	client := kalshi.NewClient()
	
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	fmt.Println("Testing Kalshi API connection...")
	fmt.Println("Base URL:", "https://api.elections.kalshi.com/trade-api/v2")
	fmt.Println()
	
	// Test 1: Get open markets
	fmt.Println("Test 1: Fetching open markets...")
	opts := &kalshi.GetMarketsOptions{
		Limit:  10,
		Status: "open",
	}
	
	marketsResp, err := client.GetMarkets(ctx, opts)
	if err != nil {
		log.Fatalf("Failed to get markets: %v", err)
	}
	
	fmt.Printf("Successfully retrieved %d markets\n", len(marketsResp.Markets))
	fmt.Println()
	
	// Display first few markets
	for i, market := range marketsResp.Markets {
		if i >= 3 {
			break
		}
		fmt.Printf("Market %d:\n", i+1)
		fmt.Println(market.Summary())
		fmt.Println()
		
		// Check for arbitrage opportunities
		arbPercent, exists := kalshi.CalculateArbitrageOpportunity(
			market.YesBid,
			market.YesAsk,
			market.NoBid,
			market.NoAsk,
		)
		
		if exists {
			fmt.Printf("⚠️  ARBITRAGE OPPORTUNITY: %.2f%%\n", arbPercent)
			fmt.Printf("   YES Bid: %d ($%s), NO Ask: %d ($%s)\n", 
				market.YesBid, market.YesBidDollars,
				market.NoAsk, market.NoAskDollars)
			fmt.Printf("   Combined: %d cents (should be 100)\n", market.YesBid+market.NoAsk)
		} else {
			fmt.Println("✓ No arbitrage opportunity")
		}
		fmt.Println("---")
	}
	
	// Test 2: Calculate implied prices
	fmt.Println("\nTest 2: Price calculations...")
	if len(marketsResp.Markets) > 0 {
		market := marketsResp.Markets[0]
		impliedNoFromYesBid := kalshi.CalculateImpliedNoPrice(market.YesBid)
		impliedNoFromYesAsk := kalshi.CalculateImpliedNoPrice(market.YesAsk)
		
		fmt.Printf("Market: %s\n", market.Ticker)
		fmt.Printf("YES Bid: %d → Implied NO Ask: %d (Actual NO Ask: %d)\n",
			market.YesBid, impliedNoFromYesBid, market.NoAsk)
		fmt.Printf("YES Ask: %d → Implied NO Bid: %d (Actual NO Bid: %d)\n",
			market.YesAsk, impliedNoFromYesAsk, market.NoBid)
		
		// Check for price discrepancies
		if impliedNoFromYesBid != market.NoAsk {
			fmt.Printf("⚠️  Price discrepancy: YES bid implies NO ask of %d, but actual is %d\n",
				impliedNoFromYesBid, market.NoAsk)
		}
		if impliedNoFromYesAsk != market.NoBid {
			fmt.Printf("⚠️  Price discrepancy: YES ask implies NO bid of %d, but actual is %d\n",
				impliedNoFromYesAsk, market.NoBid)
		}
	}
	
	// Test 3: Market statistics
	fmt.Println("\nTest 3: Market statistics...")
	if len(marketsResp.Markets) > 0 {
		var totalVolume int
		var openMarkets int
		var avgYesBid float64
		
		for _, market := range marketsResp.Markets {
			totalVolume += market.Volume
			if market.Status == "open" {
				openMarkets++
			}
			avgYesBid += float64(market.YesBid)
		}
		
		avgYesBid = avgYesBid / float64(len(marketsResp.Markets))
		
		fmt.Printf("Total markets retrieved: %d\n", len(marketsResp.Markets))
		fmt.Printf("Open markets: %d\n", openMarkets)
		fmt.Printf("Total volume across markets: %d\n", totalVolume)
		fmt.Printf("Average YES bid price: %.2f cents ($%.4f)\n", 
			avgYesBid, avgYesBid/100)
	}
	
	fmt.Println("\n✅ Kalshi API test completed successfully!")
}