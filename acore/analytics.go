package acore

import (
	"io"
)

// This interface is the main API exposed by the analytics package.
// Values that satsify this interface are returned by the client constructors
// provided by the package and provide a way to send messages via the HTTP API.
type CClient interface {
	io.Closer

	// Queues a message to be sent by the client when the conditions for a batch
	// upload are met.
	// This is the main method you'll be using, a typical flow would look like
	// this:
	//
	//	client := analytics.New(writeKey)
	//	...
	//	client.Enqueue(analytics.Track{ ... })
	//	...
	//	client.Close()
	//
	// The method returns an error if the message queue not be queued, which
	// happens if the client was already closed at the time the method was
	// called or if the message was malformed.
	Enqueue(Message) error
}

func dereferenceMessage(msg Message) Message {
	switch m := msg.(type) {
	case *Alias:
		if m == nil {
			return nil
		}

		return *m
	case *Group:
		if m == nil {
			return nil
		}

		return *m
	case *Identify:
		if m == nil {
			return nil
		}

		return *m
	case *Page:
		if m == nil {
			return nil
		}

		return *m
	case *Screen:
		if m == nil {
			return nil
		}

		return *m
	case *Track:
		if m == nil {
			return nil
		}

		return *m
	}

	return msg
}

// func (c *Client) Enqueue(msg Message) (err error) {

// 	msg = dereferenceMessage(msg)
// 	if err = msg.Validate(); err != nil {
// 		return
// 	}

// 	var id = c.uid()
// 	var ts = c.now()

// 	switch m := msg.(type) {
// 	case Alias:
// 		m.Type = "alias"
// 		m.MessageId = makeMessageId(m.MessageId, id)
// 		m.Timestamp = makeTimestamp(m.Timestamp, ts)
// 		if m.Context == nil {
// 			m.Context = makeContext()
// 		}
// 		msg = m

// 	case Group:
// 		m.Type = "group"
// 		m.MessageId = makeMessageId(m.MessageId, id)
// 		if m.AnonymousId == "" {
// 			m.AnonymousId = makeAnonymousId(m.UserId)
// 		}
// 		m.Timestamp = makeTimestamp(m.Timestamp, ts)
// 		if m.Context == nil {
// 			m.Context = makeContext()
// 		}
// 		msg = m

// 	case Identify:
// 		m.Type = "identify"
// 		m.MessageId = makeMessageId(m.MessageId, id)
// 		if m.AnonymousId == "" {
// 			m.AnonymousId = makeAnonymousId(m.UserId)
// 		}
// 		m.Timestamp = makeTimestamp(m.Timestamp, ts)
// 		if m.Context == nil {
// 			m.Context = makeContext()
// 		}
// 		msg = m

// 	case Page:
// 		m.Type = "page"
// 		m.MessageId = makeMessageId(m.MessageId, id)
// 		if m.AnonymousId == "" {
// 			m.AnonymousId = makeAnonymousId(m.UserId)
// 		}
// 		m.Timestamp = makeTimestamp(m.Timestamp, ts)
// 		if m.Context == nil {
// 			m.Context = makeContext()
// 		}
// 		msg = m

// 	case Screen:
// 		m.Type = "screen"
// 		m.MessageId = makeMessageId(m.MessageId, id)
// 		if m.AnonymousId == "" {
// 			m.AnonymousId = makeAnonymousId(m.UserId)
// 		}
// 		m.Timestamp = makeTimestamp(m.Timestamp, ts)
// 		if m.Context == nil {
// 			m.Context = makeContext()
// 		}
// 		msg = m

// 	case Track:
// 		m.Type = "track"
// 		m.MessageId = makeMessageId(m.MessageId, id)
// 		if m.AnonymousId == "" {
// 			m.AnonymousId = makeAnonymousId(m.UserId)
// 		}
// 		m.Timestamp = makeTimestamp(m.Timestamp, ts)
// 		if m.Context == nil {
// 			m.Context = makeContext()
// 		}
// 		msg = m

// 	default:
// 		err = fmt.Errorf("messages with custom types cannot be enqueued: %T", msg)
// 		return
// 	}

// 	defer func() {
// 		// When the `msgs` channel is closed writing to it will trigger a panic.
// 		// To avoid letting the panic propagate to the caller we recover from it
// 		// and instead report that the client has been closed and shouldn't be
// 		// used anymore.
// 		if recover() != nil {
// 			err = ErrClosed
// 		}
// 	}()

// 	c.msgs <- msg
// 	return
// }
