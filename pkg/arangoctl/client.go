package arangoctl

import (
	"context"
	"fmt"
	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	log "github.com/sirupsen/logrus"
)

var (
	stringTypeToCollectionType = map[string]driver.CollectionType{
		"document": driver.CollectionTypeDocument,
		"edge": driver.CollectionTypeEdge,
	}
)

// Client is a general client for interactive with Arango cluster
type Client struct {
	arangoClient 		driver.Client
	ctx 				context.Context
	db 					driver.Database
}

type ClientConfig struct {
	Db 			string
	Endpoints 	[]string
	Authentication ClusterAuthentication
	Context 	context.Context
}

func NewClient(config ClientConfig) (*Client, error) {

	// Setup authentication
	var auth driver.Authentication
	if config.Authentication != (ClusterAuthentication{}) {
		auth = driver.BasicAuthentication(
			config.Authentication.Username,
			config.Authentication.Password,
		)
	}

	// create and store the arango client
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: config.Endpoints,
	})
	if err != nil {
		return nil, err
	}

	c, err := driver.NewClient(driver.ClientConfig{
		Connection: conn,
		Authentication: auth,
	})
	if err != nil {
		return nil, err
	}


	exists, err := c.DatabaseExists(config.Context, config.Db)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, fmt.Errorf("Database %s does not exist", config.Db)
	}

	db, err := c.Database(config.Context, config.Db)
	if err != nil {
		return nil, err
	}

	client := &Client{
		arangoClient: c,
		ctx: config.Context,
		db: db,
	}

	return client, nil
}

// CheckOrCreateCollection checks if collection exists
// and if not will create it. This is mostly used by apply operation
// It will also error out if the collection exists but not same type as Arango
// really does not allow to change type of an existing collection
func (c *Client) CheckOrCreateCollection(name string, collectionType string) error {

	found, err := c.db.CollectionExists(c.ctx, name)
	if err != nil {
		return err
	}

	var collection driver.Collection
	if found {
		collection, err = c.db.Collection(c.ctx, name)
		if err != nil {
			return err
		}

		props, err := collection.Properties(c.ctx)
		if err != nil {
			return err
		}

		if props.CollectionInfo.Type != stringTypeToCollectionType[collectionType] {
			return fmt.Errorf("Existing collection %s type does not match from what is asked", name)
		}

		log.Infof("Collection %s exists", name)
	} else {

		_, ok := stringTypeToCollectionType[collectionType]
		if !ok {
			return fmt.Errorf("Invalid type %s for collection %s", collectionType, name)
		}

		options := &driver.CreateCollectionOptions{
			Type: stringTypeToCollectionType[collectionType],
		}

		collection, err = c.db.CreateCollection(c.ctx, name, options)
		if err != nil {
			return err
		}

		log.Infof("Collection %s created", name)
	}

	return nil
}

func (c *Client) GetIndexes(collection string) ([]driver.Index, error) {
	col, err := c.db.Collection(c.ctx, collection)
	if err != nil {
		return nil, err
	}

	return col.Indexes(c.ctx)
}

func (c *Client) ApplyIndex(collection string, index CollectionIndex) (driver.Index, bool, error) {
	col, err := c.db.Collection(c.ctx, collection)
	if err != nil {
		return nil, false, err
	}

	// Currently we only support persistent indexes
	persistentIndexOptions := driver.EnsurePersistentIndexOptions{
		Unique: index.Options.Unique,
		Sparse: index.Options.Sparse,
		InBackground: index.Options.InBackground,
		Name: index.Name,
	}

	return col.EnsurePersistentIndex(c.ctx, index.Fields, &persistentIndexOptions)
}

func (c *Client) ViewExists(view string) (bool, error) {
	return c.db.ViewExists(c.ctx, view)
}

func (c *Client) GetView(view string) (driver.View, error) {
	return c.db.View(c.ctx, view)
}

func (c *Client) CreateView(name string, options *driver.ArangoSearchViewProperties) (driver.View, error) {
	return c.db.CreateArangoSearchView(c.ctx, name, options)
}

func (c *Client) UpdateView(name string, options driver.ArangoSearchViewProperties) error {
	view, err := c.GetView(name)
	if err != nil {
		return err
	}

	// driver.View is an interface of different types
	// of view but right now only searchview is supported but still better to handle
	// exceptions
	searchview, err := view.ArangoSearchView()
	if err != nil {
		return err
	}

	return searchview.SetProperties(c.ctx, options)
}