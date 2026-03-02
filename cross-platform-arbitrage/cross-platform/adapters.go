package crossplatform

import (
	"fmt"
)

// KalshiMarket interface to avoid circular imports
type KalshiMarket interface {
	GetTicker() string
	GetTitle() string
	GetYesBid() int
	GetYesAsk() int
	GetNoBid() int
	GetNoAsk() int
	GetLastPrice() int
	GetStatus() string
}

// PredictItContract interface
type PredictItContract interface {
	GetID() int
	GetName() string
	GetShortName() string
	GetBestBuyYesCost() float64
	GetBestBuyNoCost() float64
	GetBestSellYesCost() float64
	GetBestSellNoCost() float64
	GetLastTradePrice() float64
	GetStatus() string
}

// PredictItMarket interface
type PredictItMarket interface {
	GetID() int
	GetName() string
	GetShortName() string
	GetURL() string
	GetContracts() []PredictItContract
}

// ConvertKalshiToUnified converts a Kalshi market to unified format
func ConvertKalshiToUnified(ticker, title string, yesBid, yesAsk, noBid, noAsk, lastPrice int, status string) UnifiedContract {
	return UnifiedContract{
		Platform:     "kalshi",
		MarketID:     ticker,
		MarketName:   title,
		ContractName: "Binary Outcome",
		YesBid:       float64(yesBid) / 100.0,   // Kalshi uses cents
		YesAsk:       float64(yesAsk) / 100.0,
		NoBid:        float64(noBid) / 100.0,
		NoAsk:        float64(noAsk) / 100.0,
		LastPrice:    float64(lastPrice) / 100.0,
		Status:       status,
		URL:          fmt.Sprintf("https://kalshi.com/markets/%s", ticker),
	}
}

// ConvertPredictItToUnified converts a PredictIt contract to unified format
func ConvertPredictItToUnified(marketID int, marketName, marketURL string, contractID int, contractName string, 
	bestBuyYesCost, bestBuyNoCost, bestSellYesCost, bestSellNoCost, lastTradePrice float64, status string) UnifiedContract {
	return UnifiedContract{
		Platform:     "predictit",
		MarketID:     fmt.Sprintf("%d-%d", marketID, contractID),
		MarketName:   marketName,
		ContractName: contractName,
		YesBid:       bestSellYesCost,  // What you receive when selling YES
		YesAsk:       bestBuyYesCost,   // What you pay to buy YES
		NoBid:        bestSellNoCost,   // What you receive when selling NO
		NoAsk:        bestBuyNoCost,    // What you pay to buy NO
		LastPrice:    lastTradePrice,
		Status:       status,
		URL:          marketURL,
	}
}