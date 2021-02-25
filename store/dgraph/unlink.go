package dgraph

import (
	"context"

	dgapi "github.com/dgraph-io/dgo/v200/protos/api"
	"github.com/milosgajdos/netscrape/pkg/store"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// unlinkRequest link from and to Objects. No other types can be linked.
func (s *Store) unlinkRequest(ctx context.Context, from, to uuid.UID, opts ...store.Option) (*dgapi.Request, error) {
	q := `
	{
		var(func: eq(xid, "` + from.Value() + `")) @filter(type(Object)) {
			from as uid
		}

		var(func: eq(xid, "` + to.Value() + `")) @filter(type(Object)) {
			to as uid
		}
	}
	`

	sopts := store.Options{}
	for _, apply := range opts {
		apply(&sopts)
	}

	link := &Object{
		UID: "uid(from)",
		Links: []Object{
			{UID: "uid(to)"},
		},
	}

	cond := `@if(gt(len(from), 0) AND gt(len(to), 0))`

	return UpsertReqJSON(UnlinkOp, link, q, cond)
}
