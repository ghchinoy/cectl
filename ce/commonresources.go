package ce

const (
	// CommonResourcesURI is the base URI of the hidden API
	// for Common Object Resources; this one provides an array of common objects
	// with the element instance IDs associated, as well as details about
	// the field's heirarchy (org, account, instance)
	CommonResourcesURI = "/common-resources"
	// CommonResourceURI is the base URI for common object resources
	// this is a simple object with keys being the common object names and no
	// details about associated elements or field level hierarchy
	CommonResourceURI = "/organizations/objects/definitions"
	// CommonResourceDefinitionsFormatURI is a string format for the URI of Common Object Resource definition, given a name of a Common Object
	CommonResourceDefinitionsFormatURI = "/organizations/objects/%s/definitions"
	// CommonResourceTransformationsFormatURI is the string format for the URI of an Element's transformation / mapping, given an element key and an object name
	CommonResourceTransformationsFormatURI = "/organizations/elements/%s/transformations/%s"
)

// CommonResource represents a normalized data object (resource)
type CommonResource struct {
	Name               string  `json:"name,omitempty"`
	ElementInstanceIDs []int   `json:"elementInstanceIds,omitempty"`
	Fields             []Field `json:"fields"`
	Level              string  `json:"level,omitempty"`
}

// Field is a set of  a common resource fields
type Field struct {
	Type            string `json:"type"`
	Path            string `json:"path"`
	AssociatedLevel string `json:"organization,omitempty"`
	AssociatedID    int    `json:"associatedId,omitempty"`
}
