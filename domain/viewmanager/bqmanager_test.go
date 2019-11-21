package viewmanager_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/googleapis/google-cloud-go-testing/bigquery/bqiface"
	"github.com/rerost/bqv/domain/viewmanager"
	"github.com/rerost/bqv/mocks/mock_bqiface"
)

type testingBQClient struct {
	viewmanager.BQClient
	bqClient *mock_bqiface.MockClient
}

type testingBQDataset struct {
	bqiface.Dataset
	bqDataset *mock_bqiface.MockDataset
}

func TestList(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	// defer ctrl.Finish()

	bqDataset := mock_bqiface.NewMockDataset(ctrl)
	bqClient := mock_bqiface.NewMockClient(ctrl)
	bqClient.EXPECT().Dataset(ctx).Return(testingBQDataset{bqDataset: bqDataset}).Times(1)

	manager := viewmanager.NewBQManager(testingBQClient{bqClient: bqClient})
	manager.List(ctx)
}
