package view

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"path/filepath"
	"strings"
	"sync"

	"github.com/go-flow/flow/render"
)

// FileHandler function
type FileHandler func(config Config, tplFile string) (content string, err error)

// Engine is view engine object
type Engine struct {
	config      Config
	tplMap      map[string]*template.Template
	tplMutex    sync.RWMutex
	fileHandler FileHandler
}

// Renderer returns go-flow render.Renderer interface instance
func (e *Engine) Renderer(name string, data map[string]interface{}, helpers template.FuncMap) render.Renderer {
	return Render{
		Engine:  e,
		Name:    name,
		Data:    data,
		Helpers: helpers,
	}
}

func (e *Engine) executeRender(out io.Writer, name string, data map[string]interface{}, helpers template.FuncMap) error {
	useMaster := true
	if filepath.Ext(name) == e.config.Extension {
		useMaster = false
		name = strings.TrimSuffix(name, e.config.Extension)

	}
	return e.executeTemplate(out, name, data, helpers, useMaster)
}

func (e *Engine) executeTemplate(out io.Writer, name string, data map[string]interface{}, helpers template.FuncMap, useMaster bool) error {
	var tpl *template.Template
	var err error
	var ok bool

	templateMaster := e.config.Master
	if val, ok := data["viewMasterTemplate"]; ok {
		templateMaster = val.(string)
	}

	allFuncs := make(template.FuncMap, 0)
	// viewEngine related functions
	allFuncs["include"] = func(layout string) (template.HTML, error) {
		buf := new(bytes.Buffer)
		err := e.executeTemplate(buf, layout, data, helpers, false)
		return template.HTML(buf.String()), err
	}

	// Get application wide functions
	for k, v := range e.config.Funcs {
		allFuncs[k] = v
	}

	// Get request Scoped functions
	for k, v := range helpers {
		allFuncs[k] = v
	}

	e.tplMutex.RLock()
	tpl, ok = e.tplMap[name]
	e.tplMutex.RUnlock()

	exeName := name
	if useMaster && templateMaster != "" {
		exeName = templateMaster
	}

	if !ok || e.config.DisableCache {
		tplList := make([]string, 0)
		if useMaster && templateMaster != "" {
			tplList = append(tplList, templateMaster)
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

	// Display the content to the screen
	err = tpl.Funcs(allFuncs).ExecuteTemplate(out, exeName, data)
	if err != nil {
		return fmt.Errorf("TemplateEngine execute template error: %v", err)
	}

	return nil
}

// SetFileHandler sets fileHandler
func (e *Engine) SetFileHandler(handle FileHandler) {
	if handle == nil {
		panic("FileHandler can't set nil!")
	}
	e.fileHandler = handle
}

// SetTemplateFuncs sets template functs for engine
//
// templateFuncs are used as view helpers
func (e *Engine) SetTemplateFuncs(funcs template.FuncMap) {
	for name, fn := range funcs {
		e.config.Funcs[name] = fn
	}
}
