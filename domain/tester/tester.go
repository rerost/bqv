package tester

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/pkg/errors"
	"github.com/rerost/bqv/domain/query"
)

var (
	TestService_IncludeWithError = errors.New("Assert query includes with clause")
)

type TestService interface {
	Test(ctx context.Context, viewQuery string, assertQuery string) error
}

type testServiceImpl struct {
	queryService query.QueryService

	targetViewName string
	tmpTableName   string
}

func NewTestService(queryService query.QueryService) TestService {
	return &testServiceImpl{
		queryService:   queryService,
		targetViewName: "$TARGET",
		tmpTableName:   "bqv_testing_table",
	}
}

func (t *testServiceImpl) Test(ctx context.Context, viewQuery string, assertQuery string) error {
	err := t.queryService.Exec(ctx, t.testQuery(viewQuery, assertQuery))
	return err
}

func (t *testServiceImpl) testQuery(viewQuery string, assertQuery string) string {
	tmpl, _ := template.New("query").Parse(`
			CREATE TEMP TABLE {{ .TmpTableName }} AS (
			  {{ .ViewQuery }}
			);

			{{ .AssertQuery }}
			  `)

	b := bytes.NewBuffer([]byte{})
	tmpl.Execute(
		b,
		struct {
			TmpTableName string
			ViewQuery    string
			AssertQuery  string
		}{
			TmpTableName: t.tmpTableName,
			ViewQuery:    viewQuery,
			AssertQuery:  assertQuery,
		})
	testQuery := b.String()
	fmt.Println(testQuery)

	return testQuery
}
