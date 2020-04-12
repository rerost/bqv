package alpha

import (
	"context"

	"github.com/rerost/bqv/cmd/alpha/query"
	dquery "github.com/rerost/bqv/domain/query"
	"github.com/spf13/cobra"
)

func NewCmd(
	ctx context.Context,
	queryService dquery.QueryService,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alpha",
		Short: "unstable functions",
	}

	cmd.AddCommand(
		query.NewCmd(ctx, queryService),
	)

	return cmd
}
