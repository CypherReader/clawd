package polymarket

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	gammaAPIBaseURL = "https://gamma-api.polymarket.com"
)

// Client represents a Polymarket API client
type Client struct {
	httpClient *resty.Client
}

// NewClient creates a new Polymarket API client
func NewClient() *Client {
	client := resty.New().
		SetBaseURL(gammaAPIBaseURL).
		SetTimeout(15 * time.Second).
		SetRetryCount(3).
		SetRetryWaitTime(1 * time.Second).
		SetRetryMaxWaitTime(5 * time.Second)

	return &Client{
		httpClient: client,
	}
}

// Event represents a Polymarket event (can contain multiple markets)
type Event struct {
	ID            string    `json:"id"`
	Ticker        string    `json:"ticker"`
	Slug          string    `json:"slug"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	StartDate     string    `json:"startDate"`
	EndDate       string    `json:"endDate"`
	Image         string    `json:"image"`
	Icon          string    `json:"icon"`
	Active        bool      `json:"active"`
	Closed        bool      `json:"closed"`
	Archived      bool      `json:"archived"`
	Volume        float64   `json:"volume"`
	Liquidity     float64   `json:"liquidity"`
	Volume24hr    float64   `json:"volume24hr"`
	Markets       []Market  `json:"markets"`
	Category      string    `json:"category"`
	ClosedTime    string    `json:"closedTime"`
}

// Market represents a Polymarket market (specific outcome pair)
type Market struct {
	ID                 string  `json:"id"`
	Question           string  `json:"question"`
	ConditionID        string  `json:"conditionId"`
	Slug               string  `json:"slug"`
	EndDate            string  `json:"endDate"`
	StartDate          string  `json:"startDate"`
	Category           string  `json:"category"`
	Image              string  `json:"image"`
	Icon               string  `json:"icon"`
	Description        string  `json:"description"`
	Outcomes           string  `json:"outcomes"`           // JSON array string: ["Yes", "No"]
	OutcomePrices      string  `json:"outcomePrices"`      // JSON array string: ["0.45", "0.55"]
	Volume             string  `json:"volume"`
	Active             bool    `json:"active"`
	Closed             bool    `json:"closed"`
	Archived           bool    `json:"archived"`
	VolumeNum          float64 `json:"volumeNum"`
	LiquidityNum       float64 `json:"liquidityNum"`
	Volume24hr         float64 `json:"volume24hr"`
	LastTradePrice     float64 `json:"lastTradePrice"`
	BestBid            float64 `json:"bestBid"`
	BestAsk            float64 `json:"bestAsk"`
	Spread             float64 `json:"spread"`
}

// EventsResponse represents the response from /events endpoint
type EventsResponse struct {
	Events []Event `json:"events,omitempty"`
}

// MarketsResponse represents the response from /markets endpoint
type MarketsResponse struct {
	Markets []Market `json:"markets,omitempty"`
}

// GetEvents retrieves events from Polymarket
func (c *Client) GetEvents(ctx context.Context, limit int, active bool) ([]Event, error) {
	var events []Event

	params := map[string]string{
		"limit":    strconv.Itoa(limit),
		"active":   strconv.FormatBool(active),
		"closed":   "false",
		"archived": "false",
	}

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetQueryParams(params).
		SetResult(&events).
		Get("/events")

	if err != nil {
		return nil, fmt.Errorf("failed to fetch events: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API returned error: %s", resp.Status())
	}

	return events, nil
}

// GetMarkets retrieves markets from Polymarket
func (c *Client) GetMarkets(ctx context.Context, limit int, active bool) ([]Market, error) {
	var markets []Market

	params := map[string]string{
		"limit":    strconv.Itoa(limit),
		"active":   strconv.FormatBool(active),
		"closed":   "false",
		"archived": "false",
	}

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetQueryParams(params).
		SetResult(&markets).
		Get("/markets")

	if err != nil {
		return nil, fmt.Errorf("failed to fetch markets: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API returned error: %s", resp.Status())
	}

	return markets, nil
}

// ParseOutcomePrices extracts Yes and No prices from the outcomePrices JSON string
// Returns (yesPrice, noPrice, error)
func ParseOutcomePrices(outcomePrices string) (float64, float64, error) {
	// outcomePrices is a JSON array string like: ["0.45", "0.55"]
	// Parse manually to avoid extra dependencies
	
	if len(outcomePrices) < 2 {
		return 0, 0, fmt.Errorf("invalid outcomePrices format")
	}

	// Simple parsing: find the two numbers
	var prices []float64
	var currentNum string
	inQuote := false

	for i := 0; i < len(outcomePrices); i++ {
		ch := outcomePrices[i]
		
		if ch == '"' {
			if inQuote && len(currentNum) > 0 {
				price, err := strconv.ParseFloat(currentNum, 64)
				if err != nil {
					return 0, 0, fmt.Errorf("failed to parse price: %w", err)
				}
				prices = append(prices, price)
				currentNum = ""
			}
			inQuote = !inQuote
		} else if inQuote {
			currentNum += string(ch)
		}
	}

	if len(prices) != 2 {
		return 0, 0, fmt.Errorf("expected 2 prices, got %d", len(prices))
	}

	return prices[0], prices[1], nil
}

// CalculateArbitrageOpportunity checks for arbitrage in a Polymarket market
// Polymarket prices are already in decimal format (0.45 = $0.45 = 45%)
func CalculateArbitrageOpportunity(market Market) (float64, bool, error) {
	yesPrice, noPrice, err := ParseOutcomePrices(market.OutcomePrices)
	if err != nil {
		return 0, false, err
	}

	// Skip markets with zero prices
	if yesPrice == 0 || noPrice == 0 {
		return 0, false, nil
	}

	// For arbitrage, we need: yesAsk + noAsk < $1.00
	// In Polymarket, the price shown is effectively the "ask" price (what you pay to buy)
	totalCost := yesPrice + noPrice

	if totalCost < 1.0 {
		arbPercent := (1.0 - totalCost) * 100
		return arbPercent, true, nil
	}

	return 0, false, nil
}

// IsMarketActive checks if a market is actively trading
func IsMarketActive(market Market) bool {
	return market.Active &&
		!market.Closed &&
		!market.Archived &&
		market.VolumeNum > 0
}
