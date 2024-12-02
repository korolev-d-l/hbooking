package app

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type configHTTP struct {
	Server struct {
		Port              string
		HandlerTimeout    time.Duration
		ReadHeaderTimeout time.Duration
		ReadTimeout       time.Duration
		IdleTimeout       time.Duration
		ShutdownTimeout   time.Duration
	}
	GinMode string // "debug" | "release" | "test"
}

func newHTTPConfig() *configHTTP {
	c := &configHTTP{}

	// Set default
	c.Server.Port = "8080"
	c.Server.HandlerTimeout = 6 * time.Second
	c.Server.ReadHeaderTimeout = 3 * time.Second
	c.Server.ReadTimeout = 3 * time.Second
	c.Server.IdleTimeout = 3 * time.Second
	c.Server.ShutdownTimeout = 10 * time.Second

	c.Load()
	c.Validate()
	return c
}

func (c *configHTTP) Load() {
	if httpPort := os.Getenv("HTTP_PORT"); httpPort != "" {
		c.Server.Port = httpPort
	}
	if httpTimeout := os.Getenv("HTTP_TIMEOUT"); httpTimeout != "" {
		handlerTimeoutDuration, err := time.ParseDuration(httpTimeout)
		if err != nil {
			panic(fmt.Errorf("incorrect env HTTP_TIMEOUT: %w", err))
		}
		c.Server.HandlerTimeout = handlerTimeoutDuration
	}
	if httpReadHeaderTimeout := os.Getenv("HTTP_READ_HEADER_TIMEOUT"); httpReadHeaderTimeout != "" {
		readHeaderTimeoutDuration, err := time.ParseDuration(httpReadHeaderTimeout)
		if err != nil {
			panic(fmt.Errorf("incorrect env HTTP_READ_HEADER_TIMEOUT: %w", err))
		}
		c.Server.ReadHeaderTimeout = readHeaderTimeoutDuration
	}
	if httpReadTimeout := os.Getenv("HTTP_READ_TIMEOUT"); httpReadTimeout != "" {
		readTimeoutDuration, err := time.ParseDuration(httpReadTimeout)
		if err != nil {
			panic(fmt.Errorf("incorrect env HTTP_READ_TIMEOUT: %w", err))
		}
		c.Server.ReadTimeout = readTimeoutDuration
	}
	if httpIdleTimeout := os.Getenv("HTTP_IDLE_TIMEOUT"); httpIdleTimeout != "" {
		idleTimeoutDuration, err := time.ParseDuration(httpIdleTimeout)
		if err != nil {
			panic(fmt.Errorf("incorrect env HTTP_IDLE_TIMEOUT: %w", err))
		}
		c.Server.IdleTimeout = idleTimeoutDuration
	}
	if httpShutdownTimeout := os.Getenv("HTTP_SHUTDOWN_TIMEOUT"); httpShutdownTimeout != "" {
		shutdownTimeoutDuration, err := time.ParseDuration(httpShutdownTimeout)
		if err != nil {
			panic(fmt.Errorf("incorrect env HTTP_SHUTDOWN_TIMEOUT: %w", err))
		}
		c.Server.ShutdownTimeout = shutdownTimeoutDuration
	}
	c.GinMode = os.Getenv(gin.EnvGinMode)
}

func (c *configHTTP) Validate() {

}
