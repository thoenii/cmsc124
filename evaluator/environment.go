package evaluator

import (
	"cmsc124/scanner"
	"fmt"
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

// Define always creates (or overwrites) the name in THIS environment only.
func (e *Environment) Define(name string, value any) {
	e.values[name] = value
}

// Get walks outward and returns the first match.
func (e *Environment) Get(name scanner.Token) (any, error) {
	if value, ok := e.values[name.Lexeme]; ok {
		return value, nil
	}
	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}
	return nil, undefinedVariable(name)
}

// Assign walks outward and changes the first match. It never creates anything.
func (e *Environment) Assign(name scanner.Token, value any) error {
	if _, ok := e.values[name.Lexeme]; ok {
		e.values[name.Lexeme] = value
		return nil
	}
	if e.enclosing != nil {
		return e.enclosing.Assign(name, value)
	}
	return undefinedVariable(name)
}

func undefinedVariable(name scanner.Token) error {
	return fmt.Errorf("[line %d] Runtime error: Undefined variable '%s'.", name.Line, name.Lexeme)
}
