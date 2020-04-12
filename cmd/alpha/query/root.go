package query

import (
	"context"
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/rerost/bqv/domain/query"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func NewCmd(
	ctx context.Context,
	queryService query.QueryService,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query",
		Short: "",
	}

	cmd.AddCommand(
		&cobra.Command{
			Use: "exec",
			RunE: func(_ *cobra.Command, args []string) error {
				var eg errgroup.Group

				queries := make([]string, len(args))
				for i, file := range args {
					i := i
					file := file
					eg.Go(func() error {
						b, err := ioutil.ReadFile(file)
						if err != nil {
							return errors.WithStack(err)
						}

						queries[i] = string(b)
						return nil
					})
				}

				if err := eg.Wait(); err != nil {
					return errors.WithStack(err)
				}

				return errors.WithStack(queryService.BulkExec(ctx, queries))
			},
			Args: cobra.MinimumNArgs(1),
		},
	)

	return cmd
}
