package scenes

import (
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/shinomontaz/hoi4_visual_modder/internal/app"
	"github.com/shinomontaz/hoi4_visual_modder/internal/domain"
	"github.com/shinomontaz/hoi4_visual_modder/internal/parser"
	"github.com/shinomontaz/hoi4_visual_modder/internal/ui/components"
)

// CountrySelectionScene allows selecting a country from bookmarks
type CountrySelectionScene struct {
	manager           *SceneManager
	state             *app.State
	countries         []*domain.BookmarkCountry
	filteredCountries []*domain.BookmarkCountry
	selectedIndex     int
	scrollOffset      int

	// Filter state
	filterMajor bool
	filterMinor bool
	filterAll   bool
	searchText  string

	// Buttons
	backButton        *components.Button
	continueButton    *components.Button
	filterAllButton   *components.Button
	filterMajorButton *components.Button
	filterMinorButton *components.Button

	errorMessage string
	loading      bool
}

// NewCountrySelectionScene creates a new country selection scene
func NewCountrySelectionScene(manager *SceneManager, state *app.State) *CountrySelectionScene {
	scene := &CountrySelectionScene{
		manager:           manager,
		state:             state,
		countries:         make([]*domain.BookmarkCountry, 0),
		filteredCountries: make([]*domain.BookmarkCountry, 0),
		selectedIndex:     -1,
		scrollOffset:      0,
		filterAll:         true,
		loading:           true,
	}

	// Create buttons
	scene.backButton = components.NewButton(50, 650, 200, 50, "← Back")
	scene.continueButton = components.NewButton(1030, 650, 200, 50, "Continue →")
	scene.filterAllButton = components.NewButton(440, 150, 120, 40, "All")
	scene.filterMajorButton = components.NewButton(570, 150, 120, 40, "Major")
	scene.filterMinorButton = components.NewButton(700, 150, 120, 40, "Minor")

	// Load countries in background
	go scene.loadCountries()

	return scene
}

// loadCountries loads countries from bookmarks
func (s *CountrySelectionScene) loadCountries() {
	defer func() {
		s.loading = false
	}()

	// Get mod and game paths
	modPath := s.state.GetModPath()
	gamePath := s.state.GetGamePath()

	if modPath == "" {
		s.errorMessage = "Mod path not set"
		return
	}

	// Parse bookmarks
	bookmarkParser := parser.NewBookmarkParser(modPath, gamePath)
	bookmarks, err := bookmarkParser.ParseBookmarks()
	if err != nil {
		s.errorMessage = "Failed to load bookmarks: " + err.Error()
		return
	}

	// Extract all countries from all bookmarks
	countryMap := make(map[string]*domain.BookmarkCountry)
	for _, bookmark := range bookmarks {
		for _, country := range bookmark.Countries {
			// Use first occurrence of each country tag
			if _, exists := countryMap[country.Tag]; !exists {
				countryMap[country.Tag] = country
			}
		}
	}

	// Convert map to slice
	s.countries = make([]*domain.BookmarkCountry, 0, len(countryMap))
	for _, country := range countryMap {
		s.countries = append(s.countries, country)
	}

	// Apply initial filter
	s.applyFilter()
}

// Update updates the country selection scene
func (s *CountrySelectionScene) Update() error {
	// Update buttons
	s.backButton.Update()
	s.continueButton.Update()
	s.filterAllButton.Update()
	s.filterMajorButton.Update()
	s.filterMinorButton.Update()

	// Handle back button
	if s.backButton.IsClicked() {
		s.manager.SwitchTo(SceneStartup)
		return nil
	}

	// Handle continue button (only if country selected)
	if s.continueButton.IsClicked() && s.selectedIndex >= 0 && s.selectedIndex < len(s.filteredCountries) {
		// Set country context in state
		selectedCountry := s.filteredCountries[s.selectedIndex]
		s.state.SetCountryContext(selectedCountry)

		// Switch to country menu scene
		countryMenu := NewCountryMenuScene(s.manager, s.state)
		s.manager.AddScene("country_menu", countryMenu)
		s.manager.SwitchToNamed("country_menu")
	}

	// Handle filter buttons
	if s.filterAllButton.IsClicked() {
		s.filterAll = true
		s.filterMajor = false
		s.filterMinor = false
		s.applyFilter()
	}
	if s.filterMajorButton.IsClicked() {
		s.filterAll = false
		s.filterMajor = true
		s.filterMinor = false
		s.applyFilter()
	}
	if s.filterMinorButton.IsClicked() {
		s.filterAll = false
		s.filterMajor = false
		s.filterMinor = true
		s.applyFilter()
	}

	// Handle mouse wheel for scrolling
	_, dy := ebiten.Wheel()
	if dy != 0 {
		s.scrollOffset -= int(dy * 3)
		if s.scrollOffset < 0 {
			s.scrollOffset = 0
		}
		maxScroll := len(s.filteredCountries) - 10
		if maxScroll < 0 {
			maxScroll = 0
		}
		if s.scrollOffset > maxScroll {
			s.scrollOffset = maxScroll
		}
	}

	// Handle mouse click on country list
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		if mx >= 440 && mx <= 840 && my >= 210 && my <= 610 {
			// Calculate which country was clicked
			relY := my - 210
			index := relY/40 + s.scrollOffset
			if index >= 0 && index < len(s.filteredCountries) {
				s.selectedIndex = index
			}
		}
	}

	return nil
}

// applyFilter applies the current filter to countries
func (s *CountrySelectionScene) applyFilter() {
	s.filteredCountries = make([]*domain.BookmarkCountry, 0)

	for _, country := range s.countries {
		// Apply major/minor filter
		if s.filterMajor && !country.IsMajor {
			continue
		}
		if s.filterMinor && country.IsMajor {
			continue
		}

		// Apply search filter (if implemented)
		if s.searchText != "" {
			if !strings.Contains(strings.ToLower(country.Tag), strings.ToLower(s.searchText)) &&
				!strings.Contains(strings.ToLower(country.Name), strings.ToLower(s.searchText)) {
				continue
			}
		}

		s.filteredCountries = append(s.filteredCountries, country)
	}

	// Reset selection if out of bounds
	if s.selectedIndex >= len(s.filteredCountries) {
		s.selectedIndex = -1
	}
}

// Draw renders the country selection scene
func (s *CountrySelectionScene) Draw(screen *ebiten.Image) {
	// Clear screen
	screen.Fill(color.RGBA{30, 30, 30, 255})

	// Draw title
	ebitenutil.DebugPrintAt(screen, "Select Country", 550, 50)

	// Draw subtitle
	if s.state.ModDescriptor != nil {
		ebitenutil.DebugPrintAt(screen, "Mod: "+s.state.ModDescriptor.Name, 440, 90)
	}

	// Draw filter buttons
	s.filterAllButton.Draw(screen)
	s.filterMajorButton.Draw(screen)
	s.filterMinorButton.Draw(screen)

	// Highlight active filter
	if s.filterAll {
		ebitenutil.DebugPrintAt(screen, "[Active]", 460, 200)
	} else if s.filterMajor {
		ebitenutil.DebugPrintAt(screen, "[Active]", 590, 200)
	} else if s.filterMinor {
		ebitenutil.DebugPrintAt(screen, "[Active]", 720, 200)
	}

	// Draw loading message
	if s.loading {
		ebitenutil.DebugPrintAt(screen, "Loading countries...", 540, 350)
		return
	}

	// Draw country list
	if len(s.filteredCountries) == 0 {
		ebitenutil.DebugPrintAt(screen, "No countries found", 540, 350)
	} else {
		// Draw list background
		listX, listY := 440, 210
		listW := 400

		// Draw countries (10 visible at a time)
		visibleCount := 10
		for i := 0; i < visibleCount && (i+s.scrollOffset) < len(s.filteredCountries); i++ {
			index := i + s.scrollOffset
			country := s.filteredCountries[index]

			y := listY + i*40

			// Highlight selected
			if index == s.selectedIndex {
				// Draw selection background
				ebitenutil.DrawRect(screen, float64(listX), float64(y), float64(listW), 35, color.RGBA{70, 70, 100, 255})
			}

			// Draw country info
			typeLabel := country.GetTypeLabel()
			text := country.Tag + " - " + country.GetDisplayName() + " (" + typeLabel + ")"
			if country.Ideology != "" {
				text += " [" + country.Ideology + "]"
			}

			ebitenutil.DebugPrintAt(screen, text, listX+10, y+10)
		}

		// Draw scroll indicator
		if len(s.filteredCountries) > visibleCount {
			scrollText := "Scroll: " + string(rune(s.scrollOffset)) + "/" + string(rune(len(s.filteredCountries)-visibleCount))
			ebitenutil.DebugPrintAt(screen, scrollText, listX+listW+10, listY)
		}
	}

	// Draw buttons
	s.backButton.Draw(screen)
	if s.selectedIndex >= 0 {
		s.continueButton.Draw(screen)
	}

	// Draw error message
	if s.errorMessage != "" {
		ebitenutil.DebugPrintAt(screen, "Error: "+s.errorMessage, 440, 620)
	}

	// Draw hint
	ebitenutil.DebugPrintAt(screen, "Click on a country to select, then click Continue", 440, 680)
}

// OnEnter is called when entering this scene
func (s *CountrySelectionScene) OnEnter() {
	// Reset state
	s.selectedIndex = -1
	s.scrollOffset = 0
	s.errorMessage = ""
}

// OnExit is called when leaving this scene
func (s *CountrySelectionScene) OnExit() {
	// Cleanup
}
