# Refactoring Plan: Country-Centric Architecture

## Overview
Переход от file-centric к country-centric архитектуре с правильной структурой навигации.

---

## Phase 1: Startup Scene Refactoring

### 1.1 Mod Selection Component
**Files to create/modify:**
- `internal/app/mod_validator.go` - валидация .mod файлов
- `internal/ui/scenes/startup.go` - обновление UI

**Tasks:**
- ✅ Парсинг .mod файлов (уже есть lexer/parser)
- ✅ Извлечение `path` переменной из .mod
- ✅ Валидация существования mod folder
- ✅ UI: File picker для .mod файлов
- ✅ UI: Отображение имени мода и версии

**Data structures:**
```go
type ModDescriptor struct {
    FilePath       string
    Name           string
    Version        string
    SupportedVersion string
    Path           string // Relative path to mod folder
    ReplacePaths   []string
    Tags           []string
}
```

### 1.2 Game Installation Component
**Files to create/modify:**
- `internal/app/game_validator.go` - валидация игры
- `internal/ui/scenes/startup.go` - UI для выбора игры

**Tasks:**
- ✅ Auto-detect функция (перенесена из icon_loader)
- ✅ Manual folder selection (валидация готова)
- ✅ Валидация: hoi4.exe + структура папок
- ✅ UI: Folder picker + "Auto-detect" кнопка
- ✅ Сохранение путей в конфиг

**Data structures:**
```go
type GameInstallation struct {
    Path      string
    Version   string
    IsValid   bool
    Executable string // hoi4.exe, hoi4.app, etc.
}
```

### 1.3 Configuration Persistence
**Files to create:**
- `internal/app/config.go` - конфигурация приложения

**Tasks:**
- ✅ Сохранение mod path
- ✅ Сохранение game path
- ✅ Сохранение последней выбранной страны
- ✅ JSON формат конфига
- ✅ Загрузка при старте

**Config structure:**
```go
type AppConfig struct {
    ModPath       string
    GamePath      string
    LastCountry   string
    WindowWidth   int
    WindowHeight  int
}
```

---

## Phase 2: Country Selection Scene

### 2.1 Bookmark Parser
**Files to create:**
- `internal/parser/bookmark_parser.go` - парсинг bookmarks
- `internal/domain/bookmark.go` - модели данных

**Tasks:**
- ⬜ Поиск файлов в `common/bookmarks/` (mod → game)
- ⬜ Парсинг bookmark структуры
- ⬜ Извлечение country tags и metadata
- ⬜ Группировка major/minor

**Data structures:**
```go
type Bookmark struct {
    Name           string
    Description    string
    Date           string
    DefaultCountry string
    Countries      []*BookmarkCountry
}

type BookmarkCountry struct {
    Tag      string // "GER", "SOV", "USA"
    Name     string // Localized or tag
    History  string // Description key
    Ideology string
    IsMajor  bool
    Ideas    []string
    Focuses  []string
}
```

### 2.2 Country Selection UI
**Files to create:**
- `internal/ui/scenes/country_selection.go` - новая сцена

**Tasks:**
- ⬜ Список стран с фильтрацией (major/minor/all)
- ⬜ Поиск по тегу/имени
- ⬜ Отображение флагов (если есть)
- ⬜ Отображение идеологии и базовой инфы
- ⬜ Переход к Country Scene

**UI Layout:**
```
┌─────────────────────────────────────┐
│ Select Country                      │
├─────────────────────────────────────┤
│ Filter: [All ▼] Search: [____]     │
├─────────────────────────────────────┤
│ ┌─────────────────────────────────┐ │
│ │ 🇩🇪 GER - Germany (Fascism)     │ │
│ │ 🇺🇸 USA - United States (Lib.) │ │
│ │ 🇷🇺 SOV - Soviet Union (Comm.) │ │
│ │ ...                             │ │
│ └─────────────────────────────────┘ │
│                                     │
│ [Back to Mod Selection]             │
└─────────────────────────────────────┘
```

---

## Phase 3: Country Scene (Main Menu)

### 3.1 Country Context
**Files to create:**
- `internal/app/country_context.go` - контекст страны

**Tasks:**
- ⬜ Хранение выбранной страны
- ⬜ Определение доступных tech folders
- ⬜ Определение focus tree файла
- ⬜ Загрузка country-specific данных

**Data structures:**
```go
type CountryContext struct {
    Tag              string
    Name             string
    ModPath          string
    GamePath         string
    
    // Resolved paths
    FocusTreeFile    string
    TechFolders      map[string]string // folder_name → file_path
    
    // Metadata
    Ideology         string
    IsMajor          bool
}
```

### 3.2 Country Scene UI
**Files to create:**
- `internal/ui/scenes/country_menu.go` - главное меню страны

**Tasks:**
- ⬜ Отображение информации о стране
- ⬜ Кнопка "National Focus Tree"
- ⬜ Кнопка "Technologies" (с подменю категорий?)
- ⬜ Кнопка "Back to Country Selection"
- ⬜ Breadcrumb навигация

**UI Layout:**
```
┌─────────────────────────────────────┐
│ Germany (GER) - Fascism             │
├─────────────────────────────────────┤
│                                     │
│   [📋 National Focus Tree]          │
│                                     │
│   [🔬 Technologies]                 │
│      ├─ Infantry                    │
│      ├─ Air Force (Luftwaffe)       │
│      ├─ Armor                       │
│      ├─ Naval                       │
│      └─ Industry                    │
│                                     │
│   [🏭 Production]                   │
│   [🗺️  Map Editor]                  │
│                                     │
│ [← Back to Country Selection]      │
└─────────────────────────────────────┘
```

---

## Phase 4: Technology Folder Resolution

### 4.1 Technology Tags Parser
**Files to create:**
- `internal/parser/tech_tags_parser.go` - парсинг technology_tags
- `internal/domain/tech_folder.go` - модели

**Tasks:**
- ⬜ Парсинг `technology_folders` блока
- ⬜ Извлечение `available` условий
- ⬜ Резолвинг folder → file mapping
- ⬜ Country-specific folder detection

**Data structures:**
```go
type TechFolder struct {
    Name      string // "luftwaffe_folder"
    Ledger    string // "air", "army", "navy", "civilian"
    Available string // Condition script
    FilePath  string // Resolved: "GER_air.txt"
}

type TechFolderResolver struct {
    CountryTag string
    Folders    []*TechFolder
}
```

### 4.2 Folder → File Mapping
**Files to create:**
- `internal/app/tech_resolver.go` - резолвер технологий

**Tasks:**
- ⬜ Создать mapping table (folder_name → file_name)
- ⬜ Логика резолвинга для country-specific
- ⬜ Fallback на generic файлы

**Mapping examples:**
```go
var folderToFileMap = map[string]string{
    "luftwaffe_folder":        "GER_air.txt",
    "sovietair_folder":        "SOV_air.txt",
    "usair_folder":            "USA_air.txt",
    "trm_armour_ger_folder":   "GER_armor.txt",
    "trm_armour_sov_folder":   "SOV_armor.txt",
    "infantry_folder":         "infantry.txt",
    "support_folder":          "support.txt",
    "industry_folder":         "industry.txt",
    // ... etc
}
```

---

## Phase 5: Integration with Existing Viewers

### 5.1 TechViewerScene Update
**Files to modify:**
- `internal/ui/scenes/tech_viewer.go`

**Tasks:**
- ⬜ Принимать CountryContext вместо filePath
- ⬜ Использовать resolved tech file path
- ⬜ Breadcrumb: Country → Technologies → [Category]
- ⬜ Кнопка "Back to Country Menu"

### 5.2 FocusViewerScene Creation
**Files to create:**
- `internal/ui/scenes/focus_viewer.go`

**Tasks:**
- ⬜ Дублировать структуру TechViewerScene
- ⬜ Использовать FocusParser
- ⬜ Relative positioning для фокусов
- ⬜ Breadcrumb навигация

---

## Phase 6: Scene Navigation System

### 6.1 Scene Flow
```
StartupScene
    ↓ (select mod + game)
CountrySelectionScene
    ↓ (select country)
CountryMenuScene
    ├─→ FocusViewerScene
    └─→ TechViewerScene (with category selection)
```

### 6.2 Navigation Stack
**Files to create:**
- `internal/ui/scenes/navigation.go` - навигационный стек

**Tasks:**
- ⬜ History stack для "Back" кнопок
- ⬜ Breadcrumb rendering
- ⬜ Context passing между сценами

---

## Implementation Order

### Sprint 1: Foundation (Week 1)
1. ✅ Создать _DATA_STRUCTURE.md
2. ✅ Создать _REFACTORING_PLAN.md
3. ✅ Создать ModDescriptor parser
4. ✅ Создать GameInstallation validator
5. ✅ Обновить StartupScene UI
6. ✅ Создать AppConfig persistence

### Sprint 2: Country Selection (Week 2)
7. ✅ Создать BookmarkParser
8. ✅ Создать CountrySelectionScene
9. ✅ Интеграция mod/game path resolution
10. ✅ UI для списка стран

### Sprint 3: Country Context (Week 3)
11. ⬜ Создать CountryContext
12. ⬜ Создать CountryMenuScene
13. ⬜ Создать TechFolderResolver
14. ⬜ Парсинг technology_tags

### Sprint 4: Integration (Week 4)
15. ⬜ Обновить TechViewerScene
16. ⬜ Создать FocusViewerScene
17. ⬜ Навигационная система
18. ⬜ Тестирование end-to-end

---

## Migration Strategy

### Backward Compatibility
- Старый file-based подход удалить после полной реализации
- Сохранить существующие parsers (TechParser, FocusParser)
- Обновить только UI flow и data loading

### Testing Checklist
- [ ] Выбор мода и валидация
- [ ] Auto-detect игры
- [ ] Загрузка bookmarks из мода
- [ ] Fallback на vanilla bookmarks
- [ ] Резолвинг country-specific tech folders
- [ ] Загрузка технологий для разных стран
- [ ] Загрузка фокусов для разных стран
- [ ] Навигация между сценами
- [ ] Сохранение/загрузка конфига

---

## Open Questions

1. **Localization:** Использовать ключи локализации или показывать теги?
2. **DLC Detection:** Проверять наличие DLC для conditional content?
3. **Multi-mod Support:** Поддержка нескольких модов одновременно?
4. **Country Flags:** Загружать флаги из `gfx/flags/`?
5. **Technology Categories UI:** Flat list или tree structure?

---

## Success Criteria

✅ Пользователь может:
1. Выбрать .mod файл и игру
2. Увидеть список стран из bookmarks
3. Выбрать страну и увидеть меню
4. Открыть дерево фокусов для страны
5. Открыть технологии (с правильными country-specific файлами)
6. Вернуться назад на любом этапе
7. Конфиг сохраняется между сессиями
