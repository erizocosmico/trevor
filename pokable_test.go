package trevor

import (
	"testing"
	"time"
)

func (p *fooPlugin) PokeEvery() time.Duration {
	return 50 * time.Millisecond
}

func (p *fooPlugin) Poke() bool {
	p.poked++

	return false
}

func (s *fooService) PokeEvery() time.Duration {
	return 50 * time.Millisecond
}

func (s *fooService) Poke() bool {
	s.poked++

	return s.poked > 2
}

func TestPokablePlugins(t *testing.T) {
	pokables := PokablePlugins([]Plugin{
		&barPlugin{},
		&fooPlugin{},
		&salutePlugin{},
	})

	if len(pokables) != 1 {
		t.Errorf("expected pokables length to be 1, %d received", len(pokables))
	}
}

func TestPokableServices(t *testing.T) {
	pokables := PokableServices([]Service{
		&barService{},
		&fooService{},
	})

	if len(pokables) != 1 {
		t.Errorf("expected pokables length to be 1, %d received", len(pokables))
	}
}

func TestRunPokeWorker(t *testing.T) {
	service := &fooService{}
	RunPokeWorker(service)
	if service.poked != 3 {
		t.Errorf("expected service to be poked 5 times but was poked just %d", service.poked)
	}
}
