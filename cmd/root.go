package cmd

import (
	"context"

	"github.com/rerost/bqv/cmd/test"
	"github.com/rerost/bqv/cmd/view"
	"github.com/rerost/bqv/domain/viewmanager"
	"github.com/rerost/bqv/domain/viewservice"
	"github.com/spf13/cobra"
)

func NewCmdRoot(
	ctx context.Context,
	viewService viewservice.ViewService,
	bqManager viewmanager.BQManager,
	fileManager viewmanager.FileManager,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bqv",
		Short: "Manage BigQuery view",
	}

	cmd.AddCommand(
		view.NewCmd(ctx, viewService, bqManager, fileManager),
		test.NewCmd(ctx),
	)

	return cmd
}
