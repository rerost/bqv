package test

import (
	"context"

	"github.com/spf13/cobra"
)

func NewCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use: "test",
		RunE: func(_ *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
