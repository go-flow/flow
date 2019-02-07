package view

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// FileHandler function
type FileHandler func(config Config, tplFile string) (content string, err error)

// HTMLEngine is html view engine implementation
type HTMLEngine struct {
	config      Config
	tplMap      map[string]*template.Template
	tplMutex    sync.RWMutex
	fileHandler FileHandler
}

// NewHTMLEngine creates HTML view Engine instance
func NewHTMLEngine(config Config) *HTMLEngine {
	return &HTMLEngine{
		config:      config,
		tplMap:      make(map[string]*template.Template),
		tplMutex:    sync.RWMutex{},
		fileHandler: DefaultFileHandler(),
	}
}

// Render renders HTML content to output writer
//
// if name contains view extension then master layout view will not be parsed
func (e *HTMLEngine) Render(out io.Writer, name string, data map[string]interface{}, viewFuncs template.FuncMap) error {
	useMaster := true
	if filepath.Ext(name) == e.config.Ext {
		useMaster = false
		name = strings.TrimSuffix(name, e.config.Ext)
	}
	return e.RenderTemplate(out, name, data, viewFuncs, useMaster)
}

// SetViewHelpers sets view helper functions to Engine
func (e *HTMLEngine) SetViewHelpers(viewFuncs template.FuncMap) {
	for name, fn := range viewFuncs {
		e.config.Funcs[name] = fn
	}
}

// RenderTemplate renders HTML Template to output writer
func (e *HTMLEngine) RenderTemplate(out io.Writer, name string, data map[string]interface{}, viewFuncs template.FuncMap, useMaster bool) error {
	var tpl *template.Template
	var err error
	var ok bool

	masterTpl := e.config.Master
	// override master template passing `viewMasterTemplate` to data
	if val, ok := data["viewMasterTemplate"]; ok {
		masterTpl = val.(string)
	}
	// create view Engine helper functions
	allFuncs := make(template.FuncMap, 0)

	// view helper function to include layout file
	allFuncs["include"] = func(layout string) (template.HTML, error) {
		buf := new(bytes.Buffer)
		err = e.RenderTemplate(buf, layout, data, viewFuncs, false)
		return template.HTML(buf.String()), err
	}

	// store viewFuncs from engine configuration
	for k, v := range e.config.Funcs {
		allFuncs[k] = v
	}

	// store render scoped functions
	for k, v := range viewFuncs {
		allFuncs[k] = v
	}

	// get template
	e.tplMutex.RLock()
	tpl, ok = e.tplMap[name]
	e.tplMutex.RUnlock()

	exeTpl := name
	if useMaster && masterTpl != "" {
		exeTpl = masterTpl
	}

	if !ok || e.config.DisableCache {
		tplList := make([]string, 0)
		if useMaster && masterTpl != "" {
			tplList = append(tplList, masterTpl)
		}
		tplList = append(tplList, name)
		tplList = append(tplList, e.config.Partials...)

		// Loop through each template and test the full path
		tpl = template.New(name).Funcs(allFuncs).Delims(e.config.Delims.Left, e.config.Delims.Right)
		for _, v := range tplList {
			var data string
			data, err = e.fileHandler(e.config, v)
			if err != nil {
				return err
			}
			var tmpl *template.Template
			if v == name {
				tmpl = tpl
			} else {
				tmpl = tpl.New(v)
			}
			_, err = tmpl.Parse(data)
			if err != nil {
				return fmt.Errorf("TemplateEngine render parser name:%v, error: %v", v, err)
			}
		}
		e.tplMutex.Lock()
		e.tplMap[name] = tpl
		e.tplMutex.Unlock()
	}

	// execute template
	return tpl.Funcs(allFuncs).ExecuteTemplate(out, exeTpl, data)
}

// DefaultFileHandler  is helper function used for loading html templates
func DefaultFileHandler() FileHandler {
	return func(config Config, tplFile string) (content string, err error) {
		// Get the absolute path of the root template
		path, err := filepath.Abs(config.Root + string(os.PathSeparator) + tplFile + config.Ext)
		if err != nil {
			return "", fmt.Errorf("ViewEngine path:%v, error: %v", path, err)
		}
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("ViewEngine render read name:%v, path:%v, error: %v", tplFile, path, err)
		}
		return string(data), nil
	}
}
