package web

import (
	"path"
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
