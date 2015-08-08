package trevor

import (
	"errors"
	"strings"
	"testing"
)

type salutePlugin struct{}

func (p *salutePlugin) Analyze(text string) (Score, interface{}) {
	if "how are you?" == strings.ToLower(text) {
		return NewScore(9, true), nil
	} else {
		return NewScore(0, false), nil
	}
}

func (p *salutePlugin) Process(text string, _ interface{}) (interface{}, error) {
	return "fine, and you?", nil
}

func (p *salutePlugin) Name() string {
	return "salute"
}

func (p *salutePlugin) Precedence() int {
	return 1
}

type fooPlugin struct{}

func (p *fooPlugin) Analyze(text string) (Score, interface{}) {
	return NewScore(5, false), nil
}

func (p *fooPlugin) Process(text string, _ interface{}) (interface{}, error) {
	return nil, errors.New("i always throw error")
}

func (p *fooPlugin) Name() string {
	return "foo"
}

func (p *fooPlugin) Precedence() int {
	return 1
}

func TestProcess(t *testing.T) {
	engine := NewEngine()
	engine.SetPlugins(dummyPlugins())

	dataType, data, err := engine.Process("how are you?")

	if err != nil {
		t.Error("unexpected error!")
	}

	if dataType != "salute" {
		t.Errorf("expected data type to be 'salute', '%s' found", dataType)
	}

	if data.(string) != "fine, and you?" {
		t.Errorf("expected data to be 'find, and you?' but was '%s'", data.(string))
	}
}

func TestProcessNoPlugins(t *testing.T) {
	engine := NewEngine()
	engine.SetPlugins(make([]Plugin, 0))

	_, _, err := engine.Process("how are you?")

	if err == nil {
		t.Error("expected error!")
	}
}

func dummyPlugins() []Plugin {
	return []Plugin{&salutePlugin{}, &fooPlugin{}}
}
