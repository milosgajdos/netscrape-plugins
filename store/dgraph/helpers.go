package dgraph

import (
	"encoding/json"
	"fmt"

	dgapi "github.com/dgraph-io/dgo/v200/protos/api"
	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/space/object"
	"github.com/milosgajdos/netscrape/pkg/space/resource"
	"github.com/milosgajdos/netscrape/pkg/store"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// AttrsToMap returns a encoded as map.
func AttrsToMap(a attrs.Attrs) map[string]string {
	attrs := make(map[string]string)
	for _, k := range a.Keys() {
		attrs[k] = a.Get(k)
	}

	return attrs
}

// contains returns true if a contains x.
// NOTE: this is a libear search byt the slice should generally be small
// as it should only contain dgraph.dtypes
func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

// UpsertReqJSON encodes e into JSON and returns Upsert mutation request with given query and cond.
// It returns error if the object failed to be encoded into JSON.
func UpsertReqJSON(op Op, e interface{}, query, cond string) (*dgapi.Request, error) {
	mu, err := MutationJSON(op, e, cond)
	if err != nil {
		return nil, err
	}

	return &dgapi.Request{
		Query:     query,
		Mutations: []*dgapi.Mutation{mu},
		CommitNow: true,
	}, nil
}

// MutationJSON returns JSON mutation for the given op with the given cond.
func MutationJSON(op Op, e interface{}, cond string) (*dgapi.Mutation, error) {
	pb, err := json.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("MutationJSON marshal: %w", err)
	}

	mu := &dgapi.Mutation{
		Cond: cond,
	}

	switch op {
	case AddOp, LinkOp:
		mu.SetJson = pb
	case DelOp, UnlinkOp:
		mu.DeleteJson = pb
	default:
		return nil, fmt.Errorf("MutationJSON unknown op: %v", op)
	}

	return mu, nil
}

func resourceToSpaceResource(r *Resource) (space.Resource, error) {
	if r == nil {
		return nil, object.ErrMissingResource
	}

	uid, err := uuid.NewFromString(r.XID)
	if err != nil {
		return nil, err
	}

	a, err := attrs.NewFromMap(r.Attrs)
	if err != nil {
		return nil, err
	}

	resOpts := []resource.Option{
		resource.WithUID(uid),
		resource.WithAttrs(a),
	}

	return resource.New(r.Name, r.Group, r.Version, r.Kind, r.Namespaced, resOpts...)
}

func objectToSpaceObject(o *Object) (space.Object, error) {
	res, err := resourceToSpaceResource(o.Resource)
	if err != nil {
		return nil, err
	}

	uid, err := uuid.NewFromString(o.XID)
	if err != nil {
		return nil, err
	}

	a, err := attrs.NewFromMap(o.Attrs)
	if err != nil {
		return nil, err
	}

	objOpts := []object.Option{
		object.WithUID(uid),
		object.WithAttrs(a),
	}

	return object.New(o.Name, o.Namespace, res, objOpts...)
}

func decodeJSONGetEntity(b []byte) ([]store.Entity, error) {
	var result struct {
		Entity []struct {
			DType []string `json:"dgraph.type,omitempty"`
		} `json:"entity"`
	}

	if err := json.Unmarshal(b, &result); err != nil {
		return nil, fmt.Errorf("decodeJSONGet %w", err)
	}

	supportedDTypes := []string{ObjectDType, ResourceDType}
	var ents []store.Entity

	for _, e := range result.Entity {
		for _, t := range e.DType {
			if !contains(supportedDTypes, t) {
				continue
			}
			switch t {
			case ObjectDType:
				var result struct {
					Objects []*Object `json:"entity"`
				}
				if err := json.Unmarshal(b, &result); err != nil {
					return nil, fmt.Errorf("decodeJSONObject %w", err)
				}
				for _, o := range result.Objects {
					obj, err := objectToSpaceObject(o)
					if err != nil {
						return nil, err
					}
					ents = append(ents, obj)
				}
			case ResourceDType:
				var result struct {
					Resources []*Resource `json:"entity"`
				}
				if err := json.Unmarshal(b, &result); err != nil {
					return nil, fmt.Errorf("decodeJSONResource %w", err)
				}
				for _, r := range result.Resources {
					res, err := resourceToSpaceResource(r)
					if err != nil {
						return nil, err
					}
					ents = append(ents, res)
				}
			}
		}
	}

	return ents, nil
}

// DecodeJSONEntity accepts JSON response and returns a slice of store.Entity
// NOTE: this is a temporary disgusting hack function; I need to take a cold shower after this.
func DecodeJSONEntity(b []byte, Op Op) ([]store.Entity, error) {
	switch Op {
	case GetOp:
		return decodeJSONGetEntity(b)
	default:
		return nil, ErrUnknownOp
	}
}
