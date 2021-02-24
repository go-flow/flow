package flow

import (
	"fmt"
	"io/ioutil"
	"path"
	"reflect"
	"regexp"
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
	routes = append(routes, *root.handle)
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
