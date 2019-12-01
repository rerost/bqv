// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package cmd

import (
	"cloud.google.com/go/bigquery"
	"context"
	"github.com/googleapis/google-cloud-go-testing/bigquery/bqiface"
	"github.com/pkg/errors"
	"github.com/rerost/bqv/domain/viewmanager"
	"github.com/rerost/bqv/domain/viewservice"
	"github.com/spf13/cobra"
)

// Injectors from wire.go:

func InitializeCmd(ctx context.Context, cfg Config) (*cobra.Command, error) {
	viewService := viewservice.NewService()
	bqClient, err := NewBQClient(ctx, cfg)
	if err != nil {
		return nil, err
	}
	bqManager := viewmanager.NewBQManager(bqClient)
	fileManager := NewFileManager(cfg)
	command := NewCmdRoot(ctx, viewService, bqManager, fileManager)
	return command, nil
}

// wire.go:

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
