package arangoctl

// CollectionConfig defines YAML for collection resource
type CollectionConfig struct {
	Kind string         `json:"Kind"`
	Meta CollectionMeta `json:"meta"`
	Spec CollectionSpec `json:"spec"`
}

// CollectionMeta defines mostly immutable parts of collections.
// We dont allow renaming of collection name through collection meta yet
type CollectionMeta struct {
	Name 	string `json:"name"`
	Type	string `json:"type"`
}

// CollectionSpec defines the mutable part of collections
// As of now we only let you update indexes
type CollectionSpec struct {
	Indexes		[]CollectionIndex `json:"indexes"`
}

// CollectionIndex defines one index
type CollectionIndex struct {
	Name    string                 `json:"name"`
	Type    string                 `json:"type"`
	Fields  []string               `json:"fields"`
	Options CollectionIndexOptions `json:"options"`
}

// CollectionIndexOptions defines options for an index
// Different index type have different set of option but this is a global
// object that is a superset of all options. We will pick and choose
// the relevant one when we create indexes based on type
type CollectionIndexOptions struct {
	Unique			bool `json:"unique"`
	Sparse 			bool `json:"sparse"`
	InBackground	bool `json:"inbackground"`
}

func (c CollectionConfig) GetResource() Resource {
	collection := &Collection{config: c}
	return collection
}