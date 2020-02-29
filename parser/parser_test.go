package parser

import (
	"reflect"
	"testing"
)

func TestParser(t *testing.T) {
	parser := New()
	for _, cas := range []struct {
		pkg string
	}{
		{"errors"},
		{"github.com/pkg/errors"},
		// module内的包（不在GOROOT或GOPATH），需要1.14及以上版本才能使用
		// {"github.com/donnol/gen/parser"},
		// {"github.com/donnol/gen/parser/testdata"},
	} {
		r, err := parser.ParsePkg(cas.pkg)
		if err != nil {
			t.Fatalf("%+v\n", err)
		}
		t.Logf("r: %+v\n", r)
	}
}

func TestParseGenCommand(t *testing.T) {
	for _, cas := range []struct {
		line string
		want []Command
	}{
		{"@gen list",
			[]Command{
				{"list", ""},
			}},
		{"@gen list column",
			[]Command{
				{"list", "column"},
			}},
		{"@gen list map",
			[]Command{
				{"list", "map"},
			}},
		{"@gen list slicemap",
			[]Command{
				{"list", "slicemap"},
			}},
		{"@gen list [map, slicemap]",
			[]Command{
				{"list", "map"}, {"list", "slicemap"},
			}},
	} {
		r := parseGenCommand(cas.line)
		if !reflect.DeepEqual(r, cas.want) {
			t.Fatalf("Bad result: %+v != %+v\n", r, cas.want)
		}
	}
}
