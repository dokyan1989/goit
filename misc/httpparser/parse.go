package httpparser

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/gorilla/schema"
)

var errNilPtr = errors.New("dst pointer is nil")

// parseParams extracts and populates url route parameters into dst
// e.g: /api/v1/todos/{id} -> /api/v1/todos/1 -> dst=1
func parseParams(r *http.Request, key string, dst any) error {

	s := chi.URLParam(r, key)
	if s == "" {
		return fmt.Errorf("%s does not exist", key)
	}

	// ref: https://cs.opensource.google/go/go/+/refs/tags/go1.22.5:src/database/sql/convert.go;l=388

	// check dst is pointer and not nil
	dpv := reflect.ValueOf(dst)
	if dpv.Kind() != reflect.Pointer {
		return errors.New("dst not a pointer")
	}
	if dpv.IsNil() {
		return errNilPtr
	}
	dv := reflect.Indirect(dpv)

	// do conversion
	switch dv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i64, err := strconv.ParseInt(s, 10, dv.Type().Bits())
		if err != nil {
			return err
		}
		dv.SetInt(i64)
		return nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u64, err := strconv.ParseUint(s, 10, dv.Type().Bits())
		if err != nil {
			return err
		}
		dv.SetUint(u64)
		return nil

	case reflect.String:
		dv.SetString(s)
		return nil

	default:
		return fmt.Errorf("converting value to %s is unsupported", dv.Kind())
	}
}

// parseQuery populates url query parameters to dst
func parseQuery(r *http.Request, dst any) error {
	dec := schema.NewDecoder()
	dec.IgnoreUnknownKeys(true)
	dec.SetAliasTag("json")

	err := dec.Decode(dst, r.URL.Query())
	if err != nil {
		return err
	}

	return nil
}

// https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body
type MalformedError struct {
	status int
	msg    string
}

func (mr *MalformedError) Error() string {
	return mr.msg
}

func (mr *MalformedError) Status() int {
	return mr.status
}

// parseBody populates http request body parameters to dst
func parseBody(w http.ResponseWriter, r *http.Request, dst any) error {
	// If the Content-Type header is present, check that it has the value
	// application/json. Note that we parse and normalize the header to remove
	// any additional parameters (like charset or boundary information) and normalize
	// it by stripping whitespace and converting to lowercase before we check the
	// value.
	ct := r.Header.Get("Content-Type")
	if ct != "" {
		mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
		if mediaType != "application/json" {
			msg := "Content-Type header is not application/json"
			return &MalformedError{status: http.StatusUnsupportedMediaType, msg: msg}
		}
	}

	// Use http.MaxBytesReader to enforce a maximum read of 1MB from the
	// response body. A request body larger than that will now result in
	// Decode() returning a "http: request body too large" error.
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	// Setup the decoder and call the DisallowUnknownFields() method on it.
	// This will cause Decode() to return a "json: unknown field ..." error
	// if it encounters any extra unexpected fields in the JSON. Strictly
	// speaking, it returns an error for "keys which do not match any
	// non-ignored, exported fields in the destination".
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		// Catch any syntax errors in the JSON and send an error message
		// which interpolates the location of the problem to make it
		// easier for the client to fix.
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &MalformedError{status: http.StatusBadRequest, msg: msg}

		// In some circumstances Decode() may also return an
		// io.ErrUnexpectedEOF error for syntax errors in the JSON. There
		// is an open issue regarding this at
		// https://github.com/golang/go/issues/25956.
		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := "Request body contains badly-formed JSON"
			return &MalformedError{status: http.StatusBadRequest, msg: msg}

		// Catch any type errors, like trying to assign a string in the
		// JSON request body to a int field in our Person struct. We can
		// interpolate the relevant field name and position into the error
		// message to make it easier for the client to fix.
		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &MalformedError{status: http.StatusBadRequest, msg: msg}

		// Catch the error caused by extra unexpected fields in the request
		// body. We extract the field name from the error message and
		// interpolate it in our custom error message. There is an open
		// issue at https://github.com/golang/go/issues/29035 regarding
		// turning this into a sentinel error.
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &MalformedError{status: http.StatusBadRequest, msg: msg}

		// An io.EOF error is returned by Decode() if the request body is
		// empty.
		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return &MalformedError{status: http.StatusBadRequest, msg: msg}

		// Catch the error caused by the request body being too large. Again
		// there is an open issue regarding turning this into a sentinel
		// error at https://github.com/golang/go/issues/30715.
		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return &MalformedError{status: http.StatusRequestEntityTooLarge, msg: msg}

		// Otherwise default to logging the error and sending a 500 Internal
		// Server Error response.
		default:
			return err
		}
	}

	// Call decode again, using a pointer to an empty anonymous struct as
	// the destination. If the request body only contained a single JSON
	// object this will return an io.EOF error. So if we get anything else,
	// we know that there is additional data in the request body.
	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		msg := "Request body must only contain a single JSON object"
		return &MalformedError{status: http.StatusBadRequest, msg: msg}
	}

	return nil
}