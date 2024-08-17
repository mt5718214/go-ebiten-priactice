package game

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth  = 3200
	ScreenHeight = 2400
)

type Color struct {
	R float64 
	G float64
	B float64
	A float64
}

type Game struct {
	playerPosition Vector
	ChangeColorTimer *Timer
	Color Color
	player *Player
	meteorSpawnTimer *Timer
	meteors []*Meteor
	bullets []*Bullet
}

func NewGame() *Game {
	g := &Game{playerPosition: Vector{X: 100, Y: 100}, ChangeColorTimer: NewTimer(5 * time.Second), meteorSpawnTimer: NewTimer(5 * time.Second)}
	g.player = NewPlayer(g)
	return g
}

func (g *Game) Update() error {
	g.ChangeColorTimer.Update()
	if g.ChangeColorTimer.IsReady() {
		g.ChangeColorTimer.Reset()

		// change the color
		print("change the color")
		g.Color.B += 0.01
	}

	// update player
	g.player.Update()

	// update meteor
	g.meteorSpawnTimer.Update()
	if g.meteorSpawnTimer.IsReady() {
		g.meteorSpawnTimer.Reset()

		// spawn a meteor
		meteor := NewMeteor()
		g.meteors = append(g.meteors, meteor)
	}

	for _, m := range g.meteors {
		m.Update()
	}

	// update bullet
	for _, b := range g.bullets {
		b.Update()
	}
	

	// speed 是每一tick（one update call）會移動的距離
	// 預設下每秒會有60tick, 所以一秒會移動300像素
	// speed := 5.0

	// Move 300 pixels per second
	// speed := float64(300 / ebiten.TPS())

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// op := &ebiten.DrawImageOptions{}
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
	// op.GeoM.Translate(g.playerPosition.X, g.playerPosition.Y)
	// screen.DrawImage(PlayerImage, op)

	// 可以同時Draw多張圖片
	// op1 := &colorm.DrawImageOptions{}
	// cm := colorm.ColorM{}
	// op1.GeoM.Translate(g.playerPosition.X, g.playerPosition.Y)
	// op1.GeoM.Scale(0.5, 0.5)
	// 該函式有四個參數, 前三個分別代表r g b 三原色(數值範圍為0~1 = 0~100%), 最後一個參數是背景透明度
	// cm.Translate(g.Color.R, g.Color.G, g.Color.B, g.Color.A)
	// colorm.DrawImage(screen, PlayerImage, cm, op1)

	// draw player
	g.player.Draw(screen)

	// draw meteor
	for _, m := range g.meteors {
		m.Draw(screen)
	}

	// draw bullet
	for _, b := range g.bullets {
		b.Draw(screen)
	}

	// op1.GeoM.Translate(200, 200)
	// op1.GeoM.Scale(0.6, 0.6)
	// cm.Translate(1.0, 1.0, 1.0, 1.0)
	// colorm.DrawImage(screen, PlayerImage, cm, op1)

	// ActualTPS returns the current TPS (ticks per second),
	// that represents how many Update function is called in a second.
	// ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) AddBullet(b *Bullet) {
	g.bullets = append(g.bullets, b)
}

func (g *Game) Reset() {
	g.player = NewPlayer(g)
	g.meteors = nil
	g.bullets = nil
}