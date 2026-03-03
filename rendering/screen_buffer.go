package rendering

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type ScreenBuffer struct {
	img   *ebiten.Image
	color color.RGBA
}

func NewScreenBuffer(width, height int, color color.RGBA) *ScreenBuffer {
	img := ebiten.NewImage(width, height)
	img.Fill(color)
	return &ScreenBuffer{
		img:   img,
		color: color,
	}
}

func (sb *ScreenBuffer) Clear() {
	sb.img.Fill(sb.color)
}
