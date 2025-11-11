package config

type mySQL struct {
	Addr     string
	Database string
	Username string
	Password string
	Charset  string
}

type server struct {
	Addr string
	Port int
}

type aiEndpoint struct {
    Url    string
    ApiKey string
    Model  string
}

type config struct {
    MySQL      mySQL
    Server     server
    AiEndpoint aiEndpoint
}
