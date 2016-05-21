package boondoggle

import (
	"fmt"
	"reflect"
	"runtime"
)

// Transformer is the type signature of functions that modify Article types
type Transformer func(*Article) error

// String returns the name of the transformer function
func (t Transformer) String() string {
	return runtime.FuncForPC(reflect.ValueOf(t).Pointer()).Name()
}

// Noop is an example no-operation Transformer
func Noop(article *Article) error {
	return nil
}

var _ = Transformer(Noop)

// CauseError is an example Transformer that returns an error
func CauseError(article *Article) error {
	return fmt.Errorf("I am an intentionally caused error")
}

var _ = Transformer(CauseError)
