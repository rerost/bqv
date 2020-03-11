package viewmanager_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/rerost/bqv/cmd"
	"github.com/rerost/bqv/domain/viewmanager"
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

func TestCreate(t *testing.T) {
	viewmanager.SetTest()
	ctx := context.Background()
	fmt.Println(os.Getenv("GOOGLE_APPLICATION_PROJECT_ID"))
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
