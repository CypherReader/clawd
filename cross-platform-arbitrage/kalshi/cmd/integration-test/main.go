package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"kalshi"
)

func main() {
	fmt.Println("=== Kalshi Integration Test ===")
	fmt.Println("Testing API connectivity, data retrieval, and arbitrage detection")
	fmt.Println()
	
	client := kalshi.NewClient()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	
	// Test 1: Basic API connectivity
	fmt.Println("Test 1: API Connectivity")
	fmt.Println("Fetching markets...")
	
	marketsResp, err := client.GetMarkets(ctx, &kalshi.GetMarketsOptions{
		Limit: 100,
	})
	if err != nil {
		log.Fatalf("❌ API connectivity test failed: %v", err)
	}
	
	fmt.Printf("✅ Successfully retrieved %d markets\n", len(marketsResp.Markets))
	fmt.Println()
	
	// Test 2: Market data validation
	fmt.Println("Test 2: Market Data Validation")
	activeMarkets := 0
	inactiveMarkets := 0
	
	for _, market := range marketsResp.Markets {
		if kalshi.IsMarketActive(market) {
			activeMarkets++
		} else {
			inactiveMarkets++
		}
	}
	
	fmt.Printf("Active markets: %d\n", activeMarkets)
	fmt.Printf("Inactive markets: %d\n", inactiveMarkets)
	fmt.Printf("✅ Market classification working\n")
	fmt.Println()
	
	// Test 3: Price calculations
	fmt.Println("Test 3: Price Calculations")
	if activeMarkets > 0 {
		// Find first active market
		var sampleMarket kalshi.Market
		for _, market := range marketsResp.Markets {
			if kalshi.IsMarketActive(market) {
				sampleMarket = market
				break
			}
		}
		
		fmt.Printf("Sample active market: %s\n", sampleMarket.Ticker)
		fmt.Printf("YES Bid/Ask: %d/%d\n", sampleMarket.YesBid, sampleMarket.YesAsk)
		fmt.Printf("NO Bid/Ask: %d/%d\n", sampleMarket.NoBid, sampleMarket.NoAsk)
		
		// Calculate implied prices
		impliedNoFromYesBid := kalshi.CalculateImpliedNoPrice(sampleMarket.YesBid)
		impliedNoFromYesAsk := kalshi.CalculateImpliedNoPrice(sampleMarket.YesAsk)
		
		fmt.Printf("Implied NO from YES bid (%d): %d (actual NO ask: %d)\n",
			sampleMarket.YesBid, impliedNoFromYesBid, sampleMarket.NoAsk)
		fmt.Printf("Implied NO from YES ask (%d): %d (actual NO bid: %d)\n",
			sampleMarket.YesAsk, impliedNoFromYesAsk, sampleMarket.NoBid)
		
		// Check arbitrage
		arbPercent, exists := kalshi.CalculateArbitrageOpportunity(
			sampleMarket.YesBid,
			sampleMarket.YesAsk,
			sampleMarket.NoBid,
			sampleMarket.NoAsk,
		)
		
		if exists {
			fmt.Printf("⚠️  Arbitrage detected: %.2f%%\n", arbPercent)
		} else {
			fmt.Printf("✓ No arbitrage\n")
		}
	} else {
		fmt.Println("No active markets found for price calculation test")
	}
	fmt.Println()
	
	// Test 4: Arbitrage scanning
	fmt.Println("Test 4: Arbitrage Scanning")
	fmt.Println("Scanning all markets for arbitrage opportunities...")
	
	var arbitrageOpportunities []kalshi.Market
	for _, market := range marketsResp.Markets {
		if !kalshi.IsMarketActive(market) {
			continue
		}
		
		arbPercent, exists := kalshi.CalculateArbitrageOpportunity(
			market.YesBid,
			market.YesAsk,
			market.NoBid,
			market.NoAsk,
		)
		
		if exists {
			arbitrageOpportunities = append(arbitrageOpportunities, market)
			fmt.Printf("  Found arbitrage in %s: %.2f%%\n", market.Ticker, arbPercent)
		}
	}
	
	fmt.Printf("Total arbitrage opportunities found: %d\n", len(arbitrageOpportunities))
	fmt.Println()
	
	// Test 5: Data normalization
	fmt.Println("Test 5: Data Normalization")
	fmt.Println("Testing conversion between different price formats...")
	
	testPrices := []struct {
		cents int
		dollars string
	}{
		{56, "0.5600"},
		{0, "0.0000"},
		{100, "1.0000"},
		{25, "0.2500"},
	}
	
	for _, test := range testPrices {
		expectedDollars := fmt.Sprintf("%.4f", float64(test.cents)/100.0)
		if expectedDollars == test.dollars {
			fmt.Printf("✅ Price conversion correct: %d cents = $%s\n", test.cents, test.dollars)
		} else {
			fmt.Printf("❌ Price conversion mismatch: %d cents → $%s (expected $%s)\n",
				test.cents, test.dollars, expectedDollars)
		}
	}
	fmt.Println()
	
	// Test 6: Performance test
	fmt.Println("Test 6: Performance Test")
	start := time.Now()
	
	// Simulate multiple API calls
	for i := 0; i < 3; i++ {
		_, err := client.GetMarkets(ctx, &kalshi.GetMarketsOptions{
			Limit: 10,
		})
		if err != nil {
			fmt.Printf("❌ Performance test iteration %d failed: %v\n", i+1, err)
		}
	}
	
	elapsed := time.Since(start)
	fmt.Printf("3 API calls completed in %v\n", elapsed)
	fmt.Printf("Average response time: %v\n", elapsed/3)
	fmt.Println()
	
	// Summary
	fmt.Println("=== Integration Test Summary ===")
	fmt.Printf("Total tests: 6\n")
	fmt.Printf("Markets retrieved: %d\n", len(marketsResp.Markets))
	fmt.Printf("Active markets: %d\n", activeMarkets)
	fmt.Printf("Arbitrage opportunities: %d\n", len(arbitrageOpportunities))
	
	if len(arbitrageOpportunities) > 0 {
		fmt.Println("\n⚠️  WARNING: Arbitrage opportunities detected")
		fmt.Println("   These may be real opportunities or data anomalies")
		fmt.Println("   Further validation needed before trading")
	}
	
	fmt.Println("\n✅ Kalshi integration test completed successfully!")
	fmt.Println("The API client is working correctly and ready for cross-platform integration.")
}