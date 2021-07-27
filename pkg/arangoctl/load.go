package arangoctl

import (
	"errors"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/mitchellh/mapstructure"
	"io/ioutil"
)

var (
	// Kind for resource config not provided
	InvalidKind = errors.New("kind missing or invalid type. Must be string")
)

// LoadClusterFile loads a ClusterConfig from a apth to YAML file
func LoadClusterFile(path string) (ClusterConfig, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return ClusterConfig{}, err
	}
	return LoadClusterBytes(contents)
}

// LoadClusterBytes loads a ClusterConfig from YAML bytes
func LoadClusterBytes(contents []byte) (ClusterConfig, error) {
	config := ClusterConfig{}
	err := yaml.Unmarshal(contents, &config)
	return config, err
}

// LoadResourceFile loads a generic map interface from a
func LoadResourceFile(path string) (ResourceConfig, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return LoadResourceBytes(contents)
}

// LoadResourceBytes loads ResourceConfig from YAML bytes
func LoadResourceBytes(contents []byte) (ResourceConfig, error) {
	var config map[string]interface{}
	err := yaml.Unmarshal(contents, &config)
	if err != nil {
		return nil, err
	}

	// Using kind we load up resource config to its concrete implementation
	kind, ok := config["Kind"].(string)
	if !ok {
		return nil, InvalidKind
	}

	switch kind {
	case "Collection":
		var c CollectionConfig
		mapstructure.Decode(config, &c)
		return c, nil
	case "View":
		var v SearchViewConfig

		// Mapstructure is unable to reflect insidde the spec
		// to map to searchviewproperties straight
		// so we do double mapping
		if _, ok := config["spec"]; ok {
			mapstructure.Decode(config, &v)
			searchProperties := config["spec"]
			mapstructure.Decode(searchProperties, &v.SearchViewProperties)
			//do something here
		} else {
			return nil, fmt.Errorf("no spec defined for searchview")
		}

		return v, nil
	default:
		return nil, InvalidKind
	}
}