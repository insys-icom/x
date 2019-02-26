package handlerx

import (
	"fmt"
	"os"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	// VueConfig defines the config for Vue middleware.
	VueConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper middleware.Skipper
		// StaticDir is the location on disk of the Vue.js app
		StaticDir string `yaml:"staticDir"`
		// PublicPath is the publicPath setting in vue.config.js. Default '/'.
		PublicPath string `yaml:"publicPath"`
	}
)

const (
	defaultPublicPath string = "/"
)

var (
	// DefaultVueConfig is the default Vue middleware config.
	DefaultVueConfig = VueConfig{
		PublicPath: defaultPublicPath,
	}
)

// Vue returns a middleware which serves a Vue.js single page app
func Vue(dir string) echo.HandlerFunc {
	c := DefaultVueConfig
	c.StaticDir = strings.TrimSuffix(dir, "/")
	return VueWithConfig(c)
}

// VueWithConfig return Vue middleware with config.
// See: `Vue()`.
func VueWithConfig(config VueConfig) echo.HandlerFunc {
	if config.StaticDir == "" {
		panic("echox: vue middleware requires a static directory")
	}
	// Remove trailing slash
	config.StaticDir = strings.TrimSuffix(config.StaticDir, "/")

	// Remove trailing slash if we have a non default public path
	if config.PublicPath != defaultPublicPath {
		config.PublicPath = strings.TrimSuffix(config.PublicPath, "/")
	}

	// TODO: populate cache or file list here!

	return func(c echo.Context) error {
		fmt.Println("enter handler func")
		p := c.Request().URL.Path

		// If the Vue.js app has a publicPath setting we have to trim the
		// path before we're checking if the file exists
		if config.PublicPath != defaultPublicPath {
			p = strings.TrimPrefix(p, config.PublicPath)
		}

		p = config.StaticDir + p

		fmt.Println("call fileExits with param: ", p)

		// If the request file exists in Vue.js folder then deliver it
		if fileExits(p) {
			return c.File(p)
		}

		// In every other case we're running the Vue.js app by service the
		// index.html file.
		return c.File(config.StaticDir + "/index.html")
	}
}

// fileExits reports if the named file exists.
func fileExits(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return !info.IsDir()
}
