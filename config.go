package flow

// Config -
type Config map[string]interface{}

// Get returns value for a given key
func (c Config) Get(key string) interface{} {
	return c.GetDefault(key, nil)
}

// GetDefault returns value for a given key, if value for a given key does not exists
// then defaultVal is returned
func (c Config) GetDefault(key string, defaultVal interface{}) interface{} {
	if data, ok := c[key]; ok {
		return data
	}
	return defaultVal
}

// String returns value as string for a given key
func (c Config) String(key string) string {
	return c.StringDefault(key, "")
}

// StringDefault returns value as string for a given key,
// if value for a given key does not exists
// then defaultVal is returned
func (c Config) StringDefault(key, defaultVal string) string {
	if val, ok := c.GetDefault(key, defaultVal).(string); ok {
		return val
	}
	return defaultVal
}

// Int returns value as int for a given key,
// if value is not found or if value can not be converted to int
// function will return 0 value
func (c Config) Int(key string) int {
	return c.IntDefault(key, 0)
}

// IntDefault returns value as int for a given key,
// if value is not found or if value can not be converted to int
// function will return defaultValue value
func (c Config) IntDefault(key string, defaultVal int) int {
	if val, ok := c.GetDefault(key, defaultVal).(int); ok {
		return val
	}
	return defaultVal
}

// Int32 returns value as int for a given key,
// if value is not found or if value can not be converted to int32
// function will return 0 value
func (c Config) Int32(key string) int32 {
	return c.Int32Default(key, 0)
}

// Int32Default returns value as int for a given key,
// if value is not found or if value can not be converted to int
// function will return defaultValue value
func (c Config) Int32Default(key string, defaultVal int32) int32 {
	if val, ok := c.GetDefault(key, defaultVal).(int32); ok {
		return val
	}
	return defaultVal
}

// Int64 returns value as int for a given key,
// if value is not found or if value can not be converted to int64
// function will return 0 value
func (c Config) Int64(key string) int64 {
	return c.Int64Default(key, 0)
}

// Int64Default returns value as int for a given key,
// if value is not found or if value can not be converted to int
// function will return defaultValue value
func (c Config) Int64Default(key string, defaultVal int64) int64 {
	if val, ok := c.GetDefault(key, defaultVal).(int64); ok {
		return val
	}
	return defaultVal
}

// Bool returns value as bool for a given key,
// if value is not found or if value can not be converted to bool
// function will return false
func (c Config) Bool(key string) bool {
	return c.BoolDefault(key, false)
}

// BoolDefault returns value as bool for a given key,
// if value is not found or if value can not be converted to bool
// function will return defaultValue value
func (c Config) BoolDefault(key string, defaultVal bool) bool {
	if val, ok := c.GetDefault(key, defaultVal).(bool); ok {
		return val
	}
	return defaultVal
}
