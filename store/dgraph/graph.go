package dgraph

import (
	"context"

	"github.com/milosgajdos/netscrape/pkg/graph"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// Graph is dgraph graph
type Graph struct {
	uid    uuid.UID
	client *Client
}

// NewGraph returns dgraph graph
func NewGraph(c *Client, opts ...Option) (*Graph, error) {
	gopts := Options{}
	for _, apply := range opts {
		apply(&gopts)
	}

	uid := gopts.UID
	if uid == nil {
		var err error
		uid, err = uuid.New()
		if err != nil {
			return nil, err
		}
	}

	return &Graph{
		uid:    uid,
		client: c,
	}, nil
}

// UID returns graph uid.
func (g Graph) UID() uuid.UID {
	return g.uid
}

// Node returns node with given uid.
func (g Graph) Node(ctx context.Context, uid uuid.UID) (graph.Node, error) {
	return nil, graph.ErrNotImplemented
}

// Nodes returns all graph nodes.
func (g Graph) Nodes(ctx context.Context) ([]graph.Node, error) {
	return nil, graph.ErrNotImplemented
}

// Edge returns the edge between from and to nodes.
func (g Graph) Edge(ctx context.Context, from, to uuid.UID) (graph.Edge, error) {
	return nil, graph.ErrNotImplemented
}

// Edges returns all graph edges.
func (g Graph) Edges(ctx context.Context) ([]graph.Edge, error) {
	return nil, graph.ErrNotImplemented
}

// SubGraph returns the maximum subgraph of a graph
// starting at node with given uid up to given depth.
func (g Graph) SubGraph(ctx context.Context, uid uuid.UID, depth int, opts ...Option) (graph.Graph, error) {
	return nil, graph.ErrNotImplemented
}
