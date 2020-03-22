package list

import (
	"testing"

	"github.com/donnol/gen/parser"
	"github.com/donnol/gen/template"
)

func TestList(t *testing.T) {
	p := parser.New()
	temp := &template.Template{}
	list := New(p, temp)
	for _, cas := range []struct {
		pkg string
	}{
		{"github.com/donnol/gen/parser"},
		{"github.com/donnol/gen/list/testdata1"},
	} {
		err := list.Parse(cas.pkg)
		if err != nil {
			t.Fatalf("%+v\n", err)
		}
	}
}
