package dgraph

import (
	"context"
	"strconv"

	dgapi "github.com/dgraph-io/dgo/v200/protos/api"
	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/store"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

const (
	// DefaultRelation is default link relation
	DefaultRelation = "Unknown"
	// DefaultWeight is default link weight
	DefaultWeight = 1.0
)

// addRequest creates a new dgraph API request for adding the given entity and returns it.
// If e is neither space.Object nor space.Resource it returns error.
func (s *Store) addRequest(ctx context.Context, e store.Entity, opts ...store.Option) (*dgapi.Request, error) {
	sopts := store.Options{}
	for _, apply := range opts {
		apply(&sopts)
	}

	switch v := e.(type) {
	case space.Object:
		return s.addObjectRequest(ctx, v)
	case space.Resource:
		return s.addResourceRequest(ctx, v)
	default:
		return nil, store.ErrUnsupported
	}
}

// addResourceRequest creates a dgraph API request for adding space.Resource and returns it.
// It returns error if r fails to be serialised as a JSON object.
func (s *Store) addResourceRequest(ctx context.Context, r space.Resource, opts ...store.Option) (*dgapi.Request, error) {
	query := `
	{
		resource(func: eq(xid, "` + r.UID().Value() + `")) {
			u as uid
	        }
	}
	`

	res := &Resource{
		UID:        "uid(u)",
		XID:        r.UID().Value(),
		Name:       r.Name(),
		Group:      r.Group(),
		Version:    r.Version(),
		Kind:       r.Kind(),
		Namespaced: r.Namespaced(),
		Attrs:      AttrsToMap(r.Attrs()),
		DType:      []string{ResourceDType},
	}

	return upsertReqJSON(AddOp, res, query, "")
}

// addResourceRequest creates a dgraph API request for adding space.Object and returns it.
// It returns error if r fails to be serialised as a JSON object.
func (s *Store) addObjectRequest(ctx context.Context, o space.Object, opts ...store.Option) (*dgapi.Request, error) {
	query := `
	{
		object(func: eq(xid, "` + o.UID().Value() + `")) {
			o as uid
		}

		resource(func: eq(xid, "` + o.Resource().UID().Value() + `")) {
			r as uid
		}
	}
	`

	obj := &Object{
		UID:       "uid(o)",
		XID:       o.UID().Value(),
		Name:      o.Name(),
		Namespace: o.Namespace(),
		Resource: &Resource{
			UID:        "uid(r)",
			XID:        o.Resource().UID().Value(),
			Name:       o.Resource().Name(),
			Group:      o.Resource().Group(),
			Version:    o.Resource().Version(),
			Kind:       o.Resource().Kind(),
			Namespaced: o.Resource().Namespaced(),
			Attrs:      AttrsToMap(o.Resource().Attrs()),
			DType:      []string{ResourceDType},
		},
		Attrs: AttrsToMap(o.Attrs()),
		DType: []string{ObjectDType},
	}

	return upsertReqJSON(AddOp, obj, query, "")
}

// getRequest creates a dgraph API request for getting entity with the given uid and returns it.
// The returned request allows for read only transactions.
func (s *Store) getRequest(ctx context.Context, uid uuid.UID, opts ...store.Option) (*dgapi.Request, error) {
	sopts := store.Options{}
	for _, apply := range opts {
		apply(&sopts)
	}

	q := `
	{
		entity(func: eq(xid, "` + uid.Value() + `")) {
			expand(_all_) {
        			expand(_all_)
				attrs
    			}
    			dgraph.type
			attrs
		}
	}
	`

	return &dgapi.Request{
		Query:    q,
		ReadOnly: true,
	}, nil
}

// deleteRequest creates a dgraph API request for deleting the entity with the given uid and returns it
// It returns error if the delete query fails to be serialized to JSON.
func (s *Store) deleteRequest(ctx context.Context, uid uuid.UID, opts ...store.Option) (*dgapi.Request, error) {
	sopts := store.Options{}
	for _, apply := range opts {
		apply(&sopts)
	}

	q := `
	{
		node(func: eq(xid, "` + uid.Value() + `")) @filter(NOT type(Resource) OR eq(count(~resource), 0)) {
			u as uid
		}
	}
	`

	node := map[string]string{"uid": "uid(u)"}

	cond := `@if(gt(len(u), 0))`

	return upsertReqJSON(DelOp, node, q, cond)
}

// linkRequest link from and to Objects. No other types can be linked.
func (s *Store) linkRequest(ctx context.Context, from, to uuid.UID, opts ...store.Option) (*dgapi.Request, error) {
	sopts := store.Options{}
	for _, apply := range opts {
		apply(&sopts)
	}

	q := `
	{
		var(func: eq(xid, "` + from.Value() + `")) @filter(type(Object)) {
			from as uid
		}

		var(func: eq(xid, "` + to.Value() + `")) @filter(type(Object)) {
			to as uid
		}
	}
	`

	weight := DefaultWeight
	relation := DefaultRelation

	if sopts.Attrs != nil {
		if w, err := strconv.ParseFloat(sopts.Attrs.Get(attrs.Weight), 64); err == nil {
			if w != 0.0 {
				weight = w
			}
		}

		if r := sopts.Attrs.Get(attrs.Relation); r != "" {
			relation = r
		}
	}

	link := &Object{
		UID:   "uid(from)",
		DType: []string{"Object"},
		Links: []Object{
			{UID: "uid(to)", DType: []string{"Object"}, Relation: relation, Weight: weight},
		},
	}

	cond := `@if(gt(len(from), 0) AND gt(len(to), 0))`

	return upsertReqJSON(LinkOp, link, q, cond)
}

// unlinkRequest link from and to Objects. No other types can be linked.
func (s *Store) unlinkRequest(ctx context.Context, from, to uuid.UID, opts ...store.Option) (*dgapi.Request, error) {
	q := `
	{
		var(func: eq(xid, "` + from.Value() + `")) @filter(type(Object)) {
			from as uid
		}

		var(func: eq(xid, "` + to.Value() + `")) @filter(type(Object)) {
			to as uid
		}
	}
	`

	sopts := store.Options{}
	for _, apply := range opts {
		apply(&sopts)
	}

	link := &Object{
		UID: "uid(from)",
		Links: []Object{
			{UID: "uid(to)"},
		},
	}

	cond := `@if(gt(len(from), 0) AND gt(len(to), 0))`

	return upsertReqJSON(UnlinkOp, link, q, cond)
}
