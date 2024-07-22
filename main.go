package main

import (
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"

	"embed"
)

type Vector struct {
	X float64
	Y float64
}

type Game struct {
	playerPosition Vector
}

func (g *Game) Update() error {
	// speed 是每一tick（one update call）會移動的距離
	// 預設下每秒會有60tick, 所以一秒會移動300像素
	// speed := 5.0

	// Move 300 pixels per second
	// speed := float64(300 / ebiten.TPS())

	speed := 5.0

	// IsKeyPressed return true if the given key is pressed in the current tick
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.playerPosition.Y += speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.playerPosition.Y -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.playerPosition.X += speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.playerPosition.X -= speed
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(200, 100) moves the image 200 pixels right and 100 pixels down.
	// op.GeoM.Rotate(45.0 * math.Pi / 180.0) rotates the image 45° clockwise. use the degrees * math.Pi / 180.0 formula
	// op.GeoM.Translate(100, 100)
	// op.GeoM.Scale(-1, 1)
	// op.GeoM.Scale(5, 5)
	// width := PlayerImage.Bounds().Dx()
	// height := PlayerImage.Bounds().Dy()
	// halw := float64(width / 2)
	// halh := float64(height / 2)
	// op.GeoM.Translate(-halw, -halh)
	// op.GeoM.Rotate(45.0 * math.Pi / 180.0)
	// op.GeoM.Translate(halw, halh)
	op.GeoM.Translate(g.playerPosition.X, g.playerPosition.Y)
	screen.DrawImage(PlayerImage, op)

	// 可以同時Draw多張圖片
	op1 := &colorm.DrawImageOptions{}
	cm := colorm.ColorM{}
	// 該函式有四個參數, 前三個分別代表r g b 三原色(數值範圍為0~1 = 0~100%), 最後一個參數是背景透明度
	cm.Translate(0, 0, 0.6, 0.0)
	colorm.DrawImage(screen, PlayerImage, cm, op1)

	// op1.GeoM.Translate(200, 200)
	// op1.GeoM.Scale(0.6, 0.6)
	// cm.Translate(1.0, 1.0, 1.0, 1.0)
	// colorm.DrawImage(screen, PlayerImage, cm, op1)
	
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

//關鍵字: go:embed，透過註解就可以直接讀取檔案
//go:embed space-shooter-extension/PNG/*
var assets embed.FS
var PlayerImage = mustLoadImage("space-shooter-extension/PNG/Sprites_X2/Ships/spaceShips_005.png")

func main() {
	g := &Game{playerPosition: Vector{X: 100, Y: 100}}

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