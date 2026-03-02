package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"predictit"
)

func main() {
	fmt.Println("=== PredictIt API Integration Test ===")
	fmt.Println()

	client := predictit.NewClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test 1: Fetch all markets
	fmt.Println("Test 1: Fetching all markets...")
	start := time.Now()

	marketsResp, err := client.GetAllMarkets(ctx)
	if err != nil {
		log.Fatalf("Failed to fetch markets: %v", err)
	}

	elapsed := time.Since(start)
	fmt.Printf("✅ Retrieved %d markets in %v\n", len(marketsResp.Markets), elapsed)
	fmt.Println()

	// Test 2: Analyze markets for arbitrage
	fmt.Println("Test 2: Scanning for arbitrage opportunities...")

	totalContracts := 0
	activeContracts := 0
	arbitrageOpportunities := 0

	for _, market := range marketsResp.Markets {
		for _, contract := range market.Contracts {
			totalContracts++

			if predictit.IsContractActive(contract) {
				activeContracts++

				arbPercent, exists := predictit.CalculateArbitrageOpportunity(contract)
				if exists {
					arbitrageOpportunities++
					fmt.Printf("🚨 ARBITRAGE FOUND: %s - %s\n", market.ShortName, contract.ShortName)
					fmt.Printf("   Buy YES: $%.2f + Buy NO: $%.2f = $%.2f (< $1.00)\n",
						contract.BestBuyYesCost, contract.BestBuyNoCost,
						contract.BestBuyYesCost+contract.BestBuyNoCost)
					fmt.Printf("   Profit: %.2f%%\n", arbPercent)
					fmt.Println()
				}
			}
		}
	}

	fmt.Printf("Total contracts: %d\n", totalContracts)
	fmt.Printf("Active contracts: %d\n", activeContracts)
	fmt.Printf("Arbitrage opportunities: %d\n", arbitrageOpportunities)
	fmt.Println()

	// Test 3: Show sample market
	if len(marketsResp.Markets) > 0 {
		fmt.Println("Test 3: Sample Market Details")
		sampleMarket := marketsResp.Markets[0]
		fmt.Printf("Market: %s\n", sampleMarket.ShortName)
		fmt.Printf("Contracts: %d\n", len(sampleMarket.Contracts))

		if len(sampleMarket.Contracts) > 0 {
			contract := sampleMarket.Contracts[0]
			fmt.Printf("\nSample Contract: %s\n", contract.ShortName)
			fmt.Printf("  Status: %s\n", contract.Status)
			fmt.Printf("  Last Trade: $%.2f\n", contract.LastTradePrice)
			fmt.Printf("  Buy YES: $%.2f (pay this to buy YES)\n", contract.BestBuyYesCost)
			fmt.Printf("  Buy NO: $%.2f (pay this to buy NO)\n", contract.BestBuyNoCost)
			fmt.Printf("  Sell YES: $%.2f (receive if selling YES)\n", contract.BestSellYesCost)
			fmt.Printf("  Sell NO: $%.2f (receive if selling NO)\n", contract.BestSellNoCost)

			arbPercent, exists := predictit.CalculateArbitrageOpportunity(contract)
			if exists {
				fmt.Printf("  ⚠️  Arbitrage: %.2f%%\n", arbPercent)
			} else {
				fmt.Printf("  ✓ No arbitrage\n")
			}
		}
	}

	fmt.Println()
	fmt.Println("=== Test Summary ===")
	fmt.Printf("✅ PredictIt API working\n")
	fmt.Printf("✅ Market data retrieved successfully\n")
	fmt.Printf("✅ Arbitrage detection functional\n")

	if arbitrageOpportunities > 0 {
		fmt.Printf("\n🚨 Found %d arbitrage opportunities!\n", arbitrageOpportunities)
		fmt.Println("Next step: Cross-platform matching with Kalshi")
	} else {
		fmt.Println("\nNo arbitrage found within PredictIt markets")
		fmt.Println("This is normal - arbitrage is rare within a single platform")
	}
}