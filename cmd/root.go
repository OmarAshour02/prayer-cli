package cmd

import (
	"fmt"
	"os"
	"prayer-cli/config"
	"prayer-cli/internal"

	"github.com/spf13/cobra"
)


func loadCity() (string, error) {

	cfg, _:= config.LoadOrInitConfig()
	
	return cfg.City, nil
}

func loadMethod() (string, error) {

	cfg, _:= config.LoadOrInitConfig()

	return cfg.Method, nil
}

func modifyCity(newCity string) error {


	err := config.UpdateConfig(func(c *config.Config) {
		c.City = newCity
	})
	
	if err != nil{
		fmt.Println("City changed successfully to", newCity)
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
			if err := modifyCity(address); err != nil {
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
		
		if err := fetchAndDisplayPrayerTimes(address); err != nil {
			fmt.Printf("%v\n", err)
		}
	},
}

func init() {
	var err error
	city, err = loadCity()
	if err != nil {
		fmt.Println("Warning: Could not load city from config, using default.")
		city = ""
	}

	method, err := loadMethod()

	if err != nil{
		fmt.Println("Could not load method")
	}
	rootCmd.Flags().StringVarP(&city, "city", "c", "", "Configure the city for prayer times")
	rootCmd.Flags().StringVarP(&method, "method", "m", "5", `Calculation method number  Possible values:

    0 - Jafari / Shia Ithna-Ashari
    1 - University of Islamic Sciences, Karachi
    2 - Islamic Society of North America
    3 - Muslim World League
    4 - Umm Al-Qura University, Makkah
    5 - Egyptian General Authority of Survey
    7 - Institute of Geophysics, University of Tehran
    8 - Gulf Region
    9 - Kuwait
    10 - Qatar
    11 - Majlis Ugama Islam Singapura, Singapore
    12 - Union Organization islamic de France
    13 - Diyanet İşleri Başkanlığı, Turkey
    14 - Spiritual Administration of Muslims of Russia
    15 - Moonsighting Committee Worldwide (also requires shafaq parameter)
    16 - Dubai (experimental)
    17 - Jabatan Kemajuan Islam Malaysia (JAKIM)
    18 - Tunisia
    19 - Algeria
    20 - KEMENAG - Kementerian Agama Republik Indonesia
    21 - Morocco
    22 - Comunidade Islamica de Lisboa
    23 - Ministry of Awqaf, Islamic Affairs and Holy Places, Jordan `)

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
