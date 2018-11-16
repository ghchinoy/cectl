package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/ghchinoy/ce-go/ce"
	"github.com/ghchinoy/cectl/output"
	"github.com/spf13/cobra"
)

// intelligenceCmd represents the elements command
var intelligenceCmd = &cobra.Command{
	Use:    "intelligence",
	Short:  "Metadata about Elements on the Platform",
	Long:   `Metadata about Elements on the Platform`,
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		// check for profile
		profilemap, err := getAuth(profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// Get metadata
		bodybytes, statuscode, curlcmd, err := ce.GetIntelligence(profilemap["base"], profilemap["auth"])
		if err != nil {
			if statuscode == -1 {
				fmt.Println("Unable to reach CE API. Please check your configuration / profile.")
			}
			fmt.Println(err)
			os.Exit(1)
		}
		// handle global options, curl
		if showCurl {
			log.Println(curlcmd)
		}

		// handle non 200
		if statuscode != 200 {
			log.Printf("HTTP Error: %v\n", statuscode)
			// handle this nicely, show error description
			os.Exit(1)
		}
		// handle global options, json
		if outputJSON {
			fmt.Printf("%s\n", bodybytes)
			return
		}
		// output
		var grepstring string
		if len(args) > 1 {
			grepstring = args[1]
		}
		output.IntelligenceMetadataTable(bodybytes, orderBy, filterBy, grepstring, outputCSV)
		fmt.Println()
	},
}

func init() {

	RootCmd.AddCommand(intelligenceCmd)

	intelligenceCmd.PersistentFlags().StringVar(&profile, "profile", "default", "profile name")
	intelligenceCmd.PersistentFlags().BoolVarP(&outputJSON, "json", "j", false, "output as json")
	intelligenceCmd.PersistentFlags().BoolVarP(&showCurl, "curl", "c", false, "show curl command")
	intelligenceCmd.PersistentFlags().BoolVarP(&outputCSV, "csv", "", false, "output as CSV")

	intelligenceCmd.Flags().StringVarP(&orderBy, "order", "", "", "order metadata (customers, instances, traffic, hub, name, api, authn)")
}
