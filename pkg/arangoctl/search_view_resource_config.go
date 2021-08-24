package arangoctl

import "github.com/arangodb/go-driver"

// SearchViewConfig defines YAML for view resource
type SearchViewConfig struct {
	Kind 	string			`json:"Kind"`
	Meta    ViewMeta		`json:"meta"`
	// I think its just easier to map the prorperties to view
	// as it lets us expand on it later on
	SearchViewProperties	driver.ArangoSearchViewProperties		`json:"spec"`
}

// ViewMeta defines immutable parts of collections
// We dont allow renaming of views yet. You will just have to delete
// and apply it again
type ViewMeta struct {
	Name 	string `json:"name"`
}

func (v SearchViewConfig) GetResource() Resource {
	view := &SearchView{config: v}
	return view
}
