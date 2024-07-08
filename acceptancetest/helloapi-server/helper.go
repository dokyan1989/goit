package at

import (
	"encoding/json"

	"github.com/google/go-cmp/cmp"
)

/**
|-------------------------------------------------------------------------
| diff helper functions
|-----------------------------------------------------------------------*/

// responseBodyCmp ...
func responseBodyCmp() cmp.Options {
	fv := cmp.FilterValues(
		func(bx, by []byte) bool {
			return json.Valid(bx) && json.Valid(by)
		},
		cmp.Transformer("ParseJSON", func(in []byte) (out any) {
			if err := json.Unmarshal(in, &out); err != nil {
				panic(err) // should never occur given previous filter to ensure valid JSON
			}
			return out
		}),
	)

	// ref: https://github.com/google/go-cmp/issues/73
	fp := cmp.FilterPath(
		func(p cmp.Path) bool {
			step, ok := p[len(p)-1].(cmp.MapIndex)
			return ok && step.Key().String() == "request_id"
		},
		cmp.Ignore(),
	)

	return cmp.Options{fv, fp}
}
