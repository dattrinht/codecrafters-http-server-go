package server

import (
	"errors"
	"fmt"
	"strings"
)

var STATUS_CODE_TO_MESSAGE = map[int]string{
	100: "Continue",
	101: "Switching Protocols",
	200: "OK",
	201: "Created",
	204: "No Content",
	301: "Moved Permanently",
	302: "Found",
	304: "Not Modified",
	400: "Bad Request",
	401: "Unauthorized",
	403: "Forbidden",
	404: "Not Found",
	405: "Method Not Allowed",
	429: "Too Many Requests",
	500: "Internal Server Error",
	502: "Bad Gateway",
	503: "Service Unavailable",
	504: "Gateway Timeout",
}

type HttpResponse struct {
	StatusCode  int
	HttpVersion string
	Headers     map[string]string
	Body        string
}

func (r *HttpResponse) ToString() (string, error) {
	statusMsg, ok := STATUS_CODE_TO_MESSAGE[r.StatusCode]
	if !ok {
		return "", errors.New("invalid status code")
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s %d %s\r\n", r.HttpVersion, r.StatusCode, statusMsg))

	sb.WriteString("\r\n")

	return sb.String(), nil
}
