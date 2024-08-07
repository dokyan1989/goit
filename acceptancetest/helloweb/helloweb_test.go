package helloweb

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/dokyan1989/goit/misc/t/httprequest"
	"github.com/dokyan1989/goit/misc/t/seeder"
	"github.com/sebdah/goldie/v2"
)

func TestWeb(t *testing.T) {
	type Input struct {
		method          string
		url             string
		httpRequestOpts httprequest.Options
		seedURL         string
	}

	type Want struct {
		statusCode int
		body       string
	}

	tests := []struct {
		name  string
		input Input
		want  Want
	}{
		{
			name: "health check",
			input: Input{
				method:          http.MethodGet,
				url:             "http://localhost:3000/health",
				httpRequestOpts: httprequest.Options{},
				seedURL:         "./sampledata",
			},
			want: Want{
				statusCode: http.StatusOK,
				body:       `{"status":"up"}`,
			},
		},
		{
			name: "todo web check",
			input: Input{
				method:  http.MethodGet,
				url:     "http://localhost:3000/v1/todos/1",
				seedURL: "./sampledata",
			},
			want: Want{
				statusCode: http.StatusOK,
				body:       `{"data":{"todos":[{"id":1,"status":"UNKNOWN","title":"title 1"},{"id":2,"status":"DOING","title":"title 2"},{"id":3,"status":"DONE","title":"title 3"}]},"message":"ok"}`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seeder.MustRun(context.Background(), t, postgresC, fmt.Sprintf("file://%s", filepath.Join(workingDir, tt.input.seedURL)))
			r := httprequest.MustNew(t, tt.input.method, tt.input.url, tt.input.httpRequestOpts)

			got, err := http.DefaultClient.Do(r)
			if err != nil {
				t.Errorf("http.DefaultClient.Do() error = %v", err)
				return
			}
			defer got.Body.Close()

			if tt.want.statusCode != got.StatusCode {
				t.Errorf("Response status code = %v, want %v", got.StatusCode, tt.want.statusCode)
				return
			}

			gotBody, err := io.ReadAll(got.Body)
			if err != nil {
				t.Errorf("io.ReadAll() error = %v", err)
				return
			}

			g := goldie.New(t)
			g.Assert(t, t.Name(), gotBody)
		})
	}

	// time.Sleep(2 * time.Second)
	// sendInterrupt()
	// t.Log("send interupt")
}

// [REQUEST]
// https://stackoverflow.com/questions/23070876/reading-body-of-http-request-without-modifying-request-state
// https://blog.flexicondev.com/read-go-http-request-body-multiple-times
// [ERROR]
// https://vladimir.varank.in/notes/2021/12/error-messages-in-go/
// https://www.reddit.com/r/golang/comments/1acx63i/how_do_you_get_stack_traces_for_errors/
// https://www.reddit.com/r/golang/comments/1acx63i/comment/kjz33yt/?utm_source=share&utm_medium=web3x&utm_name=web3xcss&utm_term=1&utm_content=share_button
// https://www.reddit.com/r/golang/comments/z870te/multiple_error_wrapping_is_coming_in_go_120/

// func TestCleanupOnPanic(t *testing.T) {
// 	cleanup := func() {
// 		println("enter cleanup")
// 		println("leave cleanup")
// 	}

// 	t.Cleanup(cleanup) // (1)
// 	// defer cleanup() // (2)
// 	t.Run("", func(t *testing.T) {
// 		panic("boom")
// 	})
// 	// cleanup() // (3)
// }
