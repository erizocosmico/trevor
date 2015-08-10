package trevor

import "net/http"

type MemoryService interface {
	// TokenForRequest returns an unique string representing the user of the request.
	TokenForRequest(*http.Request) string

	// DataForToken returns the user data associated to the token.
	DataForToken(string) (interface{}, error)

	// TokenHeader returns the name of the header used to send and receive the token.
	TokenHeader() string

	// NeededStore returns the name of the service needed to store the user data. An empty string means no store service is needed.
	NeededStore() string

	// SetStore sets the store and returns an error if the given service is not the desired one.
	SetStore(Service) error
}
