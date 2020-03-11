GO_COVER_ARGS := -coverpkg $(shell go list ./...)

PHONY: gen 
gen: mockgen
	go generate ./...

PHONY: mockgen
mockgen:
	mockgen github.com/googleapis/google-cloud-go-testing/bigquery/bqiface Client > mocks/mock_bqiface/client.go
	mockgen github.com/googleapis/google-cloud-go-testing/bigquery/bqiface Dataset > mocks/mock_bqiface/dataset.go

PHONY: ci-test
ci-test:
	echo ${GOOGLE_APPLICATION_CREDENTIALS_KEY_BASE64} | base64 -d > ${GOOGLE_APPLICATION_CREDENTIALS}
	make test

PHONY: test
test:
	go test -race -coverprofile=coverage.xml -covermode atomic ./...
