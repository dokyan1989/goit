package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dokyan1989/goit/util/httprequest"

	"github.com/gorilla/schema"
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
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return writeJSON(w, code, apiErr)
	}

	return writeJSON(w, code, NewAPIError(err))
}

// func decodeBody(resp *http.Response, dst any) error {
// 	b, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return err
// 	}

// 	return json.Unmarshal(b, dst)
// }

/*
*
|-------------------------------------------------------------------------
| Utilities for HTTP Request
|-----------------------------------------------------------------------
*/

func decodeRequestBody(w http.ResponseWriter, r *http.Request, dst any) (int, error) {
	err := httprequest.DecodeBody(w, r, dst)
	if err != nil {
		var mr *httprequest.MalformedError
		if errors.As(err, &mr) {
			return mr.Status(), mr
		}

		return http.StatusInternalServerError, err
	}

	return 0, nil
}

/*
*
|-------------------------------------------------------------------------
| Utilities for HTTP Url
|-----------------------------------------------------------------------
*/

var urlDecorder = schema.NewDecoder()

func decodeURLQuery(r *http.Request, dst any) error {
	err := urlDecorder.Decode(dst, r.URL.Query())
	if err != nil {
		return err
	}

	return nil
}

// var urlEncoder = schema.NewEncoder()

// func encodeURLQuery(src interface{}, dst url.Values) error {
// 	urlEncoder.RegisterEncoder(time.Time{}, func(v reflect.Value) string {
// 		return v.Interface().(time.Time).Format(time.RFC3339)
// 	})

// 	err := urlEncoder.Encode(src, dst)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

/*
*
|-------------------------------------------------------------------------
| Utilities for JSON
|-----------------------------------------------------------------------
*/

func jsonString(src any) string {
	b, err := json.Marshal(src)
	if err != nil {
		return err.Error()
	}

	return string(b)
}

// func jsonReader(src any) io.Reader {
// 	b, err := json.Marshal(src)
// 	if err != nil {
// 		return bytes.NewReader([]byte(err.Error()))
// 	}

// 	return bytes.NewReader(b)
// }
