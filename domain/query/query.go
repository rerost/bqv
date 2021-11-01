package query

import (
	"context"

	"github.com/googleapis/google-cloud-go-testing/bigquery/bqiface"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type QueryService interface {
	Exec(ctx context.Context, query string) error
	BulkExec(ctx context.Context, queries []string) error
}

type queryServiceImpl struct {
	bqClient bqiface.Client
}

func NewQueryService(bqClient bqiface.Client) QueryService {
	return &queryServiceImpl{
		bqClient: bqClient,
	}
}

func (q *queryServiceImpl) Exec(ctx context.Context, query string) (err error) {
	j, err := q.bqClient.Query(query).Run(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	zap.L().Debug("Exec query", zap.String("job_id", j.ID()), zap.String("query", query))

	status, err := j.Wait(ctx)
	zap.L().Debug("End query", zap.String("job_id", j.ID()))
	if err != nil {
		zap.L().Debug("Wait err", zap.String("job_id", j.ID()), zap.String("query", query))
		return errors.WithStack(err)
	}
	if err := status.Err(); err != nil {
		zap.L().Debug("Status err", zap.String("job_id", j.ID()), zap.String("query", query))
		return errors.WithStack(err)
	}

	return nil
}

func (q *queryServiceImpl) BulkExec(ctx context.Context, queries []string) error {
	var eg errgroup.Group

	for _, query := range queries {
		query := query
		eg.Go(func() error {
			return errors.WithStack(q.Exec(ctx, query))
		})
	}

	if err := eg.Wait(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
