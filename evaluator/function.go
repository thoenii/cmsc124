// function.go implements user-declared functions (fn). Each call runs in a
// fresh environment enclosed by the function's closure, giving Viper lexical
// scoping.
package evaluator

import (
	"cmsc124/parser"
	"errors"
)

// ReturnSignal carries a return value up through Exec's error returns.
// It is not a real error; UserFunction.Call intercepts it. Error handlers
// elsewhere must let it pass through untouched.
type ReturnSignal struct{ Value Value }

func (ReturnSignal) Error() string { return "return outside function" }

// UserFunction pairs a declaration with the environment it was declared in.
type UserFunction struct {
	Decl    parser.FunStmt
	Closure *Environment
}

func (f *UserFunction) Arity() int { return len(f.Decl.Params) }

func (f *UserFunction) Call(in *Interpreter, args []Value) (Value, error) {
	// Parent is the definition environment, not the call site. One
	// environment per call keeps recursive calls from sharing parameters.
	env := NewEnvironment(f.Closure)
	for i, p := range f.Decl.Params {
		env.Define(p.Lexeme, args[i])
	}

	err := in.executeBlock(f.Decl.Body, env)

	var ret ReturnSignal
	if errors.As(err, &ret) {
		return ret.Value, nil
	}
	if err != nil {
		return nil, err
	}
	return nil, nil // fell off the end
}

func (f *UserFunction) String() string {
	return "<fn " + f.Decl.Name.Lexeme + ">"
}
