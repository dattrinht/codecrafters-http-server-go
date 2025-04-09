package server

import (
	"strings"
)

type Route struct {
	routes map[string]func(*HttpRequest) *HttpResponse
}

var instance *Route

func NewRoute() *Route {
	if instance == nil {
		instance = &Route{
			routes: make(map[string]func(*HttpRequest) *HttpResponse),
		}
	}
	return instance
}

func (r *Route) AddRoute(method string, path string, handler func(*HttpRequest) *HttpResponse) {
	key := buildRouteKey(method, path)
	if _, exists := r.routes[key]; exists {
		return
	}
	r.routes[key] = handler
}

func (r *Route) GetRoute(method string, path string) (func(*HttpRequest) *HttpResponse, bool) {
	key := buildRouteKey(method, path)
	if handler, exists := r.routes[key]; exists {
		return handler, true
	}
	return nil, false
}

func (r *Route) Match(req *HttpRequest) (func(*HttpRequest) *HttpResponse, bool) {
	for routeKey := range r.routes {
		_, route := decomposeRouteKey(routeKey)
		if handler, ok := r.GetRoute(req.Method, req.Path); ok {
			return handler, true
		}

		routePaths := strings.Split(route, "/")
		pathPaths := strings.Split(req.Path, "/")

		if len(routePaths) != len(pathPaths) {
			continue
		}

		matched := true

		for i := range routePaths {
			if routePaths[i] == pathPaths[i] {
				continue
			}

			if isParam(routePaths[i]) {
				req.PathParams[routePaths[i][1:len(routePaths[i])-1]] = pathPaths[i]
				continue
			}

			matched = false
		}

		if matched {
			if handler, ok := r.GetRoute(req.Method, route); ok {
				return handler, true
			}
		}
	}

	return nil, false
}

func isParam(path string) bool {
	return strings.HasPrefix(path, "{") && strings.HasSuffix(path, "}")
}

func buildRouteKey(method string, path string) string {
	return method + "__%%__" + path
}

func decomposeRouteKey(key string) (string, string) {
	a := strings.Split(key, "__%%__")
	if len(a) != 2 {
		return "", ""
	}
	return a[0], a[1]
}
