package flow

// NewConfig creates new Configuration object for a given config data
func NewConfig(data map[string]interface{}) Config {
	return Config{
		data,
	}
}

// Config -
type Config struct {
	data map[string]interface{}
}

// Get returns value for a given key
func (c *Config) Get(key string) interface{} {
	return c.GetD(key, nil)
}

// GetD returns value for a given key, if value for a given key does not exists
// then defaultVal is returned
func (c *Config) GetD(key string, defaultVal interface{}) interface{} {
	item, found := c.data[key]
	if !found {
		return defaultVal
	}
	return item
}

// GetString returns value as string for a given key
func (c *Config) GetString(key string) string {
	return c.GetStringD(key, "")
}

// GetStringD returns value as string for a given key,
// if value for a given key does not exists
// then defaultVal is returned
func (c *Config) GetStringD(key string, defaultVal string) string {
	item, found := c.data[key]
	if !found {
		return defaultVal
	}
	return item.(string)
}

// GetInt returns value as int for a given key,
// if value is not found or if value can not be converted to int
// function will return 0 value
func (c *Config) GetInt(key string) int {
	return c.GetIntD(key, 0)
}

// GetIntD returns value as int for a given key,
// if value is not found or if value can not be converted to int
// function will return defaultValue value
func (c *Config) GetIntD(key string, defaultValue int) int {
	item, found := c.data[key]
	if !found {
		return defaultValue
	}
	v, ok := item.(int)
	if !ok {
		return defaultValue
	}
	return v
}

// GetInt32 returns value as int for a given key,
// if value is not found or if value can not be converted to int
// function will return 0 value
func (c *Config) GetInt32(key string) int32 {
	return c.GetInt32D(key, 0)
}

// GetInt32D returns value as int for a given key,
// if value is not found or if value can not be converted to int
// function will return defaultValue value
func (c *Config) GetInt32D(key string, defaultValue int32) int32 {
	item, found := c.data[key]
	if !found {
		return defaultValue
	}
	v, ok := item.(int32)
	if !ok {
		return defaultValue
	}
	return v
}

// GetInt64 returns value as int for a given key,
// if value is not found or if value can not be converted to int
// function will return 0 value
func (c *Config) GetInt64(key string) int64 {
	return c.GetInt64D(key, 0)
}

// GetInt64D returns value as int for a given key,
// if value is not found or if value can not be converted to int
// function will return defaultValue value
func (c *Config) GetInt64D(key string, defaultValue int64) int64 {
	item, found := c.data[key]
	if !found {
		return defaultValue
	}
	v, ok := item.(int64)
	if !ok {
		return defaultValue
	}
	return v
}

// GetBool returns value as bool for a given key,
// if value is not found or if value can not be converted to bool
// function will return false
func (c *Config) GetBool(key string) bool {
	return c.GetBoolD(key, false)
}

// GetBoolD returns value as bool for a given key,
// if value is not found or if value can not be converted to bool
// function will return defaultValue value
func (c *Config) GetBoolD(key string, defaultValue bool) bool {
	item, found := c.data[key]
	if !found {
		return defaultValue
	}
	v, ok := item.(bool)
	if !ok {
		return defaultValue
	}
	return v
}
