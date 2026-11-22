// callable.go defines the Callable interface shared by user-defined and
// native functions, so call evaluation and arity checks are handled uniformly.
package evaluator

// Callable is anything that can be called from Viper code.
type Callable interface {
	Arity() int

	// Call runs the function. evaluateCall has already checked
	// len(args) == Arity(), so implementations may index args freely.
	Call(in *Interpreter, args []Value) (Value, error)
}
