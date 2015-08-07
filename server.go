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

type Request struct {
	Text string `form:"text" json:"text" binding:"required"`
}

type server struct {
	engine Engine
	config Config
}

func NewServer(config Config) Server {
	engine := NewEngine()
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
	router := gin.Default()
	var endpoint = "process"
	if s.config.Endpoint != "" {
		endpoint = s.config.Endpoint
	}

	router.POST("/"+endpoint, func(c *gin.Context) {
		var json Request

		if c.BindJSON(&json) == nil {
			text := strings.TrimSpace(json.Text)
			if utf8.RuneCountInString(text) > 0 {
				dataType, data, err := s.engine.Process(text)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error":   true,
						"message": err.Error(),
					})
					return
				}

				c.JSON(http.StatusOK, gin.H{
					"error": false,
					"type":  dataType,
					"data":  data,
				})
				return
			}
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "text field is mandatory and can not be empty",
		})
	})

	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	if !s.config.Secure {
		router.Run(addr)
	} else {
		router.RunTLS(addr, s.config.CertPerm, s.config.KeyPerm)
	}
}
