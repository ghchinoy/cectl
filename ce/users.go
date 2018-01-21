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
	// UsersURI is the base uri for the Users resource
	UsersURI = "/users"
	// UserRoleURIFormat is a format string for the Roles of a user
	UserRoleURIFormat = "/users/%v/roles"
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
	Roles                []Role `json:"roles,omitempty"`
}

// Role represents a users role
type Role struct {
	Active      bool          `json:"active,omitempty"`
	Description string        `json:"description,omitempty"`
	ID          int           `json:"id,omitempty"`
	Key         string        `json:"key,omitempty"`
	Name        string        `json:"name,omitempty"`
	Features    []RoleFeature `json:"features,omitempty"`
}

// RoleFeature is a feature of a role
type RoleFeature struct {
	ID          int    `json:"id,omitempty"`
	ReadOnly    bool   `json:"read_only,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	CreateDate  string `json:"createDate,omitempty"`
	Active      bool   `json:"active,omitempty"`
}

// AddRolesToUsers appends Role array to Users
func AddRolesToUsers(base, auth string, usersbytes []byte) ([]byte, int, string, error) {

	var users []User

	err := json.Unmarshal(usersbytes, &users)
	if err != nil {
		return nil, 0, "", err
	}

	for i, u := range users {
		url := fmt.Sprintf("%s%s", base, fmt.Sprintf(UserRoleURIFormat, u.ID))
		client := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("Can't construct request", err.Error())
			os.Exit(1)
		}
		req.Header.Add("Authorization", auth)
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			break
		}
		bodybytes, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()

		var roles []Role
		err = json.Unmarshal(bodybytes, &roles)
		users[i].Roles = roles
	}

	bodybytes, err := json.Marshal(users)

	return bodybytes, 200, "", nil
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

	hasRoles := false
	err := json.Unmarshal(usersbytes, &users)
	if err != nil {
		return err
	}

	for _, u := range users {
		var roles []string
		if len(u.Roles) > 0 {
			hasRoles = true
			for _, r := range u.Roles {
				roles = append(roles, r.Key)
			}
		}
		data = append(data, []string{
			strconv.Itoa(u.ID),
			u.FullName,
			u.EMail,
			u.LastLoginDate,
			strconv.FormatBool(u.Active),
			strings.Join(roles, ","),
		})
	}

	table := tablewriter.NewWriter(os.Stdout)
	if hasRoles {
		table.SetHeader([]string{"ID", "Name", "EMail", "Last Login", "Active", "Roles"})
	} else {
		table.SetHeader([]string{"ID", "Name", "EMail", "Last Login", "Active"})
	}
	table.SetBorder(false)
	table.AppendBulk(data)
	table.Render()

	return nil
}
