package arangoctl

// Resource provides methods to identify resource and do operation on it
type Resource interface {
	Apply(client *Client)			error
	GetKind()						string
}
