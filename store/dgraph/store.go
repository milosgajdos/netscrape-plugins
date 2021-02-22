// package dgraph implements store.Store interface
// from github.com/milosgajdos/netscrape Go module.
package dgraph

import (
	"context"
	"fmt"
	"html/template"

	"github.com/milosgajdos/netscrape/pkg/store"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// Store is dgraph store
type Store struct {
	c *Client
}

// New creates new dgraph store and returns it.
func NewStore(c *Client, opts ...Option) (*Store, error) {
	sopts := Options{}
	for _, apply := range opts {
		apply(&sopts)
	}

	return &Store{
		c: c,
	}, nil
}

// Add Entity to store.
func (s *Store) Add(ctx context.Context, e store.Entity, opts ...store.Option) error {
	sopts := store.Options{}
	for _, apply := range opts {
		apply(&sopts)
	}

	req, err := s.addRequest(e)
	if err != nil {
		return err
	}

	if _, err := s.c.NewTxn().Do(ctx, req); err != nil {
		return fmt.Errorf("txn Add: %w", err)
	}

	return nil
}

// Delete Entity from store.
func (s *Store) Delete(ctx context.Context, uid uuid.UID, opts ...store.Option) error {
	sopts := store.Options{}
	for _, apply := range opts {
		apply(&sopts)
	}

	req, err := s.delRequest(uid)
	if err != nil {
		return err
	}

	if _, err := s.c.NewTxn().Do(ctx, req); err != nil {
		return fmt.Errorf("txn Delete: %w", err)
	}

	return nil
}

// Link two entities in store.
func (s *Store) Link(ctx context.Context, from, to uuid.UID, opts ...store.Option) error {
	req, err := s.linkRequest(from, to, opts...)
	if err != nil {
		return err
	}

	if _, err := s.c.NewTxn().Do(ctx, req); err != nil {
		return fmt.Errorf("txn Link: %w", err)
	}

	return nil
}

// Unlink two entities in store.
func (s *Store) Unlink(ctx context.Context, from, to uuid.UID, opts ...store.Option) error {
	req, err := s.unlinkRequest(from, to, opts...)
	if err != nil {
		return err
	}

	if _, err := s.c.NewTxn().Do(ctx, req); err != nil {
		return fmt.Errorf("txn Unlink: %w", err)
	}

	return nil
}

// Query returns entity matching the query executed with the given vars.
func (s *Store) Query(ctx context.Context, queryTpl template.Template, vars map[string]string) ([]store.Entity, error) {
	return nil, store.ErrNotImplemented
}
