package log

import (
	"fmt"
	"time"

	"go.uber.org/zap"
)

// Skip constructs a no-op field, which is often useful when handling invalid
// inputs in other Field constructors.
func Skip() zap.Field {
	return zap.Skip()
}

// Binary constructs a field that carries an opaque binary blob.
//
// Binary data is serialized in an encoding-appropriate format. For example,
// zap's JSON encoder base64-encodes binary blobs. To log UTF-8 encoded text,
// use ByteString.
func Binary(key string, val []byte) zap.Field {
	return zap.Binary(key, val)
}

// Bool constructs a field that carries a bool.
func Bool(key string, val bool) zap.Field {
	return zap.Bool(key, val)
}

// ByteString constructs a field that carries UTF-8 encoded text as a []byte.
// To log opaque binary blobs (which aren't necessarily valid UTF-8), use
// Binary.
func ByteString(key string, val []byte) zap.Field {
	return zap.ByteString(key, val)
}

// Complex128 constructs a field that carries a complex number. Unlike most
// numeric fields, this costs an allocation (to convert the complex128 to
// interface{}).
func Complex128(key string, val complex128) zap.Field {
	return Complex128(key, val)
}

// Complex64 constructs a field that carries a complex number. Unlike most
// numeric fields, this costs an allocation (to convert the complex64 to
// interface{}).
func Complex64(key string, val complex64) zap.Field {
	return Complex64(key, val)
}

// Float64 constructs a field that carries a float64. The way the
// floating-point value is represented is encoder-dependent, so marshaling is
// necessarily lazy.
func Float64(key string, val float64) zap.Field {
	return zap.Float64(key, val)
}

// Float32 constructs a field that carries a float32. The way the
// floating-point value is represented is encoder-dependent, so marshaling is
// necessarily lazy.
func Float32(key string, val float32) zap.Field {
	return zap.Float32(key, val)
}

// Int constructs a field with the given key and value.
func Int(key string, val int) zap.Field {
	return zap.Int(key, val)
}

// Int64 constructs a field with the given key and value.
func Int64(key string, val int64) zap.Field {
	return zap.Int64(key, val)
}

// Int32 constructs a field with the given key and value.
func Int32(key string, val int32) zap.Field {
	return zap.Int32(key, val)
}

// Int16 constructs a field with the given key and value.
func Int16(key string, val int16) zap.Field {
	return zap.Int16(key, val)
}

// Int8 constructs a field with the given key and value.
func Int8(key string, val int8) zap.Field {
	return zap.Int8(key, val)
}

// String constructs a field with the given key and value.
func String(key string, val string) zap.Field {
	return zap.String(key, val)
}

// Uint constructs a field with the given key and value.
func Uint(key string, val uint) zap.Field {
	return zap.Uint(key, val)
}

// Uint64 constructs a field with the given key and value.
func Uint64(key string, val uint64) zap.Field {
	return zap.Uint64(key, val)
}

// Uint32 constructs a field with the given key and value.
func Uint32(key string, val uint32) zap.Field {
	return zap.Uint32(key, val)
}

// Uint16 constructs a field with the given key and value.
func Uint16(key string, val uint16) zap.Field {
	return zap.Uint16(key, val)
}

// Uint8 constructs a field with the given key and value.
func Uint8(key string, val uint8) zap.Field {
	return zap.Uint8(key, val)
}

// Uintptr constructs a field with the given key and value.
func Uintptr(key string, val uintptr) zap.Field {
	return zap.Uintptr(key, val)
}

// Reflect constructs a field with the given key and an arbitrary object. It uses
// an encoding-appropriate, reflection-based function to lazily serialize nearly
// any object into the logging context, but it's relatively slow and
// allocation-heavy. Outside tests, Any is always a better choice.
//
// If encoding fails (e.g., trying to serialize a map[int]string to JSON), Reflect
// includes the error message in the final log output.
func Reflect(key string, val interface{}) zap.Field {
	return zap.Reflect(key, val)
}

// Namespace creates a named, isolated scope within the logger's context. All
// subsequent fields will be added to the new namespace.
//
// This helps prevent key collisions when injecting loggers into sub-components
// or third-party libraries.
func Namespace(key string) zap.Field {
	return zap.Namespace(key)
}

// Stringer constructs a field with the given key and the output of the value's
// String method. The Stringer's String method is called lazily.
func Stringer(key string, val fmt.Stringer) zap.Field {
	return zap.Stringer(key, val)
}

// Time constructs a Field with the given key and value. The encoder
// controls how the time is serialized.
func Time(key string, val time.Time) zap.Field {
	return zap.Time(key, val)
}

// Stack constructs a field that stores a stacktrace of the current goroutine
// under provided key. Keep in mind that taking a stacktrace is eager and
// expensive (relatively speaking); this function both makes an allocation and
// takes about two microseconds.
func Stack(key string) zap.Field {
	return zap.Stack(key)
}

// Duration constructs a field with the given key and value. The encoder
// controls how the duration is serialized.
func Duration(key string, val time.Duration) zap.Field {
	return zap.Duration(key, val)
}
