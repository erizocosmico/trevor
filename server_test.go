package trevor

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	server := NewServer(Config{
		Plugins:  dummyPlugins(),
		Port:     8888,
		Endpoint: "get_data",
	})

	go func() {
		server.Run()
	}()

	time.Sleep(5 * time.Millisecond)

	jsonStr := []byte(`{"text":"how are you?"}`)
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

	if resp.Status != "200 OK" {
		t.Errorf("expected status 200, got %s", resp.Status)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	if strings.TrimSpace(string(body)) != `{"data":"fine, and you?","error":false,"type":"salute"}` {
		t.Error("invalid response got")
	}
}
