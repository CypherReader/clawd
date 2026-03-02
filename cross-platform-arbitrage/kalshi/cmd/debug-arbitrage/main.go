package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"kalshi"
)

func main() {
	client := kalshi.NewClient()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	// Get markets to find the ones with arbitrage
	opts := &kalshi.GetMarketsOptions{
		Limit: 100,
	}
	
	marketsResp, err := client.GetMarkets(ctx, opts)
	if err != nil {
		log.Fatalf("Failed to get markets: %v", err)
	}
	
	fmt.Println("Debugging arbitrage calculations...")
	fmt.Printf("Total markets: %d\n\n", len(marketsResp.Markets))
	
	// Find markets that showed arbitrage
	arbitrageMarkets := []kalshi.Market{}
	for _, market := range marketsResp.Markets {
		arbPercent, exists := kalshi.CalculateArbitrageOpportunity(
			market.YesBid,
			market.YesAsk,
			market.NoBid,
			market.NoAsk,
		)
		
		if exists {
			arbitrageMarkets = append(arbitrageMarkets, market)
			fmt.Printf("Market with arbitrage: %s\n", market.Ticker)
			fmt.Printf("  YES Bid/Ask: %d/%d\n", market.YesBid, market.YesAsk)
			fmt.Printf("  NO Bid/Ask: %d/%d\n", market.NoBid, market.NoAsk)
			fmt.Printf("  Calculated arbitrage: %.2f%%\n", arbPercent)
			
			// Manual calculation check
			fmt.Printf("  Manual check:\n")
			fmt.Printf("    YES bid + NO ask = %d + %d = %d (should be >= 100)\n", 
				market.YesBid, market.NoAsk, market.YesBid+market.NoAsk)
			fmt.Printf("    NO bid + YES ask = %d + %d = %d (should be >= 100)\n",
				market.NoBid, market.YesAsk, market.NoBid+market.YesAsk)
			fmt.Println()
		}
	}
	
	fmt.Printf("\nFound %d markets with potential arbitrage\n", len(arbitrageMarkets))
	
	// Let's examine the arbitrage calculation logic
	fmt.Println("\n--- Arbitrage Calculation Logic ---")
	fmt.Println("For binary markets (YES/NO):")
	fmt.Println("1. YES price + NO price should always equal 100 cents ($1.00)")
	fmt.Println("2. Arbitrage exists if: YES bid + NO ask < 100")
	fmt.Println("   OR: NO bid + YES ask < 100")
	fmt.Println("3. This represents a risk-free profit opportunity")
	fmt.Println()
	
	// Test with specific values
	fmt.Println("Test cases:")
	testCases := []struct{
		name string
		yesBid, yesAsk, noBid, noAsk int
	}{
		{"Normal market", 45, 55, 45, 55},
		{"Wide spread", 40, 60, 40, 60},
		{"Zero prices", 0, 0, 100, 100},
		{"Edge case", 0, 1, 99, 100},
		{"Actual arbitrage", 48, 52, 47, 53}, // 48+53=101, 47+52=99 (arbitrage!)
	}
	
	for _, tc := range testCases {
		arbPercent, exists := kalshi.CalculateArbitrageOpportunity(
			tc.yesBid, tc.yesAsk, tc.noBid, tc.noAsk,
		)
		
		fmt.Printf("%s: YES %d/%d, NO %d/%d -> ", 
			tc.name, tc.yesBid, tc.yesAsk, tc.noBid, tc.noAsk)
		if exists {
			fmt.Printf("ARBITRAGE %.2f%%\n", arbPercent)
		} else {
			fmt.Printf("No arbitrage\n")
		}
	}
	
	// Check if we need to filter out markets with 0 volume
	fmt.Println("\n--- Volume Analysis ---")
	var zeroVolumeCount int
	var nonZeroVolumeCount int
	for _, market := range marketsResp.Markets {
		if market.Volume == 0 && market.Volume24h == 0 {
			zeroVolumeCount++
		} else {
			nonZeroVolumeCount++
		}
	}
	
	fmt.Printf("Markets with zero volume: %d\n", zeroVolumeCount)
	fmt.Printf("Markets with non-zero volume: %d\n", nonZeroVolumeCount)
	
	// The issue might be that markets with 0/0 and 100/100 prices
	// are being flagged as arbitrage when they're just inactive
	fmt.Println("\n--- Inactive Market Detection ---")
	fmt.Println("Markets with YES 0/0 and NO 100/100 are likely inactive")
	fmt.Println("These should be filtered out before arbitrage checking")
}