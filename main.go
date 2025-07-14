package main

import (
	"fmt"
	"os"
	"github.com/omarashour/prayer-cli/internal/location"
)



func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: prayer-cli <address>")
		return
	}

	address := os.Args[1]
	lat, lon, display, err := geocodeAddress(address)
	if err != nil {
		fmt.Println("❌ Error geocoding:", err)
		return
	}

	fmt.Println("📍 Location:", display)
	err = getPrayerTimes(lat, lon)
	if err != nil {
		fmt.Println("❌ Error fetching prayer times:", err)
	}
}
