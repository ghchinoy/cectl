package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/ghchinoy/ce-go/ce"
	"github.com/spf13/cobra"
)

var accountIDs, orgIDs, customerIDs []int

// metricsCmd represents the resources command
var metricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Obtain Platform metrics",
	Long:  `Retrieve metrics for the Platform.`,
}

// checkGroupOfRequiredFlags checks to see if one of a group of flags is present
func checkGroupOfRequiredFlags() (string, []int, error) {
	message := `must have one of either "--accounts", "--orgs", or "--customers" flag`
	present := false
	var listType string
	var values []int
	var count int
	if accountIDs != nil {
		present = true
		count++
		listType = "accountIds[]"
		values = accountIDs
	}
	if orgIDs != nil {
		present = true
		count++
		listType = "orgIds[]"
		values = orgIDs
	}
	if customerIDs != nil {
		present = true
		count++
		listType = "customerIds[]"
		values = customerIDs
	}
	if present {
		if count > 1 {
			return listType, values, fmt.Errorf("too many flags present: %v - %s", count, message)
		}
		return listType, values, nil
	}
	return listType, values, fmt.Errorf("no required flag present - %s", message)
}

var getAPIMetricsCmd = &cobra.Command{
	Use:   "api",
	Short: "get API metrics on the platform",
	Long: `Get API metrics on the platform. 
This command requires a list of either customers, orgs, or accounts 
(as list of integers, like: --accounts 123,345) for metrics.`,
	Run: func(cmd *cobra.Command, args []string) {
		// check for profile
		profilemap, err := getAuth(profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if debug {
			log.Printf("%10s %v\n", "customers:", customerIDs)
			log.Printf("%10s %v\n", "orgs:", orgIDs)
			log.Printf("%10s %v\n", "accounts:", accountIDs)
		}

		listType, values, err := checkGroupOfRequiredFlags()
		if err != nil {
			fmt.Println(err)
			cmd.Usage()
			os.Exit(1)
		}

		bodybytes, status, curlcmd, err := ce.GetMetrics(profilemap["base"], profilemap["auth"], listType, values, debug)

		if showCurl {
			log.Println(curlcmd)
		}

		if err != nil {
			fmt.Println("Unable to obtain API metrics", status)
			fmt.Printf("%s\n", bodybytes)
			return
		}

		if status != 200 {
			fmt.Print(status)
			if status == 404 {
				fmt.Println("Unable to contact CE API")
				return
			}
			fmt.Println()
		}

		//if outputJSON {
		fmt.Printf("%s\n", bodybytes)
		//return
		//}
	},
}

func init() {
	RootCmd.AddCommand(metricsCmd)

	metricsCmd.PersistentFlags().StringVar(&profile, "profile", "default", "profile name")
	metricsCmd.PersistentFlags().BoolVarP(&outputJSON, "json", "j", false, "output as json")
	metricsCmd.PersistentFlags().BoolVarP(&showCurl, "curl", "c", false, "show curl command")
	metricsCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "output debug logs")

	metricsCmd.PersistentFlags().IntSliceVarP(&customerIDs, "customers", "", nil, "customer ID list")
	metricsCmd.PersistentFlags().IntSliceVarP(&orgIDs, "orgs", "o", nil, "organization ID list")
	metricsCmd.PersistentFlags().IntSliceVarP(&accountIDs, "accounts", "a", nil, "account ID list")

	metricsCmd.AddCommand(getAPIMetricsCmd)
}
