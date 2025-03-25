package server

import (
	"fmt"
	"net"
	"os"
)

type Server struct {
	routes map[string]func(*HttpRequest) *HttpResponse
}

func NewServer() *Server {
	return &Server{
		routes: make(map[string]func(*HttpRequest) *HttpResponse),
	}
}

func (s *Server) Handle(path string, handler func(*HttpRequest) *HttpResponse) {
	s.routes[path] = handler
}

func (s *Server) Listen(port string) {
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		fmt.Printf("Failed to bind to port %s\n", port)
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	_, err = s.HandleConn(conn)
	if err != nil {
		fmt.Println("Failed to handle connection: ", err.Error())
		os.Exit(1)
	}
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

	route := NewRoute(s.routes)
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
