---

## 3. PredictIt (Politics-Focused)

### Overview
- **Type**: Regulated prediction market (US)
- **Focus**: Primarily political markets
- **Regulation**: CFTC no-action letter (specific exemption)
- **Ownership**: Academic (Victoria University of Wellington)
- **Key Feature**: $850 maximum investment per contract

### API Research - ✅ PARTIAL

**API Status**: PredictIt has an unofficial/undocumented API that community has reverse-engineered

**Key Findings**:
1. **No Official Public API**: PredictIt doesn't have official public API documentation
2. **Community Libraries**: Several Python/JavaScript libraries exist that use reverse-engineered endpoints
3. **Data Format**: Returns XML/CSV data (not JSON)
4. **Rate Limits**: Unknown, likely conservative due to academic/research focus

**Community Libraries Found**:
1. **Python**: `predictit-markets` (PyPI package)
2. **Python**: `predictit` (multiple GitHub repos)
3. **JavaScript**: `node-predict-it` (Node.js wrapper)
4. **R**: `rpredictit` (R package)

**Example Usage** (from `predictit-markets` package):
```python
from predictit_markets import market_data, market_name

# Get market data for market ID 6598 (2020 Washington presidential election)
df = market_data(6598)
market_name = market_name(6598)
```

**Data Structure** (inferred from community packages):
- Market ID-based access
- Historical data available (24h, 7, 30, 90 days)
- Contract-level data (multiple contracts per market)
- Price, volume, time series data

**Access Requirements**:
- Public market data: Available via undocumented endpoints
- Trading: Requires account (academic/research focus)
- Non-commercial use only (based on terms)

**Limitations**:
1. **Unofficial API**: No documentation, subject to change
2. **Academic Focus**: Non-commercial use emphasis
3. **Political Focus**: Primarily US political markets
4. **Investment Caps**: $850 maximum per contract

**Potential Endpoints** (inferred):
- Market data: `https://www.predictit.org/api/marketdata/markets/{id}`
- Historical: `https://www.predictit.org/api/marketdata/graph/{id}`

**Note**: Need to verify actual endpoints and test access