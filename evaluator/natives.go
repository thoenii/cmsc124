// natives.go registers all built-ins in the global environment before user
// code runs, and defines the non-game natives: str and input.
package evaluator

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

// registerNatives defines every built-in in the global environment.
func registerNatives(env *Environment, w *World, stdin *bufio.Reader) {
	natives := []*NativeFn{
		NewNative("str", 1, nativeStr),
		newInputNative(stdin),
	}
	natives = append(natives, newGameNatives(w)...)
	for _, n := range natives {
		env.Define(n.name, n)
	}
}

// nativeStr converts a number to its string form, using the same formatting
// as hiss (5.0 -> "5", 2.5 -> "2.5").
func nativeStr(args []Value) (Value, error) {
	n, ok := args[0].(float64)
	if !ok {
		return nil, errors.New("str: argument must be a number.")
	}
	return Stringify(n), nil
}

// nativeInput returns a native that reads one line from stdin, without the
// line ending, or nil at end of input so piped programs cannot loop forever.
func newInputNative(stdin *bufio.Reader) *NativeFn {
	return NewNative("input", 0, func(_ []Value) (Value, error) {
		line, err := stdin.ReadString('\n')
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return nil, err // real read failure
			}
			if line == "" {
				return nil, nil // clean EOF
			}
		}
		return strings.TrimRight(line, "\r\n"), nil
	})
}