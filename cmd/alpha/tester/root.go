package tester

import (
	"context"
	"os"

	"github.com/pkg/errors"
	"github.com/rerost/bqv/domain/tester"
	"github.com/spf13/cobra"
)

func NewCmd(
	ctx context.Context,
	testService tester.TestService,
) *cobra.Command {
	cmd := &cobra.Command{
		Use: "test",
		RunE: func(_ *cobra.Command, args []string) error {
			viewQueryFile := args[0]
			assertQueryFile := args[1]

			viewQuery, err := os.ReadFile(viewQueryFile)
			if err != nil {
				return errors.WithStack(err)
			}

			assertQuery, err := os.ReadFile(assertQueryFile)
			if err != nil {
				return errors.WithStack(err)
			}

			if err := testService.Test(ctx, string(viewQuery), string(assertQuery)); err != nil {
				return errors.WithStack(err)
			}

			return nil
		},
		Args: cobra.ExactArgs(2),
	}

	return cmd
}
