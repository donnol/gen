package parser

import "testing"

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
