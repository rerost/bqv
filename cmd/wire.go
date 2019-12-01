//+build wireinject

package cmd

import (
	"context"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/google/wire"
	"github.com/googleapis/google-cloud-go-testing/bigquery/bqiface"
	"github.com/pkg/errors"
	"github.com/rerost/bq-table-validator/domain/bqquery"
	"github.com/rerost/bq-table-validator/domain/tablemock"
	"github.com/rerost/bq-table-validator/domain/validator"
	"github.com/rerost/bqv/domain/annotateparser"
	"github.com/rerost/bqv/domain/viewmanager"
	"github.com/rerost/bqv/domain/viewservice"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func NewBQClient(ctx context.Context, cfg Config) (viewmanager.BQClient, error) {
	c, err := bigquery.NewClient(ctx, cfg.ProjectID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return viewmanager.BQClient(bqiface.AdaptClient(c)), nil
}

func NewFileManager(cfg Config) viewmanager.FileManager {
	return viewmanager.NewFileManager(cfg.Dir)
}

func NewRawBQClient(ctx context.Context, cfg Config) (bqiface.Client, error) {
	zap.L().Debug("Create BQ Client", zap.String("ProjectID", cfg.ProjectID))
	bqClient, err := bigquery.NewClient(ctx, cfg.ProjectID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return bqiface.AdaptClient(bqClient), nil
}

func NewBQMiddleware(bqClient bqiface.Client) validator.Middleware {
	return bqquery.NewBQQuery(bqClient)
}

func NewTime() time.Time {
	return time.Now()
}

func InitializeCmd(ctx context.Context, cfg Config) (*cobra.Command, error) {
	wire.Build(
		NewCmdRoot,
		viewservice.NewService,
		viewmanager.NewBQManager,
		annotateparser.NewParser,
		annotateparser.NewExtractor,
		annotateparser.NewManifests,
		NewFileManager,
		NewBQClient,
		validator.NewValidator,
		tablemock.NewTableMock,
		NewTime,
		NewRawBQClient,
		NewBQMiddleware,
	)
	return nil, nil
}
