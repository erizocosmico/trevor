package trevor

// Config is the configuration passed to start the Server.
type Config struct {
	// Plugins is a list of plugins for the trevor server
	Plugins []Plugin

	// Services is a list of services for the trevor engine
	Services []Service

	// Port is the port in which the server will be run
	Port int

	// Host is the host in which the server will be run
	Host string

	// Secure determines if the server will be run over HTTP or HTTPS
	Secure bool

	// Endpoint is the endpoint to get the processed data. e.g: http://localhost:8080/get_data
	Endpoint string

	// KeyPerm is the key for the SSL
	KeyPerm string

	// CertPerm is the cert for the SSL
	CertPerm string
}
