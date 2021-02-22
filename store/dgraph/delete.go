package dgraph

import (
	dgapi "github.com/dgraph-io/dgo/v200/protos/api"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

func (s *Store) delRequest(uid uuid.UID) (*dgapi.Request, error) {
	q := `
	{
		node(func: eq(xid, "` + uid.Value() + `")) @filter(NOT type(Resource) OR eq(count(~resource), 0)) {
			u as uid
		}
	}
	`

	node := map[string]string{"uid": "uid(u)"}

	cond := `@if(gt(len(u), 0))`

	return UpsertReqJSON(DelOp, node, q, cond)
}
