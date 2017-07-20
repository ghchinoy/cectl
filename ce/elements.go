package ce

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/moul/http2curl"
	"github.com/olekukonko/tablewriter"
)

const (
	ElementsURI                    = "/elements"
	ElementsKeysURI                = "/elements/keys"
	ElementsDocsFormatURI          = "/elements/%s/docs"
	ElementsMetadataFormatURI      = "/elements/%s/metadata"
	ElementFormatURI               = "/elements/%s"
	ElementInstancesFormatURI      = "/elements/%s/instances"
	ElementInstanceFormatURI       = "/elements/%s/instances/%s"
	ElementsOAuthTokenFormatURI    = "/elements/%s/oauth/token"
	ElementsOAuthURLTokenFormatURI = "/elements/%s/oauth/url"
)

// Element represents an Element
type Element struct {
	ID                     int                    `json:"id,omitempty"`
	Name                   string                 `json:"name,omitempty"`
	Key                    string                 `json:"key,omitempty"`
	Description            string                 `json:"description,omitempty"`
	Image                  string                 `json:"image,omitempty"`
	Active                 bool                   `json:"active,omitempty"`
	Deleted                bool                   `json:"deleted,omitempty"`
	OAuth                  bool                   `json:"typeOauth,omitempty"`
	TrialAccount           bool                   `json:"trialAccount,omitempty"`
	Configuration          []ElementConfiguration `json:"configuration,omitempty"`
	TransformationsEnabled bool                   `json:"transformationsEnabled,omitempty"`
	BulkDownloadEnabled    bool                   `json:"bulkDownloadEnabled,omitempty"`
	Cloneable              bool                   `json:"cloneable,omitempty"`
	Extendable             bool                   `json:"extendable,omitempty"`
	Beta                   bool                   `json:"beta,omitempty"`
	Authentication         ElementAuthentication  `json:"authentication,omitempty"`
	Extended               bool                   `json:"extended,omitempty"`
	Hub                    string                 `json:"hub,omitempty"`
	ProtocolType           string                 `json:"protocolType,omitempty"`
	Private                bool                   `json:"private,omitempty"`
}

// ElementConfiguration represents an element's configuration
type ElementConfiguration struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Key             string `json:"key"`
	Description     string `json:"description"`
	DefaultValue    string `json:"defaultValue"`
	ResellerConfig  bool   `json:"resellerConfig"`
	CompanyConfig   bool   `json:"companyConfig"`
	Active          bool   `json:"active"`
	Internal        bool   `json:"internal"`
	GroupControl    bool   `json:"groupControl"`
	DisplayOrder    int    `json:"displayOrder"`
	Type            string `json:"type"`
	HideFromConsole bool   `json:"hideFromConsole"`
	Required        bool   `json:"required"`
}

// ElementAuthentication represents an element's authentication
type ElementAuthentication struct {
	Type string `json:"type,omitempty"`
}

// Elements is a struct container for a list of elements, used in sorting
type Elements []Element

func (elements Elements) Len() int           { return len(elements) }
func (elements Elements) Less(i, j int) bool { return elements[i].ID < elements[j].ID }
func (elements Elements) Swap(i, j int)      { elements[i], elements[j] = elements[j], elements[i] }

// ByHub implements sort.Interface for Elements
type ByHub []Element

func (e ByHub) Len() int           { return len(e) }
func (e ByHub) Less(i, j int) bool { return e[i].Hub < e[j].Hub }
func (e ByHub) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }

// ByName implements sort.Interface for Elements
type ByName []Element

func (e ByName) Len() int           { return len(e) }
func (e ByName) Less(i, j int) bool { return strings.ToLower(e[i].Name) < strings.ToLower(e[j].Name) }
func (e ByName) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }

// GetAllElements returns all Elements as bytes
func GetAllElements(base, auth string) ([]byte, int, string, error) {

	var bodybytes []byte

	url := fmt.Sprintf("%s%s", base, ElementsURI)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Can't construct request", err.Error())
		os.Exit(1)
	}
	req.Header.Add("Authorization", auth)
	req.Header.Add("Accept", "application/json")
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

// OutputElementsTable writes out a tabular view of the elements list
func OutputElementsTable(elementsbytes []byte, orderBy string, filterBy string) {
	var elements Elements
	err := json.Unmarshal(elementsbytes, &elements)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	sort.Sort(elements)
	if orderBy == "name" {
		sort.Sort(ByName(elements))
	} else if orderBy == "hub" {
		sort.Sort(ByHub(elements))
	}
	data := [][]string{}
	for _, v := range elements {
		configcount := strconv.Itoa(len(v.Configuration))
		data = append(data, []string{
			strconv.Itoa(v.ID),
			v.Key,
			v.Name,
			v.Hub,
			configcount,
			strconv.FormatBool(v.Active),
		})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Key", "Name", "Hub", "Configs", "Active"})
	table.SetBorder(false)
	table.AppendBulk(data)
	table.Render()
}
