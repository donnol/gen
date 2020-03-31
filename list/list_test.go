package list

import (
	"reflect"
	"testing"

	"github.com/donnol/gen/list/testdata1"
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

func TestListGenResult(t *testing.T) {
	ml := testdata1.ModelList{
		{ID: 1, Name: "jd", Age: 31.0, UserID: 1, ContentID: 1, ModelID: 1, AddrID: 1},
	}

	cr := ml.ColumnID()
	cwant := []int{1}
	if !reflect.DeepEqual(cr, cwant) {
		t.Fatalf("Bad result: %+v != %+v\n", cr, cwant)
	}

	mr := ml.MapName()
	mwant := map[string]testdata1.Model{"jd": ml[0]}
	if !reflect.DeepEqual(mr, mwant) {
		t.Fatalf("Bad result: %+v != %+v\n", mr, mwant)
	}

	mlr := ml.MapListByName()
	mlwant := map[string][]testdata1.Model{"jd": ml}
	if !reflect.DeepEqual(mr, mwant) {
		t.Fatalf("Bad result: %+v != %+v\n", mlr, mlwant)
	}

	jur := ml.JoinUserByUserIDEqualID([]testdata1.User{{ID: 1, Name: "jd"}}, func(m testdata1.Model, u testdata1.User) testdata1.Model {
		m.UserName = u.Name
		return m
	})
	jurwant := testdata1.ModelList{
		{ID: 1, Name: "jd", Age: 31.0, UserID: 1, UserName: "jd", ContentID: 1, ModelID: 1, AddrID: 1},
	}
	if !reflect.DeepEqual(mr, mwant) {
		t.Fatalf("Bad result: %+v != %+v\n", jur, jurwant)
	}

	dur := ml.DeriveByUserIDEqualID([]testdata1.User{{ID: 1, Name: "jd"}}, func(m testdata1.Model, u testdata1.User) testdata1.ModelUser {
		return testdata1.ModelUser{
			ModelID:   m.ID,
			ModelName: m.Name,
			UserID:    u.ID,
			UserName:  u.Name,
		}
	})
	durwant := []testdata1.ModelUser{
		{ModelID: 1, ModelName: "jd", UserID: 1, UserName: "jd"},
	}
	if !reflect.DeepEqual(dur, durwant) {
		t.Fatalf("Bad result: %+v != %+v\n", jur, jurwant)
	}
}
