package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	assets "airplane/spaceshooterextension"
)

const (
	bulletSpeedPerSecond = 350.0
)

type Bullet struct {
	position Vector
	sprite *ebiten.Image
	rotation float64
}

func NewBullet(pos Vector, rotation float64) *Bullet {
	bounds := assets.BulletImage.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	pos.X -= halfW
	pos.Y -= halfH

	return &Bullet{
		position: pos,
		sprite: assets.BulletImage,
		rotation: rotation,
	}
}

func (b *Bullet) Update() error {
	speed := bulletSpeedPerSecond / float64(ebiten.TPS())
	b.position.X += math.Sin(b.rotation) * speed
	b.position.Y += math.Cos(b.rotation) * -speed
	// b.position.Y += math.Cos(b.rotation) * speed
	return nil
}

func (b *Bullet) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	bounds := b.sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	op.GeoM.Translate(-halfW, -halfH)
	op.GeoM.Rotate(b.rotation)
	op.GeoM.Translate(b.position.X, b.position.Y)
	screen.DrawImage(b.sprite, op)
}

func (b *Bullet) Collider() Rect {
	bounds := b.sprite.Bounds()

	return Rect{
		X: b.position.X,
		Y: b.position.Y,
		W: float64(bounds.Dx()),
		H: float64(bounds.Dy()),
	}
}