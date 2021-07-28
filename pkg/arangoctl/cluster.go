package arangoctl

import (
	"context"
)

// ClusterConfig stores information about a cluster that's referred to by one
// or more topic configs. These configs should reflect the reality of what's been
// set up externally; there's no way to "apply" these at the moment.
type ClusterConfig struct {
	Meta ClusterMeta `json:"meta"`
	Spec ClusterSpec `json:"spec"`
}

// ClusterMeta contains (mostly immutable) metadata about the cluster. Inspired
// by the meta fields in Kubernetes objects.
type ClusterMeta struct {
	Name        string `json:"name"`
	Environment string `json:"environment"`
	Description string `json:"description"`
}

// ClusterAuthentication contains auth details to connect to the cluster
type ClusterAuthentication struct {
	Type 		string `json:"type"`
	Username	string `json:"username"`
	Password 	string `json:"password"`
}

// ClusterSpec contains the details necessary to communicate with a ArangoDB cluster.
type ClusterSpec struct {
	// Database to use for this connection
	// This database already has to be present
	Db 		  string 	`json:"db"`

	Authentication ClusterAuthentication `json:"auth"`

	// BootstrapAddrs is a list of one or more broker bootstrap addresses. These can use IPs
	// or DNS names.
	Endpoints []string `json:"endpoints"`
}

// NewAdminClient returns a new admin client using parameters in the current config
func (c ClusterConfig) NewAdminClient(ctx context.Context) (*Client, error) {
	return NewClient(ClientConfig{
		Endpoints: c.Spec.Endpoints,
		Db: c.Spec.Db,
		Context: ctx,
		Authentication: c.Spec.Authentication,
	})
}