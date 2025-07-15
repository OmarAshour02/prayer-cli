package internal

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


func GetPrayerTimes(lat string, lon string) error {
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

	fmt.Printf("  %-10s %s\n", "Fajr:", result.Data.Timings["Fajr"])
	fmt.Printf("  %-10s %s\n", "Dhuhr:", result.Data.Timings["Dhuhr"])
	fmt.Printf("  %-10s %s\n", "Asr:", result.Data.Timings["Asr"])
	fmt.Printf("  %-10s %s\n", "Maghrib:", result.Data.Timings["Maghrib"])
	fmt.Printf("  %-10s %s\n", "Isha:", result.Data.Timings["Isha"])

	return nil
}

func GetPrayerTimesByCity(city string) error {
	lat, lon, _, err := GeocodeAddress(city)
	if err != nil {
		return fmt.Errorf("error geocoding city: %w", err)
	}

	return GetPrayerTimes(lat, lon)
}