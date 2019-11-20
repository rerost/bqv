package viewmanager

import (
	"context"

	"cloud.google.com/go/bigquery"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
)

type BQManager struct {
	bqClient *bigquery.Client
}

func NewBQManager(bqClient *bigquery.Client) BQManager {
	return BQManager{
		bqClient: bqClient,
	}
}

type bqView struct {
	dataSet string
	name    string
	query   string
}

func (b bqView) DataSet() string {
	return b.dataSet
}

func (b bqView) Name() string {
	return b.name
}

func (b bqView) Query() string {
	return b.query
}

func (b BQManager) List(ctx context.Context) ([]View, error) {
	datasets := b.bqClient.Datasets(ctx)
	views := []View{}
	for {
		dataset, err := datasets.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, errors.WithStack(err)
		}

		tables := dataset.Tables(ctx)
		for {
			table, err := tables.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return nil, errors.WithStack(err)
			}

			tmd, err := table.Metadata(ctx)
			if err != nil {
				return nil, errors.WithStack(err)
			}

			if tmd.Type != bigquery.ViewTable {
				continue
			}

			dsmd, err := dataset.Metadata(ctx)
			if err != nil {
				return nil, errors.WithStack(err)
			}

			views = append(views, bqView{
				dataSet: dsmd.Name,
				name:    tmd.Name,
				query:   tmd.ViewQuery,
			})
		}
	}

	return views, nil
}
func (b BQManager) Get(ctx context.Context, dataset string, name string) (View, error) {
	ds := b.bqClient.Dataset(dataset)
	t := ds.Table(name)
	tmd, err := t.Metadata(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return bqView{
		dataSet: dataset,
		name:    name,
		query:   tmd.ViewQuery,
	}, nil
}
func (b BQManager) Create(ctx context.Context, view View) (View, error) {
	ds := b.bqClient.Dataset(view.DataSet())
	t := ds.Table(view.Name())
	t.Create(ctx, &bigquery.TableMetadata{
		Name:      view.Name(),
		ViewQuery: view.Query(),
	})

	return b.Get(ctx, view.DataSet(), view.Name())
}
func (b BQManager) Update(ctx context.Context, view View) (View, error) {
	ds := b.bqClient.Dataset(view.DataSet())
	t := ds.Table(view.Name())
	t.Update(ctx, bigquery.TableMetadataToUpdate{
		ViewQuery: view.Query(),
	}, "")

	return b.Get(ctx, view.DataSet(), view.Name())
}
func (b BQManager) Delete(ctx context.Context, view View) error {
	ds := b.bqClient.Dataset(view.DataSet())
	t := ds.Table(view.Name())
	return errors.WithStack(t.Delete(ctx))
}
