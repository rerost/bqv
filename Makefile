PHONY: gen 
gen: mockgen
	go generate ./...

PHONY: mockgen
mockgen:
	mockgen github.com/googleapis/google-cloud-go-testing/bigquery/bqiface Client > mocks/mock_bqiface/client.go
	mockgen github.com/googleapis/google-cloud-go-testing/bigquery/bqiface Dataset > mocks/mock_bqiface/dataset.go
