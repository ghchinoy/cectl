package ce

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/moul/http2curl"
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
	InstancesObjectsDefinitionsURI       = "/instances/objects/definitions"
	InstancesTransformationsURI          = "/instances/transformations"
	InstanceTransformationsFormatURI     = "/instances/%s/transformations"
	InstanceDocFormatURI                 = "/instances/%s/docs"
)

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

func GetInstanceInfo(id string) {}

func GetInstanceDocs(id string) {}

func GetInstanceTransformations(id string) {}
