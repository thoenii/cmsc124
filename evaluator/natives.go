// natives.go registers all built-ins in the global environment before user
// code runs, and defines the non-game natives: str and input.
package evaluator

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

// TODO: registerNatives defines every built-in in the global environment.

// TODO: replace with the real implementation.
func nativeStr(args []Value) (Value, error) {
	return nil, nil
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
