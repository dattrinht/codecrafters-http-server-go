package main

import (
	"strconv"

	"github.com/codecrafters-io/http-server-starter-go/app/server"
)

func main() {
	s := server.NewServer()

	s.Handle("/", func(req *server.HttpRequest) *server.HttpResponse {
		return &server.HttpResponse{
			StatusCode:  200,
			HttpVersion: req.HttpVersion,
		}
	})

	s.Handle("/echo/{str}", func(req *server.HttpRequest) *server.HttpResponse {
		str := req.PathParams["str"]
		return &server.HttpResponse{
			StatusCode:  200,
			HttpVersion: req.HttpVersion,
			Body:        str,
			Headers: map[string]string{
				"Content-Type":   "text/plain",
				"Content-Length": strconv.Itoa(len(str)),
			},
		}
	})

	s.Handle("/user-agent", func(req *server.HttpRequest) *server.HttpResponse {
		userAgent := req.Headers["User-Agent"]
		return &server.HttpResponse{
			StatusCode:  200,
			HttpVersion: req.HttpVersion,
			Body:        userAgent,
			Headers: map[string]string{
				"Content-Type":   "text/plain",
				"Content-Length": strconv.Itoa(len(userAgent)),
			},
		}
	})

	s.Listen("4221")
}
