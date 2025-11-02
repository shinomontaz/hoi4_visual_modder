package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Canvas represents a drawable canvas with pan and zoom
type Canvas struct {
	Width  int
	Height int
	
	// Camera position (pan)
	OffsetX float64
	OffsetY float64
	
	// Zoom level (1.0 = 100%)
	Zoom float64
	
	// Grid settings
	GridSize    int
	GridColor   color.Color
	ShowGrid    bool
	
	// Background
	BackgroundColor color.Color
}

// NewCanvas creates a new canvas
func NewCanvas(width, height int) *Canvas {
	return &Canvas{
		Width:           width,
		Height:          height,
		OffsetX:         0,
		OffsetY:         0,
		Zoom:            1.0,
		GridSize:        50,
		GridColor:       color.RGBA{40, 40, 40, 255},
		ShowGrid:        true,
		BackgroundColor: color.RGBA{20, 20, 20, 255},
	}
}

// Update updates the canvas (handle pan/zoom input)
func (c *Canvas) Update() {
	// Pan with arrow keys
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		c.OffsetX += 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		c.OffsetX -= 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		c.OffsetY += 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		c.OffsetY -= 5
	}
	
	// Zoom with +/- keys
	if ebiten.IsKeyPressed(ebiten.KeyEqual) || ebiten.IsKeyPressed(ebiten.KeyNumpadAdd) {
		c.Zoom += 0.01
		if c.Zoom > 3.0 {
			c.Zoom = 3.0
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyMinus) || ebiten.IsKeyPressed(ebiten.KeyNumpadSubtract) {
		c.Zoom -= 0.01
		if c.Zoom < 0.3 {
			c.Zoom = 0.3
		}
	}
	
	// Reset with R key
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		c.OffsetX = 0
		c.OffsetY = 0
		c.Zoom = 1.0
	}
}

// Draw draws the canvas background and grid
func (c *Canvas) Draw(screen *ebiten.Image) {
	// Draw background
	screen.Fill(c.BackgroundColor)
	
	// Draw grid if enabled
	if c.ShowGrid {
		c.drawGrid(screen)
	}
}

// drawGrid draws the grid lines
func (c *Canvas) drawGrid(screen *ebiten.Image) {
	gridSize := float64(c.GridSize) * c.Zoom
	
	// Calculate visible grid range
	startX := int(-c.OffsetX / gridSize)
	endX := int((float64(c.Width) - c.OffsetX) / gridSize) + 1
	startY := int(-c.OffsetY / gridSize)
	endY := int((float64(c.Height) - c.OffsetY) / gridSize) + 1
	
	// Draw vertical lines
	for i := startX; i <= endX; i++ {
		x := float32(float64(i)*gridSize + c.OffsetX)
		if x >= 0 && x <= float32(c.Width) {
			vector.StrokeLine(screen, x, 0, x, float32(c.Height), 1, c.GridColor, false)
		}
	}
	
	// Draw horizontal lines
	for i := startY; i <= endY; i++ {
		y := float32(float64(i)*gridSize + c.OffsetY)
		if y >= 0 && y <= float32(c.Height) {
			vector.StrokeLine(screen, 0, y, float32(c.Width), y, 1, c.GridColor, false)
		}
	}
}

// WorldToScreen converts world coordinates to screen coordinates
func (c *Canvas) WorldToScreen(worldX, worldY int) (float64, float64) {
	screenX := float64(worldX)*c.Zoom + c.OffsetX
	screenY := float64(worldY)*c.Zoom + c.OffsetY
	return screenX, screenY
}

// ScreenToWorld converts screen coordinates to world coordinates
func (c *Canvas) ScreenToWorld(screenX, screenY float64) (int, int) {
	worldX := int((screenX - c.OffsetX) / c.Zoom)
	worldY := int((screenY - c.OffsetY) / c.Zoom)
	return worldX, worldY
}

// GridToWorld converts grid coordinates to world pixel coordinates
func (c *Canvas) GridToWorld(gridX, gridY int) (int, int) {
	worldX := gridX * c.GridSize
	worldY := gridY * c.GridSize
	return worldX, worldY
}
