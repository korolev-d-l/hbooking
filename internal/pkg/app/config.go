package app

type config struct {
	HTTP *configHTTP
	DB   *configDB
}

func (a *App) initConfig() {
	c := new(config)
	c.HTTP = newHTTPConfig()
	c.DB = newDBConfig()
	a.cfg = c
}
