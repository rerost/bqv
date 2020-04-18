package template

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rerost/bqv/domain/template"
	"github.com/spf13/cobra"
)

func NewCmd(
	ctx context.Context,
	templateService template.TemplateService,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "template",
		Short: "",
	}

	cmd.AddCommand(
		&cobra.Command{
			Use: "render",
			RunE: func(_ *cobra.Command, args []string) error {
				viewsDirPath := args[0]
				templateFilePath := args[1:]
				err := templateService.Run(ctx, viewsDirPath, templateFilePath)
				if err != nil {
					return errors.WithStack(err)
				}

				return nil
			},
			Args: cobra.MinimumNArgs(2),
		},
	)
	return cmd
}
