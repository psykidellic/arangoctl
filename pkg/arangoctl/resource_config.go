package arangoctl

// ResourceConfig provides methods to convert to actual arango Resource object
type ResourceConfig interface {
	GetResource() Resource
}
