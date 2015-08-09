package trevor

import (
	"errors"
	"fmt"
)

type Engine interface {
	// SetPlugins sets the list of plugins of the engine.
	SetPlugins([]Plugin)

	// SetServices sets the list of services of the engine.
	SetServices([]Service)

	// Process takes the text to process and returns the name of the plugin that
	// processed the text and the data returned by it.
	Process(string) (string, interface{}, error)

	// SchedulePokes schedules all pokes to run indefinitely.
	SchedulePokes()
}

type engine struct {
	plugins  []Plugin
	services map[string]Service
}

// NewEngine creates a new Engine instance
func NewEngine() Engine {
	return &engine{services: map[string]Service{}}
}

func (e *engine) SetPlugins(plugins []Plugin) {
	SortPlugins(plugins)
	e.plugins = e.injectServices(plugins)
}

func (e *engine) SetServices(services []Service) {
	for _, service := range services {
		e.services[service.Name()] = service
	}
}

func (e *engine) injectServices(plugins []Plugin) []Plugin {
	for _, plugin := range plugins {
		if injectablePlugin, ok := plugin.(InjectablePlugin); ok {
			for _, serviceName := range injectablePlugin.NeededServices() {
				service, ok := e.services[serviceName]
				if !ok {
					panic(fmt.Sprintf("unknown service with name: %s", serviceName))
				}

				injectablePlugin.SetService(serviceName, service)
			}
		}
	}

	return plugins
}

func (e *engine) Process(text string) (string, interface{}, error) {
	if len(e.plugins) == 0 {
		return "", nil, errors.New("no plugins found. can't process anything")
	}

	results := make([]analysisResult, len(e.plugins))
	for i, plugin := range e.plugins {
		score, metadata := plugin.Analyze(text)
		results[i] = newAnalysisResult(score.Score(), score.IsExactMatch(), plugin.Precedence(), plugin.Name(), metadata)
	}

	bestResult := getBestResult(results)

	var chosenPlugin Plugin
	for _, plugin := range e.plugins {
		if plugin.Name() == bestResult.name {
			chosenPlugin = plugin
			break
		}
	}

	data, err := chosenPlugin.Process(text, bestResult.metadata)
	return chosenPlugin.Name(), data, err
}

func (e *engine) SchedulePokes() {
	for _, pp := range PokablePlugins(e.plugins) {
		go RunPokeWorker(pp)
	}

	var services = make([]Service, 0, len(e.services))
	for _, service := range e.services {
		services = append(services, service)
	}

	for _, ps := range PokableServices(services) {
		go RunPokeWorker(ps)
	}
}
