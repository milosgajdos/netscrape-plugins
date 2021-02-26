package dgraph

import (
	"context"
	"flag"
	"reflect"
	"testing"

	dgapi "github.com/dgraph-io/dgo/v200/protos/api"
	"github.com/milosgajdos/netscrape/pkg/store"
	"github.com/milosgajdos/netscrape/pkg/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding/gzip"
)

var (
	host = flag.String("host", "localhost:9080", "DGrapg host")
	drop = flag.Bool("drop", false, "drop all data including schema")
)

func MustNewStore(dsn string, drop bool, t *testing.T) *Store {
	dialOpts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)),
	}

	s, err := NewStore(dsn, WithDialOpts(dialOpts...))
	if err != nil {
		t.Fatal(err)
	}

	if drop {
		op := &dgapi.Operation{
			DropOp: dgapi.Operation_ALL,
		}

		if err := s.Alter(context.Background(), op); err != nil {
			t.Fatal(err)
			defer s.Close()
		}
	}

	op := &dgapi.Operation{
		Schema: SpaceDQLSchema,
	}

	if err := s.Alter(context.Background(), op); err != nil {
		t.Fatal(err)
		defer s.Close()
	}

	return s
}

func TestAdd(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		s := MustNewStore(*host, *drop, t)
		defer s.Close()

		obj, err := newTestObject("ent1", "entNs")
		if err != nil {
			t.Fatal(err)
		}

		if err := s.Add(context.Background(), obj); err != nil {
			t.Fatal(err)
		}

		res, err := newTestResource(resName, resGroup, resVersion, resKind, true)
		if err != nil {
			t.Fatal(err)
		}

		if err := s.Add(context.Background(), res); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("ErrUnsupported", func(t *testing.T) {
		s := MustNewStore(*host, *drop, t)
		defer s.Close()

		if err := s.Add(context.Background(), nil); err != store.ErrUnsupported {
			t.Fatalf("got: %v, want: %v", err, store.ErrUnsupported)
		}
	})
}

func TestGet(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		s := MustNewStore(*host, *drop, t)
		defer s.Close()

		obj, err := newTestObject("ent1", "entNs")
		if err != nil {
			t.Fatal(err)
		}

		if err := s.Add(context.Background(), obj); err != nil {
			t.Fatal(err)
		}

		e, err := s.Get(context.Background(), obj.UID())
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(e.UID(), obj.UID()) {
			t.Fatalf("got: %s, want: %s", e.UID(), obj.UID())
		}
	})

	t.Run("ErrNodeNotFound", func(t *testing.T) {
		s := MustNewStore(*host, *drop, t)
		defer s.Close()

		uid, err := uuid.New()
		if err != nil {
			t.Fatal(err)
		}

		if _, err := s.Get(context.Background(), uid); err != store.ErrNodeNotFound {
			t.Fatalf("got: %v, want: %v", err, store.ErrNodeNotFound)
		}
	})
}

func TestDelete(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		s := MustNewStore(*host, *drop, t)
		defer s.Close()

		obj, err := newTestObject("ent1", "entNs")
		if err != nil {
			t.Fatal(err)
		}

		if err := s.Add(context.Background(), obj); err != nil {
			t.Fatal(err)
		}

		if err := s.Delete(context.Background(), obj.UID()); err != nil {
			t.Fatal(err)
		}

		if _, err := s.Get(context.Background(), obj.UID()); err != store.ErrNodeNotFound {
			t.Fatalf("got: %v, want: %v", err, store.ErrNodeNotFound)
		}
	})

	t.Run("DeleteNonExistant", func(t *testing.T) {
		s := MustNewStore(*host, *drop, t)
		defer s.Close()

		uid, err := uuid.New()
		if err != nil {
			t.Fatal(err)
		}

		if err := s.Delete(context.Background(), uid); err != nil {
			t.Fatal(err)
		}
	})
}

func TestLink(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		s := MustNewStore(*host, *drop, t)
		defer s.Close()

		obj1, err := newTestObject("ent1", "entNs")
		if err != nil {
			t.Fatal(err)
		}

		if err := s.Add(context.Background(), obj1); err != nil {
			t.Fatal(err)
		}

		obj2, err := newTestObject("ent2", "entNs")
		if err != nil {
			t.Fatal(err)
		}

		if err := s.Add(context.Background(), obj2); err != nil {
			t.Fatal(err)
		}

		if err := s.Link(context.Background(), obj1.UID(), obj2.UID()); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("LinkNonExistant", func(t *testing.T) {
		s := MustNewStore(*host, *drop, t)
		defer s.Close()

		uid, err := uuid.New()
		if err != nil {
			t.Fatal(err)
		}

		obj1, err := newTestObject("ent1", "entNs")
		if err != nil {
			t.Fatal(err)
		}

		if err := s.Add(context.Background(), obj1); err != nil {
			t.Fatal(err)
		}

		if err := s.Link(context.Background(), obj1.UID(), uid); err != nil {
			t.Fatal(err)
		}
	})
}

func TestUnlink(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		s := MustNewStore(*host, *drop, t)
		defer s.Close()

		obj1, err := newTestObject("ent1", "entNs")
		if err != nil {
			t.Fatal(err)
		}

		if err := s.Add(context.Background(), obj1); err != nil {
			t.Fatal(err)
		}

		obj2, err := newTestObject("ent2", "entNs")
		if err != nil {
			t.Fatal(err)
		}

		if err := s.Add(context.Background(), obj2); err != nil {
			t.Fatal(err)
		}

		if err := s.Link(context.Background(), obj1.UID(), obj2.UID()); err != nil {
			t.Fatal(err)
		}

		if err := s.Unlink(context.Background(), obj1.UID(), obj2.UID()); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("UnlinkNonExistant", func(t *testing.T) {
		s := MustNewStore(*host, *drop, t)
		defer s.Close()

		uid, err := uuid.New()
		if err != nil {
			t.Fatal(err)
		}

		obj1, err := newTestObject("ent1", "entNs")
		if err != nil {
			t.Fatal(err)
		}

		if err := s.Add(context.Background(), obj1); err != nil {
			t.Fatal(err)
		}

		if err := s.Unlink(context.Background(), obj1.UID(), uid); err != nil {
			t.Fatal(err)
		}
	})
}
