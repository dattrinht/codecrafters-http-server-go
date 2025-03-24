package server

import (
	"fmt"
	"net"
	"os"
)

func Start(port string) {
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

	_, err = HandleConn(conn)
	if err != nil {
		fmt.Println("Failed to handle connection: ", err.Error())
		os.Exit(1)
	}
}

func HandleConn(c net.Conn) (int, error) {
	defer c.Close()

	buffer := make([]byte, 1024)
	r, err := c.Read(buffer)
	if err != nil {
		return r, err
	}

	req, err := ParseRequest(buffer)
	if err != nil {
		return r, err
	}

	res := HttpResponse{
		HttpVersion: req.HttpVersion,
	}
	if req.Path == "/" {
		res.StatusCode = 200
	} else {
		res.StatusCode = 404
	}

	message, err := res.ToString()
	if err != nil {
		return r, err
	}

	w, err := c.Write([]byte(message))
	if err != nil {
		return w, err
	}

	return w, nil
}
