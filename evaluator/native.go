// native.go defines NativeFn, a generic wrapper for built-in functions
// Natives return plain errors without position information;
// evaluateCall attaches the line number, keeping error formatting in one place.
package evaluator

import "fmt"

// NativeFn is a built-in function. Always used by pointer: a struct with a
// func field is not comparable and would panic under ==.
type NativeFn struct {
	name  string
	arity int
	fn    func(args []Value) (Value, error)
}

func NewNative(name string, arity int, fn func(args []Value) (Value, error)) *NativeFn {
	return &NativeFn{name: name, arity: arity, fn: fn}
}

func (n *NativeFn) Arity() int { return n.arity }

// Call ignores the interpreter; natives that need shared state capture it
// in their closure.
func (n *NativeFn) Call(_ *Interpreter, args []Value) (Value, error) {
	return n.fn(args)
}

func (n *NativeFn) String() string {
	return fmt.Sprintf("<native fn %s>", n.name)
}
