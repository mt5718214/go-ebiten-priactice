package game

import (
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"

	assets "airplane/spaceshooterextension"
	// assets "airplane/space-shooter-extension/spaceshooterextension"
)

const (
	bulletSpawnOffset = 50.0
)


type Player struct {
	game *Game

	position Vector
	sprite *ebiten.Image
	rotation float64
	color Color
	// 射子彈的CD時間
	shootCooldown *Timer
}

func NewPlayer(g *Game) *Player {
	sprite := assets.PlayerImage

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
		game: g,
		position: pos,
		sprite: assets.PlayerImage,
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
		bounds := p.sprite.Bounds()
		halfW := float64(bounds.Dx()) / 2
		halfH := float64(bounds.Dy()) / 2

		spawnPos := Vector{
			p.position.X + halfW + math.Sin(p.rotation)*bulletSpawnOffset,
	    p.position.Y + halfH + math.Cos(p.rotation)*-bulletSpawnOffset,
		}

		bullet := NewBullet(spawnPos, p.rotation)
		p.game.AddBullet(bullet)
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

func (p *Player) Collider() Rect {
	bounds := p.sprite.Bounds()

	return Rect{
		X: p.position.X,
		Y: p.position.Y,
		W: float64(bounds.Dx()),
		H: float64(bounds.Dy()),
	}
}