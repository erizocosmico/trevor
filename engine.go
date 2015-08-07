package trevor

import "errors"

type Engine interface {
	// SetPlugins sets the list of plugins of the engine.
	SetPlugins([]Plugin)

	// Process takes the text to process and returns the name of the plugin that
	// processed the text and the data returned by it.
	Process(string) (string, interface{}, error)
}

type engine struct {
	plugins []Plugin
}

// NewEngine creates a new Engine instance
func NewEngine() Engine {
	return &engine{}
}

func (e *engine) SetPlugins(plugins []Plugin) {
	e.plugins = plugins
}

func (e *engine) Process(text string) (string, interface{}, error) {
	if len(e.plugins) == 0 {
		return "", nil, errors.New("no plugins found. can't process anything")
	}

	results := make([]analysisResult, len(e.plugins))
	for i, plugin := range e.plugins {
		score := plugin.Analyze(text)
		results[i] = newAnalysisResult(score.Score(), score.IsExactMatch(), plugin.Precedence(), plugin.Name())
	}

	bestResult := getBestResult(results)

	var chosenPlugin Plugin
	for _, plugin := range e.plugins {
		if plugin.Name() == bestResult.name {
			chosenPlugin = plugin
			break
		}
	}

	data, err := chosenPlugin.Process(text)
	return chosenPlugin.Name(), data, err
}
