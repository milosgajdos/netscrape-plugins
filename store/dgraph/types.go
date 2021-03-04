package dgraph

const (
	// ObjectDType is Object entity dgraph.dtype
	ObjectDType = "Object"
	// ResourceDType is Resource entity dgraph.dtype
	ResourceDType = "Resource"
)

// Op is dgraph operation.
type Op int

const (
	// AddOp is add operation
	AddOp Op = iota
	// DelOp is delete operation
	DelOp
	// GetOp is get operation
	GetOp
	// LinkOp is link operation
	LinkOp
	// UnlinkOp is unlink operation
	UnlinkOp
)

type Resource struct {
	UID        string            `json:"uid,omitempty"`
	XID        string            `json:"xid,omitempty"`
	Name       string            `json:"name,omitempty"`
	Group      string            `json:"group,omitempty"`
	Version    string            `json:"version,omitempty"`
	Kind       string            `json:"kind,omitempty"`
	Namespaced bool              `json:"namespaced,omitempty"`
	Attrs      map[string]string `json:"attrs,omitempty"`
	DType      []string          `json:"dgraph.type,omitempty"`
}

type Object struct {
	UID       string            `json:"uid,omitempty"`
	XID       string            `json:"xid,omitempty"`
	Name      string            `json:"name,omitempty"`
	Namespace string            `json:"namespace,omitempty"`
	Resource  *Resource         `json:"resource,omitempty"`
	Links     []Object          `json:"links,omitempty"`
	Attrs     map[string]string `json:"attrs,omitempty"`
	DType     []string          `json:"dgraph.type,omitempty"`

	// Links facets
	LUID     string  `json:"links|uid,omitempty"`
	Relation string  `json:"links|relation,omitempty"`
	Weight   float64 `json:"links|weight,omitempty"`
}
