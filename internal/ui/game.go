package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/shinomontaz/hoi4_visual_modder/internal/app"
	"github.com/shinomontaz/hoi4_visual_modder/internal/ui/scenes"
)

// Game represents the main game instance
type Game struct {
	sceneManager *scenes.SceneManager
	state        *app.State
	width        int
	height       int
}

// NewGame creates a new Game instance
func NewGame() *Game {
	state := app.NewState()
	
	game := &Game{
		state:  state,
		width:  1280,
		height: 720,
	}
	
	// Initialize scene manager with startup scene
	game.sceneManager = scenes.NewSceneManager(state)
	
	return game
}

// Update updates the game state
func (g *Game) Update() error {
	return g.sceneManager.Update()
}

// Draw renders the game
func (g *Game) Draw(screen *ebiten.Image) {
	g.sceneManager.Draw(screen)
}

// Layout returns the game's logical screen size
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.width = outsideWidth
	g.height = outsideHeight
	return outsideWidth, outsideHeight
}
