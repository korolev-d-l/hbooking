package app

import "os"

type configDB struct {
	URL string
}

func newDBConfig() *configDB {
	c := &configDB{}
	c.Load()
	c.Validate()
	return c
}

func (c *configDB) Load() {
	c.URL = os.Getenv("DATABASE_URL")
}

func (c *configDB) Validate() {
	if c.URL == "" {
		panic("env DATABASE_URL must be set")
	}
}
