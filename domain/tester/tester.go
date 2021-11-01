package tester

import (
	"context"
	"fmt"
	"regexp"
	"strings"

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

	targetTableName string
	tmpTableName    string
	withClause      *regexp.Regexp
}

func NewTestService(queryService query.QueryService) TestService {
	withClause, _ := regexp.Compile("(WITH|with) ([a-zA-Z]+[a-zA-Z0-9]*) (AS|as)")

	return &testServiceImpl{
		queryService:    queryService,
		targetTableName: "$TARGET",
		tmpTableName:    "bqv_testing_table",
		withClause:      withClause,
	}
}

func (t *testServiceImpl) Test(ctx context.Context, viewQuery string, assertQuery string) error {
	if res := t.withClause.Find([]byte(assertQuery)); res != nil {
		return errors.Wrapf(
			TestService_IncludeWithError,
			"Detected with clause: %s",
			res,
		)
	}

	err := t.queryService.Exec(ctx, t.testQuery(viewQuery, assertQuery))
	fmt.Println(err)
	return nil
}

func (t *testServiceImpl) testQuery(viewQuery string, assertQuery string) string {
	sql := `
WITH %s AS (
  %s
)

%s
  `

	testQuery := fmt.Sprintf(
		sql,
		t.tmpTableName,
		viewQuery,
		assertQuery,
		strings.ReplaceAll(assertQuery, t.targetTableName, t.tmpTableName),
	)

	return testQuery
}
