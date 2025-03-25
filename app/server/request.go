package server

import (
	"errors"
	"strings"
)

type HttpRequest struct {
	HttpVersion string
	Method      string
	Path        string
	PathParams  map[string]string
	Headers     map[string]string
	Body        string
}

func ParseHttpRequest(data []byte) (*HttpRequest, error) {
	requestString := string(data)
	lines := strings.Split(requestString, "\r\n")
	if len(lines) == 0 {
		return nil, errors.New("empty request")
	}

	requestLine := lines[0]
	requestParts := strings.Split(requestLine, " ")
	if len(requestParts) < 3 {
		return nil, errors.New("invalid request line")
	}

	method := requestParts[0]
	path := requestParts[1]
	httpVersion := requestParts[2]

	headers := make(map[string]string)
	i := 1
	for ; i < len(lines); i++ {
		line := lines[i]
		if line == "" {
			break
		}

		headerParts := strings.SplitN(line, ": ", 2)
		if len(headerParts) < 2 {
			continue
		}

		headers[strings.TrimSpace(headerParts[0])] = strings.TrimSpace(headerParts[1])
	}

	body := ""
	if i < len(lines)-1 {
		body = strings.Join(lines[i+1:], "\r\n")
	}

	pathParams := make(map[string]string)

	return NewRequest(
		httpVersion,
		method,
		path,
		pathParams,
		headers,
		body,
	), nil
}

func NewRequest(httpVersion string, method string, path string, pathParams map[string]string, headers map[string]string, body string) *HttpRequest {
	return &HttpRequest{
		HttpVersion: httpVersion,
		Method:      method,
		Path:        path,
		PathParams:  pathParams,
		Headers:     headers,
		Body:        body,
	}
}
