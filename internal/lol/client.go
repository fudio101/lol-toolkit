package lol

import (
	"fmt"

	"github.com/KnutZuidema/golio"
	"github.com/KnutZuidema/golio/api"
)

// Client wraps the golio client for League of Legends API
type Client struct {
	golio  *golio.Client
	region api.Region
}

// regionMap maps region strings to golio Region constants
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

// NewClient creates a new LoL API client
func NewClient(apiKey string, region string) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	r, ok := regionMap[region]
	if !ok {
		r = api.RegionVietnam // Default to Vietnam
	}

	client := golio.NewClient(apiKey, golio.WithRegion(r))

	return &Client{
		golio:  client,
		region: r,
	}, nil
}

// GetGolio returns the underlying golio client for direct access
func (c *Client) GetGolio() *golio.Client {
	return c.golio
}

// GetRegion returns the current region
func (c *Client) GetRegion() api.Region {
	return c.region
}
