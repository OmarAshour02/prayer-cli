package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)
type NominatimResult struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
	DisplayName string `json:"display_name"`
}


func GeocodeAddress(address string) (lat string, lon string, display string, err error) {
	baseURL := "https://nominatim.openstreetmap.org/search"
	params := url.Values{}
	params.Set("q", address)
	params.Set("format", "json")
	params.Set("limit", "1")

	resp, err := http.Get(fmt.Sprintf("%s?%s", baseURL, params.Encode()))
	if err != nil {
		return "", "", "", err
	}
	defer resp.Body.Close()

	var results []NominatimResult
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &results)

	if len(results) == 0 {
		return "", "", "", fmt.Errorf("no results found")
	}
	return results[0].Lat, results[0].Lon, results[0].DisplayName, nil
}
