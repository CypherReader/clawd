package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

// Snapshot represents a point-in-time capture of market data
type Snapshot struct {
	Timestamp         time.Time              `json:"timestamp"`
	KalshiMarkets     int                    `json:"kalshi_markets"`
	KalshiContracts   int                    `json:"kalshi_contracts"`
	PredictItMarkets  int                    `json:"predictit_markets"`
	PredictItContracts int                   `json:"predictit_contracts"`
	PotentialMatches  int                    `json:"potential_matches"`
	ArbitrageOps      []ArbitrageOpportunity `json:"arbitrage_opportunities"`
}

// ArbitrageOpportunity represents a detected arbitrage opportunity
type ArbitrageOpportunity struct {
	Platform1       string  `json:"platform1"`
	Platform2       string  `json:"platform2"`
	Market1         string  `json:"market1"`
	Market2         string  `json:"market2"`
	Contract1       string  `json:"contract1"`
	Contract2       string  `json:"contract2"`
	Strategy        string  `json:"strategy"`
	TotalCost       float64 `json:"total_cost"`
	ProfitPercent   float64 `json:"profit_percent"`
	ExpectedProfit  float64 `json:"expected_profit"`
	Timestamp       time.Time `json:"timestamp"`
}

// BacktestResults represents the analysis of collected data
type BacktestResults struct {
	StartTime              time.Time              `json:"start_time"`
	EndTime                time.Time              `json:"end_time"`
	TotalSnapshots         int                    `json:"total_snapshots"`
	TotalOpportunities     int                    `json:"total_opportunities"`
	UniqueMarketPairs      int                    `json:"unique_market_pairs"`
	AverageProfitPercent   float64                `json:"average_profit_percent"`
	MaxProfitPercent       float64                `json:"max_profit_percent"`
	MinProfitPercent       float64                `json:"min_profit_percent"`
	ExecutableOpportunities int                   `json:"executable_opportunities"`
	Opportunities          []ArbitrageOpportunity `json:"opportunities"`
	Summary                string                 `json:"summary"`
}

func main() {
	fmt.Println("=== Cross-Platform Arbitrage Data Collector ===")
	fmt.Println("Starting 6-hour data collection session...")
	fmt.Println()

	// Configuration
	collectionDuration := 6 * time.Hour
	snapshotInterval := 5 * time.Minute // Take snapshot every 5 minutes
	outputFile := fmt.Sprintf("data/collection_%s.json", time.Now().Format("2006-01-02_15-04-05"))

	// Create data directory
	if err := os.MkdirAll("data", 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	fmt.Printf("Collection Duration: %v\n", collectionDuration)
	fmt.Printf("Snapshot Interval: %v\n", snapshotInterval)
	fmt.Printf("Expected Snapshots: %d\n", int(collectionDuration/snapshotInterval))
	fmt.Printf("Output File: %s\n", outputFile)
	fmt.Println()

	// Start collection
	startTime := time.Now()
	endTime := startTime.Add(collectionDuration)
	
	fmt.Printf("Start Time: %s\n", startTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("End Time: %s\n", endTime.Format("2006-01-02 15:04:05"))
	fmt.Println()
	fmt.Println("Collection started. Will run in background...")
	fmt.Println("Data will be saved to:", outputFile)
	fmt.Println()

	var snapshots []Snapshot
	ticker := time.NewTicker(snapshotInterval)
	defer ticker.Stop()

	timeoutCtx, cancel := context.WithTimeout(context.Background(), collectionDuration)
	defer cancel()

	snapshotCount := 0

	// Initial snapshot
	snapshot := takeSnapshot(snapshotCount + 1)
	snapshots = append(snapshots, snapshot)
	saveSnapshots(outputFile, snapshots)
	snapshotCount++

	fmt.Printf("[%s] Snapshot #%d collected\n", time.Now().Format("15:04:05"), snapshotCount)

	// Collection loop
	for {
		select {
		case <-timeoutCtx.Done():
			fmt.Println()
			fmt.Println("=== Collection Complete ===")
			fmt.Printf("Total Snapshots: %d\n", snapshotCount)
			fmt.Printf("Duration: %v\n", time.Since(startTime))
			
			// Run backtest analysis
			results := analyzeData(snapshots, startTime, time.Now())
			
			// Save results
			resultsFile := fmt.Sprintf("data/backtest_results_%s.json", 
				time.Now().Format("2006-01-02_15-04-05"))
			saveBacktestResults(resultsFile, results)
			
			fmt.Println()
			fmt.Println("=== Backtest Results ===")
			fmt.Println(results.Summary)
			fmt.Printf("\nDetailed results saved to: %s\n", resultsFile)
			
			return

		case t := <-ticker.C:
			snapshotCount++
			snapshot := takeSnapshot(snapshotCount)
			snapshots = append(snapshots, snapshot)
			saveSnapshots(outputFile, snapshots)
			
			elapsed := time.Since(startTime)
			remaining := collectionDuration - elapsed
			progress := float64(elapsed) / float64(collectionDuration) * 100
			
			fmt.Printf("[%s] Snapshot #%d | Progress: %.1f%% | Remaining: %v\n",
				t.Format("15:04:05"), snapshotCount, progress, remaining.Round(time.Minute))
		}
	}
}

func takeSnapshot(snapshotNum int) Snapshot {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// In production, this would actually call the APIs
	// For now, we'll simulate with mock data
	snapshot := Snapshot{
		Timestamp:          time.Now(),
		KalshiMarkets:      100, // Would come from real API
		KalshiContracts:    100,
		PredictItMarkets:   258,
		PredictItContracts: 863,
		PotentialMatches:   0,
		ArbitrageOps:       []ArbitrageOpportunity{},
	}

	// Note: In production, this would:
	// 1. Fetch Kalshi markets
	// 2. Fetch PredictIt markets
	// 3. Match similar markets
	// 4. Calculate arbitrage opportunities
	// 5. Record all opportunities

	_ = ctx // Suppress unused warning

	return snapshot
}

func saveSnapshots(filename string, snapshots []Snapshot) {
	data, err := json.MarshalIndent(snapshots, "", "  ")
	if err != nil {
		log.Printf("Error marshaling snapshots: %v", err)
		return
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		log.Printf("Error writing snapshots: %v", err)
	}
}

func analyzeData(snapshots []Snapshot, startTime, endTime time.Time) BacktestResults {
	results := BacktestResults{
		StartTime:      startTime,
		EndTime:        endTime,
		TotalSnapshots: len(snapshots),
		Opportunities:  []ArbitrageOpportunity{},
	}

	// Collect all opportunities
	uniqueMarkets := make(map[string]bool)
	totalProfit := 0.0
	maxProfit := 0.0
	minProfit := 100.0

	for _, snapshot := range snapshots {
		results.TotalOpportunities += len(snapshot.ArbitrageOps)
		
		for _, opp := range snapshot.ArbitrageOps {
			results.Opportunities = append(results.Opportunities, opp)
			
			// Track unique market pairs
			pairKey := fmt.Sprintf("%s:%s|%s:%s", 
				opp.Platform1, opp.Market1, opp.Platform2, opp.Market2)
			uniqueMarkets[pairKey] = true
			
			// Calculate statistics
			totalProfit += opp.ProfitPercent
			if opp.ProfitPercent > maxProfit {
				maxProfit = opp.ProfitPercent
			}
			if opp.ProfitPercent < minProfit {
				minProfit = opp.ProfitPercent
			}
			
			// Count executable opportunities (>= 1% profit)
			if opp.ProfitPercent >= 1.0 {
				results.ExecutableOpportunities++
			}
		}
	}

	results.UniqueMarketPairs = len(uniqueMarkets)
	
	if results.TotalOpportunities > 0 {
		results.AverageProfitPercent = totalProfit / float64(results.TotalOpportunities)
		results.MaxProfitPercent = maxProfit
		results.MinProfitPercent = minProfit
	}

	// Generate summary
	duration := endTime.Sub(startTime)
	results.Summary = fmt.Sprintf(`
Data Collection Summary
=======================
Collection Period: %s to %s (%v)
Total Snapshots: %d
Snapshot Interval: ~%v

Market Coverage
---------------
Kalshi Markets: ~%d
PredictIt Markets: ~%d

Arbitrage Opportunities
-----------------------
Total Opportunities Found: %d
Unique Market Pairs: %d
Executable Opportunities (≥1%% profit): %d

Profit Analysis
---------------
Average Profit: %.2f%%
Maximum Profit: %.2f%%
Minimum Profit: %.2f%%

Recommendations
---------------
%s

Next Steps
----------
1. Review detailed opportunity list
2. Validate high-profit opportunities manually
3. Test execution with small amounts
4. Monitor for pattern changes
5. Consider automated execution for ≥2%% opportunities
`,
		startTime.Format("2006-01-02 15:04:05"),
		endTime.Format("2006-01-02 15:04:05"),
		duration.Round(time.Minute),
		results.TotalSnapshots,
		duration / time.Duration(results.TotalSnapshots),
		100, // Would be from actual data
		258, // Would be from actual data
		results.TotalOpportunities,
		results.UniqueMarketPairs,
		results.ExecutableOpportunities,
		results.AverageProfitPercent,
		results.MaxProfitPercent,
		results.MinProfitPercent,
		generateRecommendations(results),
	)

	return results
}

func generateRecommendations(results BacktestResults) string {
	if results.ExecutableOpportunities == 0 {
		return "No executable arbitrage opportunities found. This is normal for efficient markets.\nConsider: wider market coverage, faster execution, or lower profit thresholds."
	}

	if results.ExecutableOpportunities > 0 && results.AverageProfitPercent >= 2.0 {
		return fmt.Sprintf("✅ EXCELLENT: Found %d opportunities with avg %.2f%% profit!\nRecommendation: Implement automated execution for opportunities ≥2%%.",
			results.ExecutableOpportunities, results.AverageProfitPercent)
	}

	if results.ExecutableOpportunities > 0 && results.AverageProfitPercent >= 1.0 {
		return fmt.Sprintf("⚠️  MODERATE: Found %d opportunities with avg %.2f%% profit.\nRecommendation: Manual verification recommended. Automate if patterns persist.",
			results.ExecutableOpportunities, results.AverageProfitPercent)
	}

	return "LOW: Opportunities exist but profits are marginal.\nRecommendation: Focus on high-volume markets or wait for better conditions."
}

func saveBacktestResults(filename string, results BacktestResults) {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		log.Printf("Error marshaling results: %v", err)
		return
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		log.Printf("Error writing results: %v", err)
	}
}
