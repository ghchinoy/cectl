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
	UsersURI = "/users"
)

type User struct {
	ID                   int    `json:"id,omitempty"`
	CreatedDate          string `json:"createdDte,omitempty"`
	LastLoginDate        string `json:"lastLoginDate,omitempty"`
	FullName             string `json:"fullName,omitempty"`
	FirstName            string `json:"firstName,omitempty"`
	LastName             string `json:"lastName,omitempty"`
	Password             string `json:"password,omitempty"`
	EMail                string `json:"email,omitempty"`
	Active               bool   `json:"active,omitempty"`
	AccountExpired       bool   `json:"accountExpired,omitempty"`
	AccountLocked        bool   `json:"accountLocked,omitempty"`
	AccountNonExpired    bool   `json:"accountNonExpired,omitempty"`
	EmailValid           bool   `json:"emailValid,omitempty"`
	CredentialNonExpired bool   `json:"credentialNonExpired,omitempty"`
	AccountNonLocked     bool   `json:"accountNonLocked,omitempty"`
	Enabled              bool   `json:"enabled,omitempty"`
}

// GetAllUsers returns a byte stream of users, status code, curl cmd, and error (if occured)
func GetAllUsers(base, auth string) ([]byte, int, string, error) {

	var bodybytes []byte

	url := fmt.Sprintf("%s%s", base, UsersURI)

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

func FormatUserList(usersbytes []byte) error {

	data := [][]string{}

	var users []User

	err := json.Unmarshal(usersbytes, &users)
	if err != nil {
		return err
	}

	for _, u := range users {
		data = append(data, []string{
			strconv.Itoa(u.ID),
			u.FullName,
			u.EMail,
			u.LastLoginDate,
			strconv.FormatBool(u.Active),
		})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "EMail", "Last Login", "Active"})
	table.SetBorder(false)
	table.AppendBulk(data)
	table.Render()

	return nil
}
