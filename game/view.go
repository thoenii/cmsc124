package game

import (
	"fmt"
	"strings"
)

// Render returns the board as text, one row per line, ending in a newline.
// H is the head, o is the body, * is the apple, . is an empty cell.
func (g *Game) Render() string {
	var sb strings.Builder

	// Go through the board row by row (y), and within each row cell by cell (x)
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			p := Point{x, y}

			switch {
			case p == g.Snake[0]:	// head
				sb.WriteByte('H')
			case g.onSnake(p):		// body
				sb.WriteByte('o')
			case p == g.Apple && !g.Over:
				sb.WriteByte('*')	// apple ; hide if the game is over
			default:
				sb.WriteByte('.')	// empty cell
			}
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

// Turn sets the direction to "up", "down", "left" or "right". 
// A turn that would reverse a snake longer than one cell is ignored, since it would
// run the head into the neck.
func (g *Game) Turn(dir string) error {
	var d Point

	// The screen's y axis points down, so "up" is y = -1 and "down" is y = +1
	switch dir {
	case "up":
		d = Point{0, -1}
	case "down":
		d = Point{0, 1}
	case "left":
		d = Point{-1, 0}
	case "right":
		d = Point{1, 0}
	default:
		return fmt.Errorf("unknown direction %q (use up, down, left or right).", dir)
	}
	// A one-cell snake has no neck, so it may turn around freely.
	// direction reverses when the new direction is the negative of the current direction
	if len(g.Snake) > 1 && d.X == -g.Dir.X && d.Y == -g.Dir.Y {
		return nil
	}
	g.Dir = d // for next move get direction
	return nil
}