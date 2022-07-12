/*
import "github.com/hajimehoshi/ebiten/v2"

func FlipHorizontal(source *ebiten.Image) *ebiten.Image {
	result := ebiten.NewImage(source.Size())
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(-1, 1)
	op.GeoM.Translate(float64(img.Bounds().Dx()), 0)
	return result.DrawImage(source, op)
}

func FlipVertical(source *ebiten.Image) *ebiten.Image {
	result := ebiten.NewImage(source.Size())
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(1, -1)
	op.GeoM.Translate(0, float64(img.Bounds().Dy()))
	return result.DrawImage(source, op)
}
*/
