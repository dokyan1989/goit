package helloweb

import (
	"encoding/json"

	"github.com/google/go-cmp/cmp"
)

/**
|-------------------------------------------------------------------------
| Utilities for JSON
|-----------------------------------------------------------------------*/

// convertToJSON wraps json.Marshal() function. If marshalled failed, returns error message.
// It is intended for using in log function.
func convertToJSON(src any) string {
	b, err := json.Marshal(src)
	if err != nil {
		return err.Error()
	}

	return string(b)
}

/**
|-------------------------------------------------------------------------
| diff helper functions
|-----------------------------------------------------------------------*/

// responseCmp ...
func responseCmp() cmp.Option {
	return cmp.FilterValues(
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
}
