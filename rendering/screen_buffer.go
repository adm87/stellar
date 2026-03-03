package rendering

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// BufferResizeMode defines how the screen buffer should handle resizing when the game window size changes.
type BufferResizeMode uint8

const (
	BufferResizeNone BufferResizeMode = iota
	BufferResizeStretch
	BufferResizeMaintainWidth
	BufferResizeMaintainHeight
)

// ScreenBuffer represents an off-screen rendering target that can be drawn to and then applied to the main screen with various resizing options.
type ScreenBuffer struct {
	img        *ebiten.Image
	matrix     ebiten.GeoM
	options    ebiten.DrawImageOptions
	color      color.RGBA
	resizeMode BufferResizeMode
	isDirty    bool
}

// NewScreenBuffer creates a new ScreenBuffer with the specified width, height, background color, and resize mode.
// It initializes the underlying image and fills it with the background color.
func NewScreenBuffer(width, height int, color color.RGBA, resizeMode BufferResizeMode) *ScreenBuffer {
	img := ebiten.NewImage(width, height)
	img.Fill(color)
	return &ScreenBuffer{
		img:        img,
		color:      color,
		resizeMode: resizeMode,
		options:    ebiten.DrawImageOptions{},
		isDirty:    true,
	}
}

// Clear fills the screen buffer with its background color, effectively clearing any previous drawings.
func (sb *ScreenBuffer) Clear() {
	sb.img.Fill(sb.color)
}

// Image returns the underlying ebiten.Image of the screen buffer, allowing it to be drawn to the main screen or
// used as a texture in other rendering operations.
func (sb *ScreenBuffer) Image() *ebiten.Image {
	return sb.img
}

// SetFilter sets the filter mode for the screen buffer.
func (sb *ScreenBuffer) SetFilter(filter ebiten.Filter) {
	sb.options.Filter = filter
}

// SetResizeMode changes the resize mode of the screen buffer, marking it as dirty to trigger a
// recalculation of the aspect ratio matrix on the next draw.
func (sb *ScreenBuffer) SetResizeMode(mode BufferResizeMode) {
	if sb.resizeMode != mode {
		sb.resizeMode = mode
		sb.isDirty = true
	}
}

func (sb *ScreenBuffer) SetBackgroundColor(color color.RGBA) {
	sb.color = color
}

// Resize changes the dimensions of the screen buffer to the specified width and height, filling it with the background color.
// It marks the buffer as dirty, indicating that the aspect ratio matrix needs to be recalculated before the next draw.
func (sb *ScreenBuffer) Resize(newWidth, newHeight int) {
	currentWidth := sb.img.Bounds().Dx()
	currentHeight := sb.img.Bounds().Dy()

	if newWidth == currentWidth && newHeight == currentHeight {
		return
	}

	sb.img.Deallocate()

	sb.img = ebiten.NewImage(newWidth, newHeight)
	sb.img.Fill(sb.color)

	sb.isDirty = true
}

// ApplyTo draws the screen buffer onto the target image using the specified draw options,
// applying the appropriate scaling and translation based on the resize mode.
// It recalculates the aspect ratio matrix if the buffer is marked as dirty, ensuring that the buffer
// is rendered correctly according to the current window size and resize mode.
func (sb *ScreenBuffer) ApplyTo(target *ebiten.Image) {
	if sb.isDirty {
		sb.matrix = calculateAspectRatioMatrix(
			sb.img.Bounds().Dx(),
			sb.img.Bounds().Dy(),
			target.Bounds().Dx(),
			target.Bounds().Dy(),
			sb.resizeMode,
		)
		sb.isDirty = false
		sb.options.GeoM = sb.matrix
	}

	target.DrawImage(sb.img, &sb.options)
}

func calculateAspectRatioMatrix(srcWidth, srcHeight, dstWidth, dstHeight int, mode BufferResizeMode) ebiten.GeoM {
	var scaleX, scaleY float64
	var offsetX, offsetY float64

	switch mode {
	case BufferResizeStretch:
		scaleX = float64(dstWidth) / float64(srcWidth)
		scaleY = float64(dstHeight) / float64(srcHeight)
	case BufferResizeMaintainWidth:
		scaleX = float64(dstWidth) / float64(srcWidth)
		scaleY = scaleX
		offsetY = (float64(dstHeight) - float64(srcHeight)*scaleY) / 2
	case BufferResizeMaintainHeight:
		scaleY = float64(dstHeight) / float64(srcHeight)
		scaleX = scaleY
		offsetX = (float64(dstWidth) - float64(srcWidth)*scaleX) / 2
	default:
		scaleX = 1.0
		scaleY = 1.0
	}

	var matrix ebiten.GeoM
	matrix.Scale(scaleX, scaleY)
	matrix.Translate(offsetX, offsetY)

	return matrix
}
