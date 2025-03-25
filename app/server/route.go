package server

import "strings"

type Route struct {
	routes map[string]func(*HttpRequest) *HttpResponse
}

func NewRoute(routes map[string]func(*HttpRequest) *HttpResponse) *Route {
	return &Route{
		routes: routes,
	}
}

func (r *Route) Match(req *HttpRequest) (func(*HttpRequest) *HttpResponse, bool) {
	for route := range r.routes {
		if route == req.Path {
			return r.routes[req.Path], true
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
			return r.routes[route], true
		}
	}

	return nil, false
}

func isParam(path string) bool {
	return strings.HasPrefix(path, "{") && strings.HasSuffix(path, "}")
}
