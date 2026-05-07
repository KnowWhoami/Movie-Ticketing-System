package router

import (
	"net/http"
	"strings"
)

type NetHTTPRouter struct {
	handler         http.Handler
	notFoundHandler http.Handler
	router          func() *http.ServeMux
}

func (netHttpRouter *NetHTTPRouter) NotFound(handler http.Handler) {
	netHttpRouter.notFoundHandler = handler
}

func (netHttpRouter *NetHTTPRouter) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	netHttpRouter.handler.ServeHTTP(writer, request)
}

func (netHttpRouter *NetHTTPRouter) Handle(method, path string, handler http.Handler) {
	method = strings.ToLower(method)
	if isAllowedMethod(method) {
		netHttpRouter.router().Handle(path, http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if strings.ToLower(request.Method) != method {
				if netHttpRouter.notFoundHandler == nil {
					// TODO: Write Error
					return
				}
				handler = netHttpRouter.notFoundHandler
			}
			handler.ServeHTTP(writer, request)
		}))
	}
}

func NewNetHTTP(handler http.Handler) *NetHTTPRouter {
	netHttpRouter := &NetHTTPRouter{handler: handler}
	netHttpRouter.router = func() *http.ServeMux {
		return netHttpRouter.handler.(*http.ServeMux)
	}
	return netHttpRouter
}

func (netHttpRouter *NetHTTPRouter) Name() string {
	return "Net Http Router"
}
