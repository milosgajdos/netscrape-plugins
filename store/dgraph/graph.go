package dgraph

import (
	"context"

	"github.com/milosgajdos/netscrape/pkg/graph"
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// Graph is dgraph graph
type Graph struct {
	uid    uuid.UID
	client *Client
}

// NewGraph returns dgraph graph
func NewGraph(c *Client) (*Graph, error) {
	return &Graph{
		client: c,
	}, nil
}

// UID returns graph uid.
func (g Graph) UID() uuid.UID {
	return g.uid
}

// NewNode returns a new node.
func (g *Graph) NewNode(ctx context.Context, e space.Entity, opts ...Option) (graph.Node, error) {
	return nil, graph.ErrNotImplemented
}

// AddNode adds a new node to graph.
func (g *Graph) AddNode(ctx context.Context, n graph.Node) error {
	return graph.ErrNotImplemented
}

// Node returns node with given uid.
func (g Graph) Node(ctx context.Context, uid uuid.UID) (graph.Node, error) {
	return nil, graph.ErrNotImplemented
}

// Nodes returns all graph nodes.
func (g Graph) Nodes(ctx context.Context) ([]graph.Node, error) {
	return nil, graph.ErrNotImplemented
}

// RemoveNode removes node from graph.
func (g *Graph) RemoveNode(ctx context.Context, uid uuid.UID) error {
	return graph.ErrNotImplemented
}

// Link links two nodes and returns the new edge.
func (g *Graph) Link(ctx context.Context, from, to uuid.UID, opts ...Option) (graph.Edge, error) {
	return nil, graph.ErrNotImplemented
}

// Unlink removes link from graph.
func (g *Graph) Unlink(ctx context.Context, from, to uuid.UID) error {
	return graph.ErrNotImplemented
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
