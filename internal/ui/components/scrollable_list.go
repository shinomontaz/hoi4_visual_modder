package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// ScrollableList is a scrollable list component
type ScrollableList struct {
	x, y          int
	width, height int
	items         []string
	displayCount  int // max visible items
	scrollOffset  int // current scroll position
	selectedIndex int
	hoveredIndex  int

	itemHeight    int
	showScrollbar bool

	// Colors
	bgColor        color.Color
	hoverColor     color.Color
	selectedColor  color.Color
	textColor      color.Color
	scrollbarColor color.Color
}

// NewScrollableList creates a new scrollable list
func NewScrollableList(x, y, width, height int, displayCount int) *ScrollableList {
	return &ScrollableList{
		x:              x,
		y:              y,
		width:          width,
		height:         height,
		displayCount:   displayCount,
		scrollOffset:   0,
		selectedIndex:  -1,
		hoveredIndex:   -1,
		itemHeight:     40,
		showScrollbar:  true,
		bgColor:        color.RGBA{50, 50, 50, 255},
		hoverColor:     color.RGBA{70, 70, 70, 255},
		selectedColor:  color.RGBA{80, 120, 180, 255},
		textColor:      color.RGBA{255, 255, 255, 255},
		scrollbarColor: color.RGBA{100, 100, 100, 255},
	}
}

// SetItems sets the list items
func (sl *ScrollableList) SetItems(items []string) {
	sl.items = items
	sl.scrollOffset = 0
	sl.selectedIndex = -1
	sl.hoveredIndex = -1
}

// Update updates the list state
func (sl *ScrollableList) Update() {
	if len(sl.items) == 0 {
		return
	}

	// Get mouse position
	mx, my := ebiten.CursorPosition()

	// Check if mouse is over list
	if mx >= sl.x && mx <= sl.x+sl.width && my >= sl.y && my <= sl.y+sl.height {
		// Handle mouse wheel scrolling
		_, dy := ebiten.Wheel()
		if dy != 0 {
			sl.HandleMouseWheel(dy)
		}

		// Check which item is hovered
		relativeY := my - sl.y
		itemIndex := relativeY / sl.itemHeight
		absoluteIndex := sl.scrollOffset + itemIndex

		if absoluteIndex >= 0 && absoluteIndex < len(sl.items) {
			sl.hoveredIndex = absoluteIndex

			// Handle click
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				sl.selectedIndex = absoluteIndex
			}
		} else {
			sl.hoveredIndex = -1
		}
	} else {
		sl.hoveredIndex = -1
	}
}

// Draw renders the scrollable list
func (sl *ScrollableList) Draw(screen *ebiten.Image) {
	if len(sl.items) == 0 {
		// Draw empty state
		ebitenutil.DebugPrintAt(screen, "No items", sl.x+10, sl.y+10)
		return
	}

	// Calculate visible range
	maxVisible := sl.displayCount
	if len(sl.items) < maxVisible {
		maxVisible = len(sl.items)
	}

	endIndex := sl.scrollOffset + maxVisible
	if endIndex > len(sl.items) {
		endIndex = len(sl.items)
	}

	// Draw visible items
	for i := sl.scrollOffset; i < endIndex; i++ {
		itemY := sl.y + (i-sl.scrollOffset)*sl.itemHeight

		// Determine background color
		bgColor := sl.bgColor
		if i == sl.selectedIndex {
			bgColor = sl.selectedColor
		} else if i == sl.hoveredIndex {
			bgColor = sl.hoverColor
		}

		// Draw item background
		drawRect(screen, sl.x, itemY, sl.width, sl.itemHeight, bgColor)

		// Draw item text
		textY := itemY + (sl.itemHeight-16)/2 + 8 // Center text vertically
		ebitenutil.DebugPrintAt(screen, sl.items[i], sl.x+10, textY)
	}

	// Draw scrollbar if needed
	if len(sl.items) > sl.displayCount && sl.showScrollbar {
		sl.drawScrollbar(screen)
	}

	// Draw "... X more" indicator if scrolled
	if endIndex < len(sl.items) {
		remaining := len(sl.items) - endIndex
		indicatorY := sl.y + sl.height - 20
		ebitenutil.DebugPrintAt(screen, "... "+string(rune(remaining+'0'))+" more", sl.x+10, indicatorY)
	}
}

// drawScrollbar draws the scrollbar
func (sl *ScrollableList) drawScrollbar(screen *ebiten.Image) {
	scrollbarX := sl.x + sl.width - 10
	scrollbarWidth := 8
	scrollbarHeight := sl.height

	// Draw scrollbar background
	drawRect(screen, scrollbarX, sl.y, scrollbarWidth, scrollbarHeight, color.RGBA{40, 40, 40, 255})

	// Calculate thumb size and position
	visibleRatio := float64(sl.displayCount) / float64(len(sl.items))
	thumbHeight := int(float64(scrollbarHeight) * visibleRatio)
	if thumbHeight < 20 {
		thumbHeight = 20
	}

	scrollRatio := float64(sl.scrollOffset) / float64(len(sl.items)-sl.displayCount)
	thumbY := sl.y + int(float64(scrollbarHeight-thumbHeight)*scrollRatio)

	// Draw scrollbar thumb
	drawRect(screen, scrollbarX, thumbY, scrollbarWidth, thumbHeight, sl.scrollbarColor)
}

// GetSelectedItem returns the currently selected item
func (sl *ScrollableList) GetSelectedItem() string {
	if sl.selectedIndex >= 0 && sl.selectedIndex < len(sl.items) {
		return sl.items[sl.selectedIndex]
	}
	return ""
}

// GetSelectedIndex returns the currently selected index
func (sl *ScrollableList) GetSelectedIndex() int {
	return sl.selectedIndex
}

// ScrollUp scrolls the list up
func (sl *ScrollableList) ScrollUp() {
	if sl.scrollOffset > 0 {
		sl.scrollOffset--
	}
}

// ScrollDown scrolls the list down
func (sl *ScrollableList) ScrollDown() {
	maxScroll := len(sl.items) - sl.displayCount
	if maxScroll < 0 {
		maxScroll = 0
	}
	if sl.scrollOffset < maxScroll {
		sl.scrollOffset++
	}
}

// HandleMouseWheel handles mouse wheel input
func (sl *ScrollableList) HandleMouseWheel(delta float64) {
	if delta > 0 {
		sl.ScrollUp()
	} else if delta < 0 {
		sl.ScrollDown()
	}
}

// drawRect draws a filled rectangle
func drawRect(screen *ebiten.Image, x, y, width, height int, clr color.Color) {
	rect := ebiten.NewImage(width, height)
	rect.Fill(clr)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(rect, op)
}
