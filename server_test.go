package trevor

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

func makeRequest(text string, port int) (string, string) {
	return makeRequestWithMethod(text, port, "POST")
}

func makeRequestWithMethod(text string, port int, method string) (string, string) {
	server := NewServer(Config{
		Plugins:        dummyPlugins(),
		Port:           port,
		Endpoint:       "get_data",
		InputFieldName: "input",
		CORSOrigin:     "*",
	})

	go func() {
		server.Run()
	}()

	time.Sleep(5 * time.Millisecond)

	jsonStr := []byte(text)
	req, err := http.NewRequest(method, fmt.Sprintf("http://0.0.0.0:%d/get_data", port), bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	return string(body), resp.Status
}

func TestRun(t *testing.T) {
	body, status := makeRequest(`{"input":"how are you?"}`, 9091)

	if status != "200 OK" {
		t.Errorf("expected status 200, got %s", status)
	}

	if strings.TrimSpace(body) != `{"data":"fine, and you?","error":false,"type":"salute"}` {
		t.Error("invalid response got")
	}
}

func TestRunNoText(t *testing.T) {
	_, status := makeRequest(`{"foo":"bar"}`, 9092)

	if status != "400 Bad Request" {
		t.Errorf("expected status 400, got %s", status)
	}
}

func TestRunTextEmpty(t *testing.T) {
	_, status := makeRequest(`{"input":""}`, 9093)

	if status != "400 Bad Request" {
		t.Errorf("expected status 400, got %s", status)
	}
}

func TestRunPluginError(t *testing.T) {
	_, status := makeRequest(`{"input":"foo"}`, 9094)

	if status != "400 Bad Request" {
		t.Errorf("expected status 400, got %s", status)
	}
}

// This is just for code coverage
func TestGetEngine(t *testing.T) {
	server := NewServer(Config{
		Plugins:  dummyPlugins(),
		Port:     9095,
		Endpoint: "get_data",
	})
	server.GetEngine()
}

// Just for code coverage too
func TestNotFound(t *testing.T) {
	_, status := makeRequestWithMethod(`whatever`, 9097, "GET")
	if status != "404 Not Found" {
		t.Errorf("expected error 404, %s received", status)
	}
}

func TestOptions(t *testing.T) {
	_, status := makeRequestWithMethod(`whatever`, 9099, "OPTIONS")
	if status != "200 OK" {
		t.Errorf("expected status 200 OK, %s received", status)
	}
}

// Just for code coverage too
func TestRunSecure(t *testing.T) {
	server := NewServer(Config{
		Plugins:  dummyPlugins(),
		Port:     9096,
		Endpoint: "get_data",
		Secure:   true,
	})
	server.Run()
}
