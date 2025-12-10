package scenes

import (
	"fmt"
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/shinomontaz/hoi4_visual_modder/internal/domain"
	"github.com/shinomontaz/hoi4_visual_modder/internal/parser"
	"github.com/shinomontaz/hoi4_visual_modder/internal/ui/components"
)

// TechViewerScene displays technology tree visually
type TechViewerScene struct {
	manager      *SceneManager
	canvas       *components.Canvas
	nodes        []*components.Node
	technologies []*domain.Technology
	iconLoader   *components.IconLoader

	// UI state
	selectedNode *components.Node
	hoveredNode  *components.Node

	// Info panel
	showInfo bool
}

// NewTechViewerScene creates a new tech viewer scene
func NewTechViewerScene(manager *SceneManager, filePath string) *TechViewerScene {
	scene := &TechViewerScene{
		manager:  manager,
		canvas:   components.NewCanvas(1280, 720),
		nodes:    make([]*components.Node, 0),
		showInfo: true,
	}

	// Parse the technology file
	if err := scene.loadTechnologies(filePath); err != nil {
		fmt.Printf("Error loading technologies: %v\n", err)
		return scene
	}

	// Initialize icon loader with paths from manager state
	if manager.state != nil {
		modPath := manager.state.GetModPath()
		gamePath := manager.state.GetGamePath()

		scene.iconLoader = components.NewIconLoader(modPath)
		if gamePath != "" {
			scene.iconLoader.SetGamePath(gamePath)
		}
	} else {
		// Fallback: try to detect base path from file path
		scene.iconLoader = components.NewIconLoader(detectBasePath(filePath))
	}

	// Create nodes from technologies
	scene.createNodes()

	// Center view on first node
	if len(scene.nodes) > 0 {
		scene.centerOnNode(scene.nodes[0])
	}

	return scene
}

// NewTechViewerSceneWithTree creates a new tech viewer scene from a pre-loaded technology tree
func NewTechViewerSceneWithTree(manager *SceneManager, techTree *domain.TechnologyTree) *TechViewerScene {
	// Convert map to slice
	technologies := make([]*domain.Technology, 0, len(techTree.Technologies))
	for _, tech := range techTree.Technologies {
		technologies = append(technologies, tech)
	}

	scene := &TechViewerScene{
		manager:      manager,
		canvas:       components.NewCanvas(1280, 720),
		nodes:        make([]*components.Node, 0),
		showInfo:     true,
		technologies: technologies,
	}

	// Initialize icon loader with paths from manager state
	if manager.state != nil {
		modPath := manager.state.GetModPath()
		gamePath := manager.state.GetGamePath()

		scene.iconLoader = components.NewIconLoader(modPath)
		if gamePath != "" {
			scene.iconLoader.SetGamePath(gamePath)
		}
	}

	// Create nodes from technologies
	scene.createNodes()

	// Center view on first node
	if len(scene.nodes) > 0 {
		scene.centerOnNode(scene.nodes[0])
	}

	return scene
}

// detectBasePath tries to extract base path from file path
func detectBasePath(filePath string) string {
	// Simple heuristic: find "common" in path and go up one level
	// Example: C:/mods/mymod/common/technologies/file.txt -> C:/mods/mymod
	idx := -1
	for i := len(filePath) - 1; i >= 0; i-- {
		if filePath[i] == '/' || filePath[i] == '\\' {
			if idx == -1 {
				idx = i
			} else {
				// Check if this is "common"
				segment := filePath[idx+1 : i]
				if segment == "common" {
					return filePath[:i]
				}
				idx = i
			}
		}
	}
	return ""
}

// loadTechnologies loads and parses the technology file
func (s *TechViewerScene) loadTechnologies(filePath string) error {
	// Read file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Parse with lexer and parser
	p := parser.NewParser(string(content))
	program, err := p.Parse()
	if err != nil {
		return fmt.Errorf("failed to parse: %w", err)
	}

	// Parse technologies
	techParser := parser.NewTechParser()
	technologies, err := techParser.ParseTechnologies(program)
	if err != nil {
		return fmt.Errorf("failed to parse technologies: %w", err)
	}

	s.technologies = technologies
	return nil
}

// createNodes creates visual nodes from technologies
func (s *TechViewerScene) createNodes() {
	for _, tech := range s.technologies {
		node := components.NewNode(tech.ID, tech.ID, tech.Position.X, tech.Position.Y)

		// Load icon if icon loader is available
		if s.iconLoader != nil {
			// Use tech ID as icon name (standard HOI4 convention)
			node.Icon = s.iconLoader.LoadTechIcon(tech.ID)
		}

		s.nodes = append(s.nodes, node)
	}
}

// centerOnNode centers the view on a specific node
func (s *TechViewerScene) centerOnNode(node *components.Node) {
	worldX, worldY := s.canvas.GridToWorld(node.X, node.Y)
	s.canvas.OffsetX = float64(s.canvas.Width/2) - float64(worldX) - float64(node.Width/2)
	s.canvas.OffsetY = float64(s.canvas.Height/2) - float64(worldY) - float64(node.Height/2)
}

// Update updates the scene
func (s *TechViewerScene) Update() error {
	// Update canvas (pan/zoom)
	s.canvas.Update()

	// Handle mouse hover
	mouseX, mouseY := ebiten.CursorPosition()
	s.hoveredNode = nil
	for _, node := range s.nodes {
		if node.Contains(float64(mouseX), float64(mouseY), s.canvas) {
			s.hoveredNode = node
			node.IsHovered = true
		} else {
			node.IsHovered = false
		}
	}

	// Handle mouse click
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if s.hoveredNode != nil {
			if s.selectedNode != nil {
				s.selectedNode.IsSelected = false
			}
			s.selectedNode = s.hoveredNode
			s.selectedNode.IsSelected = true
		}
	}

	// ESC to go back
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		s.manager.SwitchTo(SceneStartup)
	}

	// Toggle info panel with I key
	if ebiten.IsKeyPressed(ebiten.KeyI) {
		s.showInfo = !s.showInfo
	}

	return nil
}

// Draw draws the scene
func (s *TechViewerScene) Draw(screen *ebiten.Image) {
	// Draw canvas (background + grid)
	s.canvas.Draw(screen)

	// Draw connection lines (TODO: implement later)
	// s.drawConnections(screen)

	// Draw nodes
	for _, node := range s.nodes {
		node.Draw(screen, s.canvas)
	}

	// Draw UI overlay
	s.drawUI(screen)
}

// drawUI draws the UI overlay (info panel, controls)
func (s *TechViewerScene) drawUI(screen *ebiten.Image) {
	// Draw info panel
	if s.showInfo {
		s.drawInfoPanel(screen)
	}

	// Draw controls help
	s.drawControls(screen)

	// Draw selected node info
	if s.selectedNode != nil {
		s.drawNodeInfo(screen)
	}
}

// drawInfoPanel draws the info panel
func (s *TechViewerScene) drawInfoPanel(screen *ebiten.Image) {
	panelX := float32(10)
	panelY := float32(10)
	panelWidth := float32(250)
	panelHeight := float32(100)

	// Background
	vector.DrawFilledRect(screen, panelX, panelY, panelWidth, panelHeight,
		color.RGBA{30, 30, 30, 220}, false)
	vector.StrokeRect(screen, panelX, panelY, panelWidth, panelHeight, 2,
		color.RGBA{80, 80, 80, 255}, false)

	// Text
	y := int(panelY + 20)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Technologies: %d", len(s.technologies)), int(panelX+10), y)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Zoom: %.1f%%", s.canvas.Zoom*100), int(panelX+10), y+15)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Offset: (%.0f, %.0f)", s.canvas.OffsetX, s.canvas.OffsetY), int(panelX+10), y+30)
}

// drawControls draws the controls help
func (s *TechViewerScene) drawControls(screen *ebiten.Image) {
	controlsText := "Arrow Keys: Pan | +/-: Zoom | R: Reset | I: Toggle Info | ESC: Back"

	// Draw at bottom center (approximate)
	x := s.canvas.Width/2 - len(controlsText)*3
	y := s.canvas.Height - 30

	ebitenutil.DebugPrintAt(screen, controlsText, x, y)
}

// drawNodeInfo draws info about the selected node
func (s *TechViewerScene) drawNodeInfo(screen *ebiten.Image) {
	// Find the technology for this node
	var tech *domain.Technology
	for _, t := range s.technologies {
		if t.ID == s.selectedNode.ID {
			tech = t
			break
		}
	}

	if tech == nil {
		return
	}

	// Draw panel on the right side
	panelX := float32(s.canvas.Width - 310)
	panelY := float32(10)
	panelWidth := float32(300)
	panelHeight := float32(200)

	// Background
	vector.DrawFilledRect(screen, panelX, panelY, panelWidth, panelHeight,
		color.RGBA{30, 30, 30, 220}, false)
	vector.StrokeRect(screen, panelX, panelY, panelWidth, panelHeight, 2,
		color.RGBA{80, 120, 160, 255}, false)

	// Text
	y := int(panelY + 20)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("ID: %s", tech.ID), int(panelX+10), y)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Folder: %s", tech.Folder), int(panelX+10), y+15)

	// Show position with aliases
	posStr := fmt.Sprintf("Position: (%d, %d)", tech.Position.X, tech.Position.Y)
	if tech.Position.XVar != "" || tech.Position.YVar != "" {
		xVar := tech.Position.XVar
		if xVar == "" {
			xVar = fmt.Sprintf("%d", tech.Position.X)
		}
		yVar := tech.Position.YVar
		if yVar == "" {
			yVar = fmt.Sprintf("%d", tech.Position.Y)
		}
		posStr = fmt.Sprintf("Position: { x = %s y = %s }", xVar, yVar)
	}
	ebitenutil.DebugPrintAt(screen, posStr, int(panelX+10), y+30)

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Research Cost: %.1f", tech.ResearchCost), int(panelX+10), y+45)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Paths: %d", len(tech.Paths)), int(panelX+10), y+60)
}

// OnEnter is called when entering the scene
func (s *TechViewerScene) OnEnter() {
	// Nothing to do for now
}

// OnExit is called when exiting the scene
func (s *TechViewerScene) OnExit() {
	// Nothing to do for now
}
