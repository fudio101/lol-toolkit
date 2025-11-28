package lol

import (
	"fmt"

	"github.com/KnutZuidema/golio"
	"github.com/KnutZuidema/golio/api"
)

// Default region when not specified.
const DefaultRegion = "vn2"

// Client wraps the golio client with rate limiting.
type Client struct {
	golio       *golio.Client
	region      api.Region
	rateLimiter *RateLimiter
}

// regionMap maps region codes to golio Region constants.
var regionMap = map[string]api.Region{
	"br1":  api.RegionBrasil,
	"eun1": api.RegionEuropeNorthEast,
	"euw1": api.RegionEuropeWest,
	"jp1":  api.RegionJapan,
	"kr":   api.RegionKorea,
	"la1":  api.RegionLatinAmericaNorth,
	"la2":  api.RegionLatinAmericaSouth,
	"na1":  api.RegionNorthAmerica,
	"oc1":  api.RegionOceania,
	"tr1":  api.RegionTurkey,
	"ru":   api.RegionRussia,
	"me1":  api.RegionMiddleEast,
	"sea":  api.RegionSouthEastAsia,
	"sg2":  api.RegionSouthEastAsia,
	"tw2":  api.RegionTaiwan,
	"vn2":  api.RegionVietnam,
}

// NewClient creates a new LoL API client with rate limiting.
func NewClient(apiKey, region string) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("api key is required")
	}

	r := parseRegion(region)
	client := golio.NewClient(apiKey, golio.WithRegion(r))

	return &Client{
		golio:       client,
		region:      r,
		rateLimiter: NewRateLimiter(),
	}, nil
}

// parseRegion converts region string to api.Region.
func parseRegion(region string) api.Region {
	if r, ok := regionMap[region]; ok {
		return r
	}
	return api.RegionVietnam
}

// waitForRateLimit blocks until rate limit allows a request.
func (c *Client) waitForRateLimit() {
	c.rateLimiter.Wait()
}

// GetGolio returns the underlying golio client.
func (c *Client) GetGolio() *golio.Client {
	return c.golio
}

// GetRegion returns the current region.
func (c *Client) GetRegion() api.Region {
	return c.region
}

// GetRateLimitStatus returns current rate limit usage.
func (c *Client) GetRateLimitStatus() (shortUsed, shortLimit, longUsed, longLimit int) {
	return c.rateLimiter.GetStatus()
}
