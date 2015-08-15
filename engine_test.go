package trevor

import (
	"errors"
	"strings"
	"testing"
	"time"
)

//
// Test plugins
//

type salutePlugin struct{}

func (p *salutePlugin) Analyze(req *Request) (Score, interface{}) {
	if "how are you?" == strings.ToLower(req.Text) {
		return NewScore(9, true), nil
	} else {
		return NewScore(0, false), nil
	}
}

func (p *salutePlugin) Process(req *Request, _ interface{}) (interface{}, error) {
	return "fine, and you?", nil
}

func (p *salutePlugin) Name() string {
	return "salute"
}

func (p *salutePlugin) Precedence() int {
	return 2
}

type fooPlugin struct {
	poked int
}

func (p *fooPlugin) Analyze(req *Request) (Score, interface{}) {
	return NewScore(5, false), nil
}

func (p *fooPlugin) Process(req *Request, _ interface{}) (interface{}, error) {
	return nil, errors.New("i always throw error")
}

func (p *fooPlugin) Name() string {
	return "foo"
}

func (p *fooPlugin) Precedence() int {
	return 1
}

type barPlugin struct {
	service *barService
}

func (p *barPlugin) Analyze(req *Request) (Score, interface{}) {
	return NewScore(0, false), nil
}

func (p *barPlugin) Process(req *Request, _ interface{}) (interface{}, error) {
	return nil, nil
}

func (p *barPlugin) Name() string {
	return "bar"
}

func (p *barPlugin) Precedence() int {
	return 1
}

func (p *barPlugin) NeededServices() []string {
	return []string{"bar"}
}

func (p *barPlugin) SetService(name string, service Service) {
	if name == "bar" {
		p.service = service.(*barService)
	}
}

//
// Test services
//

type fooService struct {
	poked int
}

func (s *fooService) Name() string {
	return "foo"
}

func (s *fooService) SetName(_ string) {
}

type barService struct{}

func (s *barService) Name() string {
	return "bar"
}

func (s *barService) SetName(_ string) {
}

//
// Tests
//

func TestProcess(t *testing.T) {
	engine := NewEngine()
	engine.SetPlugins(dummyPlugins())

	dataType, data, err := engine.Process(NewRequest("how are you?", nil))

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

	_, _, err := engine.Process(NewRequest("how are you?", nil))

	if err == nil {
		t.Error("expected error!")
	}
}

func TestSetServices(t *testing.T) {
	e := NewEngine().(*engine)
	e.SetServices(dummyServices())

	for _, s := range []string{"foo", "bar"} {
		if service, ok := e.services[s]; !ok || service.Name() != s {
			t.Errorf("expected to find %s service", s)
		}
	}
}

func TestInjectServices(t *testing.T) {
	e := NewEngine().(*engine)
	e.SetServices([]Service{&barService{}})
	e.SetPlugins([]Plugin{&barPlugin{}})

	plugin := e.plugins[0].(*barPlugin)
	if plugin.service == nil || plugin.service.Name() != "bar" {
		t.Errorf("expected to find service bar in bar plugin")
	}
}

func TestInjectServicesServiceUnknown(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected a panic!")
		}
	}()

	e := NewEngine()
	e.SetPlugins([]Plugin{&barPlugin{}})
}

func TestSchedulePokes(t *testing.T) {
	e := NewEngine().(*engine)
	e.SetPlugins([]Plugin{&fooPlugin{}})
	e.SetServices([]Service{&fooService{}})

	e.SchedulePokes()

	time.Sleep(250 * time.Millisecond)

	if e.plugins[0].(*fooPlugin).poked < 4 {
		t.Errorf("plugin should have been poked at least 4 times")
	}

	if e.services["foo"].(*fooService).poked < 3 {
		t.Errorf("service should have been poked 3 times")
	}
}

func TestAnalyzer(t *testing.T) {
	e := NewEngine().(*engine)
	e.SetPlugins(dummyPlugins())
	e.SetServices(dummyServices())
	e.SetAnalyzer(func(req *Request) (string, interface{}) {
		return "foo", nil
	})

	plugin, data, err := e.Process(NewRequest("how are you?", nil))
	if err == nil || data != nil || plugin != "foo" {
		t.Errorf("expected foo plugin to process but %s plugin did", plugin)
	}
}

func TestSetMemoryServiceWithoutStore(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected a panic!")
		}
	}()

	e := NewEngine().(*engine)
	e.SetPlugins(dummyPlugins())
	e.SetServices([]Service{&memoryService{}})
}

func TestMiddleware(t *testing.T) {
	afterCalled := 0
	beforeCalled := 0

	mw := func(req *Request, getService func(string) Service, next func() (string, interface{}, error)) (string, interface{}, error) {
		defer func() {
			afterCalled++
		}()

		getService("foo").(*fooService).poked++

		beforeCalled++
		return next()
	}

	e := NewEngine().(*engine)
	e.SetServices([]Service{&fooService{}})
	e.SetPlugins(dummyPlugins())
	e.SetMiddleware([]Middleware{mw, mw, mw})

	dataType, data, err := e.Process(NewRequest("how are you?", nil))

	if err != nil {
		t.Error("unexpected error!")
	}

	if dataType != "salute" {
		t.Errorf("expected data type to be 'salute', '%s' found", dataType)
	}

	if data.(string) != "fine, and you?" {
		t.Errorf("expected data to be 'find, and you?' but was '%s'", data.(string))
	}

	if afterCalled != 3 {
		t.Errorf("expected middleware to have been called 3 times, called %d instead", afterCalled)
	}

	if beforeCalled != 3 {
		t.Errorf("expected middleware to have been called 3 times, called %d instead", beforeCalled)
	}

	if e.services["foo"].(*fooService).poked != 3 {
		t.Errorf("expected foo service to have been poked 3 times by the middleware, poked %d instead", e.services["foo"].(*fooService).poked)
	}
}

//
// Helper functions
//

func dummyPlugins() []Plugin {
	return []Plugin{&fooPlugin{}, &salutePlugin{}}
}

func dummyServices() []Service {
	return []Service{&fooService{}, &barService{}}
}
