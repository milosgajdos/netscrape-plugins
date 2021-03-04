// package dgraph implements store.Store interface
// from github.com/milosgajdos/netscrape Go module.
package dgraph

import (
	"context"
	"fmt"

	"github.com/milosgajdos/netscrape/pkg/store"
	"github.com/milosgajdos/netscrape/pkg/uuid"

	dgapi "github.com/dgraph-io/dgo/v200/protos/api"
)

// Store is dgraph store
type Store struct {
	c *Client
}

// New creates new dgraph store and returns it.
func NewStore(dsn string, opts ...Option) (*Store, error) {
	sopts := Options{}
	for _, apply := range opts {
		apply(&sopts)
	}

	c, err := NewClient(dsn, opts...)
	if err != nil {
		return nil, err
	}

	return &Store{
		c: c,
	}, nil
}

// Alter alters dgraph database with the given operation.
func (s *Store) Alter(ctx context.Context, op *dgapi.Operation) error {
	return s.c.Alter(ctx, op)
}

// Close closes store.
func (s *Store) Close() error {
	return s.c.Close()
}

// Add Entity to store.
func (s *Store) Add(ctx context.Context, e store.Entity, opts ...store.Option) error {
	req, err := s.addRequest(ctx, e, opts...)
	if err != nil {
		return err
	}

	if _, err := s.c.NewTxn().Do(ctx, req); err != nil {
		return fmt.Errorf("txn.Add: %w", err)
	}

	return nil
}

// Get Entity from store.
func (s *Store) Get(ctx context.Context, uid uuid.UID, opts ...store.Option) (store.Entity, error) {
	req, err := s.getRequest(ctx, uid)
	if err != nil {
		return nil, err
	}

	resp, err := s.c.NewTxn().Do(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("txn.Get: %w", err)
	}

	ents, err := decodeJSONEntity(resp.Json, GetOp)
	if err != nil {
		return nil, err
	}

	if len(ents) == 0 {
		return nil, store.ErrNodeNotFound
	}

	if len(ents) > 2 {
		panic("duplicate nodes")
	}

	return ents[0], nil
}

// Delete Entity from store.
func (s *Store) Delete(ctx context.Context, uid uuid.UID, opts ...store.Option) error {
	req, err := s.deleteRequest(ctx, uid)
	if err != nil {
		return err
	}

	if _, err := s.c.NewTxn().Do(ctx, req); err != nil {
		return fmt.Errorf("txn.Delete: %w", err)
	}

	return nil
}

// Link two entities in store.
func (s *Store) Link(ctx context.Context, from, to uuid.UID, opts ...store.Option) error {
	req, err := s.linkRequest(ctx, from, to, opts...)
	if err != nil {
		return err
	}

	if _, err := s.c.NewTxn().Do(ctx, req); err != nil {
		return fmt.Errorf("txn.Link: %w", err)
	}

	return nil
}

// Unlink two entities in store.
func (s *Store) Unlink(ctx context.Context, from, to uuid.UID, opts ...store.Option) error {
	req, err := s.unlinkRequest(ctx, from, to, opts...)
	if err != nil {
		return err
	}

	if _, err := s.c.NewTxn().Do(ctx, req); err != nil {
		return fmt.Errorf("txn.Unlink: %w", err)
	}

	return nil
}
