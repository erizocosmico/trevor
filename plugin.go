package trevor

import "sort"

// Plugin defines the base functionality that a trevor plugin should have
type Plugin interface {
	// Analyze takes the request to process and after specific analysis returns
	// a Score and metadata that will later be passed to process
	Analyze(*Request) (Score, interface{})

	// Process takes the request to process and the metadata returned by Analyze method and processes that text to return data that will be sent to the client
	Process(*Request, interface{}) (interface{}, error)

	// Name returns the name of the plugin
	Name() string

	// Precedence gets the current precedence of the plugin
	Precedence() int
}

// InjectablePlugin is a plugin that requests dependency injection
type InjectablePlugin interface {
	// NeededServices returns an array with the name of all needed services.
	NeededServices() []string

	// SetService injects a service to the plugin
	SetService(string, Service)
}

type byPluginPrecedence []Plugin

func (b byPluginPrecedence) Len() int {
	return len(b)
}

func (b byPluginPrecedence) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b byPluginPrecedence) Less(i, j int) bool {
	return b[i].Precedence() > b[j].Precedence()
}

// SortPlugins sorts all plugins in an array by precedence
func SortPlugins(plugins []Plugin) {
	sort.Sort(byPluginPrecedence(plugins))
}
