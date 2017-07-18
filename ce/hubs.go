package ce

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/moul/http2curl"
)

const (
	hubsURI = "/hubs"
)

// Hub represents metadata about a hub
type Hub struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Key         string `json:"key"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
	VideoLink   string `json:"videoLink"`
}

// ListHubs returns a list of hubs on the platform
func ListHubs(base, user, org string, outputJson bool) ([]Hub, string, error) {
	var hubs []Hub
	var curl string

	url := fmt.Sprintf("%s%s", base, hubsURI)
	auth := fmt.Sprintf("User %s, Organization %s", user, org)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return hubs, curl, err
	}
	req.Header.Add("Authorization", auth)
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return hubs, curl, err
	}
	bodybytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Print(resp.Status)
		if resp.StatusCode == 404 {
			fmt.Printf("Unable to contact CE API, %s\n", url)
			return hubs, curl, err
		}
		fmt.Println()
	}

	curlCmd, _ := http2curl.GetCurlCommand(req)
	curl = fmt.Sprintf("%s", curlCmd)

	if outputJson {
		fmt.Printf("%s\n", bodybytes)
		return hubs, curl, nil
	}

	err = json.Unmarshal(bodybytes, &hubs)
	if err != nil {
		return hubs, curl, err
	}

	return hubs, curl, nil
}
