package viewmanager

import (
	"context"

	"cloud.google.com/go/bigquery"
)

type BQManager struct {
}

func NewBQManager(bqClient *bigquery.Client) BQManager {
	// TODO
	return BQManager{}
}

func (BQManager) List(ctx context.Context) ([]View, error) {
	return nil, nil
}
func (BQManager) Get(ctx context.Context, dataset string, name string) (View, error) {
	return nil, nil
}
func (BQManager) Create(ctx context.Context, view View) (View, error) {
	return nil, nil
}
func (BQManager) Update(ctx context.Context, view View) (View, error) {
	return nil, nil
}
func (BQManager) Delete(ctx context.Context, view View) error {
	return nil
}
