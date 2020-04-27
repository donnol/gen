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
	opt := Option{UseAnnotation: true}
	list := New(p, temp, opt)
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

	mr := ml.MapByName()
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

	wherer := ml.Where(func(m testdata1.Model) bool {
		if m.Name == "jd" {
			return false
		}
		return true
	})
	if len(wherer) != 0 {
		t.Fatalf("len isn't 0, %+v\n", wherer)
	}

	ml1 := testdata1.ModelList{
		{ID: 1, Name: "jd", Age: 31.0, UserID: 1, ContentID: 1, ModelID: 1, AddrID: 1},
		{ID: 2, Name: "je", Age: 21.0, UserID: 2, ContentID: 2, ModelID: 2, AddrID: 2},
	}
	ml1.Sort(func(i, j int) bool {
		return ml1[i].ID > ml1[j].ID
	})
	sortwant := testdata1.ModelList{
		{ID: 2, Name: "je", Age: 21.0, UserID: 2, ContentID: 2, ModelID: 2, AddrID: 2},
		{ID: 1, Name: "jd", Age: 31.0, UserID: 1, ContentID: 1, ModelID: 1, AddrID: 1},
	}
	if !reflect.DeepEqual(ml1, sortwant) {
		t.Fatalf("Bad result: %+v != %+v\n", ml1, sortwant)
	}

	ml1limit := ml1.Limit(0, 1)
	limitwant := testdata1.ModelList{
		{ID: 2, Name: "je", Age: 21.0, UserID: 2, ContentID: 2, ModelID: 2, AddrID: 2},
	}
	if !reflect.DeepEqual(ml1limit, limitwant) {
		t.Fatalf("Bad result: %+v != %+v\n", ml1limit, limitwant)
	}

	ml1limit = ml1.Limit(1, 1)
	limitwant = testdata1.ModelList{
		{ID: 1, Name: "jd", Age: 31.0, UserID: 1, ContentID: 1, ModelID: 1, AddrID: 1},
	}
	if !reflect.DeepEqual(ml1limit, limitwant) {
		t.Fatalf("Bad result: %+v != %+v\n", ml1limit, limitwant)
	}

	ml1limit = ml1.Limit(1, 2)
	limitwant = testdata1.ModelList{
		{ID: 1, Name: "jd", Age: 31.0, UserID: 1, ContentID: 1, ModelID: 1, AddrID: 1},
	}
	if !reflect.DeepEqual(ml1limit, limitwant) {
		t.Fatalf("Bad result: %+v != %+v\n", ml1limit, limitwant)
	}

	ml1reduce := ml1.Reduce(func(u testdata1.Model, nu testdata1.Model) testdata1.Model {
		return testdata1.Model{
			ID: u.ID + nu.ID,
		}
	})
	reducewant := testdata1.Model{ID: 3, Name: "", Age: 0.0, UserID: 0, ContentID: 0, ModelID: 0, AddrID: 0}
	if !reflect.DeepEqual(ml1reduce, reducewant) {
		t.Fatalf("Bad result: %+v != %+v\n", ml1reduce, reducewant)
	}

	ml1reverse := ml1.Reverse()
	reversewant := testdata1.ModelList{
		{ID: 1, Name: "jd", Age: 31.0, UserID: 1, ContentID: 1, ModelID: 1, AddrID: 1},
		{ID: 2, Name: "je", Age: 21.0, UserID: 2, ContentID: 2, ModelID: 2, AddrID: 2},
	}
	if !reflect.DeepEqual(ml1reverse, reversewant) {
		t.Fatalf("Bad result: %+v != %+v\n", ml1reverse, reversewant)
	}

	ml1first := ml1.First()
	firstwant := testdata1.Model{ID: 2, Name: "je", Age: 21.0, UserID: 2, ContentID: 2, ModelID: 2, AddrID: 2}
	if !reflect.DeepEqual(ml1first, firstwant) {
		t.Fatalf("Bad result: %+v != %+v\n", ml1first, firstwant)
	}

	ml1last := ml1.Last()
	lastwant := testdata1.Model{ID: 1, Name: "jd", Age: 31.0, UserID: 1, ContentID: 1, ModelID: 1, AddrID: 1}
	if !reflect.DeepEqual(ml1last, lastwant) {
		t.Fatalf("Bad result: %+v != %+v\n", ml1last, lastwant)
	}
}
