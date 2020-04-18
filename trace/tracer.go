package trace

import (
	"fmt"
	"io"
)

//Tracer is the interface that describes an object capable of
//tracing events throughout code

type Tracer interface {
	Trace(...interface{}) //...interface denotes Trace will accept 0 or more arguments of any type
}

type tracer struct {
	out io.Writer
}

type nilTracer struct {
}

func (t *nilTracer) Trace(a ...interface{}) {}

func Off() Tracer {
	return &nilTracer{}
}

//function Trace is a method of class tracer which writes a to out instance of tracer class
func (t *tracer) Trace(a ...interface{}) {
	fmt.Fprint(t.out, a...)
	fmt.Fprintln(t.out)
}

//Tracer here is a type of data which can have any type of value
//Thus we are using this as a return type in New func
func New(w io.Writer) Tracer {
	return &tracer{out: w}
}
