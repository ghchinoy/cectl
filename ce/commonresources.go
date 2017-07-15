package ce

// CommonResource represents a normalized data object (resource)
type CommonResource struct {
	Name               string  `json:"name"`
	ElementInstanceIDs []int   `json:"elementInstanceIds"`
	Fields             []Field `json:"fields"`
}

// Field is a set of  a common resource fields
type Field struct {
	Type            string `json:"type"`
	Path            string `json:"path"`
	AssociatedLevel string `json:"organization"`
	AssociatedID    int    `json:"associatedId"`
}
