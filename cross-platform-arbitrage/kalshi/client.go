package kalshi

import (
	"context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	// BaseURL is the Kalshi API base URL
	BaseURL = "https://api.elections.kalshi.com/trade-api/v2"
	
	// DefaultTimeout is the default timeout for API requests
	DefaultTimeout = 30 * time.Second
)

// Client represents a Kalshi API client
type Client struct {
	client *resty.Client
	baseURL string
}

// NewClient creates a new Kalshi API client
func NewClient() *Client {
	client := resty.New().
		SetBaseURL(BaseURL).
		SetTimeout(DefaultTimeout).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json")
	
	return &Client{
		client:  client,
		baseURL: BaseURL,
	}
}

// Market represents a Kalshi market
type Market struct {
	Ticker               string    `json:"ticker"`
	EventTicker          string    `json:"event_ticker"`
	MarketType           string    `json:"market_type"` // Always "binary" for Kalshi
	Title                string    `json:"title"`
	Subtitle             string    `json:"subtitle"`
	YesBid               int       `json:"yes_bid"`               // Price in cents (e.g., 56 = $0.56)
	YesBidDollars        string    `json:"yes_bid_dollars"`       // Formatted dollars (e.g., "0.5600")
	YesAsk               int       `json:"yes_ask"`               // Price in cents
	YesAskDollars        string    `json:"yes_ask_dollars"`       // Formatted dollars
	NoBid                int       `json:"no_bid"`                // Price in cents
	NoBidDollars         string    `json:"no_bid_dollars"`        // Formatted dollars
	NoAsk                int       `json:"no_ask"`                // Price in cents
	NoAskDollars         string    `json:"no_ask_dollars"`        // Formatted dollars
	LastPrice            int       `json:"last_price"`            // Last trade price in cents
	LastPriceDollars     string    `json:"last_price_dollars"`    // Formatted dollars
	Volume               int       `json:"volume"`                // Trading volume
	Volume24h            int       `json:"volume_24h"`            // 24-hour volume
	Volume24hFP          string    `json:"volume_24h_fp"`         // 24-hour volume formatted
	Status               string    `json:"status"`                // "open", "closed", "active", etc.
	OpenTime             time.Time `json:"open_time"`
	CloseTime            time.Time `json:"close_time"`
	ExpirationTime       time.Time `json:"expiration_time"`
	CreatedTime          time.Time `json:"created_time"`
	UpdatedTime          time.Time `json:"updated_time"`
	CanCloseEarly        bool      `json:"can_close_early"`
	FractionalTradingEnabled bool  `json:"fractional_trading_enabled"`
	OpenInterest         int       `json:"open_interest"`
	NotionalValue        int       `json:"notional_value"`
	NotionalValueDollars string    `json:"notional_value_dollars"`
	Liquidity            int       `json:"liquidity"`
	LiquidityDollars     string    `json:"liquidity_dollars"`
	YesSubTitle          string    `json:"yes_sub_title"`
	NoSubTitle           string    `json:"no_sub_title"`
}

// MarketsResponse represents the response from GetMarkets
type MarketsResponse struct {
	Markets []Market `json:"markets"`
	Cursor  string   `json:"cursor"`
}

// GetMarketsOptions represents options for GetMarkets
type GetMarketsOptions struct {
	Limit      int    // Maximum number of markets to return (default: 100)
	Cursor     string // Pagination cursor
	Status     string // Filter by status: "open", "closed", "all"
	SeriesTicker string // Filter by series ticker
}

// GetMarkets retrieves markets from Kalshi API
func (c *Client) GetMarkets(ctx context.Context, opts *GetMarketsOptions) (*MarketsResponse, error) {
	if opts == nil {
		opts = &GetMarketsOptions{}
	}
	
	req := c.client.R().
		SetContext(ctx).
		SetResult(&MarketsResponse{})
	
	// Add query parameters
	if opts.Limit > 0 {
		req.SetQueryParam("limit", fmt.Sprintf("%d", opts.Limit))
	}
	if opts.Cursor != "" {
		req.SetQueryParam("cursor", opts.Cursor)
	}
	if opts.Status != "" {
		req.SetQueryParam("status", opts.Status)
	}
	if opts.SeriesTicker != "" {
		req.SetQueryParam("series_ticker", opts.SeriesTicker)
	}
	
	resp, err := req.Get("/markets")
	if err != nil {
		return nil, fmt.Errorf("failed to get markets: %w", err)
	}
	
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode(), resp.String())
	}
	
	result, ok := resp.Result().(*MarketsResponse)
	if !ok {
		return nil, fmt.Errorf("failed to parse response")
	}
	
	return result, nil
}

// GetMarketByTicker retrieves a specific market by ticker
func (c *Client) GetMarketByTicker(ctx context.Context, ticker string) (*Market, error) {
	// First get all markets and filter by ticker
	opts := &GetMarketsOptions{
		Limit: 1000, // Get enough to find our market
	}
	
	resp, err := c.GetMarkets(ctx, opts)
	if err != nil {
		return nil, err
	}
	
	// Search for the market with matching ticker
	for _, market := range resp.Markets {
		if market.Ticker == ticker {
			return &market, nil
		}
	}
	
	return nil, fmt.Errorf("market not found: %s", ticker)
}

// CalculateImpliedNoPrice calculates the implied NO price from YES price
// In binary markets: NO price = 100 - YES price (in cents)
func CalculateImpliedNoPrice(yesPriceCents int) int {
	return 100 - yesPriceCents
}

// CalculateArbitrageOpportunity checks for price discrepancies
// Returns arbitrage percentage if opportunity exists
func CalculateArbitrageOpportunity(yesBid int, yesAsk int, noBid int, noAsk int) (float64, bool) {
	// Filter out inactive markets (YES 0/0, NO 100/100)
	if yesBid == 0 && yesAsk == 0 && noBid == 100 && noAsk == 100 {
		return 0.0, false
	}
	
	// Check for arbitrage between YES bid and NO ask
	// If YES bid + NO ask < 100, arbitrage exists
	yesBidNoAskSum := yesBid + noAsk
	if yesBidNoAskSum < 100 {
		arbPercent := float64(100-yesBidNoAskSum) / 100.0 * 100
		return arbPercent, true
	}
	
	// Check for arbitrage between NO bid and YES ask
	// If NO bid + YES ask < 100, arbitrage exists
	noBidYesAskSum := noBid + yesAsk
	if noBidYesAskSum < 100 {
		arbPercent := float64(100-noBidYesAskSum) / 100.0 * 100
		return arbPercent, true
	}
	
	return 0.0, false
}

// IsMarketActive checks if a market has active trading
func IsMarketActive(market Market) bool {
	// Check if market has bid/ask prices
	if market.YesBid == 0 && market.YesAsk == 0 && market.NoBid == 100 && market.NoAsk == 100 {
		return false
	}
	
	// Check if market has volume
	if market.Volume == 0 && market.Volume24h == 0 {
		return false
	}
	
	// Check if market is open/active status
	if market.Status != "open" && market.Status != "active" {
		return false
	}
	
	return true
}

// MarketSummary provides a human-readable summary of a market
func (m *Market) Summary() string {
	return fmt.Sprintf(
		"Market: %s\nTitle: %s\nStatus: %s\n"+
			"YES Bid/Ask: %d/%d ($%s/$%s)\n"+
			"NO Bid/Ask: %d/%d ($%s/$%s)\n"+
			"Last Price: %d ($%s)\n"+
			"Volume: %d (24h: %d)\n"+
			"Expires: %s",
		m.Ticker, m.Title, m.Status,
		m.YesBid, m.YesAsk, m.YesBidDollars, m.YesAskDollars,
		m.NoBid, m.NoAsk, m.NoBidDollars, m.NoAskDollars,
		m.LastPrice, m.LastPriceDollars,
		m.Volume, m.Volume24h,
		m.ExpirationTime.Format("2006-01-02 15:04:05"),
	)
}