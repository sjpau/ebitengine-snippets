package misc

import (
	"embed"
	"image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func LoadPNG(files *embed.FS, path string) (*ebiten.Image, error) {
	file, err := files.Open(path)
	Check(err)
	img, err := png.Decode(file)
	Check(err)
	return ebiten.NewImageFromImage(img), file.Close()
}

func ScaledFontSize(size float64, scaling float64) int {
	scaledSize := int(size * scaling)
	if scaledSize <= 0 {
		return 1
	}
	return scaledSize
}
