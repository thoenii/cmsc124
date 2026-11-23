package evaluator

import (
	"bufio"
	"os"
)

type Interpreter struct {
	env   *Environment
	world *World
}

func NewInterpreter() *Interpreter {
	globals := NewEnvironment(nil)
	world := &World{}
	registerNatives(globals, world, bufio.NewReader(os.Stdin))
	return &Interpreter{
		env:   globals,
		world: world,
	}
}