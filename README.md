# bqv
To manage views.
Example repository https://github.com/rerost/bqv-example

## Required permissions for apply
```
bigquery.datasets.create
bigquery.datasets.get
bigquery.tables.create
bigquery.tables.get
bigquery.tables.list
bigquery.tables.update
bigquery.tables.getData
```

## Usage
```
## Base
bq view --dir=<DATASET_DIR> --projectid=<BQ_PROJECT_ID>

## Manage view with BQ
bqv view diff # TODO not color, not formatting
bqv view apply
bqv view dump

# TODO
bqv test <DATASET_DIR>
bqv test <DATASET_DIR>/<VIEW>
```
