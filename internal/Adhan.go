package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type AladhanTimings struct {
	Data struct {
		Timings map[string]string `json:"timings"`
		Date    struct {
			Readable string `json:"readable"`
		} `json:"date"`
	} `json:"data"`
}

var prayerArt = map[string][]string{
	"morning": {
		"    \\  |  /",
		"    .--'--.",
		" ---'  ☀  '---",
		"    '-----'",
		"    /  |  \\",
	},
	"night": {
	"    *     *",
	"      *",
	"  *     )",
	"           *",
	"    *",
	"~~~~~~~~~~~~~~~~",
	},
}

func getTimeOfDayArt(timings map[string]string)  []string {
	nextPrayer := getNextPrayerName(timings)	

	if nextPrayer == "Isha" || nextPrayer == "Fajr"{
		return prayerArt["night"]
	}else{
		return prayerArt["morning"]
	}
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
	
	 art := getTimeOfDayArt(result.Data.Timings)
	
	fmt.Println()
	
	prayers := []struct {
		name  string
		time  string
	}{
		{"Fajr", result.Data.Timings["Fajr"]},
		{"Dhuhr", result.Data.Timings["Dhuhr"]},
		{"Asr", result.Data.Timings["Asr"]},
		{"Maghrib", result.Data.Timings["Maghrib"]},
		{"Isha", result.Data.Timings["Isha"]},
	}
	
	for i := 0; i < 5; i++ {
		artLine := ""
		if i < len(art) {
			artLine = art[i]
		}
		
		prayerLine := ""
		if i < len(prayers) {
			prayerLine = fmt.Sprintf("%-8s  %s", prayers[i].name, prayers[i].time)
		}
		
		fmt.Printf("  %-20s    %-33s\n", artLine, prayerLine)
	}
	
	fmt.Printf("  Date: %-52s \n", result.Data.Date.Readable)
	fmt.Println()
	
	highlightNextPrayer(result.Data.Timings)
	
	return nil
}

func getNextPrayerTime(timings map[string]string) (int, int){
	prayers := []string{"Fajr", "Dhuhr", "Asr", "Maghrib", "Isha"}
	now := time.Now()
	
	for _, prayer := range prayers {
		prayerTime, err := time.Parse("15:04", timings[prayer])
		if err != nil {
			continue
		}
		
		prayerTime = time.Date(now.Year(), now.Month(), now.Day(), 
			prayerTime.Hour(), prayerTime.Minute(), 0, 0, now.Location())
		
		if prayerTime.After(now) {
			duration := prayerTime.Sub(now)
			hours := int(duration.Hours())
			minutes := int(duration.Minutes()) % 60
			
			return hours, minutes
		}
	}
	return 0, 0
}

func getNextPrayerName(timings map[string]string) string{
	prayers := []string{"Fajr", "Dhuhr", "Asr", "Maghrib", "Isha"}
	now := time.Now()
	
	for _, prayer := range prayers {
		prayerTime, err := time.Parse("15:04", timings[prayer])
		if err != nil {
			continue
		}
		
		prayerTime = time.Date(now.Year(), now.Month(), now.Day(), 
			prayerTime.Hour(), prayerTime.Minute(), 0, 0, now.Location())
		
		if prayerTime.After(now) {
			return prayer
		}
	}
	return ""

}
func highlightNextPrayer(timings map[string]string) {
	prayer := getNextPrayerName(timings)
	hours, minutes := getNextPrayerTime(timings)
			
	fmt.Printf("Next Prayer: %s in ", prayer)
	if hours > 0 {
		fmt.Printf("%dh %dm\n", hours, minutes)
	} else {
		fmt.Printf("%dm\n", minutes)
	}
}