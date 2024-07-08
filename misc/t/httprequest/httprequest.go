package httprequest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
)

type KV []string

type Options struct {
	Params KV
	Query  url.Values
	Body   any
}

func MustNew(t *testing.T, method, url string, opts Options) *http.Request {
	t.Helper()

	// 1. constructs url parameters.
	// e.g: /api/v1/todos/{id} -> /api/v1/todos/1
	if len(opts.Params) > 0 {
		if len(opts.Params)%2 != 0 {
			t.Fatal("number of keys and values must be even not odd")
		}

		var values []any
		for i, kv := range opts.Params {
			if i%2 == 0 {
				url = strings.ReplaceAll(url, fmt.Sprintf("{%s}", kv), "%s")
				continue
			}
			values = append(values, kv)
		}
		url = fmt.Sprintf(url, values...)
	}

	// 2. constructs query parameters.
	// e.g: /api/v1/todos?f1=foo&f2=bar
	if len(opts.Query) > 0 {
		url = fmt.Sprintf("%s?%s", url, opts.Query.Encode())
	}

	// 3. constructs request body
	var body io.Reader
	if opts.Body != nil {
		b, err := json.Marshal(opts.Body)
		if err != nil {
			t.Fatalf("marshalling body: %v", err)
		}
		body = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatalf("creating new request: %v", err)
	}

	return req
}

func MustNewTest(t *testing.T, method, url string, opts Options) *http.Request {
	t.Helper()

	// 1. constructs url parameters.
	// e.g: /api/v1/todos/{id} -> /api/v1/todos/1
	// ref: https://github.com/go-chi/chi/issues/76#issuecomment-370145140
	rctx := chi.NewRouteContext()
	if len(opts.Params) > 0 {
		if len(opts.Params)%2 != 0 { // element count is not even
			t.Fatal("number of keys and values must be even")
		}

		var values []any
		var rp chi.RouteParams

		for i, kv := range opts.Params {
			if i%2 == 0 {
				rp.Keys = append(rp.Keys, kv)
				url = strings.ReplaceAll(url, fmt.Sprintf("{%s}", kv), "%s")
				continue
			}

			rp.Values = append(rp.Values, kv)
			values = append(values, kv)
		}

		rctx.URLParams = rp
		url = fmt.Sprintf(url, values...)
	}

	// 2. constructs query parameters.
	// e.g: /api/v1/todos?f1=foo&f2=bar
	if opts.Query != nil {
		url = fmt.Sprintf("%s?%s", url, opts.Query.Encode())
	}

	// 3. constructs request body
	var body io.Reader
	if opts.Body != nil {
		b, err := json.Marshal(opts.Body)
		if err != nil {
			t.Fatalf("marshalling body: %v", err)
		}
		body = bytes.NewReader(b)
	}

	r := httptest.NewRequest(method, url, body)
	// add Chi's route context to request context
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	// add zero log to request context
	r = r.WithContext(zerolog.New(os.Stdout).WithContext(r.Context()))

	return r
}
