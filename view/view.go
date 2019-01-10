package view

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

var (
	htmlContentType = []string{"text/html; charset=utf-8"}

	// DefaultConfig holds default configuration for view.Ending
	DefaultConfig = Config{
		Root:         "views",
		Extension:    ".html",
		Master:       "layouts/master",
		Partials:     []string{},
		Funcs:        make(template.FuncMap),
		DisableCache: false,
		Delims:       Delims{Left: "{{", Right: "}}"},
	}
)

// Default created view.Engine with default configuration
func Default() *Engine {
	return New(DefaultConfig)
}

// New creates view Engine instance
func New(config Config) *Engine {
	return &Engine{
		config:      config,
		tplMap:      make(map[string]*template.Template),
		tplMutex:    sync.RWMutex{},
		fileHandler: DefaultFileHandler(),
	}
}

// DefaultFileHandler -
func DefaultFileHandler() FileHandler {
	return func(config Config, tplFile string) (content string, err error) {
		// Get the absolute path of the root template
		path, err := filepath.Abs(config.Root + string(os.PathSeparator) + tplFile + config.Extension)
		if err != nil {
			return "", fmt.Errorf("TemplateEngine path:%v error: %v", path, err)
		}
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("TemplateEngine render read name:%v, path:%v, error: %v", tplFile, path, err)
		}
		return string(data), nil
	}
}
