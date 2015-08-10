package trevor

import "net/http"

// Request is the context of a request to the trevor engine.
type Request struct {
	// Text is the text that came with the request.
	Text string

	// Request is the current HTTP request.
	Request *http.Request

	// Token is the associated token from the request. Will be empty
	// on the first request. This value will be sent to the client.
	// So, if you want the client to have the token you should
	// change this value in that case on your plugin, service, ...
	Token string
}

// NewRequest creates a new request instance.
func NewRequest(text string, req *http.Request) *Request {
	return &Request{
		Text:    text,
		Request: req,
	}
}
