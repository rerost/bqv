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

func NewCmd(ctx context.Context, validate validator.Validator, ap annotateparser.Parser, ape annotateparser.Extractor) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "test",
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			targetFilePath := args[0]
			targetFile, err := ioutil.ReadFile(targetFilePath)
			if err != nil {
				return errors.WithStack(err)
			}

			rawAnnotations, err := ape.Extract(string(targetFile))
			if err != nil {
				return errors.WithStack(err)
			}

			annotations, err := ap.Parse(ctx, rawAnnotations[0])
			if err != nil {
				return errors.WithStack(err)
			}

			validates := make([]types.Validate, 0, len(annotations))
			var validatesSize = 0
			for _, annotation := range annotations {
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
