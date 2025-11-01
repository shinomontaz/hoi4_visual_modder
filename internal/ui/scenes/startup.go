package scenes

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// StartupScene is the initial scene for mod selection
type StartupScene struct {
	manager *SceneManager
}

// NewStartupScene creates a new StartupScene
func NewStartupScene(manager *SceneManager) *StartupScene {
	return &StartupScene{
		manager: manager,
	}
}

// Update updates the startup scene
func (s *StartupScene) Update() error {
	// TODO: Implement file browser and mode selection
	// For now, just a placeholder
	return nil
}

// Draw renders the startup scene
func (s *StartupScene) Draw(screen *ebiten.Image) {
	// Clear screen with dark background
	screen.Fill(color.RGBA{30, 30, 30, 255})
	
	// Draw placeholder text
	ebitenutil.DebugPrint(screen, "HOI4 Visual Modder - Startup Scene\n\nTODO: Implement mod directory selection")
}

// OnEnter is called when entering this scene
func (s *StartupScene) OnEnter() {
	// Initialize scene resources
}

// OnExit is called when leaving this scene
func (s *StartupScene) OnExit() {
	// Cleanup scene resources
}
