package geom

// Rectangle represents an axis-aligned rectangle defined by its top-left corner (X, Y) and its dimensions (Width, Height).
type Rectangle struct {
	X, Y          float64
	Width, Height float64
}

// ---------------------------------------------------------------
// Shape Interface Implementation
// ---------------------------------------------------------------

// AABB returns the minimum axis-aligned bounding box of the rectangle.
func (r *Rectangle) AABB() (x, y, width, height float64) {
	return r.X, r.Y, r.Width, r.Height
}

// Contains checks if the point (x, y) is inside the rectangle.
func (r *Rectangle) Contains(x, y float64) bool {
	return x >= r.X && x <= r.X+r.Width && y >= r.Y && y <= r.Y+r.Height
}
