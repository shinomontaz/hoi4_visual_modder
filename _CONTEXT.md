# HOI4 Visual Modder - Project Context

> **For AI Agents:** This file provides project overview, goals, and technical context. See _PLAN.md for current status, roadmap, and implementation details.

## 🎯 Project Overview

**HOI4 Visual Modder** - это инструмент для визуального редактирования деревьев национальных фокусов и технологий в модах Hearts of Iron IV.

**What it does:**
- Визуально отображает деревья фокусов и технологий на интерактивной сетке
- Парсит файлы модов HOI4 в формате Paradox scripting language
- Позволяет редактировать позиции, свойства и связи между элементами
- Генерирует корректные файлы обратно в формат игры

**Current State:** Phase 1 (MVP) - структура проекта создана, domain модели реализованы, native file picker работает, следующий шаг - реализация лексера для парсинга.

## 🚀 Goals & Objectives

### Primary Goals
- Создать интуитивный визуальный редактор для .txt файлов HOI4
- Поддержать два основных режима: национальные фокусы и технологии
- Обеспечить валидацию и проверку целостности деревьев
- Генерировать корректные файлы в формате Paradox scripting language

### Target Users
- Моддеры Hearts of Iron IV
- Разработчики модов, работающие с национальными фокусами
- Создатели технологических деревьев

## 📋 Core Features

**Note:** Detailed development roadmap with phases is in _PLAN.md

### Key Capabilities
- **Visualization:** Display focus trees and technology trees on interactive grid
- **Parsing:** Read and parse Paradox scripting language (.txt files)
- **Editing:** Drag & drop positioning, property editing, connection management
- **Icon Integration:** Load and display icons from .gfx and .dds files
- **Validation:** Real-time checking for circular dependencies, conflicts, invalid references
- **File Generation:** Save changes back to .txt and .gfx files with proper formatting

## 🎨 User Experience

### Workflow (Implemented)

**File Selection Flow:**
```
Startup Scene
    ↓ Click "Open File..." or Ctrl+O
Native File Picker Dialog
    ↓ Select .txt file from mod
ModLoader Processing
    ├─ Detect Base_path (mod root directory)
    ├─ Validate mod structure (common/, national_focus/, technologies/)
    ├─ Detect file type (Focus or Technology)
    └─ Load file content (UTF-8)
    ↓
File Viewer Scene
    └─ Display file metadata and content
```

**Planned Workflow (Future):**
1. **Запуск** → Native file picker для выбора файла ✅
2. **Парсинг** → Лексер и парсер обрабатывают файл
3. **Визуализация** → Отображение дерева на canvas
4. **Редактирование** → Визуальная работа с деревом
5. **Сохранение** → Экспорт в правильном формате

### Key Interactions
- Drag & drop для перемещения элементов
- Клик для выбора и редактирования свойств
- Визуальное создание/удаление связей
- Автоматическая валидация в реальном времени

## 📁 File Structure Context

### National Focus Files
- **Location:** `Base_path/common/national_focus/*.txt`
- **Format:** Paradox scripting language
- **Structure:** focus_tree с множественными focus блоками
- **Key Elements:** id, icon, position (x,y), prerequisites, cost, completion_reward
- **Details:** hoi4_focus_tree_documentation.md

### Technology Files  
- **Location:** `Base_path/common/technologies/*.txt`
- **Format:** Paradox scripting language
- **Structure:** technologies блок с tech определениями
- **Key Elements:** id, allow, effects, paths, research_cost, position, categories
- **Details:** hoi4_tech_structure_documentation.md

### Icon Files (Images)
- **Focus Icons Location:** `Base_path/gfx/interface/goals/*.dds`
- **Tech Icons Location:** `Base_path/gfx/interface/technologies/*.dds`
- **Focus GFX Definitions:** `Base_path/interface/goals.gfx`
- **Tech GFX Definitions:** `Base_path/interface/countrytechtreeview.gfx`
- **Format:** .dds (DirectDraw Surface) image files
- **GFX Structure:** spriteType blocks linking icon names to file paths
- **Details:** hoi4_images_rules.md

### HOI4 Scripting Language Info
- **Format:** Paradox scripting language - custom text format with nested blocks
- **Syntax:** Key-value pairs, blocks with `{ }`, comments with `#`
- **Focus Trees:** Use X-Y grid positioning system (absolute or relative)
- **Technologies:** Use folder-based positioning with X-Y coordinates
- **Prerequisites:** Support AND/OR logic through multiple blocks
- **Special Values:** Dates (1939.1.1), numbers (integers/floats), strings, identifiers
- **Nesting:** Deep nesting for complex structures (rewards, conditions, effects)

## 🔧 Technical Constraints
- Должен работать с существующими модами HOI4
- Совместимость с форматом Paradox scripting language
- Кроссплатформенность (Windows/Linux/Mac)
- Автономность (не требует интернета)

## 📊 Project Scope

### In Scope
✅ Визуальное редактирование национальных фокусов  
✅ Визуальное редактирование технологий  
✅ Парсинг и генерация .txt файлов  
✅ Валидация структур данных  
✅ Drag & drop интерфейс  

### Out of Scope (for now)
❌ Редактирование других типов файлов мода  
❌ Интеграция с Steam Workshop  
❌ Мультиплеерное редактирование  
❌ Версионирование изменений  
❌ Автоматическое тестирование в игре  


## 📝 Implementation Patterns & Standards

### Go Code Style
- **Package Structure:** Follow Go standard layout (cmd/, internal/, pkg/)
- **Naming:** Use Go conventions (PascalCase for exported, camelCase for internal)
- **Error Handling:** Always return errors, use wrapped errors with context
- **Testing:** Unit tests for all parsers and validators
- **Dependencies:** Minimal external dependencies (only Ebitengine for GUI)

### Data Models
- **Focus Structure:** Matches HOI4 focus block exactly (id, icon, position, prerequisites, etc.)
  - Prerequisites stored as `[][]string` (outer array = AND, inner array = OR)
  - MutuallyExclusive as `[]string` of focus IDs
  - Position can be absolute or relative to another focus
- **Technology Structure:** Matches HOI4 tech block (id, allow, effects, paths, etc.)
  - Effects stored as `map[string]map[string]float64` (category → modifier → value)
  - Paths as slice of structs with target ID and cost coefficient
  - XOR for mutually exclusive technologies
- **Position System:** X-Y grid coordinates as in game files (integers)
- **Validation:** Built-in validation methods on all structures
  - Circular dependency detection using recursive graph traversal
  - Prerequisite existence checks
  - Position conflict detection

### Parser Architecture
- **Two-Phase Parsing:** Lexer (tokenization) → Parser (AST building)
- **Lexer:** Converts text to tokens (identifiers, strings, numbers, delimiters, keywords)
- **Parser:** Builds Abstract Syntax Tree from tokens
- **Specialized Parsers:** Separate logic for focus_tree, technologies, and .gfx files
- **Preservation:** Keep comments and formatting for round-trip editing

### File Operations
- **Encoding:** Always UTF-8 for .txt files
- **Parsing:** Preserve comments and formatting where possible
- **Backup:** Create .bak files before overwriting (atomic: write to .bak, then rename)
- **Atomic Writes:** Write to temp file first, then rename to target
- **Error Recovery:** Graceful handling of malformed files with detailed error messages

### Implemented Components

**ModLoader (internal/app/mod_loader.go):**
- `DetectBasePath(filePath)` - Extracts mod root from file path
  - Example: `C:\...\mod\MyMod\common\national_focus\file.txt` → `C:\...\mod\MyMod`
  - Handles Windows drive letters correctly
- `ValidateModStructure(basePath)` - Validates HOI4 mod directory structure
  - Checks for `common/` directory
  - Verifies `national_focus/` or `technologies/` subdirectories exist
- `DetectFileType(filePath)` - Determines file type (Focus/Technology) from path
- `LoadModFile(filePath)` - Complete file loading with validation

**UI Components (internal/ui/):**
- **Button** (components/button.go) - Reusable button with hover/pressed states
- **StartupScene** (scenes/startup.go) - File picker with native dialog
  - "Open File..." button (centered)
  - Ctrl+O keyboard shortcut
  - Error handling with dialogs
  - Displays selected file metadata
- **FileViewerScene** (scenes/file_viewer.go) - Raw file content display
  - Scrolling with mouse wheel and arrow keys
  - Visual scrollbar indicator
  - ESC to return to startup
  - Shows file type and Base_path
- **SceneManager** (scenes/scene.go) - Scene switching and state management

**State Management (internal/app/state.go):**
- Stores: BasePath, SelectedFilePath, FileType, FileContent
- Methods: LoadFile(), SetBasePath(), SelectFile()

**Dependencies:**
- `github.com/sqweek/dialog` - Native file picker dialogs (Windows/Linux/Mac)
- `github.com/hajimehoshi/ebiten/v2` - 2D game engine for GUI

## 🎯 Anti-patterns to Avoid

### Code Anti-patterns
❌ Don't hardcode file paths - always use configurable base paths  
❌ Don't ignore parsing errors - always validate and report issues  
❌ Don't modify original files without backups  
❌ Don't use global state - pass dependencies explicitly  

### UI Anti-patterns  
❌ Don't overwhelm users with all options at once  
❌ Don't allow invalid operations (like circular dependencies)  
❌ Don't lose user work - auto-save and recovery  
❌ Don't hide validation errors - make them visible and actionable  

### Architecture Anti-patterns
❌ Don't tightly couple parser and UI components  
❌ Don't skip validation for performance  
❌ Don't assume file formats won't change  
❌ Don't build without considering large file performance