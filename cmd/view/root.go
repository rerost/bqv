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
				res, err := viewService.Diff(ctx, fileManager, bqManager)
				if err != nil {
					return errors.WithStack(err)
				}
				if len(res) == 0 {
					return nil
				}
				// TODO color & format
				fmt.Println(res)

				return nil
			},
		},
		&cobra.Command{
			Use: "apply",
			RunE: func(_ *cobra.Command, args []string) error {
				err := viewService.Copy(ctx, fileManager, bqManager)
				if err != nil {
					return errors.WithStack(err)
				}

				return nil
			},
		},
		&cobra.Command{
			Use: "dump",
			RunE: func(_ *cobra.Command, args []string) error {
				err := viewService.Copy(ctx, bqManager, fileManager)
				if err != nil {
					return errors.WithStack(err)
				}

				return nil
			},
		},
		&cobra.Command{
			Use: "flist",
			RunE: func(_ *cobra.Command, args []string) error {
				views, err := viewService.List(ctx, fileManager)
				if err != nil {
					return errors.WithStack(err)
				}
				// TODO color & format
				fmt.Println(views)

				return nil
			},
		},
		&cobra.Command{
			Use: "blist",
			RunE: func(_ *cobra.Command, args []string) error {
				views, err := viewService.List(ctx, bqManager)
				if err != nil {
					return errors.WithStack(err)
				}
				// TODO color & format
				fmt.Println(views)

				return nil
			},
		},
	)

	return cmd
}
