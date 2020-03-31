package http

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

// PingHandler is a http handler who respond to /ping endpoint
type PingHandler struct{}

func NewPingHandler(r chi.Router) *PingHandler {
	h := &PingHandler{}
	r.Get("/ping", h.ping)
	return h
}

func (h *PingHandler) ping(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "pong!")
}
