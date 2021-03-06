package dgraph

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/space/entity"
	"github.com/milosgajdos/netscrape/pkg/space/resource"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

const (
	testDir = "testdata"

	resUID     = "nodeResUID"
	resName    = "nodeResName"
	resGroup   = "nodeResGroup"
	resVersion = "nodeResVersion"
	resKind    = "nodeResKind"
)

func newTestResource(name, group, version, kind string, namespaced bool, opts ...resource.Option) (space.Resource, error) {
	uid, err := uuid.NewFromString(resUID)
	if err != nil {
		return nil, err
	}

	return resource.New(name, group, version, kind, namespaced, resource.WithUID(uid))
}

func newTestEntity(entName, entNs string) (space.Entity, error) {
	r, err := newTestResource(resName, resGroup, resVersion, resKind, true)
	if err != nil {
		return nil, err
	}

	uid, err := uuid.NewFromString(entName + "/" + entNs)
	if err != nil {
		return nil, err
	}

	return entity.New(entName, entNs, r, entity.WithUID(uid))
}

func TestAttrsToMap(t *testing.T) {
	a, err := attrs.New()
	if err != nil {
		t.Fatalf("failed to create new attrs: %v", err)
	}

	aMap := AttrsToMap(a)

	for k, v := range aMap {
		if val := a.Get(k); val != v {
			t.Errorf("expected attr val: %s, got: %s", v, val)
		}
	}
}

func TestContains(t *testing.T) {
	testCases := []struct {
		a   []string
		x   string
		exp bool
	}{
		{[]string{"a", "b", "c"}, "a", true},
		{[]string{}, "a", false},
		{[]string{"a", "a", "c"}, "a", true},
	}

	for _, tc := range testCases {
		if c := contains(tc.a, tc.x); c != tc.exp {
			t.Errorf("expected: %v, got: %v", tc.exp, c)
		}
	}
}

func TestDecodeJSONEntity(t *testing.T) {
	testFiles := []string{"entity.json", "resource.json"}

	for _, f := range testFiles {
		path := path.Join(testDir, f)

		data, err := ioutil.ReadFile(path)
		if err != nil {
			t.Fatalf("failed opening file: %v", err)
		}

		ents, err := decodeJSONEntity(data, GetOp)
		if err != nil {
			t.Fatalf("failed decoding data: %v", err)
		}

		if len(ents) == 0 {
			t.Errorf("ents count: %d", len(ents))
		}
	}
}
