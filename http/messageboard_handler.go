package http

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/guilherme-santos/messageboard"

	"github.com/go-chi/chi"
)

type MessageBoardHandler struct {
	svc messageboard.Service
}

func NewMessageBoardHandler(r chi.Router, svc messageboard.Service, creds map[string]string) *MessageBoardHandler {
	h := &MessageBoardHandler{
		svc: svc,
	}
	// Register create endpoint without authentication.
	r.Post("/v1/messages", h.create)

	// Here we're using a basic auth, but we could use a JWT token, which at least
	// will validate the token, still the permissions could be inside of each service/endpoint.

	authRouter := r.With(BasicAuth("Back's Message Board", creds))
	authRouter.Get("/v1/messages", h.list)
	authRouter.Route("/v1/messages/{id}", func(r chi.Router) {
		// Add a middleware that will be called in all following endpoints.
		r = r.With(h.loadMessage)
		r.Get("/", h.get)
		r.Put("/", h.update)
	})
	return h
}

func (h *MessageBoardHandler) list(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	opts := new(messageboard.ListOptions)
	opts.Load(req.URL.Query())

	list, err := h.svc.List(ctx, opts)
	if err != nil {
		responseError(w, err)
		return
	}
	responseJSON(w, http.StatusOK, list)
}

func (h *MessageBoardHandler) create(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var reqMsg *messageboard.Message
	err := json.NewDecoder(req.Body).Decode(&reqMsg)
	if err != nil {
		responseError(w, messageboard.NewError("invalid_json", err.Error()))
		return
	}

	msg, err := h.svc.Create(ctx, reqMsg)
	if err != nil {
		responseError(w, err)
		return
	}
	responseJSON(w, http.StatusCreated, msg)
}

type contextKey string

var msgCtxKey = contextKey("message")

func (h *MessageBoardHandler) loadMessage(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id := chi.URLParamFromCtx(ctx, "id")
		u, err := h.svc.Get(ctx, id)
		if err != nil {
			responseError(w, err)
			return
		}

		// save the message in the context, so next handlers can access it.
		ctx = context.WithValue(ctx, msgCtxKey, u)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func (h *MessageBoardHandler) get(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	msg := ctx.Value(msgCtxKey).(*messageboard.Message)
	responseJSON(w, http.StatusOK, msg)
}

func (h *MessageBoardHandler) update(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var reqMsg *messageboard.Message
	err := json.NewDecoder(req.Body).Decode(&reqMsg)
	if err != nil {
		responseError(w, messageboard.NewError("invalid_json", err.Error()))
		return
	}

	currentMsg := ctx.Value(msgCtxKey).(*messageboard.Message)

	// Here we have two approaches:
	// We can update inside of currentMsg the fields that are able to update
	// 		which means by default we do not update anything
	// Or we override the fields on reqMsg
	// What choose totally depend of your struct, in this case, if my struct grows
	// I want to update all fields, and only override that ones that could not be
	// updated, like id and creation_time (it'll be ignored anyways)
	reqMsg.ID = currentMsg.ID

	msg, err := h.svc.Update(ctx, reqMsg)
	if err != nil {
		responseError(w, err)
		return
	}
	responseJSON(w, http.StatusCreated, msg)
}

// responseError inspects the error and convert it into a meaningful status code and message.
func responseError(w http.ResponseWriter, err error) {
	var mberr *messageboard.Error
	if !errors.As(err, &mberr) {
		mberr = new(messageboard.Error)
		mberr.Code = "unknown_error"
		mberr.Message = err.Error()
	}

	var statusCode int
	switch mberr.Code {
	case "not_found":
		statusCode = http.StatusNotFound
	case "unauthorized":
		statusCode = http.StatusUnauthorized
	default:
		statusCode = http.StatusInternalServerError
	}
	responseJSON(w, statusCode, mberr)
}

// responseError responds the http call with the status code and the body as json.
func responseJSON(w http.ResponseWriter, statusCode int, body interface{}) {
	if body != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}

	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(body)
	if err != nil {
		log.Println("unable to encoding response as json:", err)
	}
}
