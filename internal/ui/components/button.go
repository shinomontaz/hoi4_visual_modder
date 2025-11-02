package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Button represents a clickable button
type Button struct {
	X, Y          int
	Width, Height int
	Text          string
	
	hovered bool
	pressed bool
	clicked bool
}

// NewButton creates a new button
func NewButton(x, y, width, height int, text string) *Button {
	return &Button{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
		Text:   text,
	}
}

// Update updates the button state
func (b *Button) Update() {
	mx, my := ebiten.CursorPosition()
	
	// Check if mouse is over button
	b.hovered = mx >= b.X && mx < b.X+b.Width &&
		my >= b.Y && my < b.Y+b.Height
	
	// Check if button is pressed
	if b.hovered && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		b.pressed = true
	} else {
		b.pressed = false
	}
	
	// Check if button was clicked (released while hovered)
	b.clicked = false
	if b.hovered && inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		b.clicked = true
	}
}

// Draw renders the button
func (b *Button) Draw(screen *ebiten.Image) {
	// Choose color based on state
	var bgColor color.RGBA
	offsetY := 0
	
	if b.pressed {
		bgColor = color.RGBA{80, 80, 100, 255} // Lighter when pressed
		offsetY = 2                             // Slight offset when pressed
	} else if b.hovered {
		bgColor = color.RGBA{60, 60, 80, 255} // Medium when hovered
	} else {
		bgColor = color.RGBA{50, 50, 70, 255} // Dark when normal
	}
	
	// Draw button background
	ebitenutil.DrawRect(screen,
		float64(b.X),
		float64(b.Y+offsetY),
		float64(b.Width),
		float64(b.Height),
		bgColor)
	
	// Draw button border
	borderColor := color.RGBA{100, 100, 120, 255}
	if b.hovered {
		borderColor = color.RGBA{120, 120, 140, 255}
	}
	
	// Top border
	ebitenutil.DrawRect(screen, float64(b.X), float64(b.Y+offsetY), float64(b.Width), 2, borderColor)
	// Bottom border
	ebitenutil.DrawRect(screen, float64(b.X), float64(b.Y+b.Height+offsetY-2), float64(b.Width), 2, borderColor)
	// Left border
	ebitenutil.DrawRect(screen, float64(b.X), float64(b.Y+offsetY), 2, float64(b.Height), borderColor)
	// Right border
	ebitenutil.DrawRect(screen, float64(b.X+b.Width-2), float64(b.Y+offsetY), 2, float64(b.Height), borderColor)
	
	// Draw button text (centered)
	textWidth := len(b.Text) * 6 // Approximate width (6 pixels per char)
	textX := b.X + (b.Width-textWidth)/2
	textY := b.Y + (b.Height-16)/2 + offsetY // 16 is approximate text height
	
	ebitenutil.DebugPrintAt(screen, b.Text, textX, textY)
}

// IsClicked returns true if the button was clicked this frame
func (b *Button) IsClicked() bool {
	return b.clicked
}

// IsHovered returns true if the mouse is over the button
func (b *Button) IsHovered() bool {
	return b.hovered
}
