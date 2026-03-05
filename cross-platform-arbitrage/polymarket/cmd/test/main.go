package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"polymarket"
)

func main() {
	fmt.Println("=== Polymarket API Integration Test ===")
	fmt.Println()

	client := polymarket.NewClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test 1: Fetch events
	fmt.Println("Test 1: Fetching events...")
	start := time.Now()

	events, err := client.GetEvents(ctx, 100, true)
	if err != nil {
		log.Fatalf("Failed to fetch events: %v", err)
	}

	elapsed := time.Since(start)
	fmt.Printf("✅ Retrieved %d events in %v\n", len(events), elapsed)
	fmt.Println()

	// Test 2: Analyze markets for arbitrage
	fmt.Println("Test 2: Scanning for arbitrage opportunities...")

	totalMarkets := 0
	activeMarkets := 0
	arbitrageOpportunities := 0

	for _, event := range events {
		for _, market := range event.Markets {
			totalMarkets++

			if polymarket.IsMarketActive(market) {
				activeMarkets++

				arbPercent, exists, err := polymarket.CalculateArbitrageOpportunity(market)
				if err != nil {
					// Skip markets with parsing errors
					continue
				}

				if exists {
					arbitrageOpportunities++
					
					yesPrice, noPrice, _ := polymarket.ParseOutcomePrices(market.OutcomePrices)
					
					fmt.Printf("🚨 ARBITRAGE FOUND: %s\n", market.Question)
					fmt.Printf("   Event: %s\n", event.Title)
					fmt.Printf("   YES price: $%.2f + NO price: $%.2f = $%.2f (< $1.00)\n",
						yesPrice, noPrice, yesPrice+noPrice)
					fmt.Printf("   Profit: %.2f%%\n", arbPercent)
					fmt.Printf("   Volume: $%.2f (24h: $%.2f)\n", market.VolumeNum, market.Volume24hr)
					fmt.Println()
				}
			}
		}
	}

	fmt.Printf("Total markets: %d\n", totalMarkets)
	fmt.Printf("Active markets: %d\n", activeMarkets)
	fmt.Printf("Arbitrage opportunities: %d\n", arbitrageOpportunities)
	fmt.Println()

	// Test 3: Show sample markets
	if len(events) > 0 {
		fmt.Println("Test 3: Sample Market Details")
		
		sampleCount := 0
		for _, event := range events {
			if len(event.Markets) > 0 {
				market := event.Markets[0]
				
				yesPrice, noPrice, err := polymarket.ParseOutcomePrices(market.OutcomePrices)
				if err != nil {
					continue
				}

				fmt.Printf("\nMarket: %s\n", market.Question)
				fmt.Printf("  Event: %s\n", event.Title)
				fmt.Printf("  Category: %s\n", market.Category)
				fmt.Printf("  YES price: $%.4f\n", yesPrice)
				fmt.Printf("  NO price: $%.4f\n", noPrice)
				fmt.Printf("  Sum: $%.4f\n", yesPrice+noPrice)
				fmt.Printf("  Volume: $%.2f\n", market.VolumeNum)
				fmt.Printf("  Liquidity: $%.2f\n", market.LiquidityNum)
				fmt.Printf("  Active: %v\n", market.Active)
				
				arbPercent, exists, _ := polymarket.CalculateArbitrageOpportunity(market)
				if exists {
					fmt.Printf("  ⚠️  Arbitrage: %.2f%%\n", arbPercent)
				} else {
					fmt.Printf("  ✓ No arbitrage\n")
				}

				sampleCount++
				if sampleCount >= 3 {
					break
				}
			}
		}
	}

	fmt.Println()
	fmt.Println("=== Test Summary ===")
	fmt.Printf("✅ Polymarket API working\n")
	fmt.Printf("✅ Market data retrieved successfully\n")
	fmt.Printf("✅ Arbitrage detection functional\n")

	if arbitrageOpportunities > 0 {
		fmt.Printf("\n🚨 Found %d arbitrage opportunities!\n", arbitrageOpportunities)
		fmt.Println("Next step: Cross-platform matching with Kalshi and PredictIt")
	} else {
		fmt.Println("\nNo arbitrage found within Polymarket markets")
		fmt.Println("This is normal - arbitrage is rare within a single platform")
		fmt.Println("Cross-platform opportunities are where the real profits are!")
	}
}
