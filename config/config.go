package config

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"text/template"

	"gopkg.in/yaml.v2"
)

var errWrongConfigurationType = errors.New("Configuration type must be a pointer to a struct")

// LoadFromPath reads configuration from path and stores it to obj interface
// The format is deduced from the file extension
//	* .yml     - is decoded as yaml
func LoadFromPath(path string, obj interface{}) error {
	_, err := os.Stat(path)
	if err != nil {
		return err
	}

	data, err := os.Open(path)
	if err != nil {
		return err
	}

	return LoadFromReader(data, obj)
}

// LoadFromReader reads configuration from reader and stores it to obj interface
// The format is deduced from the file extension
//	* .yml     - is decoded as yaml
func LoadFromReader(reader io.Reader, obj interface{}) error {
	err := checkConfigObj(obj)
	if err != nil {
		return err
	}

	tmpl := template.New("app_config")
	tmpl.Funcs(map[string]interface{}{
		"envOr": func(envKey, defaultVal string) string {
			return GetEnv(envKey, defaultVal)
		},
		"env": func(envKey string) string {
			return GetEnv(envKey, "")
		},
	})

	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	t, err := tmpl.Parse(string(b))
	if err != nil {
		return err
	}

	var bb bytes.Buffer
	err = t.Execute(&bb, nil)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(bb.Bytes(), obj)
	if err != nil {
		return err
	}
	return nil
}

// GetEnv returns environment variable value for a given key
// if value is not found defaultValue param will be returned
func GetEnv(key, defaultValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultValue
}

func checkConfigObj(obj interface{}) error {
	objVal := reflect.ValueOf(obj)
	// we allow maps
	if objVal.Kind() == reflect.Map {
		return nil
	}

	// check if type is a pointer
	if objVal.Kind() != reflect.Ptr || objVal.IsNil() {
		return errWrongConfigurationType
	}

	// get and confirm struct value
	objVal = objVal.Elem()
	if objVal.Kind() != reflect.Struct {
		return errWrongConfigurationType
	}
	return nil
}
