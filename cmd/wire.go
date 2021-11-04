//+build wireinject

package cmd

import (
	"context"

	"cloud.google.com/go/bigquery"
	"github.com/google/wire"
	"github.com/googleapis/google-cloud-go-testing/bigquery/bqiface"
	datatransfer "cloud.google.com/go/bigquery/datatransfer/apiv1"
	"github.com/pkg/errors"
	"github.com/rerost/bqv/domain/query"
	"github.com/rerost/bqv/domain/template"
	"github.com/rerost/bqv/domain/template/resolver"
	"github.com/rerost/bqv/domain/tester"
	"github.com/rerost/bqv/domain/viewmanager"
	"github.com/rerost/bqv/domain/viewservice"
	"github.com/spf13/cobra"
)

func NewRawBQClient(ctx context.Context, cfg Config) (bqiface.Client, error) {
	c, err := bigquery.NewClient(ctx, cfg.ProjectID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return bqiface.AdaptClient(c), nil
}

func NewBQClient(ctx context.Context, cfg Config) (viewmanager.BQClient, error) {
	c, err := NewRawBQClient(ctx, cfg)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return viewmanager.BQClient(c), nil
}

func NewDataTransferClient(ctx context.Context) (*datatransfer.Client, error) {
	c, err := datatransfer.NewClient(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return c, nil
}

func NewViewService(datatransferClient *datatransfer.Client, cfg Config) viewservice.ViewService {
	return viewservice.NewService(datatransferClient, cfg.ProjectID)
}

func NewFileManager(cfg Config) viewmanager.FileManager {
	return viewmanager.NewFileManager(cfg.Dir)
}

func InitializeCmd(ctx context.Context, cfg Config) (*cobra.Command, error) {
	wire.Build(
		NewCmdRoot,
		viewmanager.NewBQManager,
		NewFileManager,
		NewBQClient,
		NewRawBQClient,
		NewDataTransferClient,
		NewViewService,
		query.NewQueryService,
		template.NewTemplateService,
		resolver.NewQueryResolver,
		tester.NewTestService,
	)
	return nil, nil
}
