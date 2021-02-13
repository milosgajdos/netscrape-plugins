// package dgraph implements store.Store interface
// from github.com/milosgajdos/netscrape Go module.
package dgraph

import (
	"context"

	"github.com/milosgajdos/netscrape/pkg/graph"
	"github.com/milosgajdos/netscrape/pkg/store"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// Store is dgraph store
type Store struct {
	g *Graph
}

// New creates new dgraph store and returns it.
func NewStore(c *Client) (*Store, error) {
	g, err := NewGraph(c)
	if err != nil {
		return nil, err
	}

	return &Store{
		g: g,
	}, nil
}

// Graph returns graph handle.
func (s *Store) Graph(ctx context.Context) (graph.Graph, error) {
	return s.g, nil
}

// Add Entity to store.
func (s *Store) Add(ctx context.Context, e Entity, opts ...store.Option) error {
	return store.ErrNotImplemented
}

// Link two entities in store.
func (s *Store) Link(ctx context.Context, from, to uuid.UID, opts ...store.Option) error {
	return store.ErrNotImplemented
}

// Delete Entity from store.
func (s *Store) Delete(ctx context.Context, e Entity, opts ...store.Option) error {
	return store.ErrNotImplemented
}

// Unlink two entities in store.
func (s *Store) Unlink(ctx context.Context, from, to uuid.UID, opts ...store.Option) error {
	return store.ErrNotImplemented
}
