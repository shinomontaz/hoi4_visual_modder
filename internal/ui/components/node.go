package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Node represents a visual node (technology or focus)
type Node struct {
	ID       string
	Title    string
	X        int // Grid X position
	Y        int // Grid Y position
	Width    int
	Height   int
	
	// Visual properties
	Color       color.Color
	BorderColor color.Color
	TextColor   color.Color
	
	// State
	IsSelected bool
	IsHovered  bool
}

// NewNode creates a new node
func NewNode(id, title string, x, y int) *Node {
	return &Node{
		ID:          id,
		Title:       title,
		X:           x,
		Y:           y,
		Width:       120,
		Height:      60,
		Color:       color.RGBA{60, 60, 80, 255},
		BorderColor: color.RGBA{100, 100, 120, 255},
		TextColor:   color.RGBA{220, 220, 220, 255},
		IsSelected:  false,
		IsHovered:   false,
	}
}

// Draw draws the node on the canvas
func (n *Node) Draw(screen *ebiten.Image, canvas *Canvas) {
	// Convert grid position to world position
	worldX, worldY := canvas.GridToWorld(n.X, n.Y)
	
	// Convert world position to screen position
	screenX, screenY := canvas.WorldToScreen(worldX, worldY)
	
	// Calculate scaled dimensions
	scaledWidth := float32(n.Width) * float32(canvas.Zoom)
	scaledHeight := float32(n.Height) * float32(canvas.Zoom)
	
	// Skip if node is outside visible area
	if screenX+float64(scaledWidth) < 0 || screenX > float64(canvas.Width) ||
		screenY+float64(scaledHeight) < 0 || screenY > float64(canvas.Height) {
		return
	}
	
	x := float32(screenX)
	y := float32(screenY)
	
	// Draw node background
	nodeColor := n.Color
	if n.IsSelected {
		nodeColor = color.RGBA{80, 120, 160, 255}
	} else if n.IsHovered {
		nodeColor = color.RGBA{70, 70, 90, 255}
	}
	
	vector.DrawFilledRect(screen, x, y, scaledWidth, scaledHeight, nodeColor, false)
	
	// Draw border
	borderColor := n.BorderColor
	if n.IsSelected {
		borderColor = color.RGBA{120, 180, 220, 255}
	}
	
	borderWidth := float32(2 * canvas.Zoom)
	if borderWidth < 1 {
		borderWidth = 1
	}
	
	vector.StrokeRect(screen, x, y, scaledWidth, scaledHeight, borderWidth, borderColor, false)
	
	// Draw text (only if zoom is reasonable)
	if canvas.Zoom > 0.5 {
		n.drawText(screen, x, y, scaledWidth, scaledHeight, canvas.Zoom)
	}
}

// drawText draws the node text using ebitenutil
func (n *Node) drawText(screen *ebiten.Image, x, y, width, height float32, zoom float64) {
	// Only draw text if zoom is reasonable
	if zoom < 0.7 {
		return
	}
	
	// Calculate text position (approximately centered)
	// ebitenutil.DebugPrintAt uses a fixed-width font, so we estimate
	textLen := len(n.Title)
	charWidth := 6.0 * zoom  // Approximate character width
	charHeight := 8.0 * zoom // Approximate character height
	
	textX := int(x + (width-float32(textLen)*float32(charWidth))/2)
	textY := int(y + (height-float32(charHeight))/2)
	
	// Draw text
	ebitenutil.DebugPrintAt(screen, n.Title, textX, textY)
}

// Contains checks if a screen point is inside the node
func (n *Node) Contains(screenX, screenY float64, canvas *Canvas) bool {
	worldX, worldY := canvas.GridToWorld(n.X, n.Y)
	nodeScreenX, nodeScreenY := canvas.WorldToScreen(worldX, worldY)
	
	scaledWidth := float64(n.Width) * canvas.Zoom
	scaledHeight := float64(n.Height) * canvas.Zoom
	
	return screenX >= nodeScreenX && screenX <= nodeScreenX+scaledWidth &&
		screenY >= nodeScreenY && screenY <= nodeScreenY+scaledHeight
}
