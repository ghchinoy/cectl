package ce

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/moul/http2curl"
	"github.com/olekukonko/tablewriter"
)

const (
	InstancesURI                         = "/instances"
	InstancesFormatURI                   = "/instances/%s"
	InstanceConfigurationURI             = "/instances/configuration"
	InstanceConfigurationFormatURI       = "/instances/configuration/%s"
	InstanceDocsURI                      = "/instances/docs"
	InstanceOperationDocsFormatURI       = "/instances/docs/%s"
	InstancesEventsURI                   = "/instances/events"
	InstancesEventsAnalyticsAccountsURI  = "/instances/events/analytics/accounts"
	InstancesEventsAnalyticsInstancesURI = "/instances/events/analytics/instances"
	InstancesEventsFormatURI             = "/instances/events/%s"
	InstancesTransformationsURI          = "/instances/transformations"
	InstanceTransformationsFormatURI     = "/instances/%s/transformations"
	InstanceDocFormatURI                 = "/instances/%s/docs"

	InstanceDefinitions_ID       = "/instances/%s/objects/definitions"
	InstanceDefinitions_Token    = "/instances/objects/definitions"
	InstanceOperationOAI_ID      = "/instances/%s/docs/%s"
	InstanceOAIByOperation_ID    = "/instances/%s/docs/%s/definitions"
	InstanceOAIByOperation_Token = "/instances/docs/%s/definitions"
)

// Instance represents an Element Instance
type Instance struct {
	ID                     int
	Name                   string
	CreatedDate            string
	Token                  string
	Element                Element
	ElementID              int
	Tags                   []string
	ProvisionInteractions  interface{}
	Valid                  bool
	Disabled               bool
	MaxCacheSize           int
	CacheTimeToLive        int
	Configuration          InstanceConfiguration
	EventsEnabled          bool
	TraceLoggingEnabled    bool
	CachingEnabled         bool
	ExternalAuthentication string
	User                   User
}

// GetAllInstances returns the Element Instances for the authed user
func GetAllInstances(base, auth string) ([]byte, int, string, error) {
	var bodybytes []byte

	url := fmt.Sprintf("%s%s", base, InstancesURI)

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
		// unable to reach CE API
		return bodybytes, -1, curl, err
	}
	bodybytes, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	return bodybytes, resp.StatusCode, curl, nil
}

// GetInstanceInfo obtains details of an Instance
func GetInstanceInfo(base, auth, instanceID string) ([]byte, int, string, error) {
	var bodybytes []byte

	url := fmt.Sprintf("%s%s",
		base,
		fmt.Sprintf(InstancesFormatURI, instanceID),
	)

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
		// unable to reach CE API
		return bodybytes, -1, curl, err
	}
	bodybytes, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	return bodybytes, resp.StatusCode, curl, nil
}

// GetInstanceOAI returns the OAI Spec for an Instance ID
func GetInstanceOAI(base, auth, instanceID string) ([]byte, int, string, error) {
	var bodybytes []byte
	url := fmt.Sprintf("%s%s",
		base,
		fmt.Sprintf(InstanceDocFormatURI, instanceID),
	)
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
		return bodybytes, -1, curl, err
	}
	bodybytes, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return bodybytes, resp.StatusCode, curl, nil
}

// GetInstanceTransformations is incomplete
func GetInstanceTransformations(id string) {}

// GetInstanceObjectDefinitions returns the schema definitions for an Instance
func GetInstanceObjectDefinitions(base, auth, instanceID string) ([]byte, int, string, error) {
	var bodybytes []byte

	url := fmt.Sprintf("%s%s",
		base,
		fmt.Sprintf(InstanceDefinitions_ID, instanceID),
	)

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
		// unable to reach CE API
		return bodybytes, -1, curl, err
	}
	bodybytes, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	return bodybytes, resp.StatusCode, curl, nil
}

// GetInstanceOperationDefinition returns the bytes of a call to get Instance schema definitions
func GetInstanceOperationDefinition(base, auth, instanceID, operationName string) ([]byte, int, string, error) {
	var bodybytes []byte

	url := fmt.Sprintf("%s%s",
		base,
		fmt.Sprintf(InstanceOAIByOperation_ID, instanceID, operationName),
	)

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
		// unable to reach CE API
		return bodybytes, -1, curl, err
	}
	bodybytes, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	return bodybytes, resp.StatusCode, curl, nil
}

// OutputInstanceDetails outputs Instance details
func OutputInstanceDetails(bodybytes []byte) error {
	var i Instance
	err := json.Unmarshal(bodybytes, &i)
	if err != nil {
		return err
	}
	data := [][]string{}

	data = append(data, []string{
		strconv.Itoa(i.ID),
		i.Element.Key,
		i.Name,
		strconv.FormatBool(i.Valid),
		strconv.FormatBool(i.Disabled),
		strconv.FormatBool(i.EventsEnabled),
		fmt.Sprintf("%s", i.Tags),
		i.Token,
	})

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Key", "Name", "Valid", "Disabled", "Events", "Tags", "Token"})
	table.SetBorder(false)
	table.AppendBulk(data)
	table.Render()

	return nil
}
