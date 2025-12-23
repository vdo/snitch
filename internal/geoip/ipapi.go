package geoip

import (
	"encoding/json"
	"net/http"
	"time"
)

// ipAPIService uses ip-api.com for geolocation (free, no API key required)
type ipAPIService struct {
	client *http.Client
}

type ipAPIResponse struct {
	Status      string `json:"status"`
	CountryCode string `json:"countryCode"`
	Org         string `json:"org"`
}

func (s *ipAPIService) Lookup(ip string) IPInfo {
	if s.client == nil {
		s.client = &http.Client{
			Timeout: 2 * time.Second,
		}
	}

	// ip-api.com free tier: 45 requests per minute
	// We use the batch endpoint fields to minimize response size
	url := "http://ip-api.com/json/" + ip + "?fields=status,countryCode,org"

	resp, err := s.client.Get(url)
	if err != nil {
		return IPInfo{}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return IPInfo{}
	}

	var result ipAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return IPInfo{}
	}

	if result.Status != "success" {
		return IPInfo{}
	}

	return IPInfo{
		CountryCode: result.CountryCode,
		Org:         result.Org,
	}
}
