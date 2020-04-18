package template

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
	"github.com/rerost/bqv/domain/template/resolver"
	"golang.org/x/sync/errgroup"
)

type TemplateService interface {
	Run(ctx context.Context, viewDirPath string, templateFilePaths []string) error
}

type templateServiceImpl struct {
	queryResolver resolver.QueryResolver
}

func NewTemplateService(
	queryResolver resolver.QueryResolver,
) TemplateService {
	return &templateServiceImpl{
		queryResolver: queryResolver,
	}
}

func (t *templateServiceImpl) Run(ctx context.Context, viewDirPath string, templateFilePaths []string) error {
	var eg errgroup.Group

	for _, templateFilePath := range templateFilePaths {
		templateFilePath := templateFilePath
		eg.Go(func() error {
			err := t.run(ctx, viewDirPath, templateFilePath)
			return errors.WithStack(err)
		})
	}

	if err := eg.Wait(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (t *templateServiceImpl) run(ctx context.Context, viewDirPath string, templateFilePath string) error {
	queries, err := t.queryResolver.Resolve(ctx, templateFilePath)
	if err != nil {
		return errors.WithStack(err)
	}

	var dataset string
	{
		paths := strings.Split(templateFilePath, "/")
		pathLength := len(templateFilePath)
		if pathLength < 2 {
			return errors.New("Not valid template. template path must be `foo/<dataset_name>/<template_name>.sql`")
		}

		dataset = paths[pathLength-2]
	}

	for _, query := range queries {
		err := t.save(ctx, templateFilePath, dataset, query)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (t *templateServiceImpl) save(ctx context.Context, viewDirPath string, dataset string, query resolver.Query) error {
	err := ioutil.WriteFile(fmt.Sprintf("%s/%s/%s", viewDirPath, dataset, query.Name), []byte(query.Query), 0644)
	return errors.WithStack(err)
}
