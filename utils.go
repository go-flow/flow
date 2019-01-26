package flow

import (
	"fmt"
	"io/ioutil"
	"path"
	"reflect"
	"runtime"
	"strings"
)

func filterFlags(content string) string {
	for i, char := range content {
		if char == ' ' || char == ';' {
			return content[:i]
		}
	}
	return content
}

func iterate(path, method string, routes Routes, root *node) Routes {
	path += root.path
	if len(root.handler) > 0 {
		handlerFunc := root.handler.Last()
		routes = append(routes, Route{
			Method:        method,
			Path:          path,
			HandlersChain: root.handler,
			HandlerName:   nameOfFunction(handlerFunc),
			HandlerFunc:   handlerFunc,
		})
	}
	for _, child := range root.children {
		routes = iterate(path, method, routes, child)
	}
	return routes
}

func nameOfFunction(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func lastChar(str string) uint8 {
	if str == "" {
		panic("The length of the string can't be 0")
	}
	return str[len(str)-1]
}

func joinPaths(absolutePath, relativePath string) string {
	if relativePath == "" {
		return absolutePath
	}

	finalPath := path.Join(absolutePath, relativePath)
	appendSlash := lastChar(relativePath) == '/' && lastChar(finalPath) != '/'
	if appendSlash {
		return finalPath + "/"
	}
	return finalPath
}

func byteCountDecimal(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

func byteCountBinary(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}

func loadPartials(viewsRoot, partialsRoot, ext string) ([]string, error) {
	dirname := path.Join(viewsRoot, partialsRoot)
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}
	partials := []string{}
	for _, f := range files {
		partial := f.Name()
		if strings.HasSuffix(partial, ext) {
			// remove ext from file
			partial = strings.TrimSuffix(partial, ext)
			// join file with folder name
			partial = path.Join(partialsRoot, partial)

			// add to partials
			partials = append(partials, partial)
		}
	}
	return partials, nil
}
