package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/ghchinoy/ce-go/ce"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// elementsCmd represents the elements command
var transformationsCmd = &cobra.Command{
	Use:   "transformations",
	Short: "Manage Transformations on the Platform",
	Long:  `Manage Transformations on the Platform`,
}

var withElementAssociations bool

var listTransformationsCmd = &cobra.Command{
	Use:   "list",
	Short: "List Transformations",
	Long:  "List Transformations on the Platform",
	Run: func(cmd *cobra.Command, args []string) {
		// check for profile
		profilemap, err := getAuth(profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		bodybytes, statuscode, curlcmd, err := ce.GetTransformations(profilemap["base"], profilemap["auth"])
		if err != nil {
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
		}
		if outputJSON {
			// todo uplift to output package, output/FormattedJSON
			var transformations interface{}
			err = json.Unmarshal(bodybytes, &transformations)
			if err != nil {
				fmt.Println("Can't unmarshal")
				os.Exit(1)
			}
			formattedbytes, err := json.MarshalIndent(transformations, "", "    ")
			if err != nil {
				fmt.Println("Can't format json")
				os.Exit(1)
			}
			fmt.Printf("%s", formattedbytes)
			os.Exit(0)
		}

		txs := make(map[string]ce.Transformation)
		err = json.Unmarshal(bodybytes, &txs)
		if err != nil {
			fmt.Println("Unable to parse Transformations", err.Error())
			os.Exit(1)
		}
		// sort by key
		var keys []string
		for k := range txs {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		elementAssociations := make(map[string][]string)
		if withElementAssociations {
			for _, k := range keys {
				bodybytes, status, _, err := ce.GetTransformationAssocation(profilemap["base"], profilemap["auth"], k)
				if err != nil {
					break
				}
				if status != 200 {
					break
				}
				var associations []ce.AccountElement
				err = json.Unmarshal(bodybytes, &associations)
				if err != nil {
					break
				}
				var elements []string
				for _, e := range associations {
					elements = append(elements, e.Element.Key)
				}
				elementAssociations[k] = elements
			}
		}

		data := [][]string{}
		for _, k := range keys {
			v := txs[k]
			if withElementAssociations {
				data = append(data, []string{
					k,
					v.Level,
					fmt.Sprintf("%v", len(v.Fields)),
					fmt.Sprintf("%v %s", len(elementAssociations[k]), elementAssociations[k]),
				})
			} else {
				data = append(data, []string{
					k,
					v.Level,
					fmt.Sprintf("%v", len(v.Fields)),
				})
			}
		}
		table := tablewriter.NewWriter(os.Stdout)
		//table.SetHeader([]string{"Resource", "Vendor", "Level", "# Fields", "# Configs", "Legacy", "Start Date"})
		if withElementAssociations {
			table.SetHeader([]string{"Resource", "Level", "# Fields", "Elements"})
		} else {
			table.SetHeader([]string{"Resource", "Level", "# Fields"})
		}
		table.SetBorder(false)
		table.AppendBulk(data)
		table.Render()
	},
}

func init() {
	RootCmd.AddCommand(transformationsCmd)

	transformationsCmd.PersistentFlags().StringVar(&profile, "profile", "default", "profile name")
	transformationsCmd.PersistentFlags().BoolVarP(&outputJSON, "json", "j", false, "output as json")
	transformationsCmd.PersistentFlags().BoolVarP(&showCurl, "curl", "c", false, "show curl command")
	//transformationsCmd.PersistentFlags().BoolVarP(&outputCSV, "csv", "", false, "output as CSV")
	transformationsCmd.AddCommand(listTransformationsCmd)
	listTransformationsCmd.PersistentFlags().BoolVarP(&withElementAssociations, "with-elements", "", false, "show Element associations")
}
