package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	g := &Game{}

	// RunGame starts the main loop and runs the game.
	// game's Update function is called every tick to update the game logic.
	// game's Draw function is called every frame to draw the screen.
	// game's Layout function is called when necessary, and can specify the logical screen size by the function.
	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}