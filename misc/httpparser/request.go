package httpparser

import (
	"errors"
	"net/http"
)

func RouteParams(r *http.Request, key string, dst any) error {
	return parseParams(r, key, dst)
}

func RequestQuery(r *http.Request, dst any) error {
	return parseQuery(r, dst)
}

func RequestBody(w http.ResponseWriter, r *http.Request, dst any) (int, error) {
	err := parseBody(w, r, dst)
	if err == nil {
		return 0, nil
	}

	var mr *MalformedError
	if errors.As(err, &mr) {
		return mr.Status(), mr
	}

	return http.StatusInternalServerError, err
}
