package dgraph

import (
	"context"

	dgapi "github.com/dgraph-io/dgo/v200/protos/api"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

func (s *Store) getRequest(ctx context.Context, uid uuid.UID) (*dgapi.Request, error) {
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
