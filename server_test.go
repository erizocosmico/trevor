package trevor

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

func makeRequest(text string) (string, string) {
	server := NewServer(Config{
		Plugins:  dummyPlugins(),
		Port:     8888,
		Endpoint: "get_data",
	})

	go func() {
		server.Run()
	}()

	time.Sleep(5 * time.Millisecond)

	jsonStr := []byte(text)
	req, err := http.NewRequest("POST", "http://0.0.0.0:8888/get_data", bytes.NewBuffer(jsonStr))
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
	body, status := makeRequest(`{"text":"how are you?"}`)

	if status != "200 OK" {
		t.Errorf("expected status 200, got %s", status)
	}

	if strings.TrimSpace(body) != `{"data":"fine, and you?","error":false,"type":"salute"}` {
		t.Error("invalid response got")
	}
}

func TestRunNoText(t *testing.T) {
	_, status := makeRequest(`{"foo":"bar"}`)

	if status != "400 Bad Request" {
		t.Errorf("expected status 400, got %s", status)
	}
}

func TestRunTextEmpty(t *testing.T) {
	_, status := makeRequest(`{"text":""}`)

	if status != "400 Bad Request" {
		t.Errorf("expected status 400, got %s", status)
	}
}
