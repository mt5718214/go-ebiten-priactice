package game

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"

	assets "airplane/spaceshooterextension"
)


type Meteor struct {
	position Vector
	sprite *ebiten.Image
	movement Vector
	rotation float64
	rotationSpeed float64
}

func NewMeteor() *Meteor {
	sprite := assets.MeteorSprites[rand.Intn(len(assets.MeteorSprites))]

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