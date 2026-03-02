package crossplatform

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// UnifiedContract represents a normalized contract across platforms
type UnifiedContract struct {
	Platform     string  // "kalshi" or "predictit"
	MarketID     string  // Platform-specific market ID
	MarketName   string  // Human-readable market name
	ContractName string  // Contract/outcome name
	YesBid       float64 // Price you receive when selling YES (in dollars)
	YesAsk       float64 // Price you pay to buy YES (in dollars)
	NoBid        float64 // Price you receive when selling NO (in dollars)
	NoAsk        float64 // Price you pay to buy NO (in dollars)
	LastPrice    float64 // Last traded price
	Status       string  // Market status
	URL          string  // Link to market
}

// ArbitrageOpportunity represents a cross-platform arbitrage opportunity
type ArbitrageOpportunity struct {
	Platform1   string
	Platform2   string
	Contract1   UnifiedContract
	Contract2   UnifiedContract
	Strategy    string  // Description of the arbitrage strategy
	ProfitPercent float64
	TotalCost   float64 // Cost to execute arbitrage
	ExpectedProfit float64 // Expected profit in dollars
}

// NormalizeMarketName attempts to normalize market names for matching
func NormalizeMarketName(name string) string {
	// Convert to lowercase
	normalized := strings.ToLower(name)

	// Remove common punctuation
	normalized = strings.ReplaceAll(normalized, "?", "")
	normalized = strings.ReplaceAll(normalized, "!", "")
	normalized = strings.ReplaceAll(normalized, ".", "")
	normalized = strings.ReplaceAll(normalized, ",", "")

	// Remove extra whitespace
	normalized = strings.TrimSpace(normalized)

	return normalized
}

// FindCrossPlatformArbitrage identifies arbitrage opportunities between two platforms
func FindCrossPlatformArbitrage(contracts1, contracts2 []UnifiedContract) []ArbitrageOpportunity {
	var opportunities []ArbitrageOpportunity

	// Build a map of normalized market names to contracts for faster lookup
	marketMap := make(map[string][]UnifiedContract)
	for _, c := range contracts2 {
		normalized := NormalizeMarketName(c.MarketName + " " + c.ContractName)
		marketMap[normalized] = append(marketMap[normalized], c)
	}

	// Check each contract from platform 1
	for _, c1 := range contracts1 {
		normalized1 := NormalizeMarketName(c1.MarketName + " " + c1.ContractName)

		// Try to find matching contracts in platform 2
		if matches, exists := marketMap[normalized1]; exists {
			for _, c2 := range matches {
				// Check for arbitrage opportunities

				// Strategy 1: Buy YES on platform 1, Buy NO on platform 2
				cost1 := c1.YesAsk + c2.NoAsk
				if cost1 < 1.0 && c1.YesAsk > 0 && c2.NoAsk > 0 {
					profit := 1.0 - cost1
					opportunities = append(opportunities, ArbitrageOpportunity{
						Platform1:      c1.Platform,
						Platform2:      c2.Platform,
						Contract1:      c1,
						Contract2:      c2,
						Strategy:       fmt.Sprintf("Buy YES on %s @ $%.2f + Buy NO on %s @ $%.2f", c1.Platform, c1.YesAsk, c2.Platform, c2.NoAsk),
						ProfitPercent:  profit * 100,
						TotalCost:      cost1,
						ExpectedProfit: profit,
					})
				}

				// Strategy 2: Buy NO on platform 1, Buy YES on platform 2
				cost2 := c1.NoAsk + c2.YesAsk
				if cost2 < 1.0 && c1.NoAsk > 0 && c2.YesAsk > 0 {
					profit := 1.0 - cost2
					opportunities = append(opportunities, ArbitrageOpportunity{
						Platform1:      c1.Platform,
						Platform2:      c2.Platform,
						Contract1:      c1,
						Contract2:      c2,
						Strategy:       fmt.Sprintf("Buy NO on %s @ $%.2f + Buy YES on %s @ $%.2f", c1.Platform, c1.NoAsk, c2.Platform, c2.YesAsk),
						ProfitPercent:  profit * 100,
						TotalCost:      cost2,
						ExpectedProfit: profit,
					})
				}
			}
		}
	}

	return opportunities
}

// MatchMarkets attempts to find similar markets across platforms
// Returns a map of matched market names
func MatchMarkets(contracts1, contracts2 []UnifiedContract) map[string][]string {
	matches := make(map[string][]string)

	// Build unique market names from both platforms
	markets1 := make(map[string]bool)
	markets2 := make(map[string]bool)

	for _, c := range contracts1 {
		normalized := NormalizeMarketName(c.MarketName)
		markets1[normalized] = true
	}

	for _, c := range contracts2 {
		normalized := NormalizeMarketName(c.MarketName)
		markets2[normalized] = true
	}

	// Find common markets
	for market1 := range markets1 {
		for market2 := range markets2 {
			// Simple substring matching (can be improved with fuzzy matching)
			if strings.Contains(market1, market2) || strings.Contains(market2, market1) {
				matches[market1] = append(matches[market1], market2)
			}
		}
	}

	return matches
}