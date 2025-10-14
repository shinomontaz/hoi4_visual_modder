package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"hoi4_visual_modder/internal/app"
)

//go:embed all:web
var assets embed.FS

// main запускает приложение
func main() {
	// Создаем экземпляр приложения
	appInstance := app.NewApp()

	// Создаем приложение с опциями
	err := wails.Run(&options.App{
		Title:            "HOI4 Visual Modder",
		Width:            1200,
		Height:           800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        appInstance.OnStartup,
		Bind: []interface{}{
			appInstance,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
