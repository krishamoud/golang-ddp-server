package server

import "errors"

type MethodContext struct {
	ID      string
	Params  []interface{}
	Conn    Connection
	Done    bool
	Updated bool
}

// NewMethodContext takes a message and connection then assigns the message fields to
// the MethodContext fields and returns the MethodContext
func NewMethodContext(m *Message, conn Connection) MethodContext {
	ctx := MethodContext{}
	ctx.ID = m.ID
	ctx.Params = m.Params
	ctx.Conn = conn
	return ctx
}

// SendResult takes the result of the method and sends it over a websocket
// back to the caller
func (ctx *MethodContext) SendResult(result interface{}) error {
	if ctx.Done {
		err := errors.New("results already sent")
		return err
	}

	ctx.Done = true
	msg := map[string]interface{}{
		"msg":    "result",
		"id":     ctx.ID,
		"result": result,
	}

	return ctx.Conn.WriteMessage(msg)
}

// SendError sends an error back to the method caller via websocket
func (ctx *MethodContext) SendError(e string) error {
	if ctx.Done {
		err := errors.New("already sent results for method")
		return err
	}

	ctx.Done = true
	msg := map[string]interface{}{
		"msg": "result",
		"id":  ctx.ID,
		"error": map[string]string{
			"error": e,
		},
	}

	return ctx.Conn.WriteMessage(msg)
}

// Method calls can affect data that the client is subscribed to.
// Once the server has finished sending the client all the relevant data messages
// based on this procedure call, the server should send an updated message to the 
// client with this method's ID.
func (ctx *MethodContext) SendUpdated() error {
	if ctx.Updated {
		err := errors.New("already sent updated for method")
		return err
	}

	ctx.Updated = true
	msg := map[string]interface{}{
		"msg":     "updated",
		"methods": []string{ctx.ID},
	}

	return ctx.Conn.WriteMessage(msg)
}
