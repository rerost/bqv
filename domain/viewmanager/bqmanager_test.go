package viewmanager_test

import (
	"cloud.google.com/go/bigquery"
	"context"
	"fmt"
	"github.com/googleapis/google-cloud-go-testing/bigquery/bqiface"
	"github.com/magiconair/properties/assert"
	"os"
	"testing"
	"time"

	"github.com/rerost/bqv/cmd"
	"github.com/rerost/bqv/domain/viewmanager"
)

var (
	GcpProjectId = os.Getenv("GOOGLE_APPLICATION_PROJECT_ID")
)

type dummyView struct {
	dataset string
	name    string
	query   string
	settings dummyViewSetting
}

func (dv dummyView) DataSet() string {
	return dv.dataset
}

func (dv dummyView) Name() string {
	return dv.name
}

func (dv dummyView) Query() string {
	return dv.query
}

func (dv dummyView) Setting() viewmanager.Setting {
	return dv.settings
}

type dummyViewSetting struct {
	metadata map[string]interface{}
}

func (dvs dummyViewSetting) Metadata() map[string]interface{} {
	return dvs.metadata
}

func TestCreateView(t *testing.T) {
	viewmanager.SetTest()
	ctx := context.Background()
	bqClient, err := cmd.NewBQClient(ctx, cmd.Config{
		ProjectID: GcpProjectId,
	})
	if err != nil {
		t.Error(err)
		return
	}

	bqManager := viewmanager.NewBQManager(bqClient)
	_, err = bqManager.Create(ctx, dummyView{
		dataset: "test",
		name: "test_create_view",
		query: "SELECT 1",
	})
	if err != nil {
		t.Error(err)
		return
	}
}

func TestCreateViewWithMetadata(t *testing.T) {
	viewmanager.SetTest()
	ctx := context.Background()
	bqClient, err := cmd.NewBQClient(ctx, cmd.Config{
		ProjectID: GcpProjectId,
	})
	if err != nil {
		t.Error(err)
		return
	}

	bqManager := viewmanager.NewBQManager(bqClient)
	view, err := bqManager.Create(ctx, dummyView{
		dataset: "test",
		name: "test_create_view_with_metadata",
		query: "SELECT 1",
		settings: dummyViewSetting{
			metadata: map[string]interface{}{
				"description": "test view",
				"labels": map[interface{}]interface{}{
					"app": "test",
					"env": "test",
				},
			},
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	resultTable := bqClient.Dataset(view.DataSet()).Table("test_create_view_with_metadata")
	resultMetadata, _ := resultTable.Metadata(ctx)
	assert.Equal(t, resultMetadata.Description,"test view")
	assert.Equal(
		t,
		resultMetadata.Labels,
		map[string]string{
			"app": "test",
			"env": "test",
		},
	)
}

func TestCreateMaterializedView(t *testing.T) {
	viewmanager.SetTest()
	ctx := context.Background()
	bqClient, err := cmd.NewBQClient(ctx, cmd.Config{
		ProjectID: GcpProjectId,
	})
	if err != nil {
		t.Error(err)
		return
	}

	// Prepare view source table, Materialized views must belong to the same project or organization as the tables they reference.
	mvSrcDatasetName := "rerost_bqv_test_create_materialized_view"
	mvSrcDTableName := "test_empty_table"
	mvSrcDataset := bqClient.Dataset(mvSrcDatasetName)
	mvSrcDataset.Delete(ctx)
	mvSrcDataset.Create(
		ctx,
		&bqiface.DatasetMetadata{DatasetMetadata: bigquery.DatasetMetadata{Location: "US"}})
	mvSrcDataset.Table(mvSrcDTableName).Create(ctx, &bigquery.TableMetadata{
		Schema: []*bigquery.FieldSchema{
			&bigquery.FieldSchema{
				Name: "id",
				Type: bigquery.StringFieldType,
			},
		},
	})

	expectedQuery := fmt.Sprintf("SELECT 1 as NUM FROM %s.%s", mvSrcDatasetName, mvSrcDTableName)
	bqManager := viewmanager.NewBQManager(bqClient)
	view, err := bqManager.Create(ctx, dummyView{
		dataset: "test",
		name: "test_create_materialized_view",
		query: expectedQuery,
		settings: dummyViewSetting{
			metadata: map[string]interface{}{
				"materializedView": map[interface{}]interface{}{
					"refreshInterval": 900000,
					"enableRefresh": true,
				},
			},
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	resultTable := bqClient.Dataset(view.DataSet()).Table("test_create_materialized_view")
	resultMetadata, _ := resultTable.Metadata(ctx)
	assert.Equal(t, resultMetadata.MaterializedView.Query, expectedQuery)
	assert.Equal(t, resultMetadata.MaterializedView.RefreshInterval, time.Duration(900000) * time.Millisecond)
	assert.Equal(t, resultMetadata.MaterializedView.EnableRefresh, true)

}
