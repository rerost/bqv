package alpha

import (
	"context"

	"github.com/rerost/bqv/cmd/alpha/query"
	"github.com/rerost/bqv/cmd/alpha/template"
	"github.com/rerost/bqv/cmd/alpha/tester"
	dquery "github.com/rerost/bqv/domain/query"
	dtemplate "github.com/rerost/bqv/domain/template"
	dtester "github.com/rerost/bqv/domain/tester"
	"github.com/spf13/cobra"
)

func NewCmd(
	ctx context.Context,
	queryService dquery.QueryService,
	templateService dtemplate.TemplateService,
	testService dtester.TestService,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alpha",
		Short: "unstable functions",
	}

	cmd.AddCommand(
		query.NewCmd(ctx, queryService),
		template.NewCmd(ctx, templateService),
		tester.NewCmd(ctx, testService),
	)

	return cmd
}
