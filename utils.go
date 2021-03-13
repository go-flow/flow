package flow

import (
	"path"
	"regexp"
	"strings"
)

func iterate(path, method string, routes Routes, root *node) Routes {
	path += root.path
	routes = append(routes, *root.handle)
	for _, child := range root.children {
		routes = iterate(path, method, routes, child)
	}
	return routes
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

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}-${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}-${2}")
	return strings.ToLower(snake)
}

func contentTypeFromString(ct string) string {
	for i, char := range ct {
		if char == ' ' || char == ';' {
			return ct[:i]
		}
	}
	return ct
}
