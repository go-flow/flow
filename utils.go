package flow

import (
	"path"
	"reflect"
	"runtime"
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
