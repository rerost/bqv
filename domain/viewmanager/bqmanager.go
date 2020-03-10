package viewmanager

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/bigquery"
	"github.com/googleapis/google-cloud-go-testing/bigquery/bqiface"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iterator"
)

type BQManager struct {
	bqClient BQClient
}

type BQClient interface {
	bqiface.Client
}

func NewBQManager(bqClient BQClient) BQManager {
	return BQManager{
		bqClient: bqClient,
	}
}

type bqView struct {
	dataSet string
	name    string
	query   string
	setting bqSetting
}

type bqSetting struct {
	metadata map[string]interface{}
}

func (b bqSetting) Metadata() map[string]interface{} {
	return b.metadata
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

func (b bqView) Setting() Setting {
	return Setting(b.setting)
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

			views = append(views, bqView{
				dataSet: dataset.DatasetID(),
				name:    table.TableID(),
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
		zap.L().Debug("Error when get metadata", zap.String("err", err.Error()))
		if e, ok := err.(*googleapi.Error); ok && e.Code == 404 {
			return nil, NotFoundError
		}
		return nil, errors.WithStack(err)
	}

	metadata, err := b.convertTmdToMetadata(tmd)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return bqView{
		dataSet: dataset,
		name:    name,
		query:   tmd.ViewQuery,
		setting: bqSetting{
			metadata: metadata,
		},
	}, nil
}
func (b BQManager) Create(ctx context.Context, view View) (View, error) {
	ds := b.bqClient.Dataset(view.DataSet())
	_, err := ds.Metadata(ctx)
	if err != nil {
		zap.L().Debug("Failed to create dataset", zap.String("err", err.Error()))
		if e, ok := err.(*googleapi.Error); ok && e.Code == 404 {
			err := ds.Create(
				ctx,
				&bqiface.DatasetMetadata{DatasetMetadata: bigquery.DatasetMetadata{Location: "US"}})
			if err != nil {
				return nil, errors.WithStack(err)
			}
		}
	}
	t := ds.Table(view.Name())
	tmd, err := b.converToTmd(view)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = t.Create(
		ctx,
		&tmd,
	)
	if err != nil {
		zap.L().Debug("Failed to create table", zap.String("Err", err.Error()))
		return nil, errors.WithStack(err)
	}

	return b.Get(ctx, view.DataSet(), view.Name())
}
func (b BQManager) Update(ctx context.Context, view View) (View, error) {
	ds := b.bqClient.Dataset(view.DataSet())
	t := ds.Table(view.Name())
	tmd, err := b.converToTmd(view)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	tmdForUpdate, err := b.convertTmdToForUpdate(tmd)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	_, err = t.Update(ctx, tmdForUpdate, "")
	if err != nil {
		zap.L().Debug("Failed to update view", zap.String("err", err.Error()))
		if e, ok := err.(*googleapi.Error); ok && e.Code == 404 {
			return nil, NotFoundError
		}
		return nil, errors.WithStack(err)
	}

	view, err = b.Get(ctx, view.DataSet(), view.Name())
	if err != nil {
		zap.L().Debug("Failed to get view", zap.String("err", err.Error()))
		if err == NotFoundError {
			return nil, NotFoundError
		} else {
			return nil, errors.WithStack(err)
		}
	}

	return view, nil
}
func (b BQManager) Delete(ctx context.Context, view View) error {
	ds := b.bqClient.Dataset(view.DataSet())
	t := ds.Table(view.Name())
	return errors.WithStack(t.Delete(ctx))
}

func (b BQManager) convertTmdToMetadata(tmd *bigquery.TableMetadata) (map[string]interface{}, error) {
	out, err := json.Marshal(tmd)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var res map[string]interface{}
	if err := json.Unmarshal(out, &res); err != nil {
		return nil, errors.WithStack(err)
	}

	return res, nil
}

func (b BQManager) converToTmd(view View) (bigquery.TableMetadata, error) {
	return bigquery.TableMetadata{
		Name:        view.Name(),
		ViewQuery:   view.Query(),
		Description: view.Setting().Metadata()["description"].(string),
		Labels:      view.Setting().Metadata()["labels"].(map[string]string),
	}, nil
}

func (b BQManager) convertTmdToForUpdate(tmd bigquery.TableMetadata) (bigquery.TableMetadataToUpdate, error) {
	tmdForUpdate := bigquery.TableMetadataToUpdate{
		Name:        tmd.Name,
		ViewQuery:   tmd.ViewQuery,
		Description: tmd.Description,
	}

	for k, v := range tmd.Labels {
		tmdForUpdate.SetLabel(k, v)
	}
	return tmdForUpdate, nil
}
