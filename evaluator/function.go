package evaluator

import (
	"cmsc124/parser"
	"errors"
)

// Callable is anything that can be called: user functions now, natives later.
type Callable interface {
	Arity() int
	Call(in *Interpreter, args []Value) (Value, error)
}

// ReturnSignal carries a return value up through Exec's error returns.
// It is not a real error; Call intercepts it.
type ReturnSignal struct{ Value Value }

func (ReturnSignal) Error() string { return "return outside function" }

// UserFunction pairs a declaration with the environment it was declared in.
type UserFunction struct {
	Decl    parser.FunStmt
	Closure *Environment
}

func (f *UserFunction) Arity() int { return len(f.Decl.Params) }

func (f *UserFunction) Call(in *Interpreter, args []Value) (Value, error) {
	// Fresh environment per call, parented to the DEFINITION environment
	// (not the call site). This is what makes Viper lexically scoped.
	env := NewEnvironment(f.Closure)
	for idx, p := range f.Decl.Params {
		env.Define(p.Lexeme, args[idx]) // ASSUMED: same Define signature as VarStmt
	}

	err := in.executeBlock(f.Decl.Body, env)

	var ret ReturnSignal
	if errors.As(err, &ret) {
		return ret.Value, nil
	}
	if err != nil {
		return nil, err
	}
	return nil, nil // fell off the end: returns dead (nil)
}