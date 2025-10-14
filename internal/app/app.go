package app

import (
	"context"
	"fmt"

	"hoi4_visual_modder/internal/models"
	"hoi4_visual_modder/internal/validator"
)

// App структура содержит данные приложения
type App struct {
	ctx         context.Context
	validator   *validator.ModValidator
	currentMod  *models.ModInfo
}

// NewApp создает новый экземпляр приложения
func NewApp() *App {
	return &App{
		validator: validator.NewModValidator(),
	}
}

// OnStartup вызывается при запуске приложения
func (a *App) OnStartup(ctx context.Context) {
	a.ctx = ctx
}

// SelectModDirectory позволяет пользователю выбрать каталог мода
func (a *App) SelectModDirectory(path string) (*models.ValidationResult, error) {
	if path == "" {
		return &models.ValidationResult{
			IsValid: false,
			Errors:  []string{"Путь не может быть пустым"},
		}, nil
	}

	result, err := a.validator.ValidateModDirectory(path)
	if err != nil {
		return nil, fmt.Errorf("ошибка валидации каталога мода: %w", err)
	}

	if result.IsValid {
		a.currentMod = result.ModInfo
	}

	return result, nil
}

// GetCurrentMod возвращает информацию о текущем выбранном моде
func (a *App) GetCurrentMod() *models.ModInfo {
	return a.currentMod
}

// GetModSummary возвращает краткую информацию о моде
func (a *App) GetModSummary() string {
	if a.currentMod == nil {
		return "Мод не выбран"
	}
	return a.validator.GetModSummary(a.currentMod)
}

// GetNationalFocusFiles возвращает список файлов национальных фокусов
func (a *App) GetNationalFocusFiles() []string {
	if a.currentMod == nil {
		return []string{}
	}
	return a.currentMod.NationalFocusFiles
}

// GetTechnologyFiles возвращает список файлов технологий
func (a *App) GetTechnologyFiles() []string {
	if a.currentMod == nil {
		return []string{}
	}
	return a.currentMod.TechnologyFiles
}

// ValidateCurrentMod повторно валидирует текущий мод
func (a *App) ValidateCurrentMod() (*models.ValidationResult, error) {
	if a.currentMod == nil {
		return &models.ValidationResult{
			IsValid: false,
			Errors:  []string{"Мод не выбран"},
		}, nil
	}

	return a.validator.ValidateModDirectory(a.currentMod.BasePath)
}

// GetAppInfo возвращает информацию о приложении
func (a *App) GetAppInfo() map[string]interface{} {
	return map[string]interface{}{
		"name":    "HOI4 Visual Modder",
		"version": "1.0.0",
		"author":  "HOI4 Tools",
	}
}
