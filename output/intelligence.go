package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/ghchinoy/ce-go/ce"
	"github.com/olekukonko/tablewriter"
)

// IntelligenceMetadataTable writes out either a tabular or csv view of the metadata
func IntelligenceMetadataTable(metadatabytes []byte, orderBy string, filterBy string, selectonly string, asCsv bool) {
	var intelligence ce.Intelligence
	err := json.Unmarshal(metadatabytes, &intelligence)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	sort.Sort(intelligence)
	if orderBy == "customers" {
		sort.Sort(ce.ByCustomerCount(intelligence))
	} else if orderBy == "hub" {
		sort.Sort(ce.ByIntHub(intelligence))
	} else if orderBy == "name" {
		sort.Sort(ce.ByIntName(intelligence))
	} else if orderBy == "instances" {
		sort.Sort(ce.ByInstanceCount(intelligence))
	} else if orderBy == "traffic" {
		sort.Sort(ce.ByTraffic(intelligence))
	} else if orderBy == "api" {
		sort.Sort(ce.ByAPIType(intelligence))
	} else if orderBy == "authn" {
		sort.Sort(ce.ByAuthn(intelligence))
	}

	if selectonly != "" {
		intelligence = selectOut(selectonly, intelligence)
	}

	data := [][]string{}
	for _, v := range intelligence {
		//configcount := strconv.Itoa(len(v.Configuration))
		data = append(data, []string{
			strconv.Itoa(v.ID),
			v.Key,
			v.Name,
			v.Hub,
			v.API.Type,
			v.AuthenticationType,
			strconv.FormatBool(v.Transformations),
			strconv.FormatBool(v.Active),
			strconv.FormatBool(v.Beta),
			v.ElementClass,
			strconv.Itoa(v.Usage.Traffic),
			strconv.Itoa(v.Usage.CustomerCount),
			strconv.Itoa(v.Usage.InstanceCount),
		})
	}

	if asCsv == true {
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
		table.SetHeader([]string{"ID", "Key", "Name", "Hub", "API", "Authn", "Transforms", "Hidden", "Beta", "Class", "Traffic", "Customers", "Instances"})
		table.SetBorder(false)
		table.AppendBulk(data)
		table.Render()
	}
}

// returns only keys maching filter
func selectOut(filter string, intelligence ce.Intelligence) ce.Intelligence {

	var i ce.Intelligence
	for _, v := range intelligence {
		if strings.ToLower(v.Key) == strings.ToLower(filter) {
			i = append(i, v)
		}
	}

	return i

}
