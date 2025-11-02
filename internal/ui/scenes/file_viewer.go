package scenes

import (
	"image/color"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/shinomontaz/hoi4_visual_modder/internal/app"
)

// FileViewerScene displays the content of a selected file
type FileViewerScene struct {
	manager      *SceneManager
	state        *app.State
	scrollOffset int
	lines        []string
}

// NewFileViewerScene creates a new FileViewerScene
func NewFileViewerScene(manager *SceneManager, state *app.State) *FileViewerScene {
	return &FileViewerScene{
		manager: manager,
		state:   state,
	}
}

// Update updates the file viewer scene
func (s *FileViewerScene) Update() error {
	// Handle back button (Escape key)
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		s.manager.SwitchTo(SceneStartup)
		return nil
	}
	
	// Handle scrolling
	_, dy := ebiten.Wheel()
	if dy != 0 {
		s.scrollOffset -= int(dy * 20)
		if s.scrollOffset < 0 {
			s.scrollOffset = 0
		}
		maxScroll := len(s.lines)*16 - 600
		if maxScroll < 0 {
			maxScroll = 0
		}
		if s.scrollOffset > maxScroll {
			s.scrollOffset = maxScroll
		}
	}
	
	// Handle arrow key scrolling
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		s.scrollOffset -= 2
		if s.scrollOffset < 0 {
			s.scrollOffset = 0
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		s.scrollOffset += 2
		maxScroll := len(s.lines)*16 - 600
		if maxScroll < 0 {
			maxScroll = 0
		}
		if s.scrollOffset > maxScroll {
			s.scrollOffset = maxScroll
		}
	}
	
	return nil
}

// Draw renders the file viewer scene
func (s *FileViewerScene) Draw(screen *ebiten.Image) {
	// Clear screen with dark background
	screen.Fill(color.RGBA{30, 30, 30, 255})
	
	// Draw header
	ebitenutil.DebugPrintAt(screen, "File Viewer", 20, 20)

	// Display file metadata
	if s.state.SelectedFilePath != "" {
		ebitenutil.DebugPrintAt(screen, "File: "+filepath.Base(s.state.SelectedFilePath), 20, 40)
		ebitenutil.DebugPrintAt(screen, "Type: "+s.state.FileType.String(), 20, 60)
		ebitenutil.DebugPrintAt(screen, "Base Path: "+s.state.BasePath, 20, 80)
	}

	ebitenutil.DebugPrintAt(screen, "Press ESC to go back", 20, 100)
	
	// Draw file content
	if len(s.lines) > 0 {
		contentY := 140
		lineHeight := 16
		
		// Calculate visible range
		startLine := s.scrollOffset / lineHeight
		endLine := startLine + 40 // Show ~40 lines
		
		if startLine < 0 {
			startLine = 0
		}
		if endLine > len(s.lines) {
			endLine = len(s.lines)
		}
		
		// Draw visible lines
		for i := startLine; i < endLine; i++ {
			y := contentY + (i-startLine)*lineHeight
			if y >= contentY && y < 700 {
				// Truncate long lines
				line := s.lines[i]
				if len(line) > 150 {
					line = line[:150] + "..."
				}
				ebitenutil.DebugPrintAt(screen, line, 30, y)
			}
		}
		
		// Draw scrollbar indicator
		if len(s.lines) > 40 {
			scrollbarHeight := 500
			scrollbarY := 120
			thumbHeight := 50
			thumbY := scrollbarY + int(float64(s.scrollOffset)/float64(len(s.lines)*lineHeight)*float64(scrollbarHeight))
			
			ebitenutil.DrawRect(screen, 1250, float64(scrollbarY), 10, float64(scrollbarHeight), color.RGBA{50, 50, 50, 255})
			ebitenutil.DrawRect(screen, 1250, float64(thumbY), 10, float64(thumbHeight), color.RGBA{100, 100, 120, 255})
		}
	} else {
		ebitenutil.DebugPrintAt(screen, "No content to display", 30, 120)
	}
}

// OnEnter is called when entering this scene
func (s *FileViewerScene) OnEnter() {
	// Split content into lines for display
	s.lines = strings.Split(s.state.FileContent, "\n")
	s.scrollOffset = 0
}

// OnExit is called when leaving this scene
func (s *FileViewerScene) OnExit() {
	// Cleanup
}
