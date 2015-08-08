package trevor

// Service is a type of component that is passed to the plugins.
// A plugin can request 0 or more services that will be passed to it
// on startup time.
type Service interface {
	// Name returns the name of the service
	Name() string

	// SetName sets the name to the service
	SetName(string)
}
