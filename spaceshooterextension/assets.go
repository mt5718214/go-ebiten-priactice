package assets

import (
	"embed"
	"image"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

//關鍵字: go:embed，透過註解就可以直接讀取檔案
//go:embed PNG/*
var assets embed.FS
var PlayerImage = mustLoadImage("PNG/Sprites_X2/Ships/spaceShips_005.png")
var MeteorSprites = mustLoadImages("PNG/Sprites_X2/Meteors")
var BulletImage = mustLoadImage("PNG/Sprites_X2/Rocket parts/spaceRocketParts_023.png")

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