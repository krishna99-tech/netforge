package ip

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// IPInfo represents the geolocation and ISP information for an IP.
type IPInfo struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	ISP         string  `json:"isp"`
	Org         string  `json:"org"`
	AS          string  `json:"as"`
	Query       string  `json:"query"`
}

// GetIPInfo fetches geolocation details for a given IP or domain.
// If target is empty, it fetches details for the caller's public IP.
func GetIPInfo(target string) (*IPInfo, error) {
	url := fmt.Sprintf("http://ip-api.com/json/%s", target)
	
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch IP info: %v", err)
	}
	defer resp.Body.Close()

	var info IPInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if info.Status == "fail" {
		return nil, fmt.Errorf("API error: failed to lookup IP/Domain")
	}

	return &info, nil
}
