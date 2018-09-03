package tokens

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Config holds the CE configuration info
type Config struct {
	Username    string `survey:"un"`
	Environment string `survey:"env"`
	Password    string `survey:"pwd"`
	Output      string
}

// Token is a temp structure for a Cloud Elements token
type Token struct {
	User         string `json:"userSecret"`
	Organization string `json:"organizationSecret"`
}

// ObtainCEToken returns a Token struct given a Config struct
func ObtainCEToken(config Config) (Token, error) {
	var token Token

	// POST to /authentication
	payload := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		config.Username,
		config.Password,
	}
	payloadBytes, err := json.Marshal(payload)
	url := fmt.Sprintf("%s/%s", config.Environment, "authentication")
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payloadBytes))
	if err != nil {
		return token, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return token, err
	}
	defer res.Body.Close()
	bodybytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return token, err
	}

	var response map[string]interface{}
	err = json.Unmarshal(bodybytes, &response)
	if err != nil {
		return token, err
	}
	if _, ok := response["token"]; !ok {
		return token, fmt.Errorf("empty token received for %s @ %s", config.Username, url)
	}

	// GET /authentication/secrets
	url = fmt.Sprintf("%s/%s", config.Environment, "/authentication/secrets")
	req, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return token, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", response["token"].(string)))
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		return token, err
	}
	defer res.Body.Close()
	bodybytes, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return token, err
	}

	err = json.Unmarshal(bodybytes, &token)
	if err != nil {
		return token, err
	}

	return token, nil
}

// OutputCectlTOML outputs TOML
func OutputCectlTOML(config Config, token Token) {
	fmt.Printf("[%s]\n", strings.Replace(
		strings.Split(config.Username, "@")[0],
		"+",
		"_",
		-1,
	),
	)
	fmt.Printf("  base = \"%s\"\n", config.Environment)
	fmt.Printf("  org =\"%s\"\n", token.Organization)
	fmt.Printf("  user = \"%s\"\n", token.User)
}

// PostmanEnvironment is a structure for a Postman Environment json config
type PostmanEnvironment struct {
	Name   string                    `json:"name"`
	Values []PostmanEnvironmentValue `json:"values"`
}

// PostmanEnvironmentValue is a value for a Postman Environment json config
type PostmanEnvironmentValue struct {
	Enabled bool   `json:"enabled"`
	Key     string `json:"key"`
	Value   string `json:"value"`
	Type    string `json:"type"`
}

// OutputPostmanEnvJSON outputs JSON for a Postman Environment file
func OutputPostmanEnvJSON(config Config, token Token) error {
	envjson := PostmanEnvironment{
		Name: config.Username,
		Values: []PostmanEnvironmentValue{
			{Enabled: true, Type: "text", Key: "PlatformBaseURI", Value: config.Environment},
			{Enabled: true, Type: "text", Key: "OrganizationID", Value: token.Organization},
			{Enabled: true, Type: "text", Key: "AdminUserID", Value: token.User},
		},
	}

	outputbytes, err := json.MarshalIndent(envjson, "", "  ")
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", outputbytes)

	return nil
}
