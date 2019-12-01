package test

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/rerost/bq-table-validator/domain/validator"
	"github.com/rerost/bq-table-validator/types"
	"github.com/rerost/bqv/domain/annotateparser"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

func NewCmd(ctx context.Context, validate validator.Validator, apm annotateparser.Manifests) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "test",
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			targetFilePath := args[0]
			targetFile, err := ioutil.ReadFile(targetFilePath)
			if err != nil {
				return errors.WithStack(err)
			}

			manifests, err := apm.Manifests(ctx, string(targetFile))
			if err != nil {
				return errors.WithStack(err)
			}

			validates := make([]types.Validate, 0, len(manifests))
			var validatesSize = 0
			for _, annotation := range manifests {
				if annotation.Type() != "bqv:TEST" {
					continue
				}
				validatesSize++

				var validate types.Validate
				if err := yaml.Unmarshal([]byte(annotation.Body()), &validate); err != nil {
					return errors.WithStack(err)
				}
				validates = append(validates, validate)
			}
			validates = validates[:validatesSize]

			for _, v := range validates {
				out, err := validate.Valid(ctx, v)
				if err != nil {
					return errors.WithStack(err)
				}

				// TODO format & color
				fmt.Printf(out)
			}
			return nil
		},
	}

	return cmd
}
