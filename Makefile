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
	echo $(GOOGLE_APPLICATION_CREDENTIALS_KEY) > /tmp/key.json
	export GOOGLE_APPLICATION_CREDENTIALS=/tmp/key.json
	make test

PHONY: test
test:
	go test -race -coverprofile=profile.out -covermode atomic ./...
