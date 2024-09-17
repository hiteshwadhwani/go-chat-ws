package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/websocket"
)

type Server struct {
	clients map[*websocket.Conn]bool
}

func newServer() *Server {
	return &Server{
		clients: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) handleWS(client *websocket.Conn) {
	s.clients[client] = true
	s.readLoop(client)
}

func (s *Server) readLoop(client *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		fmt.Print("Reading from client")
		n, err := client.Read(buf)
		if err != nil {
			fmt.Println("Error reading from client")
			continue
		}
		msg := buf[:n]
		fmt.Printf("Received: %v", msg)

		s.broadcast(client, msg)
	}
}

func (s *Server) broadcast(currentClient *websocket.Conn, msg []byte) {
	for client := range s.clients {
		go func(conn *websocket.Conn) {
			if currentClient != client {
				if _, err := client.Write(msg); err != nil {
					fmt.Println("Error broadcasting")
				}
			}
		}(client)
	}
}

func main() {
	server := newServer()
	http.Handle("/ws", websocket.Handler(server.handleWS))
	http.ListenAndServe(":8080", nil)
}
