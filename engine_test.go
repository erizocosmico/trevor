package trevor

import (
	"strings"
	"testing"
)

type salutePlugin struct{}

func (p *salutePlugin) Analyze(text string) Score {
	if "how are you?" == strings.ToLower(text) {
		return NewScore(10, true)
	} else {
		return NewScore(0, false)
	}
}

func (p *salutePlugin) Process(text string) (interface{}, error) {
	return "fine, and you?", nil
}

func (p *salutePlugin) Name() string {
	return "salute"
}

func (p *salutePlugin) Precedence() int {
	return 1
}

type fooPlugin struct{}

func (p *fooPlugin) Analyze(text string) Score {
	return NewScore(5, false)
}

func (p *fooPlugin) Process(text string) (interface{}, error) {
	return "foo", nil
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

func dummyPlugins() []Plugin {
	return []Plugin{&salutePlugin{}, &fooPlugin{}}
}
