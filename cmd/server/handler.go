package server

import (
	"encoding/json"
	"net/http"

	"KnowWhoami/movie-ticketing/pkgs/api"
)

func defaultHandler(handlerType string) http.Handler {
	if handlerType == "healthz" {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			healthCheckStatus(request, writer)
		})
	}
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		NotFoundHandler(request, writer)
	})
}

func wrapHandler(handler api.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		handler(request, writer)
	})
}

func healthCheckStatus(_ *http.Request, writer http.ResponseWriter) {
	body, _ := json.Marshal(map[string]string{"name": "apiServer", "status": "OK"})
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write(body)
}

func NotFoundHandler(_ *http.Request, writer http.ResponseWriter) {
	body, _ := json.Marshal(map[string]string{"name": "NotFound", "status": "404"})
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusNotFound)
	_, _ = writer.Write(body)
}
