package http

import (
	"net/http"
)

type ProbeHandler struct {
}

func NewProbeHandler() *ProbeHandler {
	return &ProbeHandler{}
}

func (h *ProbeHandler) Live(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	w.Write([]byte("Auth Service is alive"))
}

func (h *ProbeHandler) Ready(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	w.Write([]byte("Auth Service is ready"))
}

func (h *ProbeHandler) StartUp(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	w.Write([]byte("Auth Service is start up"))
}
