package view

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/rerost/bqv/domain/viewmanager"
	"github.com/rerost/bqv/domain/viewservice"
	"github.com/spf13/cobra"
)

func NewCmd(ctx context.Context, viewService viewservice.ViewService, bqManager viewmanager.BQManager, fileManager viewmanager.FileManager) *cobra.Command {
	cmd := &cobra.Command{
		Use: "view",
	}

	cmd.AddCommand(
		&cobra.Command{
			Use: "diff",
			RunE: func(_ *cobra.Command, args []string) error {
				res, err := viewService.Diff(ctx)
				if err != nil {
					return errors.WithStack(err)
				}
				fmt.Println(res)

				return nil
			},
		},
		&cobra.Command{
			Use: "sync",
			RunE: func(_ *cobra.Command, args []string) error {
				err := viewService.Sync(ctx)
				if err != nil {
					return errors.WithStack(err)
				}

				return nil
			},
		},
	)

	return cmd
}
