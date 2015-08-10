package trevor

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"unicode/utf8"
)

// Server is a Trevor server ready to run
type Server interface {
	// Run starts the server.
	Run()

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

	return &server{
		engine: engine,
		config: config,
	}
}

func (s *server) GetEngine() Engine {
	return s.engine
}

func (s *server) Run() {
	var (
		router    = gin.Default()
		endpoint  = "process"
		inputName = "text"
		errorText = inputName + " field is mandatory and can not be empty"
	)

	if s.config.Endpoint != "" {
		endpoint = s.config.Endpoint
	}

	if s.config.InputFieldName != "" {
		inputName = s.config.InputFieldName
	}

	router.POST("/"+endpoint, func(c *gin.Context) {
		var json map[string]string

		if c.BindJSON(&json) == nil {
			text, ok := json[inputName]
			if ok && utf8.RuneCountInString(strings.TrimSpace(text)) > 0 {
				req := NewRequest(strings.TrimSpace(text), c.Request)
				if s.engine.Memory() != nil {
					req.Token = c.Request.Header.Get(s.engine.Memory().TokenHeader())
				}

				dataType, data, err := s.engine.Process(req)
				if err != nil {
					errorText = err.Error()
				} else {
					if s.engine.Memory() != nil {
						c.Header(s.engine.Memory().TokenHeader(), req.Token)
					}

					c.JSON(http.StatusOK, gin.H{
						"error": false,
						"type":  dataType,
						"data":  data,
					})
					return
				}
			}
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": errorText,
		})
	})

	s.engine.SchedulePokes()

	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	if !s.config.Secure {
		router.Run(addr)
	} else {
		router.RunTLS(addr, s.config.CertPerm, s.config.KeyPerm)
	}
}
