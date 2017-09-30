package ce

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/moul/http2curl"
	"github.com/olekukonko/tablewriter"
)

const (
	// CommonResourcesURI is the base URI of the hidden API
	// for Common Object Resources; this one provides an array of common objects
	// with the element instance IDs associated, as well as details about
	// the field's heirarchy (org, account, instance)
	CommonResourcesURI = "/common-resources"
	// CommonResourceURI is the base URI for common object resources
	// this is a simple object with keys being the common object names and no
	// details about associated elements or field level hierarchy
	CommonResourceURI = "/organizations/objects/definitions"
	// CommonResourceDefinitionsFormatURI is a string format for the URI of Common Object Resource definition, given a name of a Common Object
	CommonResourceDefinitionsFormatURI = "/organizations/objects/%s/definitions"
	// CommonResourceTransformationsFormatURI is the string format for the URI of an Element's transformation / mapping, given an element key and an object name
	CommonResourceTransformationsFormatURI = "/organizations/elements/%s/transformations/%s"
)

// CommonResource represents a normalized data object (resource)
type CommonResource struct {
	Name               string  `json:"name,omitempty"`
	ElementInstanceIDs []int   `json:"elementInstanceIds,omitempty"`
	Fields             []Field `json:"fields"`
	Level              string  `json:"level,omitempty"`
}

// Field is a set of  a common resource fields
type Field struct {
	Type            string `json:"type"`
	Path            string `json:"path"`
	AssociatedLevel string `json:"organization,omitempty"`
	AssociatedID    int    `json:"associatedId,omitempty"`
}

// ResourcesList retruns a list of common resource objects
func ResourcesList(base, auth string) ([]byte, int, string, error) {
	var bodybytes []byte
	url := fmt.Sprintf("%s%s", base, CommonResourcesURI)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// cant construct request
		return bodybytes, -1, "", err
	}
	req.Header.Add("Authorization", auth)
	req.Header.Add("Accpet", "application/json")
	req.Header.Add("Content-Type", "application/json")
	curlCmd, _ := http2curl.GetCurlCommand(req)
	curl := fmt.Sprintf("%s", curlCmd)
	resp, err := client.Do(req)
	if err != nil {
		return bodybytes, resp.StatusCode, curl, err
	}
	bodybytes, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return bodybytes, resp.StatusCode, curl, nil
}

// OutputResourcesList prints a nicely formatted table to stdout
func OutputResourcesList(resourcesbytes []byte) error {
	data := [][]string{}

	var commonResources []CommonResource
	err := json.Unmarshal(resourcesbytes, &commonResources)
	if err != nil {
		fmt.Printf("Response not a list of Common Resources, %s", err.Error())
		return err
	}

	for _, v := range commonResources {

		var fieldList string
		if len(v.Fields) > 0 {
			var fields []string
			for _, f := range v.Fields {
				fields = append(fields, f.Path)
			}
			fieldList = strings.Join(fields[:], ", ")
			fieldList = " [" + fieldList + "]"
		}

		var instanceList string
		if len(v.ElementInstanceIDs) > 0 {
			var ids []string
			for _, i := range v.ElementInstanceIDs {
				ids = append(ids, strconv.Itoa(i))
			}
			instanceList = strings.Join(ids[:], ", ")
			instanceList = " [" + instanceList + "]"
		}

		data = append(data, []string{
			v.Name,
			strconv.Itoa(len(v.ElementInstanceIDs)) + instanceList,
			strconv.Itoa(len(v.Fields)),
			fieldList,
		})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Mapped Instances", "#", "Fields"})
	table.SetBorder(false)
	table.SetColWidth(40)
	table.AppendBulk(data)
	table.Render()

	return nil
}
