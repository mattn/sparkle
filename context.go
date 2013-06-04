package sparkle

type Context struct {
	data map[string]interface{}
}

func newContext() *Context {
	return &Context{
		data: make(map[string]interface{}),
	}
}

func (c *Context) Set(key string, value interface{}) {
	c.data[key] = value
}

func (c *Context) Get(key string) interface{} {
	result, ok := c.data[key]

	if !ok {
		return nil
	}


	return result
}