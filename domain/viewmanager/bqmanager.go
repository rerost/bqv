package viewmanager

import "cloud.google.com/go/bigquery"

type BQManager interface {
	ViewReadWriter
}

func NewBQManager(bqClient *bigquery.Client) BQManager {
	// TODO
	return nil
}
