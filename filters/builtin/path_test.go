package builtin

import (
	"net/http"
	"testing"

	"github.com/ardaguclu/skipper/filters"
	"github.com/ardaguclu/skipper/filters/filtertest"
)

type createTestItem struct {
	msg  string
	args []interface{}
	err  bool
}

func TestModifyPath(t *testing.T) {
	spec := NewModPath()
	f, err := spec.CreateFilter([]interface{}{"/replace-this/", "/with-this/"})
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest("GET", "https://www.example.org/path/replace-this/yo", nil)
	if err != nil {
		t.Error(err)
	}

	ctx := &filtertest.Context{FRequest: req}
	f.Request(ctx)
	if req.URL.Path != "/path/with-this/yo" {
		t.Error("failed to replace path")
	}
}

func TestModifyPathWithInvalidExpression(t *testing.T) {
	spec := NewModPath()
	if f, err := spec.CreateFilter([]interface{}{"(?=;)", "foo"}); err == nil || f != nil {
		t.Error("Expected error for invalid regular expression parameter")
	}
}

func TestSetPath(t *testing.T) {
	spec := NewSetPath()
	f, err := spec.CreateFilter([]interface{}{"/baz/qux"})
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest("GET", "https://www.example.org/foo/bar", nil)
	if err != nil {
		t.Error(err)
	}

	ctx := &filtertest.Context{FRequest: req}
	f.Request(ctx)
	if req.URL.Path != "/baz/qux" {
		t.Error("failed to replace path")
	}
}

func TestSetPathWithTemplate(t *testing.T) {
	spec := NewSetPath()
	f, err := spec.CreateFilter([]interface{}{"/path/${param2}/${param1}"})
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest("GET", "https://www.example.org/foo/bar", nil)
	if err != nil {
		t.Error(err)
	}

	ctx := &filtertest.Context{FRequest: req, FParams: map[string]string{
		"param1": "foo",
		"param2": "bar",
	}}

	f.Request(ctx)
	if req.URL.Path != "/path/bar/foo" {
		t.Error("failed to transform path")
	}
}

func testCreate(t *testing.T, spec func() filters.Spec, items []createTestItem) {
	for _, ti := range items {
		func() {
			f, err := spec().CreateFilter(ti.args)
			switch {
			case ti.err && err == nil:
				t.Error(ti.msg, "failed to fail")
			case !ti.err && err != nil:
				t.Error(ti.msg, err)
			case err == nil && f == nil:
				t.Error(ti.msg, "failed to create filter")
			}
		}()
	}
}

func TestCreateModPath(t *testing.T) {
	testCreate(t, NewModPath, []createTestItem{{
		"no args",
		nil,
		true,
	}, {
		"single arg",
		[]interface{}{".*"},
		true,
	}, {
		"non-string arg, pos 1",
		[]interface{}{3.14, "/foo"},
		true,
	}, {
		"non-string arg, pos 2",
		[]interface{}{".*", 2.72},
		true,
	}, {
		"more than two args",
		[]interface{}{".*", "/foo", "/bar"},
		true,
	}, {
		"create",
		[]interface{}{".*", "/foo"},
		false,
	}})
}

func TestCreateSetPath(t *testing.T) {
	testCreate(t, NewSetPath, []createTestItem{{
		"no args",
		nil,
		true,
	}, {
		"non-string arg",
		[]interface{}{3.14},
		true,
	}, {
		"more than one args",
		[]interface{}{"/foo", "/bar"},
		true,
	}, {
		"create",
		[]interface{}{"/foo"},
		false,
	}})
}
