package trevor

// Plugin defines the base functionality that a trevor plugin should have
type Plugin interface {
	// Analyze takes a string with the text and after specific analysis returns
	// a Score
	Analyze(string) Score

	// Process takes a string with the text and processes that text to return data that will be sent to the client
	Process(string) (interface{}, error)

	// Name returns the name of the plugin
	Name() string

	// Precedence gets the current precedence of the plugin
	Precedence() int
}
