// natives_game.go connects Viper to package game. World holds a swappable
// pointer to the current game; the natives are thin wrappers that validate
// arguments and call one game method. They read World.G at call time so
// setup() can replace the game.
package evaluator

import (
	"errors"
	"fmt"
	"math"
	"time"	// for the random number geneator seed

	"cmsc124/game"
)

// World holds the current game. G is nil until setup() is called.
// NewInterpreter creates one World and hands it to all the natives.
type World struct {
	G *game.Game
}

// errNoGame is the shared message for natives called before setup()
var errNoGame = errors.New("no game yet; call setup(w, h) first.")

// newGameNatives returns setup, render, turn, move, score and over, bound to w.
func newGameNatives(w *World) []*NativeFn {
	return []*NativeFn{

		// setup(w, h): create a new w x h board, replacing any existing game
		NewNative("setup", 2, func(args []Value) (Value, error) {
			// Viper numbers are float64, so check if they are whole and >= 1.
			width, err := positiveInt("setup", "w", args[0])
			if err != nil {
				return nil, err
			}
			height, err := positiveInt("setup", "h", args[1])
			if err != nil {
				return nil, err
			}
			w.G = game.NewGame(width, height, time.Now().UnixNano())	//for seeding the apple
			return nil, nil
		}),

		// render(): prints the board to stdout, using H for the head, o for the body, * for the apple, and . for empty cells.
		NewNative("render", 0, func(_ []Value) (Value, error) {
			if w.G == nil {
				return nil, fmt.Errorf("render: %w", errNoGame)
			}
			fmt.Print(w.G.Render())
			return nil, nil
		}),

		// turn(dir): set the direction to "up", "down", "left" or "right".
		NewNative("turn", 1, func(args []Value) (Value, error) {
			dir, ok := args[0].(string)
			if !ok {
				return nil, errors.New("turn: dir must be a string.")
			}
			if w.G == nil {
				return nil, fmt.Errorf("turn: %w", errNoGame)
			}
			// Turn rejects unknown directions and ignores 180-degree reversals.
			if err := w.G.Turn(dir); err != nil {
				return nil, fmt.Errorf("turn: %w", err)
			}
			return nil, nil
		}),

		// move(): advance the snake one cell. Handles eating, walls and self-collision inside game
		NewNative("move", 0, func(_ []Value) (Value, error) {
			if w.G == nil {
				return nil, fmt.Errorf("move: %w", errNoGame)
			}
			w.G.Move()
			return nil, nil
		}),

		// score(): number of apples eaten
		NewNative("score", 0, func(_ []Value) (Value, error) {
			if w.G == nil {
				return nil, fmt.Errorf("score: %w", errNoGame)
			}
			return float64(w.G.Score), nil
		}),

		// over(): true once the snake has hit a wall or itself, or the board is full.
		NewNative("over", 0, func(_ []Value) (Value, error) {
			if w.G == nil {
				return nil, fmt.Errorf("over: %w", errNoGame)
			}
			return w.G.Over, nil
		}),
	}
}

// positiveInt checks that v is a whole number >= 1 and returns it as an int.
func positiveInt(fn, param string, v Value) (int, error) {
	n, ok := v.(float64)
	if !ok {
		return 0, fmt.Errorf("%s: %s must be a number.", fn, param)
	}
	// n != math.Trunc(n) catches fractions
	// MaxInt32 keeps the conversion to int safe and stops very big nums for board creation
	if n != math.Trunc(n) || n < 1 || n > math.MaxInt32 {
		return 0, fmt.Errorf("%s: %s must be a whole number of at least 1.", fn, param)
	}
	return int(n), nil
}