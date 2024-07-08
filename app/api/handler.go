package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
)

/**
|-------------------------------------------------------------------------
| createTodo
|-----------------------------------------------------------------------*/

func (s *server) createTodo(w http.ResponseWriter, r *http.Request) {
	logger := zerolog.Ctx(r.Context()).With().Str("[handler]", "createTodo").Logger()

	var req CreateTodoRequest
	status, err := decodeRequestBody(w, r, &req)
	if err != nil {
		logger.Error().Err(err).Msg("failed to decode request body")
		writeError(w, status, err)
		return
	}
	logger.Info().Str("request", jsonString(req)).Send()

	res, err := s.todoSvc.CreateTodo(r.Context(), &req)
	if err != nil {
		logger.Error().Err(err).Msg("[Todo Service] failed to call CreateTodo()")
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, NewAPIResponse(MsgSuccess, res))
}

/**
|-------------------------------------------------------------------------
| updateTodo
|-----------------------------------------------------------------------*/

func (s *server) updateTodo(w http.ResponseWriter, r *http.Request) {
	logger := zerolog.Ctx(r.Context()).With().Str("[handler]", "updateTodo").Logger()

	var req UpdateTodoRequest

	status, err := decodeRequestBody(w, r, &req)
	if err != nil {
		logger.Error().Err(err).Msg("failed to decode request body")
		writeError(w, status, err)
		return
	}

	paramID := chi.URLParam(r, "id")
	req.ID, err = strconv.ParseInt(paramID, 10, 64)
	if err != nil {
		logger.Error().Err(err).Msg("failed to parse 'id' parameter in url to integer")
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	logger.Info().Str("request", jsonString(req)).Send()

	res, err := s.todoSvc.UpdateTodo(r.Context(), &req)
	if err != nil {
		logger.Error().Err(err).Msg("[Todo Service] failed to call UpdateTodo()")
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, NewAPIResponse(MsgSuccess, res))
}

/**
|-------------------------------------------------------------------------
| getTodo
|-----------------------------------------------------------------------*/

func (s *server) getTodo(w http.ResponseWriter, r *http.Request) {
	logger := zerolog.Ctx(r.Context()).With().Str("[handler]", "getTodo").Logger()

	var (
		req GetTodoRequest
		err error
	)

	paramID := chi.URLParam(r, "id")
	req.ID, err = strconv.ParseInt(paramID, 10, 64)
	if err != nil {
		logger.Error().Err(err).Msg("failed to parse 'id' parameter in url to integer")
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	logger.Info().Str("request", jsonString(req)).Send()

	res, err := s.todoSvc.GetTodo(r.Context(), &req)
	if err != nil {
		logger.Error().Err(err).Msg("[Todo Service] failed to call ListTodos()")
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, NewAPIResponse(MsgSuccess, res))
}

/**
|-------------------------------------------------------------------------
| listTodos
|-----------------------------------------------------------------------*/

func (s *server) listTodos(w http.ResponseWriter, r *http.Request) {
	logger := zerolog.Ctx(r.Context()).With().Str("[handler]", "listTodos").Logger()

	var req ListTodosRequest
	err := decodeURLQuery(r, &req)
	if err != nil {
		logger.Error().Err(err).Msg("failed to decode URL query parameters")
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	logger.Info().Str("request", jsonString(req)).Send()

	res, err := s.todoSvc.ListTodos(r.Context(), &req)
	if err != nil {
		logger.Error().Err(err).Msg("[Todo Service] failed to call ListTodos()")
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, NewAPIResponse(MsgSuccess, res))
}

/**
|-------------------------------------------------------------------------
| deleteTodo
|-----------------------------------------------------------------------*/

func (s *server) deleteTodo(w http.ResponseWriter, r *http.Request) {
	logger := zerolog.Ctx(r.Context()).With().Str("[handler]", "deleteTodo").Logger()

	var (
		req DeleteTodoRequest
		err error
	)

	paramID := chi.URLParam(r, "id")
	req.ID, err = strconv.ParseInt(paramID, 10, 64)
	if err != nil {
		logger.Error().Err(err).Msg("failed to parse 'id' parameter in url to integer")
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	logger.Info().Str("request", jsonString(req)).Send()

	res, err := s.todoSvc.DeleteTodo(r.Context(), &req)
	if err != nil {
		logger.Error().Err(err).Msg("[Todo Service] failed to call DeleteTodo()")
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, NewAPIResponse(MsgSuccess, res))
}
