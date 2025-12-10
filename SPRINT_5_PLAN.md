# Sprint 5: Technology Display Enhancement - Detailed Plan

## 📋 Overview

**Goal:** Implement advanced technology display features based on `HOI4_Technology_Advanced_Rules.md`

**Key Features:**
1. Country flags system for folder availability
2. Technology folder filtering with overlay logic
3. Sub-tree detection within folders
4. Scrollable UI for technology lists
5. Flag editor for managing country-specific folders

---

## 🎯 Task 21: Implement Country Flags System

### Goal
Load and manage country flags from `history/countries/*.txt` files to determine technology folder availability.

### Background
Country flags control which technology folders are available. Examples:
- `UNLOCK:infantry_folder` - unlocks infantry technologies
- `GER_air` - German-specific air tech tree
- `coastal_state` - enables naval technologies

### Implementation Steps

#### Step 1: Create CountryFlagsParser

**File:** `internal/parser/country_flags_parser.go`

```go
package parser

type CountryFlagsParser struct {
    gamePath string
    modPath  string
}

func NewCountryFlagsParser(gamePath, modPath string) *CountryFlagsParser

// ParseCountryFlags loads flags from history/countries/<TAG>.txt
func (p *CountryFlagsParser) ParseCountryFlags(countryTag string) ([]string, error)
```

**Algorithm:**
1. Try mod path first: `modPath/history/countries/<TAG> - <Name>.txt`
2. Fallback to game path: `gamePath/history/countries/<TAG> - <Name>.txt`
3. Parse file with existing Parser
4. Find all `set_country_flag = FLAG_NAME` statements
5. Return list of flag names

**Example file structure:**
```
capital = 219
set_country_flag = UNLOCK:infantry_folder
set_country_flag = UNLOCK:support_folder
set_country_flag = GER_air
set_technology = {
    infantry_weapons = 1
}
```

#### Step 2: Update CountryContext

**File:** `internal/app/country_context.go`

**Add fields:**
```go
type CountryContext struct {
    // ... existing fields
    CountryFlags []string // Loaded flags for this country
}
```

**Add methods:**
```go
// loadCountryFlags loads flags from history files
func (ctx *CountryContext) loadCountryFlags()

// HasFlag checks if country has a specific flag
func (ctx *CountryContext) HasFlag(flagName string) bool

// HasAnyFlag checks if country has any of the flags
func (ctx *CountryContext) HasAnyFlag(flags []string) bool
```

**Update NewCountryContext:**
```go
func NewCountryContext(...) *CountryContext {
    // ... existing code
    
    // Load country flags
    ctx.loadCountryFlags()
    
    return ctx
}
```

#### Step 3: Testing

**Test cases:**
1. Load flags for GER - should have `GER_air`, `GER_armor`
2. Load flags for SOV - should have `SOV_air`, `SOV_armor`
3. Load flags for minor country - should have basic unlock flags
4. Handle missing history file gracefully

**Files to test:**
- `history/countries/GER - Germany.txt`
- `history/countries/SOV - Soviet Union.txt`

---

## 🎯 Task 22: Implement Technology Folder Filtering

### Goal
Filter technology folders based on `available = { ... }` conditions and hide overlay folders when appropriate.

### Background

**Folder types:**
1. **Regular folders** - contain technologies (e.g., `infantry_folder`)
2. **Overlay folders** - block access until unlocked (e.g., `infantry_overlay_folder`)
3. **Country-specific folders** - only for certain countries (e.g., `luftwaffe_folder` for Germany)

**Overlay logic:**
- Overlay shown when `available` condition is FALSE
- Overlay hidden when condition is TRUE
- Usually: `NOT = { has_country_flag = UNLOCK:folder_name }`

### Implementation Steps

#### Step 1: Update TechnologyTagsParser

**File:** `internal/parser/technology_tags_parser.go`

**Extend TechFolder struct:**
```go
type TechFolder struct {
    Name       string
    Ledger     string // army, navy, air, civilian
    Available  *AvailableCondition // nil if always available
    IsOverlay  bool   // true if name ends with "_overlay_folder"
}

type AvailableCondition struct {
    Conditions []Condition
}

type Condition struct {
    Type     string // "has_country_flag", "NOT", "has_dlc", "major_country"
    Value    string
    Negated  bool
    Children []Condition // for nested conditions
}
```

**Update parsing:**
```go
func (p *TechnologyTagsParser) parseFolderBlock(name string, block *AssignmentStatement) *TechFolder {
    folder := &TechFolder{
        Name:      name,
        IsOverlay: strings.HasSuffix(name, "_overlay_folder"),
    }
    
    // Parse ledger, available conditions
    // ...
    
    return folder
}
```

#### Step 2: Create ConditionEvaluator

**File:** `internal/app/condition_evaluator.go`

```go
package app

type ConditionEvaluator struct {
    countryFlags []string
    dlcs         []string // for future
    isMajor      bool     // for future
}

func NewConditionEvaluator(flags []string) *ConditionEvaluator

// Evaluate checks if condition is met
func (e *ConditionEvaluator) Evaluate(condition *parser.AvailableCondition) bool

// evaluateCondition recursively evaluates a single condition
func (e *ConditionEvaluator) evaluateCondition(cond *parser.Condition) bool
```

**Logic:**
```go
func (e *ConditionEvaluator) evaluateCondition(cond *parser.Condition) bool {
    switch cond.Type {
    case "has_country_flag":
        hasFlag := e.hasFlag(cond.Value)
        if cond.Negated {
            return !hasFlag
        }
        return hasFlag
        
    case "NOT":
        // Evaluate all children, return true if ALL are false
        for _, child := range cond.Children {
            if e.evaluateCondition(&child) {
                return false // One is true, so NOT fails
            }
        }
        return true
        
    case "has_dlc":
        // For now, assume no DLCs
        return false
        
    default:
        return true // Unknown condition = allow
    }
}
```

#### Step 3: Update CountryContext Filtering

**File:** `internal/app/country_context.go`

```go
func (ctx *CountryContext) resolveTechFolders() {
    tagsParser := parser.NewTechnologyTagsParser(ctx.GamePath, ctx.ModPath)
    allFolders, err := tagsParser.ParseTechnologyFoldersDetailed()
    if err != nil {
        // ... error handling
        return
    }
    
    evaluator := NewConditionEvaluator(ctx.CountryFlags)
    
    // Filter folders
    availableFolders := make([]string, 0)
    overlayMap := make(map[string]bool) // base_folder -> overlay_active
    
    for _, folder := range allFolders {
        if folder.IsOverlay {
            // Check if overlay should be shown
            baseName := strings.TrimSuffix(folder.Name, "_overlay_folder")
            if folder.Available != nil {
                // Overlay shown when condition is FALSE
                overlayActive := !evaluator.Evaluate(folder.Available)
                overlayMap[baseName] = overlayActive
            }
        } else {
            // Regular folder - check if available
            if folder.Available == nil || evaluator.Evaluate(folder.Available) {
                availableFolders = append(availableFolders, folder.Name)
            }
        }
    }
    
    // Remove folders that have active overlays
    filtered := make([]string, 0)
    for _, folderName := range availableFolders {
        if !overlayMap[folderName] {
            filtered = append(filtered, folderName)
        }
    }
    
    ctx.TechFolders = filtered
}
```

#### Step 4: Testing

**Test scenarios:**
1. **Germany with GER_air flag:**
   - Should see `luftwaffe_folder`
   - Should NOT see generic `air_techs_folder`

2. **Minor country without unlock flags:**
   - Should see overlay folders blocking access
   - Should NOT see locked folders

3. **Country with UNLOCK:infantry_folder:**
   - Should see `infantry_folder`
   - Should NOT see `infantry_overlay_folder`

---

## 🎯 Task 23: Implement Sub-tree Detection

### Goal
Detect and display multiple technology sub-trees within a single folder (e.g., electronics_folder has 4 sub-trees).

### Background

**Example: electronics_folder**
- Electronic Engineering (X: 0-5)
- Experimental Rockets (X: 6-11)
- Jets & Aircraft Engines (X: 12-17)
- Atomic Research (X: 18-23)

**Detection algorithm:**
1. Group technologies by X-coordinate
2. Find gaps > 5 units
3. Identify sub-trees by categories

### Implementation Steps

#### Step 1: Update Domain Model

**File:** `internal/domain/technology.go`

```go
type SubTree struct {
    Name         string   // "Electronic Engineering"
    XMin         int      // 0
    XMax         int      // 5
    Categories   []string // ["electronics", "computing_tech"]
    Technologies []*Technology
}

type TechnologyTree struct {
    Technologies map[string]*Technology
    Folders      map[string][]string
    SubTrees     map[string][]*SubTree // folder_name -> sub-trees
}
```

#### Step 2: Create Sub-tree Detector

**File:** `internal/app/technology_loader.go`

```go
// DetectSubTrees groups technologies into sub-trees by X-coordinate
func (tl *TechnologyLoader) DetectSubTrees(folderName string, technologies []*domain.Technology) []*domain.SubTree {
    if len(technologies) == 0 {
        return nil
    }
    
    // Sort by X coordinate
    sort.Slice(technologies, func(i, j int) bool {
        return technologies[i].Position.X < technologies[j].Position.X
    })
    
    subTrees := make([]*domain.SubTree, 0)
    currentTree := &domain.SubTree{
        XMin:         technologies[0].Position.X,
        XMax:         technologies[0].Position.X,
        Technologies: make([]*domain.Technology, 0),
        Categories:   make(map[string]bool),
    }
    
    for _, tech := range technologies {
        // Check if this tech starts a new sub-tree (gap > 5)
        if tech.Position.X - currentTree.XMax > 5 {
            // Finalize current tree
            currentTree.Name = identifySubTreeName(currentTree.Categories)
            subTrees = append(subTrees, currentTree)
            
            // Start new tree
            currentTree = &domain.SubTree{
                XMin:         tech.Position.X,
                XMax:         tech.Position.X,
                Technologies: make([]*domain.Technology, 0),
                Categories:   make(map[string]bool),
            }
        }
        
        // Add tech to current tree
        currentTree.Technologies = append(currentTree.Technologies, tech)
        currentTree.XMax = max(currentTree.XMax, tech.Position.X)
        
        // Collect categories
        for _, cat := range tech.Categories {
            currentTree.Categories[cat] = true
        }
    }
    
    // Finalize last tree
    if len(currentTree.Technologies) > 0 {
        currentTree.Name = identifySubTreeName(currentTree.Categories)
        subTrees = append(subTrees, currentTree)
    }
    
    return subTrees
}

func identifySubTreeName(categories map[string]bool) string {
    categoryNames := map[string]string{
        "electronics":     "Electronic Engineering",
        "computing_tech":  "Electronic Engineering",
        "radar_tech":      "Electronic Engineering",
        "rocketry":        "Experimental Rockets",
        "mot_rockets":     "Experimental Rockets",
        "jet_technology":  "Jets & Aircraft Engines",
        "jet_engine":      "Jets & Aircraft Engines",
        "nuclear":         "Atomic Research",
    }
    
    for cat := range categories {
        if name, ok := categoryNames[cat]; ok {
            return name
        }
    }
    
    return "Technology Tree"
}
```

#### Step 3: Update TechViewerScene

**File:** `internal/ui/scenes/tech_viewer.go`

**Add sub-tree headers:**
```go
func (s *TechViewerScene) Draw(screen *ebiten.Image) {
    // ... existing code
    
    // Draw sub-tree headers if present
    if s.tree != nil && len(s.tree.SubTrees) > 1 {
        for _, subTree := range s.tree.SubTrees {
            headerX := (subTree.XMin + subTree.XMax) / 2
            headerY := -1 // Above the tree
            
            pos := s.canvas.GridToScreen(headerX, headerY)
            ebitenutil.DebugPrintAt(screen, subTree.Name, int(pos.X), int(pos.Y))
        }
    }
    
    // ... draw technologies
}
```

---

## 🎯 Task 24: Add Scrollable Technology List UI

### Goal
Create a scrollable list component for displaying technology folders with proper UX.

### Implementation Steps

#### Step 1: Create ScrollableList Component

**File:** `internal/ui/components/scrollable_list.go`

```go
package components

type ScrollableList struct {
    x, y          int
    width, height int
    items         []string
    displayCount  int // max visible items
    scrollOffset  int // current scroll position
    selectedIndex int
    
    itemHeight    int
    showScrollbar bool
}

func NewScrollableList(x, y, width, height int, displayCount int) *ScrollableList

func (sl *ScrollableList) SetItems(items []string)
func (sl *ScrollableList) Update()
func (sl *ScrollableList) Draw(screen *ebiten.Image)
func (sl *ScrollableList) GetSelectedItem() string
func (sl *ScrollableList) ScrollUp()
func (sl *ScrollableList) ScrollDown()
func (sl *ScrollableList) HandleMouseWheel(delta float64)
```

**Features:**
- Mouse wheel scrolling
- Click to select item
- Visual scrollbar indicator
- "... X more" text at bottom
- Hover highlighting

#### Step 2: Update CountryMenuScene

**File:** `internal/ui/scenes/country_menu.go`

**Replace button list with ScrollableList:**
```go
type CountryMenuScene struct {
    // ... existing fields
    techList *components.ScrollableList // Replace techCategoryButtons
}

func NewCountryMenuScene(manager *SceneManager, state *app.State) *CountryMenuScene {
    // ...
    scene.techList = components.NewScrollableList(440, 420, 400, 300, 6)
    scene.loadTechCategories()
    return scene
}

func (s *CountryMenuScene) loadTechCategories() {
    ctx := s.state.GetCountryContext()
    if ctx == nil {
        return
    }
    
    // Get localized names
    displayNames := make([]string, len(ctx.TechFolders))
    for i, folder := range ctx.TechFolders {
        displayNames[i] = ctx.GetLocalizedFolderName(folder)
    }
    
    s.techList.SetItems(displayNames)
}
```

---

## 🎯 Task 25: Implement Technology Folder Editor UI

### Goal
Create UI for managing country flags to unlock/lock technology folders.

### Implementation Steps

#### Step 1: Create FlagEditorScene

**File:** `internal/ui/scenes/flag_editor.go`

```go
type FlagEditorScene struct {
    manager      *SceneManager
    state        *app.State
    
    currentFlags []string
    allFlags     []string // All possible flags
    
    flagList     *components.ScrollableList
    addButton    *components.Button
    removeButton *components.Button
    saveButton   *components.Button
    backButton   *components.Button
}
```

**Features:**
- Display current flags
- Add new flags from predefined list
- Remove existing flags
- Save changes to history file
- Preview affected folders

#### Step 2: Update CountryMenuScene

Add "Edit Flags" button to open FlagEditorScene.

---

## 📊 Summary

### Priority Order
1. **Task 21** (Country Flags) - HIGH - Foundation for filtering
2. **Task 22** (Folder Filtering) - HIGH - Core functionality
3. **Task 24** (Scrollable UI) - HIGH - Usability improvement
4. **Task 23** (Sub-trees) - MEDIUM - UX enhancement
5. **Task 25** (Flag Editor) - MEDIUM - Future feature

### Estimated Time
- Task 21: 2-3 hours
- Task 22: 3-4 hours
- Task 23: 2-3 hours
- Task 24: 2-3 hours
- Task 25: 3-4 hours

**Total: 12-17 hours**

### Testing Checklist
- [ ] Flags loaded correctly for major countries
- [ ] Overlay folders hidden when unlocked
- [ ] Country-specific folders shown only for correct countries
- [ ] Sub-trees detected in electronics_folder
- [ ] Scrollable list works with mouse wheel
- [ ] UI responsive and smooth
