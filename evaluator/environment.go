package evaluator

import (
	"cmsc124/scanner"
)

type Environment struct {
	values    map[string]any
	enclosing *Environment
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		values:    make(map[string]any),
		enclosing: enclosing,
	}
}

type uninitialized struct{}

// Define always creates (or overwrites) the name in THIS environment only.
func (e *Environment) Define(name string, value any) {
	e.values[name] = value
}

// Get walks outward and returns the first match.
func (e *Environment) Get(name scanner.Token) (any, error) {
	value, found := e.values[name.Lexeme]
	if !found {
		if e.enclosing != nil {
			return e.enclosing.Get(name)
		}
		return nil, undefinedVariable(name)
	}

	switch value.(type) {
	case uninitialized:
		return nil, RuntimeError{name, "Variable '" + name.Lexeme + "' is not initialized."}
	default:
		return value, nil
	}
}

// Assign walks outward and changes the first match. It never creates anything.
func (e *Environment) Assign(name scanner.Token, value any) error {
	_, found := e.values[name.Lexeme]
	if found {
		e.values[name.Lexeme] = value
		return nil
	}

	if e.enclosing == nil {
		return undefinedVariable(name)
	}
	return e.enclosing.Assign(name, value)
}

func undefinedVariable(name scanner.Token) error {
	return RuntimeError{name, "Undefined variable '" + name.Lexeme + "'."}
}
