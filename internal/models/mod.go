package models

import (
	"fmt"
	"path/filepath"
)

// ModInfo содержит информацию о структуре мода HOI4
type ModInfo struct {
	BasePath             string                `json:"basePath"`
	Name                 string                `json:"name"`
	IsValid              bool                  `json:"isValid"`
	ValidationErrors     []string              `json:"validationErrors"`
	NationalFocusFiles   []string              `json:"nationalFocusFiles"`
	TechnologyFiles      []string              `json:"technologyFiles"`
	InterfaceFiles       InterfaceFiles        `json:"interfaceFiles"`
	GraphicsFiles        GraphicsFiles         `json:"graphicsFiles"`
}

// InterfaceFiles содержит пути к файлам интерфейса
type InterfaceFiles struct {
	GoalsGfx              string   `json:"goalsGfx"`              // interface/goals.gfx
	CountryTechTreeViewGfx string  `json:"countryTechTreeViewGfx"` // interface/countrytechtreeview.gfx
	OtherGfxFiles         []string `json:"otherGfxFiles"`
}

// GraphicsFiles содержит пути к каталогам с графикой
type GraphicsFiles struct {
	GoalsDirectory        string `json:"goalsDirectory"`        // gfx/interface/goals/
	TechnologiesDirectory string `json:"technologiesDirectory"` // gfx/interface/technologies/
}

// ModStructureValidator проверяет структуру каталога мода
type ModStructureValidator struct {
	RequiredDirectories []string
	RequiredFiles       []string
	OptionalDirectories []string
}

// NewModStructureValidator создает новый валидатор структуры мода
func NewModStructureValidator() *ModStructureValidator {
	return &ModStructureValidator{
		RequiredDirectories: []string{
			"common",
			"common/national_focus",
			"common/technologies",
		},
		OptionalDirectories: []string{
			"interface",
			"gfx",
			"gfx/interface",
			"gfx/interface/goals",
			"gfx/interface/technologies",
		},
		RequiredFiles: []string{
			// Базовые файлы мода могут отсутствовать в некоторых модах
		},
	}
}

// ValidationResult содержит результат валидации
type ValidationResult struct {
	IsValid    bool     `json:"isValid"`
	Errors     []string `json:"errors"`
	Warnings   []string `json:"warnings"`
	ModInfo    *ModInfo `json:"modInfo,omitempty"`
}

// FileInfo содержит информацию о файле
type FileInfo struct {
	Path         string `json:"path"`
	RelativePath string `json:"relativePath"`
	Name         string `json:"name"`
	Size         int64  `json:"size"`
	IsDirectory  bool   `json:"isDirectory"`
}

// GetNationalFocusPath возвращает путь к каталогу национальных фокусов
func (m *ModInfo) GetNationalFocusPath() string {
	return filepath.Join(m.BasePath, "common", "national_focus")
}

// GetTechnologiesPath возвращает путь к каталогу технологий
func (m *ModInfo) GetTechnologiesPath() string {
	return filepath.Join(m.BasePath, "common", "technologies")
}

// GetInterfacePath возвращает путь к каталогу интерфейса
func (m *ModInfo) GetInterfacePath() string {
	return filepath.Join(m.BasePath, "interface")
}

// GetGraphicsPath возвращает путь к каталогу графики
func (m *ModInfo) GetGraphicsPath() string {
	return filepath.Join(m.BasePath, "gfx")
}

// AddValidationError добавляет ошибку валидации
func (m *ModInfo) AddValidationError(err string) {
	m.ValidationErrors = append(m.ValidationErrors, err)
	m.IsValid = false
}

// AddValidationErrorf добавляет отформатированную ошибку валидации
func (m *ModInfo) AddValidationErrorf(format string, args ...interface{}) {
	m.AddValidationError(fmt.Sprintf(format, args...))
}
