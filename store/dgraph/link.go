package dgraph

import (
	"strconv"

	dgapi "github.com/dgraph-io/dgo/v200/protos/api"
	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/store"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

const (
	DefaultRelation = "Unknown"
	DefaultWeight   = 1.0
)

// linkRequest link from and to Objects. No other types can be linked.
func (s *Store) linkRequest(from, to uuid.UID, opts ...store.Option) (*dgapi.Request, error) {
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

	weight := DefaultWeight
	relation := DefaultRelation

	if sopts.Attrs != nil {
		if w, err := strconv.ParseFloat(sopts.Attrs.Get(attrs.Weight), 64); err == nil {
			if w != 0.0 {
				weight = w
			}
		}

		if r := sopts.Attrs.Get(attrs.Relation); r != "" {
			relation = r
		}
	}

	link := &Object{
		UID:   "uid(from)",
		DType: []string{"Object"},
		Links: []Object{
			{UID: "uid(to)", DType: []string{"Object"}, Relation: relation, Weight: weight},
		},
	}

	cond := `@if(gt(len(from), 0) AND gt(len(to), 0))`

	return UpsertReqJSON(LinkOp, link, q, cond)
}
