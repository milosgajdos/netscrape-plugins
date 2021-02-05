package dgraph

import (
	"context"

	dgo "github.com/dgraph-io/dgo/v200"
	"github.com/dgraph-io/dgo/v200/protos/api"
	"google.golang.org/grpc"
)

// Auth for dgraph authentication
type Auth struct {
	User   string
	Passwd string
}

// Client is dgraph client
type Client struct {
	*dgo.Dgraph
	conn *grpc.ClientConn
}

// NewClient creates new dgraph client and returns it.
func NewClient(target string, opts ...Option) (*Client, error) {
	dopts := Options{}
	for _, apply := range opts {
		apply(&dopts)
	}

	conn, err := grpc.Dial(target, dopts.DialOpts...)
	if err != nil {
		return nil, err
	}

	dg := dgo.NewDgraphClient(api.NewDgraphClient(conn))

	ctx := context.Background()

	if dopts.Auth != nil {
		err = dg.Login(ctx, dopts.Auth.User, dopts.Auth.Passwd)
		if err != nil {
			return nil, err
		}
		// TODO: implement retry
		//if err == nil || !strings.Contains(err.Error(), "Please retry") {
		//	break
		//}
	}

	return &Client{
		Dgraph: dg,
		conn:   conn,
	}, nil
}

// Close closes dgraph connection
func (c *Client) Close() error {
	return c.conn.Close()
}
