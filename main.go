package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

// Game implements ebiten.Game interface
type Game struct {
	// Do stuff, e.g. game state vars: current word, guesses, etc.
}

func (g *Game) Update(*ebiten.Image) error {
	// Handle inputs, game logic, etc
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Che Guordel!")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("CheGuordle!")
	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
	}
}
