package server

import (
	"fmt"
	"io"
	"net/http"
	"golang.org/x/net/websocket"
)

// Declare the Server type which has a map of strings that correlate to
// methods that will be written by the application writer
type Server struct {
	methods map[string]MethodHandler
}

// This will instantiate a server, instantiate the methods field,
// then return the server
func New() Server {
	s := Server{}
	s.methods = make(map[string]MethodHandler)
	return s
}

// Listen takes a string which is the address it will listen for messages from
// A websocket server is instantiated with the handler defined below.
// Currently the handshake accepts all connections
// http.Handle reads the messages being sent to the /websocket route and handles
// them properly
// ListenAndServe starts an HTTP server with the given address
func (s *Server) Listen(addr string) error {
	wsServer := websocket.Server{Handler: s.wsHandler, Handshake: s.handshake}
	http.Handle("/websocket", wsServer)
	return http.ListenAndServe(addr, nil)
}

// Method takes a string to create a key with the value being the function that
// needs to be execute
func (s *Server) Method(name string, fn MethodHandler) {
	s.methods[name] = fn
}

// wsHandler takes a websocket connection and for as long as it doesn't err
// it will handle the messages appropriately
// if an err occurs it will close the connection
func (s *Server) wsHandler(ws *websocket.Conn) {
	conn := WSConn{ws}
	defer ws.Close()

	for {
		msg, err := conn.ReadMessage()

		if err != nil {
			if err != io.EOF {
				fmt.Println("Error (Read Error):", err, msg)
			}

			break
		}

		s.handleMessage(&conn, &msg)
	}
}

// Based on the DDP spec outlined here: https://github.com/meteor/meteor/blob/devel/packages/ddp/DDP.md
// the server can handle "connect", "ping", and "method" messages.
// Server Messages not supported currently are "failed", "pong", "nosub",
// "added", "changed", "removed", "ready", "addedBefore", "movedBefore", "result",
// "updated"
func (s *Server) handleMessage(conn Connection, msg *Message) {
	switch msg.Msg {
	case "connect":
		s.handleConnect(conn, msg)
	case "ping":
		s.handlePing(conn, msg)
	case "method":
		s.handleMethod(conn, msg)
	default:
		fmt.Println("Error (Unknown Message Type):", msg)
		// TODO => send "error" ddp message
		break
	}
}

// handleConnect recieves the message "connect" and responds with "connected"
func (s *Server) handleConnect(conn Connection, m *Message) {
	msg := map[string]string{
		"msg":     "connected",
		"session": s.Id(17),
	}

	conn.WriteMessage(msg)
}

// handleMethod recieves the message "method", finds the function correlated with
// the message method name, creates a new method context then executes the method in a
// goroutine
// if the method is not found it will return an err
func (s *Server) handleMethod(conn Connection, m *Message) {
	fn, ok := s.methods[m.Method]
	if !ok {
		fmt.Println("Error: (Method Not Found)", m.Method)
		return
	}

	ctx := NewMethodContext(m, conn)
	go fn(ctx)
}

// handlePing recieves the message "ping" and responds "pong"
func (s *Server) handlePing(conn Connection, m *Message) {
	msg := map[string]string{
		"msg": "pong",
	}
	if m.ID != "" {
		msg["id"] = m.ID
	}

	conn.WriteMessage(msg)
}

// handshake recieves a request and responds nil
func (s *Server) handshake(config *websocket.Config, req *http.Request) error {
	// accept all connections
	return nil
}
