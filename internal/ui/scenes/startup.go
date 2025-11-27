package scenes

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/shinomontaz/hoi4_visual_modder/internal/app"
	"github.com/shinomontaz/hoi4_visual_modder/internal/ui/components"
	"github.com/sqweek/dialog"
)

// StartupScene is the initial scene for mod selection
type StartupScene struct {
	manager          *SceneManager
	state            *app.State
	selectModButton  *components.Button
	selectGameButton *components.Button
	autoDetectButton *components.Button
	continueButton   *components.Button
	errorMessage     string
	infoMessage      string
}

// NewStartupScene creates a new StartupScene
func NewStartupScene(manager *SceneManager, state *app.State) *StartupScene {
	// Create buttons
	selectModButton := components.NewButton(440, 200, 400, 50, "Select Mod (.mod file)")
	selectGameButton := components.NewButton(440, 270, 250, 50, "Select Game Folder")
	autoDetectButton := components.NewButton(700, 270, 140, 50, "Auto-detect")
	continueButton := components.NewButton(440, 500, 400, 60, "Continue →")

	return &StartupScene{
		manager:          manager,
		state:            state,
		selectModButton:  selectModButton,
		selectGameButton: selectGameButton,
		autoDetectButton: autoDetectButton,
		continueButton:   continueButton,
	}
}

// Update updates the startup scene
func (s *StartupScene) Update() error {
	// Update all buttons
	s.selectModButton.Update()
	s.selectGameButton.Update()
	s.autoDetectButton.Update()
	s.continueButton.Update()

	// Handle "Select Mod" button
	if s.selectModButton.IsClicked() {
		s.handleModSelection()
	}

	// Handle "Select Game" button
	if s.selectGameButton.IsClicked() {
		s.handleGameSelection()
	}

	// Handle "Auto-detect" button
	if s.autoDetectButton.IsClicked() {
		s.handleGameAutoDetect()
	}

	// Handle "Continue" button (only if both mod and game are selected)
	if s.continueButton.IsClicked() && s.canContinue() {
		// Switch to country selection scene
		countryScene := NewCountrySelectionScene(s.manager, s.state)
		s.manager.AddScene("country_selection", countryScene)
		s.manager.SwitchToNamed("country_selection")
	}

	return nil
}

// handleModSelection opens file picker for .mod file
func (s *StartupScene) handleModSelection() {
	s.errorMessage = ""
	s.infoMessage = ""

	// Open file picker for .mod files
	filePath, err := dialog.File().
		Filter("HOI4 Mod Files", "mod").
		Title("Select Mod Descriptor File").
		Load()

	if err != nil {
		if err.Error() != "Cancelled" {
			s.errorMessage = "Error: " + err.Error()
		}
		return
	}

	// Load and parse mod descriptor
	modDesc, err := app.LoadModDescriptor(filePath)
	if err != nil {
		s.errorMessage = "Failed to load mod: " + err.Error()
		dialog.Message("%s", err.Error()).Title("Mod Load Error").Error()
		return
	}

	// Save to state
	if err := s.state.SetModDescriptor(modDesc); err != nil {
		s.errorMessage = "Failed to save mod config: " + err.Error()
		return
	}

	s.infoMessage = "Mod loaded: " + modDesc.Name
}

// handleGameSelection opens folder picker for game installation
func (s *StartupScene) handleGameSelection() {
	s.errorMessage = ""
	s.infoMessage = ""

	// Open folder picker
	folderPath, err := dialog.Directory().
		Title("Select HOI4 Installation Folder").
		Browse()

	if err != nil {
		if err.Error() != "Cancelled" {
			s.errorMessage = "Error: " + err.Error()
		}
		return
	}

	// Validate game installation
	game, err := app.ValidateGameInstallation(folderPath)
	if err != nil {
		s.errorMessage = "Invalid game folder: " + err.Error()
		dialog.Message("%s", err.Error()).Title("Game Validation Error").Error()
		return
	}

	// Save to state
	if err := s.state.SetGameInstallation(game); err != nil {
		s.errorMessage = "Failed to save game config: " + err.Error()
		return
	}

	s.infoMessage = "Game found: " + game.Path
}

// handleGameAutoDetect tries to auto-detect game installation
func (s *StartupScene) handleGameAutoDetect() {
	s.errorMessage = ""
	s.infoMessage = "Searching for HOI4 installation..."

	// Try auto-detect
	gamePath, err := app.AutoDetectGamePath()
	if err != nil {
		s.errorMessage = "Auto-detect failed: " + err.Error()
		return
	}

	// Validate
	game, err := app.ValidateGameInstallation(gamePath)
	if err != nil {
		s.errorMessage = "Validation failed: " + err.Error()
		return
	}

	// Save to state
	if err := s.state.SetGameInstallation(game); err != nil {
		s.errorMessage = "Failed to save game config: " + err.Error()
		return
	}

	s.infoMessage = "Game auto-detected: " + game.Path
}

// canContinue checks if both mod and game are selected
func (s *StartupScene) canContinue() bool {
	return s.state.ModDescriptor != nil && s.state.GameInstallation != nil
}

// Draw renders the startup scene
func (s *StartupScene) Draw(screen *ebiten.Image) {
	// Clear screen with dark background
	screen.Fill(color.RGBA{30, 30, 30, 255})

	// Draw title
	ebitenutil.DebugPrintAt(screen, "HOI4 Visual Modder", 520, 80)

	// Draw subtitle
	ebitenutil.DebugPrintAt(screen, "Setup: Select Mod and Game", 500, 120)

	// Draw instructions
	ebitenutil.DebugPrintAt(screen, "1. Select your mod's .mod file", 440, 170)
	ebitenutil.DebugPrintAt(screen, "2. Select or auto-detect HOI4 installation", 440, 240)

	// Draw buttons
	s.selectModButton.Draw(screen)
	s.selectGameButton.Draw(screen)
	s.autoDetectButton.Draw(screen)

	// Draw mod info if loaded
	if s.state.ModDescriptor != nil {
		y := 350
		mod := s.state.ModDescriptor
		ebitenutil.DebugPrintAt(screen, "✓ Mod: "+mod.Name, 440, y)
		ebitenutil.DebugPrintAt(screen, "  Version: "+mod.Version+" (Game: "+mod.SupportedVersion+")", 440, y+20)
		ebitenutil.DebugPrintAt(screen, "  Path: "+mod.ModFolderPath, 440, y+40)
	}

	// Draw game info if loaded
	if s.state.GameInstallation != nil {
		y := 350
		if s.state.ModDescriptor != nil {
			y = 420 // Move down if mod is also shown
		}
		game := s.state.GameInstallation
		ebitenutil.DebugPrintAt(screen, "✓ Game: Hearts of Iron IV", 440, y)
		ebitenutil.DebugPrintAt(screen, "  Path: "+game.Path, 440, y+20)
	}

	// Draw continue button if both are selected
	if s.canContinue() {
		s.continueButton.Draw(screen)
	}

	// Draw info message
	if s.infoMessage != "" {
		ebitenutil.DebugPrintAt(screen, s.infoMessage, 440, 580)
	}

	// Draw error message if any
	if s.errorMessage != "" {
		ebitenutil.DebugPrintAt(screen, "Error: "+s.errorMessage, 440, 600)
	}
}

// OnEnter is called when entering this scene
func (s *StartupScene) OnEnter() {
	// Initialize scene resources
}

// OnExit is called when leaving this scene
func (s *StartupScene) OnExit() {
	// Cleanup scene resources
}
