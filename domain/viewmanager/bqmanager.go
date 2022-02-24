package viewmanager

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/googleapis/google-cloud-go-testing/bigquery/bqiface"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iterator"
)

var (
	datasetPrefixForTest string
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
	var views []View
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

			if tmd.Type != bigquery.ViewTable && tmd.Type != bigquery.MaterializedView {
				continue
			}

			metadata, err := b.convertTmdToMetadata(tmd)
			if err != nil {
				return nil, errors.WithStack(err)
			}

			views = append(views, bqView{
				dataSet: dataset.DatasetID(),
				name:    table.TableID(),
				query:   tmd.ViewQuery,
				setting: bqSetting{
					metadata: metadata,
				},
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
	fmt.Println(datasetPrefixForTest)
	ds := b.bqClient.Dataset(datasetPrefixForTest + view.DataSet())
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
	tmd, err := b.convertToTmd(view)
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

	return b.Get(ctx, datasetPrefixForTest+view.DataSet(), view.Name())
}
func (b BQManager) Update(ctx context.Context, view View) (View, error) {
	ds := b.bqClient.Dataset(datasetPrefixForTest + view.DataSet())
	t := ds.Table(view.Name())
	oldMeta, err := t.Metadata(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	tmd, err := b.convertToTmd(view)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	tmdForUpdate, err := b.convertTmdToForUpdate(tmd)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if reflect.ValueOf(oldMeta.MaterializedView).IsNil() {
		// Update view
		_, err = t.Update(ctx, tmdForUpdate, "")
	} else {
		if oldMeta.MaterializedView.Query == view.Query() {
			// Update materialized view metadata
			_, err = t.Update(ctx, tmdForUpdate, "")
		} else {
			// Update materialized view query and metadata
			zap.L().Warn(
				"Updating materialized view query is not yet supported. Updating with Delete and Create...",
				zap.String("view name", view.Name()),
			)
			err = t.Delete(ctx)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			err = t.Create(ctx, &tmd)
			if err != nil {
				return nil, errors.WithStack(err)
			}
		}
	}

	if err != nil {
		zap.L().Debug("Failed to update view", zap.String("err", err.Error()))
		if e, ok := err.(*googleapi.Error); ok && e.Code == 404 {
			return nil, NotFoundError
		}
		return nil, errors.WithStack(err)
	}

	view, err = b.Get(ctx, datasetPrefixForTest+view.DataSet(), view.Name())
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
	ds := b.bqClient.Dataset(datasetPrefixForTest + view.DataSet())
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

func (b BQManager) convertToTmd(view View) (bigquery.TableMetadata, error) {
	var description string
	{
		d := view.Setting().Metadata()["description"]
		if d != nil {
			description = d.(string)
		}
	}
	var labels map[string]string
	{
		l := view.Setting().Metadata()["labels"]
		var ls map[interface{}]interface{}
		if l != nil {
			ls = l.(map[interface{}]interface{})
		}
		labels = map[string]string{}
		for k, v := range ls {
			if v != nil {
				labels[fmt.Sprint(k)] = fmt.Sprint(v)
			}
		}
	}

	metadata := bigquery.TableMetadata{
		Name:        view.Name(),
		Description: description,
		Labels:      labels,
	}

	_, existMvDefData := view.Setting().Metadata()["materializedView"]
	if existMvDefData {
		mvDef, err := b.convertToMvd(view)
		if err != nil {
			panic(err)
		}
		metadata.MaterializedView = &mvDef
	} else {
		metadata.ViewQuery = view.Query()
	}

	return metadata, nil
}

func (b BQManager) convertTmdToForUpdate(tmd bigquery.TableMetadata) (bigquery.TableMetadataToUpdate, error) {
	tmdForUpdate := bigquery.TableMetadataToUpdate{
		Name:        tmd.Name,
		Description: tmd.Description,
	}

	for k, v := range tmd.Labels {
		tmdForUpdate.SetLabel(k, v)
	}
	if tmd.ViewQuery == "" {
		// Update MaterializedView
		tmdForUpdate.MaterializedView = tmd.MaterializedView
	} else {
		// Update View
		tmdForUpdate.ViewQuery = tmd.ViewQuery
	}

	return tmdForUpdate, nil
}

func (b BQManager) convertToMvd(view View) (bigquery.MaterializedViewDefinition, error) {
	mvDefData, existMvDefData := view.Setting().Metadata()["materializedView"]
	if !existMvDefData {
		return bigquery.MaterializedViewDefinition{} ,nil
	}
	zap.L().Debug("MaterializedView info", zap.String("mterializedView info", fmt.Sprintf("%v", mvDefData)))
	mvDef := bigquery.MaterializedViewDefinition{
		Query: view.Query(),
		EnableRefresh: mvDefData.(map[interface {}]interface {})["enableRefresh"].(bool),
		RefreshInterval:  time.Duration(mvDefData.(map[interface {}]interface {})["refreshInterval"].(int)) * time.Millisecond,
	}
	zap.L().Debug("Converted MaterializedViewDef", zap.String("MaterializedViewDef", fmt.Sprintf("%v", mvDef)))
	return mvDef, nil
}
