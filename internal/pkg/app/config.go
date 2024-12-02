package app

type config struct {
	HTTP *configHTTP
	DB   *configDB
}

func newConfig() *config {
	c := new(config)
	c.HTTP = newHTTPConfig()
	c.DB = newDBConfig()
	return c
}
