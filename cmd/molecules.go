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
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/ghchinoy/ce-go/ce"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	profileSource, profileTarget string
	exportCombined               bool
)

// moleculesCmd is the top level command for managing integration assets
var moleculesCmd = &cobra.Command{
	Use:    "molecules",
	Short:  "Manage integration molecules from the Platform",
	Hidden: true,
	Long:   `Manage the integration assets of the Platform`,
}

// exportCmd is the command to export assets
var exportCmd = &cobra.Command{
	Use:   "export [formulas|resources|transformations|all (default)]",
	Short: "exports assets from the platform",
	Long:  "Exports a set of assets",
	Run: func(cmd *cobra.Command, args []string) {

		// check for profile
		profilemap, err := getAuth(profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		scope := []string{"formulas", "resources", "transformations"}
		if len(args) > 0 {
			// args[0] should be either "formulas" | "resources" | "transformations"
			if args[0] == "formulas" {
				scope = []string{"formulas"}
			}
			if exportCombined { // exportCombined combines both resources & transformations
				scope = []string{"resources", "transformations"}
			} else {
				if args[0] == "resources" {
					scope = []string{"resources"}
				}
				if args[0] == "transformations" {
					scope = []string{"transformations"}
				}
			}
		}

		if exportCombined {
			vdr, err := CombineVirtualDataResourcesForExport(profilemap["base"], profilemap["auth"])
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			vdrbytes, err := json.Marshal(vdr)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			//fmt.Printf("%s", vdrbytes)
			name := fmt.Sprintf("%s.combined.vdr.json", strings.Replace(profile, " ", "", -1))
			fmt.Printf("Exporting '%s' to %s/%s\n", "combined vdr", ".", name)
			err = ioutil.WriteFile(fmt.Sprintf("%s/%s", ".", name), vdrbytes, 0644)
		}

		for _, v := range scope {
			if v == "formulas" {
				err = ExportAllFormulasToDir(profilemap["base"], profilemap["auth"], "./formulas")
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
			}
			if !exportCombined {
				if v == "resources" {
					err = ExportAllResourcesToDir(profilemap["base"], profilemap["auth"], "./resources")
					if err != nil {
						fmt.Println(err.Error())
						os.Exit(1)
					}
				}
				if v == "transformations" {
					err = ExportAllTransformationsToDir(profilemap["base"], profilemap["auth"], "./transformations")
					if err != nil {
						fmt.Println(err.Error())
						os.Exit(1)
					}
				}
			}
		}

	},
}

// AllVDR is the combination of Resources and Transformations
type AllVDR struct {
	ObjectDefinitions map[string]ce.CommonResource `json:"objectDefinitions"`
	Transformations   interface{}                  `json:"transformations"`
}

// CombineVirtualDataResourcesForExport creates a single JSON object comprising
// both the Resources and Transformations for the account.
// Two top level keys of the JSON object are: objectDefinitions for Resources
// and transformations for Transformations
// This is very slow.
func CombineVirtualDataResourcesForExport(base, auth string) (AllVDR, error) {
	var vdr AllVDR

	// todo: Gather Resources and Gather Transformations, goroutines, each

	// Gather Resources
	objs := make(map[string]ce.CommonResource)
	resourcesListBytes, _, _, err := ce.ResourcesList(base, auth)
	if err != nil {
		return vdr, err
	}
	var resources []ce.CommonResource
	err = json.Unmarshal(resourcesListBytes, &resources)
	if err != nil {
		return vdr, err
	}
	for _, r := range resources { // todo: goroutine
		//log.Println("exporting", r.Name)
		resourceBytes, _, _, err := ce.GetResourceDefinition(base, auth, r.Name, false)
		if err != nil {
			log.Println(err.Error())
			break
		}
		var obj ce.CommonResource
		err = json.Unmarshal(resourceBytes, &obj)
		objs[r.Name] = obj
	}
	vdr.ObjectDefinitions = objs

	// Gather Transformations
	txs := make(map[string]map[string]interface{})
	// Get all available transformations
	log.Println("Getting all available Transformations")
	bodybytes, status, _, err := ce.GetTransformations(base, auth)
	if err != nil {
		log.Println("Couldn't find any Transformations")
		return vdr, nil
	}
	if status != 200 {
		log.Println("No Transformations present")
		return vdr, nil
	}
	transformationnames := make(map[string]ce.Transformation)
	err = json.Unmarshal(bodybytes, &transformationnames)
	var elementids []int
	namemap := make(map[int]string)
	for k := range transformationnames {
		bodybytes, status, _, err := ce.GetTransformationAssocation(base, auth, k)
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
		for _, v := range associations {
			//fmt.Printf("%s: %s (%v)\n", k, v.Element.Key, v.Element.ID)
			if _, ok := namemap[v.Element.ID]; !ok {
				namemap[v.Element.ID] = v.Element.Key
				elementids = append(elementids, v.Element.ID)
			}
		}
	}
	for _, v := range elementids {
		transforms := make(map[string]interface{})
		idstr := strconv.Itoa(v)
		bodybytes, status, _, err := ce.GetTransformationsPerElement(base, auth, idstr)
		if err != nil {
			break
		}
		if status != 200 {
			break
		}
		err = json.Unmarshal(bodybytes, &transforms)
		txs[namemap[v]] = transforms
	}
	vdr.Transformations = txs

	return vdr, nil
}

// ExportAllTransformationsToDir creates a directory given a dirname and iterates through all
// Elements with associated transformations and creates a single JSON file
func ExportAllTransformationsToDir(base, auth string, dirname string) error {

	err := os.MkdirAll(dirname, os.ModePerm)
	if err != nil {
		return err
	}
	log.Println("Finding all Transformations")

	// Get all available transformations
	log.Println("Getting all available Transformations")
	bodybytes, status, _, err := ce.GetTransformations(base, auth)
	if err != nil {
		log.Println("Couldn't find any Transformations")
		return err
	}
	if status != 200 {
		log.Println("No Transformations present")
		return nil
	}

	log.Println("Assembling unique Element keys")
	transformationnames := make(map[string]ce.Transformation)
	err = json.Unmarshal(bodybytes, &transformationnames)
	var elementids []int
	namemap := make(map[int]string)
	for k := range transformationnames {
		bodybytes, status, _, err := ce.GetTransformationAssocation(base, auth, k)
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
		for _, v := range associations {
			//fmt.Printf("%s: %s\n", k, v.Element.Key)
			if _, ok := namemap[v.Element.ID]; !ok {
				namemap[v.Element.ID] = v.Element.Key
				elementids = append(elementids, v.Element.ID)
			}
		}
	}
	fmt.Printf("%v", elementids)
	log.Println("Exporting Transformations per Element")
	for _, v := range elementids {
		transforms := make(map[string]interface{})
		idstr := strconv.Itoa(v)
		bodybytes, status, _, err := ce.GetTransformationsPerElement(base, auth, idstr)
		if err != nil {
			break
		}
		log.Printf("%s (%s)", namemap[v], idstr)
		//log.Printf("%s\n", bodybytes)
		if status != 200 {
			break
		}
		err = json.Unmarshal(bodybytes, &transforms)
		if err != nil {
			log.Println("unable to umarshal Transformation JSON", err.Error())
			break
		}

		for n, t := range transforms {
			filename := fmt.Sprintf("%s_%s.transformation.json", namemap[v], n)

			b, err := json.Marshal(t)
			if err != nil {
				log.Println("Couldn't convert to bytes", err.Error())
			}
			//log.Printf("%s\n%s\n", n, b)
			log.Printf("Exporting %s", filename)
			err = ioutil.WriteFile(fmt.Sprintf("%s/%s", dirname, filename), b, 0644)
			if err != nil {
				log.Println("Error writing file")
				break
			}
		}
	}

	return nil
}

func interfaceToByte(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ExportAllFormulasToDir creates a directory given and exports all Formula JSON files
func ExportAllFormulasToDir(base, auth string, dirname string) error {
	formulaListByes, _, _, err := ce.FormulasList(base, auth)
	if err != nil {
		return err
	}
	var formulas []ce.Formula
	err = json.Unmarshal(formulaListByes, &formulas)
	if err != nil {
		return err
	}

	// create formulas dir
	err = os.MkdirAll(dirname, os.ModePerm)
	if err != nil {
		return err
	}
	for _, f := range formulas {
		name := fmt.Sprintf("%s.formula.json", strings.Replace(f.Name, " ", "", -1))
		formulaBytes, err := json.Marshal(f)
		if err != nil {
			break
		}
		fmt.Printf("Exporting '%s' to %s/%s\n", f.Name, dirname, name)
		err = ioutil.WriteFile(fmt.Sprintf("%s/%s", dirname, name), formulaBytes, 0644)
	}

	return nil
}

// ExportAllResourcesToDir writes out all the resources to the speceified irectory
func ExportAllResourcesToDir(base, auth string, dirname string) error {
	resourcesListBytes, _, _, err := ce.ResourcesList(base, auth)
	if err != nil {
		return err
	}
	var resources []ce.CommonResource
	err = json.Unmarshal(resourcesListBytes, &resources)
	if err != nil {
		return err
	}
	err = os.MkdirAll(dirname, os.ModePerm)
	if err != nil {
		return err
	}
	for _, r := range resources {
		resourceBytes, _, _, err := ce.GetResourceDefinition(base, auth, r.Name, false)
		if err != nil {
			log.Println(err.Error())
			break
		}
		name := fmt.Sprintf("%s.obj.json", r.Name)
		fmt.Printf("Exporting %s to %s/%s\n", r.Name, dirname, name)
		err = ioutil.WriteFile(fmt.Sprintf("%s/%s", dirname, name), resourceBytes, 0644)
	}

	return nil
}

// cloneCmd is the command to clone assets between accounts
var cloneCmd = &cobra.Command{
	Use:    "clone",
	Short:  "clone assets from one profile to another",
	Long:   "Clone exports assets from one account profile (--from) and imports them into another profile (--to)",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {

		/*
			// check for profile
			profilemap, err := getAuth(profile)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		*/

		if viper.IsSet(profileTarget) && viper.IsSet(profileSource) {
			fmt.Printf("Exporting from profile '%s' into profile '%s'\n", profileSource, profileTarget)
		} else {
			if !viper.IsSet(profileSource) {
				fmt.Println("Cannot find profile named:", profileSource)
			}
			if !viper.IsSet(profileTarget) {
				fmt.Println("Cannot find profile named:", profileTarget)
			}
		}

	},
}

func init() {
	RootCmd.AddCommand(moleculesCmd)

	moleculesCmd.PersistentFlags().StringVar(&profile, "profile", "default", "profile name")

	moleculesCmd.AddCommand(exportCmd)
	exportCmd.PersistentFlags().BoolVar(&exportCombined, "combined", false, "export resources+transformations as one file")
	moleculesCmd.AddCommand(cloneCmd)
	cloneCmd.PersistentFlags().StringVar(&profileSource, "from", "default", "source profile name")
	cloneCmd.PersistentFlags().StringVar(&profileTarget, "to", "", "target profile name")
	cloneCmd.MarkPersistentFlagRequired("to")
}
