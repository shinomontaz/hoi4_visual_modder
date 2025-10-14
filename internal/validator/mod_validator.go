package validator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"hoi4_visual_modder/internal/models"
)

// ModValidator проверяет структуру каталога мода HOI4
type ModValidator struct {
	validator *models.ModStructureValidator
}

// NewModValidator создает новый валидатор мода
func NewModValidator() *ModValidator {
	return &ModValidator{
		validator: models.NewModStructureValidator(),
	}
}

// ValidateModDirectory проверяет каталог мода и возвращает информацию о нем
func (mv *ModValidator) ValidateModDirectory(basePath string) (*models.ValidationResult, error) {
	if basePath == "" {
		return &models.ValidationResult{
			IsValid: false,
			Errors:  []string{"Путь к каталогу мода не может быть пустым"},
		}, nil
	}

	// Проверяем, существует ли каталог
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		return &models.ValidationResult{
			IsValid: false,
			Errors:  []string{fmt.Sprintf("Каталог не существует: %s", basePath)},
		}, nil
	}

	modInfo := &models.ModInfo{
		BasePath:           basePath,
		Name:               filepath.Base(basePath),
		IsValid:            true,
		ValidationErrors:   []string{},
		NationalFocusFiles: []string{},
		TechnologyFiles:    []string{},
	}

	result := &models.ValidationResult{
		IsValid:  true,
		Errors:   []string{},
		Warnings: []string{},
		ModInfo:  modInfo,
	}

	// Проверяем обязательные каталоги
	mv.validateRequiredDirectories(modInfo, result)

	// Сканируем файлы национальных фокусов
	mv.scanNationalFocusFiles(modInfo, result)

	// Сканируем файлы технологий
	mv.scanTechnologyFiles(modInfo, result)

	// Проверяем файлы интерфейса
	mv.validateInterfaceFiles(modInfo, result)

	// Проверяем каталоги графики
	mv.validateGraphicsDirectories(modInfo, result)

	// Финальная проверка валидности
	if len(modInfo.ValidationErrors) > 0 {
		modInfo.IsValid = false
		result.IsValid = false
		result.Errors = modInfo.ValidationErrors
	}

	return result, nil
}

// validateRequiredDirectories проверяет наличие обязательных каталогов
func (mv *ModValidator) validateRequiredDirectories(modInfo *models.ModInfo, result *models.ValidationResult) {
	for _, dir := range mv.validator.RequiredDirectories {
		fullPath := filepath.Join(modInfo.BasePath, dir)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			modInfo.AddValidationErrorf("Отсутствует обязательный каталог: %s", dir)
		}
	}
}

// scanNationalFocusFiles сканирует файлы национальных фокусов
func (mv *ModValidator) scanNationalFocusFiles(modInfo *models.ModInfo, result *models.ValidationResult) {
	focusPath := modInfo.GetNationalFocusPath()
	
	if _, err := os.Stat(focusPath); os.IsNotExist(err) {
		result.Warnings = append(result.Warnings, "Каталог national_focus не найден")
		return
	}

	files, err := os.ReadDir(focusPath)
	if err != nil {
		modInfo.AddValidationErrorf("Ошибка чтения каталога national_focus: %v", err)
		return
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".txt") {
			modInfo.NationalFocusFiles = append(modInfo.NationalFocusFiles, file.Name())
		}
	}

	if len(modInfo.NationalFocusFiles) == 0 {
		result.Warnings = append(result.Warnings, "В каталоге national_focus не найдено .txt файлов")
	}
}

// scanTechnologyFiles сканирует файлы технологий
func (mv *ModValidator) scanTechnologyFiles(modInfo *models.ModInfo, result *models.ValidationResult) {
	techPath := modInfo.GetTechnologiesPath()
	
	if _, err := os.Stat(techPath); os.IsNotExist(err) {
		result.Warnings = append(result.Warnings, "Каталог technologies не найден")
		return
	}

	files, err := os.ReadDir(techPath)
	if err != nil {
		modInfo.AddValidationErrorf("Ошибка чтения каталога technologies: %v", err)
		return
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".txt") {
			modInfo.TechnologyFiles = append(modInfo.TechnologyFiles, file.Name())
		}
	}

	if len(modInfo.TechnologyFiles) == 0 {
		result.Warnings = append(result.Warnings, "В каталоге technologies не найдено .txt файлов")
	}
}

// validateInterfaceFiles проверяет файлы интерфейса
func (mv *ModValidator) validateInterfaceFiles(modInfo *models.ModInfo, result *models.ValidationResult) {
	interfacePath := modInfo.GetInterfacePath()
	
	if _, err := os.Stat(interfacePath); os.IsNotExist(err) {
		result.Warnings = append(result.Warnings, "Каталог interface не найден")
		return
	}

	// Проверяем goals.gfx
	goalsGfxPath := filepath.Join(interfacePath, "goals.gfx")
	if _, err := os.Stat(goalsGfxPath); err == nil {
		modInfo.InterfaceFiles.GoalsGfx = goalsGfxPath
	} else {
		result.Warnings = append(result.Warnings, "Файл interface/goals.gfx не найден")
	}

	// Проверяем countrytechtreeview.gfx
	techGfxPath := filepath.Join(interfacePath, "countrytechtreeview.gfx")
	if _, err := os.Stat(techGfxPath); err == nil {
		modInfo.InterfaceFiles.CountryTechTreeViewGfx = techGfxPath
	} else {
		result.Warnings = append(result.Warnings, "Файл interface/countrytechtreeview.gfx не найден")
	}

	// Сканируем другие .gfx файлы
	files, err := os.ReadDir(interfacePath)
	if err == nil {
		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".gfx") {
				fileName := file.Name()
				if fileName != "goals.gfx" && fileName != "countrytechtreeview.gfx" {
					modInfo.InterfaceFiles.OtherGfxFiles = append(modInfo.InterfaceFiles.OtherGfxFiles, fileName)
				}
			}
		}
	}
}

// validateGraphicsDirectories проверяет каталоги графики
func (mv *ModValidator) validateGraphicsDirectories(modInfo *models.ModInfo, result *models.ValidationResult) {
	gfxPath := modInfo.GetGraphicsPath()
	
	if _, err := os.Stat(gfxPath); os.IsNotExist(err) {
		result.Warnings = append(result.Warnings, "Каталог gfx не найден")
		return
	}

	// Проверяем каталог gfx/interface/goals
	goalsDir := filepath.Join(gfxPath, "interface", "goals")
	if _, err := os.Stat(goalsDir); err == nil {
		modInfo.GraphicsFiles.GoalsDirectory = goalsDir
	} else {
		result.Warnings = append(result.Warnings, "Каталог gfx/interface/goals не найден")
	}

	// Проверяем каталог gfx/interface/technologies
	techDir := filepath.Join(gfxPath, "interface", "technologies")
	if _, err := os.Stat(techDir); err == nil {
		modInfo.GraphicsFiles.TechnologiesDirectory = techDir
	} else {
		result.Warnings = append(result.Warnings, "Каталог gfx/interface/technologies не найден")
	}
}

// GetModSummary возвращает краткую сводку о моде
func (mv *ModValidator) GetModSummary(modInfo *models.ModInfo) string {
	if !modInfo.IsValid {
		return fmt.Sprintf("Мод '%s' содержит ошибки: %d", modInfo.Name, len(modInfo.ValidationErrors))
	}

	return fmt.Sprintf("Мод '%s': %d файлов фокусов, %d файлов технологий", 
		modInfo.Name, 
		len(modInfo.NationalFocusFiles), 
		len(modInfo.TechnologyFiles))
}
