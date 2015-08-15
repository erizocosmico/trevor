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

	// SetAnalyzer sets the Analyzer function of the engine.
	SetAnalyzer(Analyzer)

	// Process takes the current request to process and returns the name of the plugin that
	// processed the text and the data returned by it.
	Process(*Request) (string, interface{}, error)

	// SchedulePokes schedules all pokes to run indefinitely.
	SchedulePokes()

	// Memory returns the memory service if any.
	Memory() MemoryService
}

// Analyzer is a function that takes the current request to process and returns the name of the plugin that should process it and metadata.
type Analyzer func(*Request) (string, interface{})

type engine struct {
	plugins   []Plugin
	pluginMap map[string]int
	services  map[string]Service
	analyzer  Analyzer
	memory    MemoryService
}

// NewEngine creates a new Engine instance
func NewEngine() Engine {
	return &engine{
		services:  map[string]Service{},
		pluginMap: map[string]int{},
	}
}

func (e *engine) SetPlugins(plugins []Plugin) {
	SortPlugins(plugins)
	e.plugins = e.injectServices(plugins)

	for i, p := range e.plugins {
		e.pluginMap[p.Name()] = i
	}
}

func (e *engine) getPlugin(name string) Plugin {
	return e.plugins[e.pluginMap[name]]
}

func (e *engine) SetServices(services []Service) {
	for _, service := range services {
		e.services[service.Name()] = service
	}

	e.setMemoryService()
}

func (e *engine) setMemoryService() {
	if service, ok := e.services["memory"]; ok {
		if memoryService, isMemoryService := service.(MemoryService); isMemoryService && service.Name() == "memory" {
			storeName := memoryService.NeededStore()
			if store, ok := e.services[storeName]; ok || storeName == "" {
				memoryService.SetStore(store)
			} else {
				panic(errors.New("service " + storeName + " not found but is required by memory service"))
			}

			e.memory = memoryService
		}
	}
}

func (e *engine) SetAnalyzer(analyzer Analyzer) {
	e.analyzer = analyzer
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

func (e *engine) Process(req *Request) (string, interface{}, error) {
	if len(e.plugins) == 0 {
		return "", nil, errors.New("no plugins found. can't process anything")
	}

	var bestResult analysisResult
	if e.analyzer == nil {
		bestResult = getBestResult(getResults(e.plugins, req))
	} else {
		name, metadata := e.analyzer(req)
		bestResult = analysisResult{name: name, metadata: metadata}
	}

	var chosenPlugin = e.getPlugin(bestResult.name)
	data, err := chosenPlugin.Process(req, bestResult.metadata)
	return chosenPlugin.Name(), data, err
}

func (e *engine) Memory() MemoryService {
	return e.memory
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
