package plugin

import "strconv"

// Context 插件上下文
type Context map[string]string

func (c Context) Add(key, value string) {
	c[key] = value
}

func (c Context) GetString(key string) (string, bool) {
	val, ok := c[key]
	return val, ok
}

func (c Context) GetStringOrDefault(key string, defaultValue string) string {
	val, ok := c[key]
	if !ok {
		val = defaultValue
	}
	return val
}

func (c Context) GetInt(key string) (int, bool) {
	val, ok := c.GetString(key)
	if !ok {
		return 0, false
	}
	if iVal, err := strconv.Atoi(val); err == nil {
		return iVal, true
	}
	return 0, false
}

func (c Context) GetIntOrDefault(key string, defaultValue int) int {
	val, ok := c.GetInt(key)
	if !ok {
		val = defaultValue
	}
	return val
}

func (c Context) GetBoolOrDefault(key string, defaultValue bool) bool {
	val, ok := c[key]
	if !ok {
		return defaultValue
	}

	if res, err := strconv.ParseBool(val); err == nil {
		return res
	}
	return defaultValue
}
