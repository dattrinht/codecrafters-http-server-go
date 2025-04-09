package server

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type Server struct {
	routes map[string]func(*HttpRequest) *HttpResponse
}

func NewServer() *Server {
	return &Server{
		routes: make(map[string]func(*HttpRequest) *HttpResponse),
	}
}

var route = NewRoute()

func (s *Server) Handle(method string, path string, handler func(*HttpRequest) *HttpResponse) {
	route.AddRoute(method, path, handler)
}

func (s *Server) Listen(port string) {
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		fmt.Printf("Failed to bind to port %s\n", port)
		os.Exit(1)
	}

	tp := NewThreadPool(10)
	tp.Start()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			break
		}

		tp.Submit(func() {
			go func(c net.Conn) {
				_, err := s.HandleConn(c)
				if err != nil {
					fmt.Println("Failed to handle connection: ", err.Error())
				}
			}(conn)
		})
	}

	tp.Stop()
}

func (s *Server) HandleConn(c net.Conn) (int, error) {
	defer c.Close()

	buffer := make([]byte, 1024)
	r, err := c.Read(buffer)
	if err != nil {
		return r, err
	}

	req, err := ParseHttpRequest(buffer)
	if err != nil {
		return r, err
	}

	handler, ok := route.Match(req)
	var res *HttpResponse
	if !ok {
		res = &HttpResponse{
			StatusCode:  404,
			HttpVersion: req.HttpVersion,
		}
	} else {
		res = handler(req)
	}

	if res.StatusCode == 200 {
		encodings := strings.SplitSeq(req.Headers["Accept-Encoding"], ",")
		for encoding := range encodings {
			if strings.TrimSpace(encoding) == "gzip" {
				var buf bytes.Buffer

				gz := gzip.NewWriter(&buf)
				gz.Write([]byte(res.Body))
				gz.Close()

				res.Body = buf.String()
				res.Headers["Content-Length"] = strconv.Itoa(len(res.Body))
				res.Headers["Content-Encoding"] = "gzip"

				break
			}
		}
	}

	message, err := res.Stringify()
	if err != nil {
		return r, err
	}

	w, err := c.Write([]byte(message))
	if err != nil {
		return w, err
	}

	return w, nil
}
