package main

import (
	"image"
	_ "image/png"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"

	// "github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"embed"
)

type Vector struct {
	X float64
	Y float64
}

func (v *Vector) Normalize() Vector {
	// 計算向量長度
	length := math.Sqrt(v.X*v.X + v.Y*v.Y)
	if length == 0 {
		return Vector{0, 0} // Avoid division by zero
	}
	return Vector{
		X: v.X / length,
		Y: v.Y / length,
	}
}

type Color struct {
	R float64 
	G float64
	B float64
	A float64
}

type Timer struct {
	currentTicks int
	targetTicks int
}

func NewTimer(d time.Duration) *Timer {
	return &Timer{
		currentTicks: 0,
		targetTicks: int(d.Milliseconds()) * ebiten.TPS() / 1000,
	}
}

func (t *Timer) Update() {
	if t.currentTicks < t.targetTicks {
		t.currentTicks++
	}
}

func (t *Timer) IsReady() bool {
	return t.currentTicks >= t.targetTicks
}

func (t *Timer) Reset() {
	t.currentTicks = 0
}

type Player struct {
	position Vector
	sprite *ebiten.Image
	rotation float64
	color Color
	// 射子彈的CD時間
	shootCooldown *Timer
}

func NewPlayer() *Player {
	sprite := PlayerImage

	// 初始化player的位置在畫面的中心
	// 假設sprite的寬度是2x, 畫面寬度是4y, 中心的位置會是2y, 要讓sprite的中心在畫面中心的話會是2y (screenWidth/2) - x (x = sprite的寬度/2) 
	bounds := sprite.Bounds()
	halW := float64(bounds.Dx()) / 2
	halH := float64(bounds.Dy()) / 2
	pos := Vector{
		X: ScreenWidth/2 - halW,
		Y: ScreenHeight/2 - halH,
	}
	return &Player{
		position: pos,
		sprite: PlayerImage,
		shootCooldown: NewTimer(1 * time.Second),
	}
}

func (p *Player) Update() error {
	speed := math.Pi / float64(ebiten.TPS())

	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		p.rotation += speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		p.rotation -= speed
	}

	p.shootCooldown.Update()
	if ebiten.IsKeyPressed(ebiten.KeySpace) && p.shootCooldown.IsReady() {
		p.shootCooldown.Reset()
		// shoot the bullet
	}
	return nil
}

func (p *Player) Draw(screen *ebiten.Image) {
	// 取得sprite的中心位置資訊
	bounds := p.sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	op1 := &colorm.DrawImageOptions{}
	cm := colorm.ColorM{}
	// 將sprite先移動到pivot point的位置後在旋轉sprite, 之後再將sprite移回原位
	op1.GeoM.Translate(-halfW, -halfH)
	op1.GeoM.Rotate(p.rotation)
	op1.GeoM.Translate(halfW, halfH)

	// 將sprite平移動到screen中心的位置
	op1.GeoM.Translate(p.position.X, p.position.Y)
	cm.Translate(p.color.R, p.color.G, p.color.B, p.color.A)
	colorm.DrawImage(screen, p.sprite, cm, op1)
}

type Meteor struct {
	position Vector
	sprite *ebiten.Image
	movement Vector
	rotation float64
	rotationSpeed float64
}

func NewMeteor() *Meteor {
	sprite := MeteorSprites[rand.Intn(len(MeteorSprites))]

	// 取得screen的中心位置資訊
	target := Vector{
		X: ScreenWidth / 2,
		Y: ScreenHeight / 2,
	}
	// 產生meteor時距中心的距離
	r := ScreenWidth / 2.0
	// 隨機一個角度 （2π is 360°）
	angle := rand.Float64() * 2 * math.Pi
	pos := Vector{
		X: target.X + r*math.Cos(angle), // r*math.Cos(angle) = 以某個點為中心且半徑為r, 從正x軸開始移動angle角度的點(0度逆時針旋轉)
		Y: target.Y + r*math.Sin(angle),
	}
	
	// 隨機速度
	velocity := 0.25 + rand.Float64() * 1.5

	// 方向 = 目標位置 - 現在位置
	direction := Vector{
		X: target.X - pos.X,
		Y: target.Y - pos.Y,
	}
	// 標準化向量 - 只取得方向不取得長度
	normalizedDirection := direction.Normalize()

	// 將方向乘以速度
	movement := Vector{
		X: normalizedDirection.X * velocity,
		Y: normalizedDirection.Y * velocity,
	}

	return &Meteor{
		position: pos,
		sprite: sprite,
		movement: movement,
		rotation: 0,
		rotationSpeed: -0.02 + rand.Float64() * 0.04,
	}
}

func (m *Meteor) Update() error {
	m.position.X += m.movement.X
	m.position.Y += m.movement.Y
	m.rotation += m.rotationSpeed
	return nil
}

func (m *Meteor) Draw(screen *ebiten.Image) {
	bounds := m.sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-halfW, -halfH)
	op.GeoM.Rotate(m.rotation)
	op.GeoM.Translate(halfW, halfH)

	op.GeoM.Translate(m.position.X, m.position.Y)
	screen.DrawImage(m.sprite, op)
}

type Game struct {
	playerPosition Vector
	ChangeColorTimer *Timer
	Color Color
	player *Player
	meteorSpawnTimer *Timer
	meteors []*Meteor
}

func (g *Game) Update() error {
	g.ChangeColorTimer.Update()
	if g.ChangeColorTimer.IsReady() {
		g.ChangeColorTimer.Reset()

		// change the color
		print("change the color")
		g.Color.B += 0.01
	}

	g.player.Update()

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

	// op1.GeoM.Translate(200, 200)
	// op1.GeoM.Scale(0.6, 0.6)
	// cm.Translate(1.0, 1.0, 1.0, 1.0)
	// colorm.DrawImage(screen, PlayerImage, cm, op1)

	// ActualTPS returns the current TPS (ticks per second),
	// that represents how many Update function is called in a second.
	// ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()))
}


const (
	ScreenWidth  = 3200
	ScreenHeight = 2400
)
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

//關鍵字: go:embed，透過註解就可以直接讀取檔案
//go:embed space-shooter-extension/PNG/*
var assets embed.FS
var PlayerImage = mustLoadImage("space-shooter-extension/PNG/Sprites_X2/Ships/spaceShips_005.png")
var MeteorSprites = mustLoadImages("space-shooter-extension/PNG/Sprites_X2/Meteors")

func main() {
	g := &Game{playerPosition: Vector{X: 100, Y: 100}, ChangeColorTimer: NewTimer(5 * time.Second), player: NewPlayer(), meteorSpawnTimer: NewTimer(5 * time.Second), meteors: []*Meteor{}}

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

func mustLoadImages(name string) ([]*ebiten.Image) {
	var images []*ebiten.Image

	// 讀取整個資料夾
	entries, err := assets.ReadDir(name)
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		// 如果不是directory，並且檔案須為.png則讀取檔案
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".png") {
			file, err := assets.Open(name + "/" + entry.Name())
			if err != nil {
				panic(err)
			}

			img, _, err := image.Decode(file)
			if err != nil {
				panic(err)
			}

			images = append(images, ebiten.NewImageFromImage(img))
		}
	}

	return images
}