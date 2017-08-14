package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/ghchinoy/cectl/ce"
	"github.com/spf13/cobra"
)

var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "Manage users on the platform",
	Long:  "Allows management of users on the platform",
}

var listUsersCmd = &cobra.Command{
	Use:   "list",
	Short: "list the users on the platform",
	Long:  "List users on the platform",
	Run: func(cmd *cobra.Command, args []string) {
		// check for profile
		profilemap, err := getAuth(profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		bodybytes, status, curlcmd, err := ce.GetAllUsers(profilemap["base"], profilemap["auth"])
		if err != nil {
			fmt.Println("Unable to even", status)
			return
		}

		if showCurl {
			log.Println(curlcmd)
		}

		if outputJSON {
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

		err = ce.FormatUserList(bodybytes)
		if err != nil {
			fmt.Println("Unable to format")
		}
	},
}

func init() {
	RootCmd.AddCommand(usersCmd)

	usersCmd.AddCommand(listUsersCmd)

	usersCmd.PersistentFlags().StringVar(&profile, "profile", "default", "profile name")
	usersCmd.PersistentFlags().BoolVarP(&outputJSON, "json", "j", false, "output as json")
	usersCmd.PersistentFlags().BoolVarP(&showCurl, "curl", "c", false, "show curl command")

}
