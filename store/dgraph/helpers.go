package dgraph

import (
	"encoding/json"
	"fmt"

	dgapi "github.com/dgraph-io/dgo/v200/protos/api"
	"github.com/milosgajdos/netscrape/pkg/attrs"
)

// Op is dgraph operation.
type Op int

const (
	// AddOp is add operation
	AddOp Op = iota
	// DelOp is delete operation
	DelOp
	// LinkOp is link operation
	LinkOp
	// UnlinkOp is unlink operation
	UnlinkOp
)

// AttrsToMap returns a encoded as map.
func AttrsToMap(a attrs.Attrs) map[string]string {
	attrs := make(map[string]string)
	for _, k := range a.Keys() {
		attrs[k] = a.Get(k)
	}

	return attrs
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
		return nil, fmt.Errorf("Unknown op: %v", op)
	}

	return mu, nil
}
