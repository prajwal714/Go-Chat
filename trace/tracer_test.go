package trace

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf) // here we are testing a New function
	if tracer == nil {
		t.Error("return from new should not be nil")
	} else {
		tracer.Trace("Hello trace package \n")
		if buf.String() != "Hello trace package \n" {
			t.Errorf("Trace should not write '%s' ", buf.String())
		}
	}
}

func TestOff(t *testing.T) {
	var silentTracer Tracer = Off()
	silentTracer.Trace("something")
}
