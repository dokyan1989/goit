package helloapi

import (
	"encoding/json"
	"net/http"

	"github.com/google/go-cmp/cmp"
)

/*
*
|-------------------------------------------------------------------------
| Utilities for HTTP Response
|-----------------------------------------------------------------------
*/
func writeJSON(w http.ResponseWriter, code int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, code int, err error) error {
	return writeJSON(w, code, err)
}

/**
|-------------------------------------------------------------------------
| Utilities for HTTP Url
|-----------------------------------------------------------------------*/

// var urlEncoder = schema.NewEncoder()

// func encodeQueryParams(src interface{}, dst url.Values) error {
// 	urlEncoder.RegisterEncoder(time.Time{}, func(v reflect.Value) string {
// 		return v.Interface().(time.Time).Format(time.RFC3339)
// 	})

// 	err := urlEncoder.Encode(src, dst)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

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
