package domain

// Position represents a coordinate on the focus/tech tree grid
type Position struct {
	X int // Horizontal position
	Y int // Vertical position
}

// NewPosition creates a new Position
func NewPosition(x, y int) Position {
	return Position{X: x, Y: y}
}

// Equals checks if two positions are equal
func (p Position) Equals(other Position) bool {
	return p.X == other.X && p.Y == other.Y
}

// Add returns a new position with added offsets
func (p Position) Add(dx, dy int) Position {
	return Position{X: p.X + dx, Y: p.Y + dy}
}
