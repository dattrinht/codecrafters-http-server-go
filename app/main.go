package main

import (
	"os"
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

	s.Handle("/files/{filename}", func(req *server.HttpRequest) *server.HttpResponse {
		fileName := req.PathParams["filename"]

		dir := ""
		if len(os.Args) > 1 && os.Args[1] == "--directory" && len(os.Args) > 2 {
			dir = os.Args[2]
		}

		if dir != "" {
			fileName = dir + "/" + fileName
		}

		if _, err := os.Stat(fileName); err == nil {
			if content, err := os.ReadFile(fileName); err == nil {
				return &server.HttpResponse{
					StatusCode:  200,
					HttpVersion: req.HttpVersion,
					Body:        string(content),
					Headers: map[string]string{
						"Content-Type":   "application/octet-stream",
						"Content-Length": strconv.Itoa(len(content)),
					},
				}
			} else {
				return &server.HttpResponse{
					StatusCode:  500,
					HttpVersion: req.HttpVersion,
				}
			}
		} else {
			return &server.HttpResponse{
				StatusCode:  404,
				HttpVersion: req.HttpVersion,
			}
		}
	})

	s.Listen("4221")
}
