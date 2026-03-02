package predictit

import (
	"context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	baseURL = "https://www.predictit.org/api/marketdata"
)

// Client represents a PredictIt API client
type Client struct {
	httpClient *resty.Client
}

// NewClient creates a new PredictIt API client
func NewClient() *Client {
	client := resty.New().
		SetBaseURL(baseURL).
		SetTimeout(10 * time.Second).
		SetRetryCount(3).
		SetRetryWaitTime(1 * time.Second).
		SetRetryMaxWaitTime(5 * time.Second)

	return &Client{
		httpClient: client,
	}
}

// Market represents a PredictIt market
type Market struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	ShortName string     `json:"shortName"`
	Image     string     `json:"image"`
	URL       string     `json:"url"`
	Contracts []Contract `json:"contracts"`
	TimeStamp string     `json:"timeStamp"` // Using string to avoid parsing issues
	Status    string     `json:"status"`
}

// Contract represents a contract within a market
type Contract struct {
	ID               int     `json:"id"`
	DateEnd          string  `json:"dateEnd"`
	Image            string  `json:"image"`
	Name             string  `json:"name"`
	ShortName        string  `json:"shortName"`
	Status           string  `json:"status"`
	LastTradePrice   float64 `json:"lastTradePrice"`
	BestBuyYesCost   float64 `json:"bestBuyYesCost"`   // Price to BUY YES shares
	BestBuyNoCost    float64 `json:"bestBuyNoCost"`    // Price to BUY NO shares
	BestSellYesCost  float64 `json:"bestSellYesCost"`  // Price to SELL YES shares (bid)
	BestSellNoCost   float64 `json:"bestSellNoCost"`   // Price to SELL NO shares (bid)
	LastClosePrice   float64 `json:"lastClosePrice"`
	DisplayOrder     int     `json:"displayOrder"`
}

// MarketsResponse represents the response from /marketdata/all/
type MarketsResponse struct {
	Markets []Market `json:"markets"`
}

// GetAllMarkets retrieves all markets from PredictIt
func (c *Client) GetAllMarkets(ctx context.Context) (*MarketsResponse, error) {
	var result MarketsResponse

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetResult(&result).
		Get("/all/")

	if err != nil {
		return nil, fmt.Errorf("failed to fetch markets: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API returned error: %s", resp.Status())
	}

	return &result, nil
}

// GetMarket retrieves a specific market by ID
func (c *Client) GetMarket(ctx context.Context, marketID int) (*Market, error) {
	var result Market

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetResult(&result).
		Get(fmt.Sprintf("/markets/%d", marketID))

	if err != nil {
		return nil, fmt.Errorf("failed to fetch market %d: %w", marketID, err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API returned error: %s", resp.Status())
	}

	return &result, nil
}

// CalculateArbitrageOpportunity checks for arbitrage in a binary contract
// PredictIt pricing: bestSellYesCost is the bid (what you can sell YES for)
//                    bestBuyYesCost is the ask (what you pay to buy YES)
// Returns arbitrage percentage if opportunity exists
func CalculateArbitrageOpportunity(contract Contract) (float64, bool) {
	// Skip contracts without valid pricing
	if contract.BestBuyYesCost == 0 || contract.BestBuyNoCost == 0 {
		return 0.0, false
	}

	// In PredictIt:
	// - bestBuyYesCost = price to BUY YES (ask price for YES)
	// - bestBuyNoCost = price to BUY NO (ask price for NO)
	// For arbitrage, we want: buy YES + buy NO < $1.00

	totalCost := contract.BestBuyYesCost + contract.BestBuyNoCost

	if totalCost < 1.0 {
		arbPercent := (1.0 - totalCost) * 100
		return arbPercent, true
	}

	return 0.0, false
}

// IsContractActive checks if a contract is actively trading
func IsContractActive(contract Contract) bool {
	return contract.Status == "Open" &&
		contract.BestBuyYesCost > 0 &&
		contract.BestBuyNoCost > 0
}

// IsMarketActive checks if a market has active contracts
func IsMarketActive(market Market) bool {
	for _, contract := range market.Contracts {
		if IsContractActive(contract) {
			return true
		}
	}
	return false
}