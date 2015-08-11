package trevor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func startServer() {
	server := NewServer(Config{
		Plugins:        []Plugin{&rememberPlugin{}},
		Services:       []Service{&memoryService{}, &storeService{map[string]int{}}},
		Port:           8884,
		Endpoint:       "get_data",
		InputFieldName: "input",
	})

	go func() {
		server.Run()
	}()

	time.Sleep(5 * time.Millisecond)
}

func makeRequestWithHeader(header string) (string, string) {
	jsonStr := []byte(`{"input":"how are you?"}`)
	req, err := http.NewRequest("POST", "http://0.0.0.0:8884/get_data", bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	if header != "" {
		req.Header.Set("X-Memory-Token", header)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var data map[string]interface{}
	if err = json.Unmarshal(body, &data); err != nil {
		panic(err)
	}

	return data["data"].(string), resp.Header.Get("X-Memory-Token")
}

type rememberPlugin struct {
	memory MemoryService
}

func (p *rememberPlugin) Analyze(req *Request) (Score, interface{}) {
	return NewScore(10, false), nil
}

func (p *rememberPlugin) Process(req *Request, _ interface{}) (interface{}, error) {
	count, _ := p.memory.DataForToken(req.Token)
	if intCount, ok := count.(int); ok {
		return fmt.Sprintf("visit number %d", intCount), nil
	}

	req.Token = p.memory.TokenForRequest(req.Request)

	return "hello new visitor", nil
}

func (p *rememberPlugin) Name() string {
	return "remember"
}

func (p *rememberPlugin) Precedence() int {
	return 1
}

func (p *rememberPlugin) NeededServices() []string {
	return []string{"memory"}
}

func (p *rememberPlugin) SetService(name string, service Service) {
	if name == "memory" {
		p.memory = service.(MemoryService)
	}
}

type storeService struct {
	storage map[string]int
}

func (s *storeService) Name() string {
	return "store"
}

func (s *storeService) SetName(name string) {
}

type memoryService struct {
	store     *storeService
	userCount int
}

func (s *memoryService) Name() string {
	return "memory"
}

func (s *memoryService) SetName(name string) {
}

func (s *memoryService) TokenForRequest(req *http.Request) string {
	if _, ok := s.store.storage[req.Header.Get(s.TokenHeader())]; ok {
		return req.Header.Get(s.TokenHeader())
	} else {
		s.userCount++
		token := fmt.Sprintf("token_%d", s.userCount)
		s.store.storage[token] = 0
		return token
	}
}

func (s *memoryService) DataForToken(token string) (interface{}, error) {
	if _, ok := s.store.storage[token]; ok {
		s.store.storage[token]++
		return s.store.storage[token], nil
	} else {
		return nil, nil
	}
}

func (s *memoryService) TokenHeader() string {
	return "X-Memory-Token"
}

func (s *memoryService) NeededStore() string {
	return "store"
}

func (s *memoryService) SetStore(store Service) error {
	s.store = store.(*storeService)
	return nil
}

//
// Tests
//

func TestMemoryService(t *testing.T) {
	startServer()
	assertRememberRequest("", "hello new visitor", "token_1", t)
	assertRememberRequest("token_1", "visit number 1", "token_1", t)
	assertRememberRequest("token_1", "visit number 2", "token_1", t)
	assertRememberRequest("", "hello new visitor", "token_2", t)
	assertRememberRequest("token_1", "visit number 3", "token_1", t)
}

func assertRememberRequest(token, expectedData, expectedToken string, t *testing.T) {
	data, token := makeRequestWithHeader(token)
	if data != expectedData || token != expectedToken {
		t.Errorf("expecting data '%s', got '%s'. Expecting token '%a', got '%s'", expectedData, data, expectedToken, token)
	}
}
