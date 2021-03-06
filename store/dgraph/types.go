package dgraph

type Resource struct {
	UID        string            `json:"uid,omitempty"`
	XID        string            `json:"xid,omitempty"`
	Type       string            `json:"type,omitempty"`
	Name       string            `json:"name,omitempty"`
	Group      string            `json:"group,omitempty"`
	Version    string            `json:"version,omitempty"`
	Kind       string            `json:"kind,omitempty"`
	Namespaced bool              `json:"namespaced,omitempty"`
	Attrs      map[string]string `json:"attrs,omitempty"`
	DType      []string          `json:"dgraph.type,omitempty"`
}

type Entity struct {
	UID       string            `json:"uid,omitempty"`
	XID       string            `json:"xid,omitempty"`
	Type      string            `json:"type,omitempty"`
	Name      string            `json:"name,omitempty"`
	Namespace string            `json:"namespace,omitempty"`
	Resource  *Resource         `json:"resource,omitempty"`
	Links     []Entity          `json:"links,omitempty"`
	Attrs     map[string]string `json:"attrs,omitempty"`
	DType     []string          `json:"dgraph.type,omitempty"`

	// Links facets
	LUID     string  `json:"links|uid,omitempty"`
	Relation string  `json:"links|relation,omitempty"`
	Weight   float64 `json:"links|weight,omitempty"`
}
