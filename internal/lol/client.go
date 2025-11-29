package lol

import (
	"encoding/json"
	"fmt"

	"github.com/KnutZuidema/golio"
	"github.com/KnutZuidema/golio/api"
)

// Default region when not specified.
const DefaultRegion = "vn2"

// Client wraps the golio client.
type Client struct {
	golio  *golio.Client
	region api.Region
	apiKey string // stored for logging headers
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

// NewClient creates a new LoL API client.
func NewClient(apiKey, region string) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("api key is required")
	}

	r := parseRegion(region)
	client := golio.NewClient(apiKey, golio.WithRegion(r))

	return &Client{
		golio:  client,
		region: r,
		apiKey: apiKey,
	}, nil
}

// parseRegion converts region string to api.Region.
func parseRegion(region string) api.Region {
	if r, ok := regionMap[region]; ok {
		return r
	}
	return api.RegionVietnam
}

// GetGolio returns the underlying golio client.
func (c *Client) GetGolio() *golio.Client {
	return c.golio
}

// GetRegion returns the current region.
func (c *Client) GetRegion() api.Region {
	return c.region
}

// getHeaders returns the standard headers for Riot API requests.
func (c *Client) getHeaders() map[string]string {
	headers := make(map[string]string)
	if c.apiKey != "" {
		headers["X-Riot-Token"] = c.apiKey
	}
	headers["Accept"] = "application/json"
	return headers
}

// marshalResponse converts a response object to JSON string for logging.
func (c *Client) marshalResponse(data interface{}) string {
	if data == nil {
		return ""
	}
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Sprintf("(failed to marshal: %v)", err)
	}
	return string(jsonData)
}
