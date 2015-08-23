package trevor

// Config is the configuration passed to start the Server.
type Config struct {
	// Plugins is a list of plugins for the trevor server
	Plugins []Plugin

	// Services is a list of services for the trevor engine
	Services []Service

	// Middleware is a list of middlewares for the trevor engine
	Middleware []Middleware

	// Port is the port in which the server will be run
	Port int

	// Host is the host in which the server will be run
	Host string

	// Secure determines if the server will be run over HTTP or HTTPS
	Secure bool

	// Endpoint is the endpoint to get the processed data. e.g: http://localhost:8080/get_data
	Endpoint string

	// InputFieldName is the key of the JSON object passed to the endpoint that contains the input data.
	InputFieldName string

	// CORSOrigin is a comma separated list of origins allowed for CORS.
	CORSOrigin string

	// KeyPerm is the key for the SSL
	KeyPerm string

	// CertPerm is the cert for the SSL
	CertPerm string

	// Analyzer is the function used as a analyzer for choosing the adequate plugin for the request
	Analyzer Analyzer
}
