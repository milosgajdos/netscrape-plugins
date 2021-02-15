// package dgraph implements store.Store interface
// from github.com/milosgajdos/netscrape Go module.
package dgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	dgapi "github.com/dgraph-io/dgo/v200/protos/api"

	"github.com/milosgajdos/netscrape/pkg/graph"
	"github.com/milosgajdos/netscrape/pkg/store"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

const (
	DefaultRelation = "Unknown"
)

// Store is dgraph store
type Store struct {
	g *Graph
	c *Client
}

// New creates new dgraph store and returns it.
func NewStore(c *Client, opts ...Option) (*Store, error) {
	sopts := Options{}
	for _, apply := range opts {
		apply(&sopts)
	}

	g, err := NewGraph(c)
	if err != nil {
		return nil, err
	}

	return &Store{
		g: g,
		c: c,
	}, nil
}

// Graph returns graph handle.
func (s *Store) Graph(ctx context.Context) (graph.Graph, error) {
	return s.g, nil
}

// Add Entity to store.
func (s *Store) Add(ctx context.Context, e store.Entity, opts ...store.Option) error {
	sopts := store.Options{}
	for _, apply := range opts {
		apply(&sopts)
	}

	query := `
	{
		entity(func: eq(xid, "` + e.UID().Value() + `")) {
			u as uid
		}
	}
	`

	a := make(map[string]string)
	for _, k := range e.Attrs().Keys() {
		a[k] = e.Attrs().Get(k)
	}

	ent := &Entity{
		UID:       "uid(u)",
		XID:       e.UID().Value(),
		Name:      e.Name(),
		Namespace: e.Namespace(),
		Resource: Resource{
			Name:       e.Resource().Name(),
			Group:      e.Resource().Group(),
			Version:    e.Resource().Version(),
			Kind:       e.Resource().Kind(),
			Namespaced: e.Resource().Namespaced(),
			DType:      []string{"Resource"},
		},
		Attrs: a,
		DType: []string{"Entity"},
	}

	pb, err := json.Marshal(ent)
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	mu := &dgapi.Mutation{
		SetJson: pb,
	}

	req := &dgapi.Request{
		Query:     query,
		Mutations: []*dgapi.Mutation{mu},
		CommitNow: true,
	}

	txn := s.c.NewTxn()

	var txnErr error
	defer func() { txnErr = txn.Discard(ctx) }()

	if _, err := txn.Do(ctx, req); err != nil {
		return fmt.Errorf("txn Add: %w", err)
	}

	return txnErr
}

// Delete Entity from store.
func (s *Store) Delete(ctx context.Context, e store.Entity, opts ...store.Option) error {
	sopts := store.Options{}
	for _, apply := range opts {
		apply(&sopts)
	}

	q := `
		query Entity($xid: string) {
			entity(func: eq(xid, $xid)) {
				uid
				xid
			}
		 }
		`

	txn := s.c.NewTxn()

	var txnErr error
	defer func() { txnErr = txn.Discard(ctx) }()

	resp, err := txn.QueryWithVars(ctx, q, map[string]string{"$xid": e.UID().Value()})
	if err != nil {
		return err
	}

	var r struct {
		Result []Entity `json:"entity"`
	}

	if err = json.Unmarshal(resp.Json, &r); err != nil {
		return err
	}

	res := len(r.Result)

	switch {
	case res == 0:
		return nil
	case res > 1:
		return fmt.Errorf("txn query: ErrDuplicateNode")
	}

	n := r.Result[0]

	node := map[string]string{"uid": n.UID}
	pb, err := json.Marshal(node)
	if err != nil {
		return fmt.Errorf("delete json marshal: %w", err)
	}

	mu := &dgapi.Mutation{
		CommitNow:  true,
		DeleteJson: pb,
	}

	ctx = context.Background()
	_, err = s.c.NewTxn().Mutate(ctx, mu)
	if err != nil {
		return fmt.Errorf("txn Delete: %w", err)
	}

	return txnErr
}

// Link two entities in store.
func (s *Store) Link(ctx context.Context, from, to uuid.UID, opts ...store.Option) error {
	sopts := store.Options{}
	for _, apply := range opts {
		apply(&sopts)
	}

	query := `
	{
		from as var(func: eq(xid, "` + from.Value() + `")) {
			fid as uid
		}

		to as var(func: eq(xid, "` + to.Value() + `")) {
			tid as uid
		}
	}
	`
	weight := graph.DefaultWeight
	relation := DefaultRelation

	if sopts.Attrs != nil {
		if w, err := strconv.ParseFloat(sopts.Attrs.Get("weight"), 64); err == nil {
			if w != 0.0 {
				weight = w
			}
		}

		if r := sopts.Attrs.Get("relation"); r != "" {
			relation = r
		}
	}

	node := &Entity{
		UID:   "uid(fid)",
		DType: []string{"Entity"},
		Links: []Entity{
			{UID: "uid(tid)", DType: []string{"Entity"}, Relation: relation, Weight: weight},
		},
	}

	pb, err := json.Marshal(node)
	if err != nil {
		return fmt.Errorf("link json marshal: %w", err)
	}

	mu := &dgapi.Mutation{
		Cond:    `@if(gt(len(from), 0) AND gt(len(to), 0))`,
		SetJson: pb,
	}

	req := &dgapi.Request{
		Query:     query,
		Mutations: []*dgapi.Mutation{mu},
		CommitNow: true,
	}

	txn := s.c.NewTxn()

	var txnErr error
	defer func() { txnErr = txn.Discard(ctx) }()

	if _, err := txn.Do(ctx, req); err != nil {
		return fmt.Errorf("txn Link: %w", err)
	}

	return txnErr
}

// Unlink two entities in store.
func (s *Store) Unlink(ctx context.Context, from, to uuid.UID, opts ...store.Option) error {
	sopts := store.Options{}
	for _, apply := range opts {
		apply(&sopts)
	}

	query := `
		{
			from as var(func: eq(xid, "` + from.Value() + `")) {
				fid as uid
			}

			to as var(func: eq(xid, "` + to.Value() + `")) {
				tid as uid
			}
		}
		`

	node := &Entity{
		UID:   "uid(fid)",
		DType: []string{"Entity"},
		Links: []Entity{
			{UID: "uid(tid)", DType: []string{"Entity"}},
		},
	}

	pb, err := json.Marshal(node)
	if err != nil {
		return err
	}

	mu := &dgapi.Mutation{
		Cond:       `@if(gt(len(from), 0) AND gt(len(to), 0))`,
		DeleteJson: pb,
	}

	req := &dgapi.Request{
		Query:     query,
		Mutations: []*dgapi.Mutation{mu},
		CommitNow: true,
	}

	txn := s.c.NewTxn()

	var txnErr error
	defer func() { txnErr = txn.Discard(ctx) }()

	if _, err := txn.Do(ctx, req); err != nil {
		return err
	}

	return txnErr
}
