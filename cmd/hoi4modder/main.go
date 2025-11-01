package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/shinomontaz/hoi4_visual_modder/internal/ui"
)

const (
	screenWidth  = 1280
	screenHeight = 720
	windowTitle  = "HOI4 Visual Modder"
)

func main() {
	// Create the main game instance
	game := ui.NewGame()

	// Configure window
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle(windowTitle)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// Run the game loop
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
