package types

type Rect struct {
	X, Y, Width, Height int
}

func (r Rect) Overlaps(other Rect) bool {
	return r.X < other.X+other.Width &&
		r.X+r.Width > other.X &&
		r.Y < other.Y+other.Height &&
		r.Y+r.Height > other.Y
}
