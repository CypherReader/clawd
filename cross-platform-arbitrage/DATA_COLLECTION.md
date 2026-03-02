# 6-Hour Data Collection & Backtesting

This document describes the automated data collection and backtesting system for cross-platform arbitrage opportunities.

## Overview

The system collects market data from Kalshi and PredictIt every 5 minutes for 6 hours, identifies arbitrage opportunities, and generates a comprehensive backtest report.

## Quick Start

### Start Data Collection

```bash
cd /root/clawd/cross-platform-arbitrage
./run-data-collection.sh
```

This will:
- Start a 6-hour background collection process
- Take snapshots every 5 minutes (72 total snapshots)
- Save all data to `cmd/data-collector/data/`
- Generate logs in `logs/`

### Check Status

```bash
./check-collection-status.sh
```

Shows:
- Current progress (%)
- Elapsed time
- Remaining time
- PID and log locations

### View Live Progress

```bash
tail -f logs/collection_*.log
```

## Data Collection Process

### What Gets Collected

Every 5 minutes, the system captures:
1. **Kalshi Markets** - All active markets and contracts
2. **PredictIt Markets** - All active markets and contracts
3. **Market Matching** - Similar markets across platforms
4. **Arbitrage Opportunities** - All detected opportunities with:
   - Platform pair
   - Market names
   - Contract details
   - Buy/sell prices
   - Profit percentage
   - Expected profit ($)
   - Timestamp

### Data Structure

```json
{
  "timestamp": "2026-03-02T21:30:00Z",
  "kalshi_markets": 100,
  "kalshi_contracts": 100,
  "predictit_markets": 258,
  "predictit_contracts": 863,
  "potential_matches": 15,
  "arbitrage_opportunities": [
    {
      "platform1": "kalshi",
      "platform2": "predictit",
      "market1": "Will X win?",
      "market2": "X to win?",
      "strategy": "Buy YES on kalshi @ $0.45 + Buy NO on predictit @ $0.48",
      "total_cost": 0.93,
      "profit_percent": 7.0,
      "expected_profit": 0.07,
      "timestamp": "2026-03-02T21:30:00Z"
    }
  ]
}
```

## Backtest Analysis

After 6 hours, the system automatically generates a comprehensive backtest report.

### Metrics Analyzed

1. **Coverage**
   - Total snapshots collected
   - Markets monitored
   - Contracts analyzed

2. **Opportunities**
   - Total opportunities found
   - Unique market pairs
   - Executable opportunities (≥1% profit)

3. **Profitability**
   - Average profit percentage
   - Maximum profit
   - Minimum profit
   - Distribution of profits

4. **Timing**
   - Opportunity frequency
   - Duration of opportunities
   - Best time windows

5. **Trade Feasibility**
   - Which opportunities could be executed
   - Transaction cost impact
   - Liquidity considerations

### Report Format

```
Data Collection Summary
=======================
Collection Period: 2026-03-02 21:30:00 to 2026-03-03 03:30:00 (6h 0m)
Total Snapshots: 72
Snapshot Interval: ~5m

Market Coverage
---------------
Kalshi Markets: ~100
PredictIt Markets: ~258

Arbitrage Opportunities
-----------------------
Total Opportunities Found: 45
Unique Market Pairs: 12
Executable Opportunities (≥1% profit): 23

Profit Analysis
---------------
Average Profit: 2.3%
Maximum Profit: 7.5%
Minimum Profit: 0.5%

Recommendations
---------------
✅ EXCELLENT: Found 23 opportunities with avg 2.3% profit!
Recommendation: Implement automated execution for opportunities ≥2%.

Next Steps
----------
1. Review detailed opportunity list
2. Validate high-profit opportunities manually
3. Test execution with small amounts
4. Monitor for pattern changes
5. Consider automated execution for ≥2% opportunities
```

## Output Files

### During Collection

- `logs/collection_TIMESTAMP.log` - Live collection log
- `logs/collector.pid` - Process ID
- `logs/collector_status.json` - Current status

### After Collection

- `cmd/data-collector/data/collection_TIMESTAMP.json` - Raw snapshot data (all 72 snapshots)
- `cmd/data-collector/data/backtest_results_TIMESTAMP.json` - Analysis results

## Trade Accuracy Assessment

The system evaluates trade accuracy by:

1. **Opportunity Persistence**
   - How long do opportunities last?
   - Can you act fast enough?

2. **Price Stability**
   - Do prices stay stable during execution?
   - What's the slippage risk?

3. **Market Liquidity**
   - Can you fill orders at displayed prices?
   - What's the volume at bid/ask?

4. **Transaction Costs**
   - Trading fees (typically 0-2%)
   - Withdrawal fees
   - Opportunity cost

5. **Execution Risk**
   - One side fills, other doesn't
   - Price moves during execution
   - Platform downtime

## Example Backtest Interpretation

### Scenario 1: High Frequency, Low Profit
```
Opportunities: 100
Average Profit: 0.8%
Max Profit: 1.5%
```
**Assessment**: Many opportunities but profits barely cover costs. Wait for better conditions.

### Scenario 2: Moderate Frequency, Good Profit
```
Opportunities: 25
Average Profit: 2.5%
Max Profit: 5.0%
```
**Assessment**: Worth pursuing! Focus on opportunities >2% after costs.

### Scenario 3: Low Frequency, High Profit
```
Opportunities: 5
Average Profit: 6.0%
Max Profit: 10.0%
```
**Assessment**: Rare but lucrative. Manual execution recommended. Investigate why they occur.

## Stopping Collection Early

```bash
# Find PID
cat logs/collector.pid

# Stop process
kill $(cat logs/collector.pid)

# Or use status script
./check-collection-status.sh
```

The system will generate a partial backtest report with data collected so far.

## Troubleshooting

### Collection Not Starting

```bash
# Check logs
cat logs/collection_*.log

# Check if port/resources are available
ps aux | grep data-collector
```

### No Opportunities Found

This is normal! Reasons:
- Markets are efficient
- No matching markets across platforms
- Spreads are too wide
- Low liquidity period

### High Opportunities But Low Profit

Common causes:
- Transaction costs not accounted for
- Bid/ask spread too wide
- Low liquidity
- Stale data

## Next Steps After Backtesting

Based on results:

1. **If profitable opportunities exist**:
   - Validate top opportunities manually
   - Test with small amounts
   - Build execution engine

2. **If few opportunities**:
   - Expand to more platforms (Polymarket)
   - Reduce snapshot interval (faster detection)
   - Improve market matching algorithm

3. **If marginal profits**:
   - Focus on high-volume markets
   - Negotiate lower fees
   - Improve execution speed

## Advanced Configuration

Edit `cmd/data-collector/main.go`:

```go
// Change collection duration
collectionDuration := 6 * time.Hour  // Modify this

// Change snapshot frequency
snapshotInterval := 5 * time.Minute  // Modify this

// Change profit threshold for "executable"
if opp.ProfitPercent >= 1.0 {  // Modify this
    results.ExecutableOpportunities++
}
```

## Integration with Oracle Service

This data collection feeds into the broader Oracle service:
- Market predictions
- Price forecasts
- Trading recommendations
- Portfolio optimization

See main Oracle documentation for full integration details.
