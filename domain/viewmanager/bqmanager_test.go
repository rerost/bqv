package viewmanager_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/rerost/bqv/cmd"
	"github.com/rerost/bqv/domain/viewmanager"
	"google.golang.org/api/iterator"
)

type dummyView struct {
	dataset string
	name    string
	query   string
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
	return dummyViewSetting{}
}

type dummyViewSetting struct {
}

func (dvs dummyViewSetting) Metadata() map[string]interface{} {
	return map[string]interface{}{}
}

func TestMain(m *testing.M) {
	bqDatasetPrefix := viewmanager.SetTest()

	ctx := context.Background()
	bqClient, err := cmd.NewBQClient(ctx, cmd.Config{
		ProjectID: os.Getenv("GOOGLE_APPLICATION_PROJECT_ID"),
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	exitCode := m.Run()

	for it := bqClient.Datasets(ctx); ; {
		ds, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if !strings.HasPrefix(ds.DatasetID(), bqDatasetPrefix) {
			continue
		}
		ds.Delete(ctx)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	os.Exit(exitCode)
}

func TestCreate(t *testing.T) {
	ctx := context.Background()
	bqClient, err := cmd.NewBQClient(ctx, cmd.Config{
		ProjectID: os.Getenv("GOOGLE_APPLICATION_PROJECT_ID"),
	})
	if err != nil {
		t.Error(err)
		return
	}

	bqManager := viewmanager.NewBQManager(bqClient)
	_, err = bqManager.Create(ctx, dummyView{dataset: "test", name: "test", query: "SELECT 1"})
	if err != nil {
		t.Error(err)
		return
	}
}
