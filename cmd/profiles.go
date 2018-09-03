// Copyright © 2017 G. Hussain Chinoy <ghchinoy@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/olekukonko/tablewriter"
	toml "github.com/pelletier/go-toml"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ghchinoy/cectl/tokens"
)

// profilesCmd represents the profile command
var profilesCmd = &cobra.Command{
	Use:   "profiles",
	Short: "Manage profiles",
	Long:  `Add, remove, list profiles to manage Cloud Elements access`,
}

var useLoginFlow bool

// addProfileCmd represents the addProfile command
var addProfileCmd = &cobra.Command{
	Use:   "add <profile>",
	Short: "add a new profile",
	Long: `Adds a new profile to the available profiles. Provide a name to get started.
Use the flag --login or -l to log in to CE and create profile.`,
	Run: func(cmd *cobra.Command, args []string) {
		// check for args, if arg, then list the details for that particular profile
		if len(args) > 0 {
			profile = args[0]
		} else {
			fmt.Println("please provide a profile name to add, profile add <name>")
			os.Exit(1)
		}
		fmt.Printf("%7s: %s\n", "profile", profile)
		if viper.IsSet(profile) {
			fmt.Printf("Profile %s exists.\n", profile)
			os.Exit(1)
		}

		var base, org, user string

		// Login flow
		if useLoginFlow {
			var err error
			base, org, user, err = tokens.LoginInquiry()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else { // manual entry flow
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("base URI: ")
			base, _ = reader.ReadString('\n')
			base = strings.Replace(base, "\n", "", -1)
			fmt.Print("user token: ")
			user, _ = reader.ReadString('\n')
			user = strings.Replace(user, "\n", "", -1)
			fmt.Print("org token: ")
			org, _ = reader.ReadString('\n')
			org = strings.Replace(org, "\n", "", -1)
		}
		viper.Set(profile+".base", base)
		viper.Set(profile+".org", org)
		viper.Set(profile+".user", user)

		err := writeConfigAs(viper.ConfigFileUsed(), true)
		if err != nil {
			fmt.Println("Unable to write config file", err.Error())
			fmt.Printf("Config file %s unchanged.\n", viper.ConfigFileUsed())
		}
		fmt.Printf("Added profile %s\n", profile)
	},
}

var longProfile bool

// listProfilesCmd represents the listProfiles command
var listProfilesCmd = &cobra.Command{
	Use:   "list",
	Short: "lists available profiles",
	Long:  `Lists available profiles`,
	Run: func(cmd *cobra.Command, args []string) {
		// check for args, if arg, then list the details for that particular profile
		if len(args) > 0 {
			profile = args[0]
		}

		if !(outputCSV) {
			fmt.Printf("%7s: %s", "profile", profile)
			if viper.IsSet(profile) {
				p := viper.GetStringMap(profile)

				if val, ok := p["label"]; ok {
					fmt.Printf(" (%s)\n", val)
				} else {
					fmt.Println()
				}

				for k, v := range p {
					if k == "base" {
						fmt.Printf("%7s: %s\n", k, v)
					}
				}
			} else {
				fmt.Printf("No %s profile exists in config file %s.", profile, cfgFile)
			}
			fmt.Println()
		}

		settings := viper.AllSettings()
		var profiles []string
		for k := range settings {
			profiles = append(profiles, k)
		}
		sort.Strings(profiles)
		if longProfile {
			data := [][]string{}
			if !(outputCSV) {
				fmt.Printf("%v profiles\n", len(profiles))
			}
			for _, k := range profiles {
				if k == "profile" {
					continue
				}
				v := settings[k].(map[string]interface{})
				if v["base"] == nil {
					continue
				}
				data = append(data, []string{
					k,
					fmt.Sprintf("%s", v["base"]),
				})
			}
			if outputCSV {
				w := csv.NewWriter(os.Stdout)
				for _, record := range data {
					if err := w.Write(record); err != nil {
						log.Fatalln("error writing record to csv:", err)
					}
				}
				w.Flush()
				if err := w.Error(); err != nil {
					log.Fatal(err)
				}
			} else {
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"profile", "base url"})
				table.SetBorder(false)
				table.SetAlignment(tablewriter.ALIGN_LEFT)
				table.AppendBulk(data)
				table.Render()
			}
			os.Exit(0)
		}
		posn := sort.SearchStrings(profiles, "profile")
		profiles = append(profiles[:posn], profiles[posn+1:]...)
		fmt.Println("Valid profiles:", strings.Join(profiles, ", "))
	},
}

// setProfileCmd represents the setProfile command
var setProfileCmd = &cobra.Command{
	Use:   "set <profile>",
	Short: "sets a profile to be the default profile",
	Long:  `Sets given profile name as the default profile`,
	Run: func(cmd *cobra.Command, args []string) {

		// check for args, if arg, then list the details for that particular profile
		if len(args) > 0 {
			profile = args[0]
		}
		fmt.Printf("%7s: %s\n", "profile", profile)
		// check if specified profile exists
		if viper.IsSet(profile) {
			p := viper.GetStringMap(profile)

			// copy profile data to default, update label
			// indicate change has occurred

			viper.Set("default.base", p["base"])
			viper.Set("default.org", p["org"])
			viper.Set("default.user", p["user"])
			viper.Set("default.label", profile)

			// writing back has a PR to make this more formal: https://github.com/spf13/viper/pull/287
			// Once that PR is merged, replace writeConfigFile with
			/*
				if err := viper.WriteConfigAs(viper.ConfigFileUsed()); err != nil {
					fmt.Fatal(err)
				}
			*/
			err := writeConfigAs(viper.ConfigFileUsed(), true)
			if err != nil {
				fmt.Println("Unable to write config file", err.Error())
				fmt.Printf("Config file %s unchanged.\n", viper.ConfigFileUsed())
			}
			fmt.Printf("Default profile set to %s\n", profile)

			for k, v := range p {
				if k == "base" {
					fmt.Printf("%7s: %s\n", k, v)
				}
				if k == "label" {
					fmt.Printf("%7s: %s\n", k, v)
				}
			}
		} else { // if not, end
			fmt.Printf("No %s profile exists in config file %s.\n", profile, cfgFile)
			fmt.Printf("Cannot set %s as default profile.\n", profile)
			fmt.Println()
			settings := viper.AllSettings()
			var profiles []string
			for k := range settings {
				profiles = append(profiles, k)
			}
			sort.Strings(profiles)
			posn := sort.SearchStrings(profiles, "profile")
			profiles = append(profiles[:posn], profiles[posn+1:]...)
			fmt.Println("Valid profiles:", strings.Join(profiles, ", "))
		}
	},
}

func writeConfigAs(filename string, force bool) error {

	t, err := toml.TreeFromMap(viper.AllSettings())
	if err != nil {
		return err
	}
	s := t.String()

	var flags int
	if force == true {
		flags = os.O_CREATE | os.O_TRUNC | os.O_WRONLY
	} else {
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			flags = os.O_WRONLY
		} else {
			return fmt.Errorf("file: %s exists - cannot overwrite, use force option", filename)
		}
	}

	var AppFs = afero.NewOsFs()
	f, err := AppFs.OpenFile(filename, flags, os.FileMode(0644))
	if err != nil {
		return err
	}

	_, err = f.WriteString(s)
	if err != nil { // issue writing string to file
		return err
	}

	return nil
}

var profilesEnvCmd = &cobra.Command{
	Use:   "env",
	Short: "Output bash env vars for profile",
	Long: `Outputs bash environment variables for current profile
Do the following to add env variables
source <(cectl profiles env)`,
	Run: func(cmd *cobra.Command, args []string) {
		// check for profile
		profilemap, err := getAuth(profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("export CE_AUTH=\"%s\"\n", profilemap["auth"])
		fmt.Printf("export CE_BASE=%s\n", profilemap["base"])
		fmt.Printf("export CE_ORG=\"%s\"\n", profilemap["org"])
		fmt.Printf("export CE_USER=\"%s\"\n", profilemap["user"])

	},
}

func init() {
	RootCmd.AddCommand(profilesCmd)

	profilesCmd.PersistentFlags().StringVar(&cfgFile, "config", "", cfgHelp)
	profilesCmd.PersistentFlags().StringVar(&profile, "profile", "default", "profile name")
	// Set bash-completion
	validConfigFilenames := []string{"toml", ""}
	profilesCmd.PersistentFlags().SetAnnotation("config", cobra.BashCompFilenameExt, validConfigFilenames)

	profilesCmd.AddCommand(listProfilesCmd)
	listProfilesCmd.PersistentFlags().BoolVarP(&longProfile, "long", "l", false, "show long profile")
	listProfilesCmd.PersistentFlags().BoolVarP(&outputCSV, "csv", "", false, "output as CSV")

	profilesCmd.AddCommand(addProfileCmd)
	addProfileCmd.PersistentFlags().BoolVarP(&useLoginFlow, "login", "l", false, "prompt for login")

	profilesCmd.AddCommand(setProfileCmd)
	profilesCmd.AddCommand(profilesEnvCmd)

}
