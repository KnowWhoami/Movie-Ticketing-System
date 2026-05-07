package router

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type HTTPRouter struct {
	handler http.Handler
	router  func() *httprouter.Router
}

func NewHTTPRouter(handler http.Handler) *HTTPRouter {
	httpRouter := &HTTPRouter{handler: handler}
	httpRouter.router = func() *httprouter.Router {
		return httpRouter.handler.(*httprouter.Router)
	}
	return httpRouter
}

func (httpRouter *HTTPRouter) Name() string {
	return "HTTP Router"
}

func (httpRouter *HTTPRouter) Handle(method, path string, handler http.Handler) {
	if isAllowedMethod(method) {
		httpRouter.router().Handler(method, path, handler)
	}
}

func (httpRouter *HTTPRouter) NotFound(handler http.Handler) {
	httpRouter.router().NotFound = handler
}

func (httpRouter *HTTPRouter) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	httpRouter.handler.ServeHTTP(writer, request)
}
