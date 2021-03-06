package dgraph

import (
	"context"
	"strconv"

	dgapi "github.com/dgraph-io/dgo/v200/protos/api"
	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/entity"
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
// It returns err if e is neither space.Entity nor space.Resource or if they fail to serialised to JSON.
func (s *Store) addRequest(ctx context.Context, e store.Entity, opts ...store.Option) (*dgapi.Request, error) {
	sopts := store.Options{}
	for _, apply := range opts {
		apply(&sopts)
	}

	switch v := e.(type) {
	case space.Entity:
		return s.addEntityRequest(ctx, v)
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
			r as uid
	        }
	}
	`

	res := &Resource{
		UID:        "uid(r)",
		XID:        r.UID().Value(),
		Type:       r.Type().String(),
		Name:       r.Name(),
		Group:      r.Group(),
		Version:    r.Version(),
		Kind:       r.Kind(),
		Namespaced: r.Namespaced(),
		Attrs:      AttrsToMap(r.Attrs()),
		DType:      []string{entity.ResourceType.String()},
	}

	return upsertReqJSON(AddOp, res, query, "")
}

// addResourceRequest creates a dgraph API request for adding space.Entity and returns it.
// It returns error if entity fails to be serialised into JSON.
func (s *Store) addEntityRequest(ctx context.Context, e space.Entity, opts ...store.Option) (*dgapi.Request, error) {
	query := `
	{
		entity(func: eq(xid, "` + e.UID().Value() + `")) {
			e as uid
		}

		resource(func: eq(xid, "` + e.Resource().UID().Value() + `")) {
			r as uid
		}
	}
	`

	obj := &Entity{
		UID:       "uid(e)",
		XID:       e.UID().Value(),
		Type:      e.Type().String(),
		Name:      e.Name(),
		Namespace: e.Namespace(),
		Resource: &Resource{
			UID:        "uid(r)",
			XID:        e.Resource().UID().Value(),
			Type:       e.Resource().Type().String(),
			Name:       e.Resource().Name(),
			Group:      e.Resource().Group(),
			Version:    e.Resource().Version(),
			Kind:       e.Resource().Kind(),
			Namespaced: e.Resource().Namespaced(),
			Attrs:      AttrsToMap(e.Resource().Attrs()),
			DType:      []string{entity.ResourceType.String()},
		},
		Attrs: AttrsToMap(e.Attrs()),
		DType: []string{entity.EntityType.String()},
	}

	return upsertReqJSON(AddOp, obj, query, "")
}

// getRequest creates a dgraph API request for getting entity with the given uid and returns it.
// The returned request allows for read only transactions.
func (s *Store) getRequest(ctx context.Context, uid uuid.UID, opts ...store.Option) (*dgapi.Request, error) {
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

// linkRequest creates dgraph API request to link from and to entities stored in the dgraph database.
// The link is created only if both from and to nodes exist and are both of Entity types.
func (s *Store) linkRequest(ctx context.Context, from, to uuid.UID, opts ...store.Option) (*dgapi.Request, error) {
	sopts := store.Options{}
	for _, apply := range opts {
		apply(&sopts)
	}

	q := `
	{
		var(func: eq(xid, "` + from.Value() + `")) @filter(type(Entity)) {
			from as uid
		}

		var(func: eq(xid, "` + to.Value() + `")) @filter(type(Entity)) {
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

	link := &Entity{
		UID:   "uid(from)",
		DType: []string{entity.EntityType.String()},
		Links: []Entity{
			{UID: "uid(to)", DType: []string{entity.EntityType.String()}, Relation: relation, Weight: weight},
		},
	}

	cond := `@if(gt(len(from), 0) AND gt(len(to), 0))`

	return upsertReqJSON(LinkOp, link, q, cond)
}

// unlinkRequest creates dgraph API request to remove the link between and to entities stored in the dgraph database.
// The link is removed only if both from and to nodes exist, are both of Entity types and there is a link between them.
func (s *Store) unlinkRequest(ctx context.Context, from, to uuid.UID, opts ...store.Option) (*dgapi.Request, error) {
	q := `
	{
		var(func: eq(xid, "` + from.Value() + `")) @filter(type(Entity)) {
			from as uid
		}

		var(func: eq(xid, "` + to.Value() + `")) @filter(type(Entity)) {
			to as uid
		}
	}
	`

	sopts := store.Options{}
	for _, apply := range opts {
		apply(&sopts)
	}

	link := &Entity{
		UID: "uid(from)",
		Links: []Entity{
			{UID: "uid(to)"},
		},
	}

	cond := `@if(gt(len(from), 0) AND gt(len(to), 0))`

	return upsertReqJSON(UnlinkOp, link, q, cond)
}
