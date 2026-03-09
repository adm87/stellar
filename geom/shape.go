package geom

// Shape defines a geometric shape that can be used for collision detection, spatial queries, and other geometric operations.
type Shape interface {
	AABB() (x, y, width, height float64) // AABB returns the minimum axis-aligned bounding box of the shape.
	Contains(x, y float64) bool          // Contains checks if the point (x, y) is inside the shape.
}
