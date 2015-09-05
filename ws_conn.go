package server

import "golang.org/x/net/websocket"

type WSConn struct {
	ws *websocket.Conn
}

// ReadMessage takes the json message and returns the values or error
func (c *WSConn) ReadMessage() (Message, error) {
	msg := Message{}
	err := websocket.JSON.Receive(c.ws, &msg)
	return msg, err
}

// WriteMessage takes a json value and returns it to the caller over a websocket
func (c *WSConn) WriteMessage(msg interface{}) error {
	return websocket.JSON.Send(c.ws, msg)
}
