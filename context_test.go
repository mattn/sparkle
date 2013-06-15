package sparkle

import "testing"

func Test_ContextDataAttachment(t *testing.T) {
	context := newContext()

	context.Set("Test", 10)
	i := context.Get("Test")

	if i == nil {
		t.Error("Failed to retrieve inserted data from context")
		return
	}

	r, ok := i.(int)
	if !ok {
		t.Error("Returned value was not the same type as passed value")
	}

	if r != 10 {
		t.Error("Return value was not the same")
	}
}
