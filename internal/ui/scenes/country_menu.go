package scenes

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/shinomontaz/hoi4_visual_modder/internal/app"
	"github.com/shinomontaz/hoi4_visual_modder/internal/domain"
	"github.com/shinomontaz/hoi4_visual_modder/internal/ui/components"
)

// CountryMenuScene displays the main menu for a selected country
type CountryMenuScene struct {
	manager *SceneManager
	state   *app.State

	// Buttons
	focusTreeButton *components.Button
	techButton      *components.Button
	backButton      *components.Button

	// Tech category buttons (shown when tech button clicked)
	showTechCategories  bool
	techCategoryButtons []*components.Button
	techCategories      []string

	errorMessage string
}

// NewCountryMenuScene creates a new country menu scene
func NewCountryMenuScene(manager *SceneManager, state *app.State) *CountryMenuScene {
	scene := &CountryMenuScene{
		manager: manager,
		state:   state,
	}

	// Create main buttons
	scene.focusTreeButton = components.NewButton(440, 250, 400, 60, "National Focus Tree")
	scene.techButton = components.NewButton(440, 330, 400, 60, "Technologies")
	scene.backButton = components.NewButton(50, 650, 200, 50, "← Back")

	// Load tech categories
	scene.loadTechCategories()

	return scene
}

// loadTechCategories loads available technology categories
func (s *CountryMenuScene) loadTechCategories() {
	ctx := s.state.GetCountryContext()
	if ctx == nil {
		return
	}

	s.techCategories = ctx.TechFolders

	// Create buttons for each category
	s.techCategoryButtons = make([]*components.Button, 0)

	startY := 420
	for i, category := range s.techCategories {
		// Get localized display name from context
		displayName := ctx.GetLocalizedFolderName(category)
		btn := components.NewButton(440, startY+i*50, 400, 45, "  • "+displayName)
		s.techCategoryButtons = append(s.techCategoryButtons, btn)
	}
}

// Update updates the country menu scene
func (s *CountryMenuScene) Update() error {
	// Update main buttons
	s.focusTreeButton.Update()
	s.techButton.Update()
	s.backButton.Update()

	// Handle back button
	if s.backButton.IsClicked() {
		s.manager.SwitchToNamed("country_selection")
		return nil
	}

	// Handle focus tree button
	if s.focusTreeButton.IsClicked() {
		s.handleFocusTreeClick()
	}

	// Handle tech button (toggle categories)
	if s.techButton.IsClicked() {
		s.showTechCategories = !s.showTechCategories
	}

	// Handle tech category buttons
	if s.showTechCategories {
		for i, btn := range s.techCategoryButtons {
			btn.Update()
			if btn.IsClicked() {
				s.handleTechCategoryClick(s.techCategories[i])
			}
		}
	}

	return nil
}

// handleFocusTreeClick handles clicking the focus tree button
func (s *CountryMenuScene) handleFocusTreeClick() {
	ctx := s.state.GetCountryContext()
	if ctx == nil {
		s.errorMessage = "No country selected"
		return
	}

	focusPath, err := ctx.GetFocusPath()
	if err != nil {
		s.errorMessage = "No focus tree available for " + ctx.GetTag()
		return
	}

	// TODO: Create FocusViewerScene and switch to it
	// For now, just show message
	s.errorMessage = "Focus tree: " + focusPath + " (Viewer not yet implemented)"
}

// handleTechCategoryClick handles clicking a tech category
func (s *CountryMenuScene) handleTechCategoryClick(category string) {
	ctx := s.state.GetCountryContext()
	if ctx == nil {
		s.errorMessage = "No country selected"
		return
	}

	// Load technologies for this folder
	technologies, err := ctx.LoadTechnologiesForFolder(category)
	if err != nil {
		s.errorMessage = "Failed to load tech: " + err.Error()
		return
	}

	// Create technology tree from loaded technologies
	techTree := domain.NewTechnologyTree()
	for _, tech := range technologies {
		techTree.AddTechnology(tech)
	}

	// Set in state
	s.state.TechnologyTree = techTree

	// Create and switch to tech viewer
	techViewer := NewTechViewerSceneWithTree(s.manager, techTree)
	s.manager.AddScene("tech_viewer", techViewer)
	s.manager.SwitchToNamed("tech_viewer")
}

// Draw renders the country menu scene
func (s *CountryMenuScene) Draw(screen *ebiten.Image) {
	// Clear screen
	screen.Fill(color.RGBA{30, 30, 30, 255})

	ctx := s.state.GetCountryContext()
	if ctx == nil {
		ebitenutil.DebugPrintAt(screen, "Error: No country selected", 440, 300)
		s.backButton.Draw(screen)
		return
	}

	// Draw title
	title := ctx.GetDisplayName() + " (" + ctx.GetTag() + ")"
	ebitenutil.DebugPrintAt(screen, title, 550, 80)

	// Draw subtitle
	if s.state.ModDescriptor != nil {
		ebitenutil.DebugPrintAt(screen, "Mod: "+s.state.ModDescriptor.Name, 440, 120)
	}

	// Draw instructions
	ebitenutil.DebugPrintAt(screen, "Select what to view:", 440, 200)

	// Draw main buttons
	s.focusTreeButton.Draw(screen)
	s.techButton.Draw(screen)

	// Show focus tree availability
	if !ctx.HasFocusTree() {
		ebitenutil.DebugPrintAt(screen, "(Using generic focus tree)", 850, 270)
	}

	// Show tech categories count
	techInfo := "(" + string(rune(len(s.techCategories))) + " categories available)"
	ebitenutil.DebugPrintAt(screen, techInfo, 850, 350)

	// Draw tech categories if shown
	if s.showTechCategories && len(s.techCategoryButtons) > 0 {
		ebitenutil.DebugPrintAt(screen, "Technology Categories:", 440, 390)

		// Draw category buttons (max 5 visible)
		maxVisible := 5
		for i := 0; i < maxVisible && i < len(s.techCategoryButtons); i++ {
			s.techCategoryButtons[i].Draw(screen)
		}

		if len(s.techCategories) > maxVisible {
			remaining := len(s.techCategories) - maxVisible
			ebitenutil.DebugPrintAt(screen, "... and "+string(rune(remaining))+" more", 460, 420+maxVisible*50)
		}
	}

	// Draw back button
	s.backButton.Draw(screen)

	// Draw error message
	if s.errorMessage != "" {
		ebitenutil.DebugPrintAt(screen, s.errorMessage, 440, 620)
	}

	// Draw hint
	ebitenutil.DebugPrintAt(screen, "Click on an option to view", 440, 680)
}

// OnEnter is called when entering this scene
func (s *CountryMenuScene) OnEnter() {
	s.errorMessage = ""
	s.showTechCategories = false
}

// OnExit is called when leaving this scene
func (s *CountryMenuScene) OnExit() {
	// Cleanup
}
