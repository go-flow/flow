package binding

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestYAMLBindingBindBody(t *testing.T) {
	var s struct {
		Foo string `yaml:"foo"`
	}
	err := yamlBinding{}.BindBody([]byte("foo: FOO"), &s)
	require.NoError(t, err)
	assert.Equal(t, "FOO", s.Foo)
}
