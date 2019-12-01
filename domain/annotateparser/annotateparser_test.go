package annotateparser_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rerost/bqv/domain/annotateparser"
)

func TestParse(t *testing.T) {
	type Out struct {
		Body string
		Type string
	}

	inOutPairs := []struct {
		name string
		in   string
		out  []Out
	}{
		{
			in: `
[bqv:TEST]
- test_count
  - target: |
    SELECT owner_user_id
    FROM dataset.view
    GROUP BY owner_user_id
			`,
			out: []Out{
				{
					Type: "bqv:TEST",
					Body: `- test_count
  - target: |
    SELECT owner_user_id
    FROM dataset.view
    GROUP BY owner_user_id
			`,
				},
			},
		},
	}

	for _, inOutPair := range inOutPairs {
		inOutPair := inOutPair
		t.Run(inOutPair.name, func(t *testing.T) {
			p := annotateparser.NewParser()
			ctx := context.Background()

			outs, err := p.Parse(ctx, inOutPair.in)
			if err != nil {
				t.Error(err)
				return
			}

			if diff := cmp.Diff(len(inOutPair.out), len(outs)); diff != "" {
				t.Error(diff)
				return
			}
			for i, target := range inOutPair.out {
				if diff := cmp.Diff(target.Body, outs[i].Body()); diff != "" {
					t.Error(diff)
				}
				if diff := cmp.Diff(target.Type, outs[i].Type()); diff != "" {
					t.Error(diff)
				}
			}
		})
	}
}
