package sparkle

type Context struct {
	data map[string]interface{}
}

func newContext() *Context {
	return &Context{
		data: make(map[string]interface{}),
	}
}

// Sets a value against a Context with the given key
//
// This is primarily there so that extensions can store data
// against a context
func (c *Context) Set(key string, value interface{}) {
	c.data[key] = value
}

// Gets a value stored against a Context by it's key
//
// If the key does not exist, nil is returned
func (c *Context) Get(key string) interface{} {
	result, ok := c.data[key]
	if !ok {
		return nil
	}
	return result
}
