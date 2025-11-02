package scenes

import (
	"image/color"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/shinomontaz/hoi4_visual_modder/internal/app"
	"github.com/shinomontaz/hoi4_visual_modder/internal/ui/components"
	"github.com/sqweek/dialog"
)

// StartupScene is the initial scene for mod selection
type StartupScene struct {
	manager      *SceneManager
	state        *app.State
	openButton   *components.Button
	errorMessage string
}

// NewStartupScene creates a new StartupScene
func NewStartupScene(manager *SceneManager, state *app.State) *StartupScene {
	// Create "Open File..." button (centered on screen)
	openButton := components.NewButton(440, 300, 400, 60, "Open File...")
	
	return &StartupScene{
		manager:    manager,
		state:      state,
		openButton: openButton,
	}
}

// Update updates the startup scene
func (s *StartupScene) Update() error {
	// Update button
	s.openButton.Update()
	
	// Handle button click or Ctrl+O shortcut
	openFilePicker := s.openButton.IsClicked()
	if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyO) {
		openFilePicker = true
	}
	
	if openFilePicker {
		s.errorMessage = "" // Clear previous errors
		
		// Open native file picker dialog
		filePath, err := dialog.File().
			Filter("HOI4 Files", "txt").
			Title("Select Focus or Technology File").
			Load()
		
		if err != nil {
			// User cancelled or error occurred
			if err.Error() != "Cancelled" {
				s.errorMessage = "Error opening file picker: " + err.Error()
			}
			return nil
		}
		
		// Load file using ModLoader
		err = s.state.LoadFile(filePath)
		if err != nil {
			s.errorMessage = err.Error()
			// Show error dialog
			dialog.Message("%s", err.Error()).Title("Error").Error()
			return nil
		}
		
		// Success! Switch to file viewer
		s.manager.SwitchTo(SceneFileViewer)
	}
	
	return nil
}

// Draw renders the startup scene
func (s *StartupScene) Draw(screen *ebiten.Image) {
	// Clear screen with dark background
	screen.Fill(color.RGBA{30, 30, 30, 255})
	
	// Draw title
	ebitenutil.DebugPrintAt(screen, "HOI4 Visual Modder", 520, 100)
	
	// Draw subtitle
	ebitenutil.DebugPrintAt(screen, "Visual editor for Hearts of Iron IV mod files", 440, 140)
	
	// Draw instructions
	ebitenutil.DebugPrintAt(screen, "Click the button below or press Ctrl+O to open a file", 420, 220)
	
	// Draw button
	s.openButton.Draw(screen)
	
	// Draw file info if loaded
	if s.state.BasePath != "" {
		y := 400
		ebitenutil.DebugPrintAt(screen, "Selected File:", 440, y)
		
		fileName := filepath.Base(s.state.SelectedFilePath)
		ebitenutil.DebugPrintAt(screen, "  File: "+fileName, 440, y+20)
		ebitenutil.DebugPrintAt(screen, "  Type: "+s.state.FileType.String(), 440, y+40)
		ebitenutil.DebugPrintAt(screen, "  Base: "+s.state.BasePath, 440, y+60)
	}
	
	// Draw error message if any
	if s.errorMessage != "" {
		ebitenutil.DebugPrintAt(screen, "Error: "+s.errorMessage, 440, 500)
	}
	
	// Draw hint
	ebitenutil.DebugPrintAt(screen, "Supported: .txt files from common/national_focus/ or common/technologies/", 340, 650)
}

// OnEnter is called when entering this scene
func (s *StartupScene) OnEnter() {
	// Initialize scene resources
}

// OnExit is called when leaving this scene
func (s *StartupScene) OnExit() {
	// Cleanup scene resources
}
