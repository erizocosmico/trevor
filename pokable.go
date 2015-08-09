package trevor

import "time"

// Pokable is a component (plugin or service) that needs to be poked every X time.
type Pokable interface {
	// PokeEvery returns the lapse between a poke and another.
	PokeEvery() time.Duration

	// Poke is the method that will be invoked after the desired time between pokes has passed. If Poke returns true the Pokable will stop being poked forever.
	Poke() bool
}

// PokablePlugins returns a list of Pokable with all the pokable plugins in the given list of plugins.
func PokablePlugins(plugins []Plugin) []Pokable {
	pokables := make([]Pokable, 0, len(plugins))

	for _, plugin := range plugins {
		if pokable, ok := plugin.(Pokable); ok {
			pokables = append(pokables, pokable)
		}
	}

	return pokables
}

// PokableServices returns a list of Pokable with all the pokable services in the given list of services.
func PokableServices(services []Service) []Pokable {
	pokables := make([]Pokable, 0, len(services))

	for _, service := range services {
		if pokable, ok := service.(Pokable); ok {
			pokables = append(pokables, pokable)
		}
	}

	return pokables
}

// RunPokeWorker runs a new worker that will run indefinitely poking the Pokable until it tells the worker to stop.
func RunPokeWorker(pokable Pokable) {
	for {
		time.Sleep(pokable.PokeEvery())

		if pokable.Poke() {
			break
		}
	}
}
