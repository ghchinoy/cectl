package ce

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
)

const (
	// FormulasURI is the base URI for Formulas
	FormulasURI = "/formulas"
	// FormulaCancelExecutionURIFormat is the API URI for cancelling a formula instance execution
	FormulaCancelExecutionURIFormat = "/formulas/instances/executions/%s"
	// FormulaExecutionsURIFormat is the URI to obtain executions of a Formula Instance
	FormulaExecutionsURIFormat = "/formulas/instances/%s/executions"
	// FormulaRetryExecutionURI is the URI to retry a Formula execution
	FormulaRetryExecutionURI = "/formulas/instances/executions/%s/retries"
	// FormulaURIFormat is the main, partial API URI for Formula
	FormulaURIFormat = "/formulas/%s"
	// FormulaInstancesURI is the main API URI for Formula Instances
	FormulaInstancesURI = "/formulas/instances"
	// FormulaInstancesURIFormat is the URI to obtain instances of a Formula template
	FormulaInstancesURIFormat = "/formulas/%s/instances"
)

// Formula represents the structure of a CE Formula
type Formula struct {
	ID             int             `json:"id"`
	Name           string          `json:"name"`
	UserID         int             `json:"userId"`
	AccountID      int             `json:"accountId"`
	CreatedDate    time.Time       `json:"createdDate"`
	Steps          []Step          `json:"steps"`
	Triggers       []Trigger       `json:"triggers"`
	Active         bool            `json:"active"`
	SingleThreaded bool            `json:"singleThreaded"`
	Configuration  []Configuration `json:"configuration"`
	API            string          `json:"api"`
}

// Step represents a Formula step
type Step struct {
	ID         int         `json:"id"`
	OnSuccess  []string    `json:"onSuccess"`
	OnFailure  []string    `json:"onFailure"`
	Name       string      `json:"name"`
	Type       string      `json:"type"`
	Properties interface{} `json:"properties"`
}

// Trigger represents an action that starts a Formula
type Trigger struct {
	ID         int         `json:"id"`
	Type       string      `json:"type"`
	OnSuccess  []string    `json:"onSuccess"`
	OnFailure  []string    `json:"onFailure"`
	Async      bool        `json:"async"`
	Name       string      `json:"name"`
	Properties interface{} `json:"properties"`
}

// Configuration represents a configuration for a formula
type Configuration struct {
	ID       int `json:"id"`
	Key      string
	Name     string
	Type     string
	Required bool
}

// FormulaInstance represents a configured instance of a Formula
type FormulaInstance struct {
	ID            int         `json:"id"`
	Formula       Formula     `json:"formula"`
	Name          string      `json:"name"`
	CreatedDate   time.Time   `json:"createdDate"`
	Settings      interface{} `json:"settings"`
	Active        bool        `json:"active"`
	Configuration interface{} `json:"configuration"`
}

// FormulaInstanceConfig represents a configuration used when creating an Instance of a Formula
type FormulaInstanceConfig struct {
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

// FormulaInstanceCreationResponse is the response returned when a Formula Instance is triggered
type FormulaInstanceCreationResponse struct {
	ID        int    `json:"id"`
	RequestID string `json:"requestId"`
	Message   string `json:"message"`
}

// FormulaInstanceExecution is a brief info about an instance Execution
type FormulaInstanceExecution struct {
	ID                int       `json:"id"`
	FormulaInstanceID int       `json:"formulaInstanceId"`
	Status            string    `json:"status"`
	CreateDate        time.Time `json:"createdDate"`
	UpdatedDate       time.Time `json:"updatedDate"`
}

// GetInstancesOfFormula returns an Instance array, given a Formula ID and an Auth header
func GetInstancesOfFormula(id int, baseurl string, auth string) ([]FormulaInstance, error) {
	var instances []FormulaInstance

	url := fmt.Sprintf("%s%s", baseurl,
		fmt.Sprintf("/formulas/%v/instances", id))

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		//fmt.Println("Can't construct request", err.Error())
		return instances, err
	}
	req.Header.Add("Authorization", auth)
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		//fmt.Println("Cannot process response", err.Error())
		return instances, err
	}
	bodybytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	err = json.Unmarshal(bodybytes, &instances)
	if err != nil {
		return instances, err
	}

	return instances, nil
}

// FormulaDetailsTableOutput prints to stdout an ASCII rendered table of the details of a Formula
func FormulaDetailsTableOutput(f Formula) error {

	// basic formula info
	data := [][]string{}

	if len(f.Triggers) < 1 {
		fmt.Printf("Formula %v is malformed, no trigger present\n", f.ID)

	} else {
		data = append(data, []string{
			strconv.Itoa(f.ID),
			f.Name,
			strconv.FormatBool(f.Active),
			strconv.Itoa(len(f.Steps)),
			f.Triggers[0].Type,
		})

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Name", "active", "steps", "trigger"})
		table.SetBorder(false)
		table.AppendBulk(data)
		table.Render()

		fmt.Println()

		// Triggers

		data = [][]string{}

		for _, v := range f.Triggers {
			data = append(data, []string{
				strconv.Itoa(v.ID),
				v.Name,
				v.Type,
				strconv.FormatBool(v.Async),
				fmt.Sprintf("%s", v.OnSuccess),
			})
		}

		table = tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Name", "Type", "Async", "Success"})
		table.SetBorder(false)
		table.AppendBulk(data)
		table.Render()

		fmt.Println()

		// Steps

		data = [][]string{}

		for _, v := range f.Steps {
			data = append(data, []string{
				strconv.Itoa(v.ID),
				v.Name,
				v.Type,
				fmt.Sprintf("%s", v.OnSuccess),
				fmt.Sprintf("%s", v.OnFailure),
			})
		}

		table = tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Name", "Type", "Success", "Failure"})
		table.SetBorder(false)
		table.AppendBulk(data)
		table.Render()
	}

	return nil
}
