package alpha

import (
	"context"

	"github.com/rerost/bqv/cmd/alpha/query"
	"github.com/rerost/bqv/cmd/alpha/template"
	dquery "github.com/rerost/bqv/domain/query"
	dtemplate "github.com/rerost/bqv/domain/template"
	"github.com/spf13/cobra"
)

func NewCmd(
	ctx context.Context,
	queryService dquery.QueryService,
	templateService dtemplate.TemplateService,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alpha",
		Short: "unstable functions",
	}

	cmd.AddCommand(
		query.NewCmd(ctx, queryService),
		template.NewCmd(ctx, templateService),
	)

	return cmd
}
