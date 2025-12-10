package domain

// Position represents a coordinate on the focus/tech tree grid
type Position struct {
	X    int    // Horizontal position (resolved numeric value)
	Y    int    // Vertical position (resolved numeric value)
	XVar string // Variable name for X (e.g., "@RADAR", "@HQ")
	YVar string // Variable name for Y (e.g., "@1940", "@1936")
}

// NewPosition creates a new Position with numeric coordinates
func NewPosition(x, y int) Position {
	return Position{X: x, Y: y, XVar: "", YVar: ""}
}

// NewPositionWithVars creates a new Position with both numeric and variable names
func NewPositionWithVars(x, y int, xVar, yVar string) Position {
	return Position{X: x, Y: y, XVar: xVar, YVar: yVar}
}

// Equals checks if two positions are equal
func (p Position) Equals(other Position) bool {
	return p.X == other.X && p.Y == other.Y
}

// Add returns a new position with added offsets
func (p Position) Add(dx, dy int) Position {
	return Position{X: p.X + dx, Y: p.Y + dy}
}
