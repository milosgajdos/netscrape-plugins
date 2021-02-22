package dgraph

import (
	dgapi "github.com/dgraph-io/dgo/v200/protos/api"
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/store"
)

func (s *Store) addRequest(e store.Entity) (*dgapi.Request, error) {
	switch v := e.(type) {
	case space.Object:
		return s.addObjectRequest(v)
	case space.Resource:
		return s.addResourceRequest(v)
	default:
		return nil, store.ErrUnsupported
	}
}

func (s *Store) addResourceRequest(r space.Resource) (*dgapi.Request, error) {
	query := `
	{
		resource(func: eq(xid, "` + r.UID().Value() + `")) {
			u as uid
	        }
	}
	`

	res := &Resource{
		UID:        "uid(u)",
		XID:        r.UID().Value(),
		Name:       r.Name(),
		Group:      r.Group(),
		Version:    r.Version(),
		Kind:       r.Kind(),
		Namespaced: r.Namespaced(),
		Attrs:      AttrsToMap(r.Attrs()),
		DType:      []string{"Resource"},
	}

	return UpsertReqJSON(AddOp, res, query, "")
}

func (s *Store) addObjectRequest(o space.Object) (*dgapi.Request, error) {
	query := `
	{
		object(func: eq(xid, "` + o.UID().Value() + `")) {
			o as uid
		}

		resource(func: eq(xid, "` + o.Resource().UID().Value() + `")) {
			r as uid
		}
	}
	`

	obj := &Object{
		UID:       "uid(o)",
		XID:       o.UID().Value(),
		Name:      o.Name(),
		Namespace: o.Namespace(),
		Resource: &Resource{
			UID:        "uid(r)",
			XID:        o.Resource().UID().Value(),
			Name:       o.Resource().Name(),
			Group:      o.Resource().Group(),
			Version:    o.Resource().Version(),
			Kind:       o.Resource().Kind(),
			Namespaced: o.Resource().Namespaced(),
			Attrs:      AttrsToMap(o.Resource().Attrs()),
			DType:      []string{"Resource"},
		},
		Attrs: AttrsToMap(o.Attrs()),
		DType: []string{"Object"},
	}

	return UpsertReqJSON(AddOp, obj, query, "")
}
