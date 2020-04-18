package resolver

import (
	"context"
	"io/ioutil"

	"github.com/googleapis/google-cloud-go-testing/bigquery/bqiface"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
)

type Query struct {
	Name  string `bigquery:"name"`
	Query string `bigquery:"query"`
}

type QueryResolver interface {
	Resolve(ctx context.Context, templateFilePath string) ([]Query, error)
}

type queryResolverImpl struct {
	bqClient bqiface.Client
}

func NewQueryResolver(bqClient bqiface.Client) QueryResolver {
	return &queryResolverImpl{
		bqClient: bqClient,
	}
}

func (qr queryResolverImpl) Resolve(ctx context.Context, templateFilePath string) ([]Query, error) {
	var templateFile string
	{
		b, err := ioutil.ReadFile(templateFilePath)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		templateFile = string(b)
	}

	q := qr.bqClient.Query(templateFile)
	rowIterator, err := q.Read(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	queries := []Query{}
	for {
		res := Query{}
		err := rowIterator.Next(&res)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, errors.WithStack(err)
		}
		queries = append(queries, res)
	}

	return queries, nil
}
