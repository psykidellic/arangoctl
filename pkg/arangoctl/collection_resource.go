package arangoctl

import (
	"fmt"
	"github.com/arangodb/go-driver"
	log "github.com/sirupsen/logrus"
	"strings"
)

var (
	validIndexTypes = map[string]bool{
		"persistent": true,
	}
)

// Collection implementation of resource
type Collection struct {
	config CollectionConfig
}

// GetKind implements Resource.Getkind interface
func (c *Collection) GetKind() string {
	return "collection"
}

// Apply applies collection based on config specs
func (c *Collection) Apply(client *Client) error {

	// Check if collection exists. If not create it
	// with the specs provided
	err := client.CheckOrCreateCollection(c.config.Meta.Name, strings.ToLower(c.config.Meta.Type))
	if err != nil {
		return err
	}

	err = c.validate()
	if err != nil {
		return err
	}

	// Currently, only thing we do is manage persistent indexes
	return c.applyIndexes(client)
}

// validate Validates user provided collection config is correct with
// known values
func (c *Collection) validate() error {

	for _, index := range c.config.Spec.Indexes {
		_, ok := validIndexTypes[index.Type]
		if !ok {
			return fmt.Errorf("Index %s not of valid type", index.Name)
		}
	}

	return nil
}

// applyIndexes applies indexes. it will delete indexes which are present in collection
// but not provided in config. will update with new properties for existing indexes with the
// same name. We cannot rename an index yet. Its a delete -> create operation if you rename
// index.
func (c *Collection) applyIndexes (client *Client) error {
	err := c.deleteIndexes(client)
	if err != nil {
		return err
	}

	return c.ensureIndexes(client)
}

// deleteIndexes will delete any indexes that are present but not
// asked in the yaml
func (c *Collection) deleteIndexes (client *Client) error {
	// Get indexes present
	indexes, err := client.GetIndexes(c.config.Meta.Name)
	if err != nil {
		return err
	}

	// get existing and asked so we can do set maths
	existing := make([]string, 0)
	existingid := make(map[string]driver.Index)
	asked := make([]string, 0)

	for _, existingindex := range indexes {
		// We cannot work on primary and edge indexes
		if existingindex.Type() == driver.PrimaryIndex || existingindex.Type() == driver.EdgeIndex {
			continue
		}
		existing = append(existing, existingindex.UserName())
		// add actual index object to map so we can delete it later
		existingid[existingindex.UserName()] = existingindex
	}

	for _, askedindex := range c.config.Spec.Indexes {
		asked = append(asked, askedindex.Name)
	}

	log.Infof("%s : Existing Indexes: %v, Asked Indexes: %v", c.config.Meta.Name, existing, asked)

	// all indxes to be deleted
	// there is no generic set subtraction so we
	// implement our own
	diff := make([]string, 0)
	m := make(map[string]bool)
	for _, item := range asked {
		m[item] = true
	}

	for _, item := range existing {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}

	// Delete extra indexes
	if len(diff) > 0 {
		for _, indexname := range diff {
			log.Infof("Deleting indexes: %v", indexname)
			index := existingid[indexname]
			err := index.Remove(client.ctx)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ensureIndexes is basically a wapper over driver.Ensureindex
// to create/update indexes from yaml specified
func (c *Collection) ensureIndexes (client *Client) error {

	for _, index := range c.config.Spec.Indexes {
		_, created, err := client.ApplyIndex(c.config.Meta.Name, index)
		if err != nil {
			return err
		}

		var operation string
		if created {
			operation = "created"
		} else {
			operation = "updated"
		}
		log.Infof("Index %s %s", index.Name, operation)
	}

	return nil
}


