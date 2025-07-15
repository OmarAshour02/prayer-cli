package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"prayer-cli/internal"
)

type Config struct {
	City string `json:"city"`
}

func loadCity() (string, error) {
	data, err := os.ReadFile("config.json")
	if err != nil {
		return "", err
	}
	
	var config map[string]string
	if err := json.Unmarshal(data, &config); err != nil {
		return "", err
	}
	
	return config["city"], nil
}

func modifyCity(filename string, newCity string) error {

	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	config.City = newCity

	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(filename, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func fetchAndDisplayPrayerTimes(address string) error {
	
	lat, lon, display, err := internal.GeocodeAddress(address)
	if err != nil {
		return fmt.Errorf("error geocoding: %w", err)
	}
	
	fmt.Println("Location:", display)
	
	err = internal.GetPrayerTimes(lat, lon)
	if err != nil {
		return fmt.Errorf("error fetching prayer times: %w", err)
	}
	
	return nil
}


var city string

var rootCmd = &cobra.Command{
	Use:   "prayers",
	Short: "Get today's prayer times",
	Long: `Prayer Times CLI - Get Islamic prayer times for any city.

			Examples:
			prayers --city "Cairo, Egypt"
			prayers -c "Alexandria"
			prayers  (uses previously configured city, Default: Makkah, Saudi Arabia)`,


	Run: func(cmd *cobra.Command, args []string) {
		var address string
		
		cityFlag, _ := cmd.Flags().GetString("city")
		
		if cityFlag != "" {
			address = cityFlag
			var configFile = "config.json"

			if err := modifyCity(configFile, address); err != nil {
				fmt.Printf("Warning: Failed to save city to config: %v\n", err)
			}
		} else {
			savedCity, err := loadCity()
			if err != nil || savedCity == "" {
				fmt.Println("No city configured and none provided.")
				fmt.Println("Usage:")
				fmt.Println("prayers --city \"Cairo, Egypt\"")
				fmt.Println("prayers -c \"Alexandria\"")
				return
			}
			address = savedCity
		}
		
		// Fetch and display prayer times
		if err := fetchAndDisplayPrayerTimes(address); err != nil {
			fmt.Printf("%v\n", err)
		}
	},}

func init() {
	var err error
	city, err = loadCity()
	if err != nil {
		fmt.Println("Warning: Could not load city from config, using default.")
		city = ""
	}
	rootCmd.Flags().StringVarP(&city, "city", "c", "", "Configure the city for prayer times")


}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
