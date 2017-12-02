package bootstrap

import (
	"encoding/json"

	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
)

// MessageOut represents a message going out
type MessageOut struct {
	Name    string      `json:"name"`
	Payload interface{} `json:"payload,omitempty"`
}

// MessageIn represents a message going in
type MessageIn struct {
	Name    string          `json:"name"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

// handleMessages handles messages
func handleMessages(w *astilectron.Window, messageHandler MessageHandler) astilectron.ListenerMessage {
	return func(e astilectron.Event) (v interface{}) {
		// Unmarshal message
		var i MessageIn
		var err error
		if err = e.Message.Unmarshal(&i); err != nil {
			astilog.Error(errors.Wrapf(err, "unmarshaling message %+v failed", *e.Message))
			return
		}

		// Handle message
		var p interface{}
		if p, err = messageHandler(w, i); err != nil {
			astilog.Error(errors.Wrapf(err, "handling message %+v failed", i))
		}

		// Return message
		if p != nil {
			o := &MessageOut{Name: i.Name + ".callback", Payload: p}
			if err != nil {
				o.Name = "error"
			}
			v = o
		}
		return
	}
}
