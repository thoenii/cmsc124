// Package game is the model for Viper's Snake game: board, snake, apple, and
// the movement and collision rules.
package game

import "math/rand"

type Point struct {
	X, Y int
}

type Game struct {
	Width, Height int
	Snake         []Point // Snake[0] is the head
	Dir           Point
	Apple         Point
	Score         int
	Over          bool
	rng           *rand.Rand
}

func NewGame(w, h int, seed int64) *Game {
	source := rand.NewSource(seed)

	g := &Game{
		Width: w, Height: h,
		Snake: []Point{{w / 2, h / 2}},
		Dir:   Point{1, 0},
		rng:   rand.New(source),
	}
	g.placeApple()
	return g
}

// placeApple picks from the list of free cells instead of retrying randomly,
// so a nearly full grid cannot loop. A full grid ends the game (a win).
func (g *Game) placeApple() {
	var free []Point
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			p := Point{x, y}
			if !g.onSnake(p) {
				free = append(free, p)
			}
		}
	}
	if len(free) == 0 {
		g.Over = true
		return
	}
	count := len(free)        // how many free cells there are
	pick := g.rng.Intn(count) // a random number from 0 to count-1
	g.Apple = free[pick]      // the free cell at that index
}

func (g *Game) onSnake(p Point) bool {
	for _, s := range g.Snake {
		if s == p {
			return true
		}
	}
	return false
}

// Move advances the snake one cell in its current direction.
func (g *Game) Move() {
	if g.Over {
		return
	}

	head := g.Snake[0]
	next := Point{head.X + g.Dir.X, head.Y + g.Dir.Y}

	if next.X < 0 || next.X >= g.Width || next.Y < 0 || next.Y >= g.Height {
		g.Over = true
		return
	}

	eating := next == g.Apple

	// When not eating, the tail cell is vacated this turn, so moving into
	// it is legal. When eating, the tail stays, so the whole body counts.
	body := g.Snake
	if !eating {
		tail := len(g.Snake) - 1
		body = g.Snake[:tail] // everything except the tail
	}
	for _, s := range body {
		if s == next {
			g.Over = true
			return
		}
	}

	g.Snake = append([]Point{next}, g.Snake...)
	if eating {
		g.Score++
		g.placeApple() // after the new head is in, so the apple can't spawn under it
	} else {
		last := len(g.Snake) - 1
		g.Snake = g.Snake[:last] // drop the tail
	}
}
