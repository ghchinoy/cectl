package output

import (
	"encoding/json"

	"github.com/ghchinoy/ce-go/ce"
)

// ROIData structure for ROI data
type ROIData struct {
	Name    string `json:"name"`
	Key     string `json:"key"`
	Hub     string `json:"hub"`
	APIType string `json:"apiType"`
	Beta    bool   `json:"beta"`
	Active  bool   `json:"active"`
}

// ElementsForROICalculator outputs JSON as used in https://github.com/cdoelling/ce-roi-calc
func ElementsForROICalculator(elementData, intelligenceData []byte) ([]byte, error) {
	var elements []ce.Element
	var intelligence []ce.Metadata

	var outputs []ROIData

	err := json.Unmarshal(elementData, &elements)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(intelligenceData, &intelligence)
	if err != nil {
		return nil, err
	}

	metadata := intelligenceMap(intelligence)

	for _, v := range elements {
		if v.Active == true { // only Active
			if v.Private == false { // only not Private elements
				m := metadata[v.Key]
				if m.API.Type != "" { // must have an API type
					data := ROIData{v.Name, v.Key, v.Hub, m.API.Type, v.Beta, v.Active}
					outputs = append(outputs, data)
				}
			}
		}
	}
	outbytes, err := json.Marshal(outputs)
	if err != nil {
		return nil, err
	}
	return outbytes, nil
}

func intelligenceMap(metadata []ce.Metadata) map[string]ce.Metadata {
	intmetadata := make(map[string]ce.Metadata)
	for _, v := range metadata {
		intmetadata[v.Key] = v
	}
	return intmetadata
}
