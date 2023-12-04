package main

import (
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"

	"embed"
)

type Game struct {}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(PlayerImage, nil)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

//關鍵字: go:embed，透過註解就可以直接讀取檔案
//go:embed space-shooter-extension/PNG/*
var assets embed.FS
var PlayerImage = mustLoadImage("space-shooter-extension/PNG/Sprites_X2/Ships/spaceShips_005.png")

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

func mustLoadImage(name string) *ebiten.Image {
	f, err := assets.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	
	//image/png used for decoding
	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	return ebiten.NewImageFromImage(img)
}