package storage_test

import (
	"reflect"
	"testing"

	"github.com/lucapette/deloominator/pkg/db"
	"github.com/lucapette/deloominator/pkg/testutil"
)

func TestStorage_Question_InsertAndFind(t *testing.T) {
	type args struct{ title, query string }

	sources := testutil.CreateDSN(t)
	for _, source := range sources {
		s := testutil.NewStorage(t, source)
		defer s.Close()

		ds, err := db.NewDataSource(source)
		if err != nil {
			t.Errorf("could not create data source: %v", err)
		}
		defer ds.Close()

		insertQuestion, err := s.InsertQuestion("New Question", `select 42 from answer`)
		if err != nil {
			t.Errorf("expected to insert question, but did not: %v", err)
		}

		findQuestion, err := s.FindQuestion(insertQuestion.ID)
		if err != nil {
			t.Fatalf("expected to insert question with id=%s, but did not: %v", insertQuestion.ID, err)
		}

		if !reflect.DeepEqual(insertQuestion, findQuestion) {
			t.Fatalf("unexpected result, diff: %v", testutil.Diff(insertQuestion, findQuestion))
		}
	}
}
