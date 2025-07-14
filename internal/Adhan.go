package location

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)
type AladhanTimings struct {
	Data struct {
		Timings map[string]string `json:"timings"`
		Date    struct {
			Readable string `json:"readable"`
		} `json:"date"`
	} `json:"data"`
}


func getPrayerTimes(lat string, lon string) error {
	apiURL := fmt.Sprintf("https://api.aladhan.com/v1/timings?latitude=%s&longitude=%s&method=5", lat, lon)
	resp, err := http.Get(apiURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result AladhanTimings
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	fmt.Println("📅 Date:", result.Data.Date.Readable)
	fmt.Println("🕌 Prayer Times:")
	for name, time := range result.Data.Timings {
		fmt.Printf("  %-10s %s\n", name+":", time)
	}

	return nil
}
