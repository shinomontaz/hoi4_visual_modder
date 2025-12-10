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

	// Tech categories list (shown when tech button clicked)
	showTechCategories bool
	techList           *components.ScrollableList
	techCategories     []string

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

	// Create scrollable list for tech categories
	scene.techList = components.NewScrollableList(440, 420, 400, 300, 6)

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

	// Create display names for scrollable list
	displayNames := make([]string, len(s.techCategories))
	for i, category := range s.techCategories {
		displayNames[i] = ctx.GetLocalizedFolderName(category)
	}

	s.techList.SetItems(displayNames)
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

	// Handle tech categories list
	if s.showTechCategories {
		prevIndex := s.techList.GetSelectedIndex()
		s.techList.Update()

		// Check if an item was just selected (index changed)
		selectedIndex := s.techList.GetSelectedIndex()
		if selectedIndex >= 0 && selectedIndex < len(s.techCategories) && selectedIndex != prevIndex {
			s.handleTechCategoryClick(s.techCategories[selectedIndex])
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

	// Detect sub-trees for this folder
	loader := app.NewTechnologyLoader(ctx.ModPath, ctx.GamePath)
	subTrees := loader.DetectSubTrees(category, technologies)
	if len(subTrees) > 0 {
		techTree.SubTrees[category] = subTrees
		println("Detected", len(subTrees), "sub-trees in", category)
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

	// Draw tech categories list if shown
	if s.showTechCategories {
		ebitenutil.DebugPrintAt(screen, "Technology Categories:", 440, 390)
		s.techList.Draw(screen)
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
