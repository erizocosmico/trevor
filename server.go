package trevor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"unicode/utf8"
)

// Server is a Trevor server ready to run
type Server interface {
	// Run starts the server.
	Run() error

	// GetEngine returns the current Engine being used on the server.
	GetEngine() Engine
}

type server struct {
	engine Engine
	config Config
}

func NewServer(config Config) Server {
	engine := NewEngine()
	engine.SetServices(config.Services)
	engine.SetPlugins(config.Plugins)
	engine.SetMiddleware(config.Middleware)

	return &server{
		engine: engine,
		config: config,
	}
}

func (s *server) GetEngine() Engine {
	return s.engine
}

func (s *server) Run() error {
	var (
		endpoint  = "process"
		inputName = "text"
	)

	if s.config.Endpoint != "" {
		endpoint = s.config.Endpoint
	}

	if s.config.InputFieldName != "" {
		inputName = s.config.InputFieldName
	}

	router := http.NewServeMux()
	router.HandleFunc("/"+endpoint, processHandler(inputName, endpoint, s))

	s.engine.SchedulePokes()

	var err error
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	if !s.config.Secure {
		err = http.ListenAndServe(addr, router)
	} else {
		err = http.ListenAndServeTLS(addr, s.config.CertPerm, s.config.KeyPerm, router)
	}

	return err
}

func processHandler(inputName, endpoint string, s *server) func(http.ResponseWriter, *http.Request) {
	var errorText = inputName + " field is mandatory and can not be empty"

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			var (
				jsonInput map[string]string
				response  map[string]interface{}
				status    int
			)

			content, err := ioutil.ReadAll(r.Body)
			if err = json.Unmarshal(content, &jsonInput); err == nil {
				text, ok := jsonInput[inputName]
				if ok && utf8.RuneCountInString(strings.TrimSpace(text)) > 0 {
					req := NewRequest(strings.TrimSpace(text), r)
					if s.engine.Memory() != nil {
						req.Token = r.Header.Get(s.engine.Memory().TokenHeader())
					}

					dataType, data, err := s.engine.Process(req)
					if err != nil {
						errorText = err.Error()
					} else {
						if s.engine.Memory() != nil {
							w.Header().Set(s.engine.Memory().TokenHeader(), req.Token)
						}

						response = map[string]interface{}{
							"error": false,
							"type":  dataType,
							"data":  data,
						}
						status = http.StatusOK
					}
				}
			}

			if status == 0 {
				response = map[string]interface{}{
					"error":   true,
					"message": errorText,
				}
				status = http.StatusBadRequest
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)
			resp, _ := json.Marshal(response)
			w.Write(resp)
		} else {
			http.NotFound(w, r)
		}
	}
}
