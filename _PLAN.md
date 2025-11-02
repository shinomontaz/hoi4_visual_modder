# HOI4 Visual Modder - Project plan

## 🎯 Current Focus

**Current Phase:** Phase 1 - MVP (Basic visualization and parsing)

**What exists now:**
- ✅ Complete project structure with all directories
- ✅ Domain models fully implemented (Focus, Technology, Tree, Position)
- ✅ Basic Ebitengine window with scene manager
- ✅ Lexer and Parser placeholder files created
- ✅ Application compiles and runs (bin/modder.exe)
- 📄 Test data available: test_tech.txt in project root

**What's missing:**
- ❌ File browser UI for mod selection
- ❌ File content display (show selected .txt file)
- ❌ Lexer implementation (tokenization logic)
- ❌ Parser implementation (AST building)
- ❌ Canvas rendering (grid, nodes, connections)

### Next Steps
1. ✅ ~~Create initial project structure (directories, go.mod, main.go)~~
2. ✅ ~~Implement domain models (Focus, Technology, Position, Tree)~~
3. **Implement file browser UI** ← NEXT
   - Directory selection dialog (Base_path)
   - Scan for .txt files in common/national_focus/ and common/technologies/
   - Display file list in startup scene
4. **File selection and display**
   - Click on file to select
   - Load file content (UTF-8)
   - Display raw text content on screen
5. **Build Paradox script lexer** (tokenization)
6. **Implement parser** for focus and technology files
7. **Test parser** with selected files
8. **Implement canvas rendering** (grid, nodes, connections)

### Active Work
- [x] Project structure setup ✅
- [x] Domain layer implementation ✅
- [x] File browser UI (basic) ✅
  - ✅ File scanning (common/national_focus/*.txt, common/technologies/*.txt)
  - ✅ File list display with clickable items and hover effect
  - ✅ File content loading and display
  - ✅ File viewer scene with scrolling
- [x] **Native File Picker Integration** ✅
  - ✅ Add file dialog library (github.com/sqweek/dialog)
  - ✅ Create file picker UI in StartupScene
  - ✅ Implement "Open File" button
  - ✅ File type filters (.txt files only)
  - ✅ Auto-detect Base_path from selected file
  - ✅ Validate mod structure
  - ✅ Store Base_path in State
  - ✅ Update UI to show selected file and Base_path
  - ✅ Fixed Base_path detection bug (duplicate drive letter)
- [ ] **Parser implementation** ← **NEXT**
  - [ ] Implement Lexer (tokenization)
  - [ ] Implement Parser (AST building)
  - [ ] Test with real mod files

---

## 🗺️ Development Roadmap

### Phase 1: MVP (Minimum Viable Product) ⬅️ **Current Phase**
**Goal:** Basic visualization and parsing - read-only viewer

**Features:**
- ✅ Architecture design
- ✅ Technology stack selection
- ✅ **Project Structure**
  - ✅ Created directory structure (cmd/, internal/, pkg/, assets/, test_data/)
  - ✅ Initialized go.mod with Ebitengine v2.9.3
  - ✅ Setup main.go entry point
  - ✅ Created .gitignore
  - ✅ Binary builds successfully (bin/modder.exe)
- ✅ **Domain Models**
  - ✅ Focus struct with all properties (ID, Icon, Position, Prerequisites, etc.)
  - ✅ Technology struct with all properties (ID, Effects, Paths, etc.)
  - ✅ FocusTree and TechnologyTree structures
  - ✅ Position and grid system
  - ✅ Validation methods (circular dependencies, prerequisites, position conflicts)
- ✅ **Basic GUI (Ebitengine)** - Completed
  - ✅ Window setup and game loop (main.go)
  - ✅ Scene manager (scene switching)
  - ✅ Startup scene with native file picker
  - ✅ Native file picker dialog with .txt filter
  - ✅ File scanner (scan common/national_focus/ and common/technologies/)
  - ✅ Button component (reusable UI element)
  - ✅ File selection handling (mouse click detection)
  - ✅ File content loading (read UTF-8 text file)
  - ✅ File viewer scene (show raw file content with scrolling)
  - ✅ ModLoader (Base_path detection and validation)
- [ ] **Paradox Script Parser**
  - Lexer: tokenize Paradox scripting language
  - Parser: build AST from tokens
  - Focus parser: parse `focus_tree` and `focus` blocks
  - Technology parser: parse `technologies` and tech blocks
  - Error handling and reporting
- [ ] **Visual Editor Canvas**
  - Grid rendering with coordinates
  - Render nodes as white squares with ID text
  - Render connection lines (prerequisites/paths)
  - Zoom functionality (mouse wheel)
  - Pan functionality (drag canvas)
  - Camera system for viewport management
- [ ] **Read-Only Mode**
  - Load and parse existing .txt files
  - Display tree on canvas
  - Navigate and explore the tree
  - View node positions and connections

**Deliverable:** Application that can load and visualize existing focus trees and technology trees

---

### Phase 2: Extended Features
**Goal:** Icon integration and basic editing

**Features:**
- [ ] **GFX Integration**
  - Parse .gfx files (goals.gfx, countrytechtreeview.gfx)
  - Load .dds icon files
  - Display actual icons instead of white squares
  - Icon caching system
- [ ] **Position Editing**
  - Drag & drop nodes to new positions
  - Snap to grid functionality
  - Real-time position updates
  - Visual feedback during drag
- [ ] **Properties Panel**
  - Display detailed node information
  - Show all focus/tech properties
  - View completion rewards / effects
  - View prerequisites and conditions
  - Read-only property display
- [ ] **Validation Feedback**
  - Highlight circular dependencies
  - Show position conflicts
  - Display invalid references
  - Warning indicators on nodes

**Deliverable:** Application with icon display and drag-drop editing of positions

---

### Phase 3: Advanced Editing
**Goal:** Full editing capabilities with file generation

**Features:**
- [ ] **Property Editing**
  - Edit focus/tech properties in panel
  - Modify completion rewards
  - Edit availability conditions
  - Change costs and research values
- [ ] **Connection Editing**
  - Visual creation of prerequisites
  - Visual creation of technology paths
  - Mutual exclusivity setup
  - Delete connections
- [ ] **Icon Management**
  - Upload new icon images
  - Automatic .dds conversion
  - Auto-generate GFX sprite entries
  - Icon preview and selection
- [ ] **File Operations**
  - Save changes back to .txt files
  - Update .gfx files automatically
  - Create .bak backups
  - Atomic file writes
  - Export to new files
- [ ] **Advanced Features**
  - Undo/redo support
  - Copy/paste nodes
  - Duplicate branches
  - Auto-layout algorithms
  - Search and filter nodes

**Deliverable:** Full-featured editor with complete read/write capabilities

---

## 📐 Architecture & Technical Decisions

### [2025-01-02] - Application Architecture Design

**Decision:** Layered architecture with clear separation of concerns using Go + Ebitengine

**Tech Stack:**
- **Language:** Go 1.21+
- **GUI Framework:** Ebitengine v2 (2D game engine for cross-platform GUI)
- **No external dependencies** for parsing (custom Paradox script parser)

**Architecture Layers:**

```
┌─────────────────────────────────────────────────────────────┐
│                     PRESENTATION LAYER                       │
│  (Ebitengine-based GUI, User Interactions, Rendering)       │
├─────────────────────────────────────────────────────────────┤
│                     APPLICATION LAYER                        │
│  (Business Logic, Validation, State Management)             │
├─────────────────────────────────────────────────────────────┤
│                        DOMAIN LAYER                          │
│  (Core Models: Focus, Technology, Tree structures)          │
├─────────────────────────────────────────────────────────────┤
│                    INFRASTRUCTURE LAYER                      │
│  (File I/O, Parser, Serializer, GFX Integration)           │
└─────────────────────────────────────────────────────────────┘
```

**Project Structure:**

```
hoi4_visual_modder/
├── cmd/
│   └── hoi4modder/
│       └── main.go                 # Application entry point
│
├── internal/
│   ├── domain/                     # Domain models (core entities)
│   │   ├── focus.go               # Focus structure & methods
│   │   ├── technology.go          # Technology structure & methods
│   │   ├── tree.go                # Tree/Graph structures
│   │   ├── position.go            # Position & grid system
│   │   └── validation.go          # Domain validation rules
│   │
│   ├── parser/                     # Paradox script parsing
│   │   ├── lexer.go               # Tokenization
│   │   ├── parser.go              # AST building
│   │   ├── focus_parser.go        # Focus-specific parsing
│   │   ├── tech_parser.go         # Technology-specific parsing
│   │   └── gfx_parser.go          # GFX file parsing
│   │
│   ├── serializer/                 # File generation
│   │   ├── focus_writer.go        # Focus tree serialization
│   │   ├── tech_writer.go         # Technology serialization
│   │   └── gfx_writer.go          # GFX file generation
│   │
│   ├── app/                        # Application logic
│   │   ├── project.go             # Project/mod management
│   │   ├── validator.go           # Cross-entity validation
│   │   ├── state.go               # Application state
│   │   └── commands.go            # User actions/commands
│   │
│   └── ui/                         # Presentation layer
│       ├── game.go                # Main Ebitengine game struct
│       ├── scenes/                # Different UI screens
│       │   ├── scene.go           # Scene interface
│       │   ├── startup.go         # Mod selection screen
│       │   ├── focus_editor.go    # Focus tree editor
│       │   └── tech_editor.go     # Technology editor
│       ├── components/            # Reusable UI components
│       │   ├── canvas.go          # Main editing canvas
│       │   ├── node.go            # Visual node (focus/tech)
│       │   ├── connection.go      # Visual connections/arrows
│       │   ├── properties.go      # Properties panel
│       │   ├── toolbar.go         # Toolbar component
│       │   └── dialog.go          # Modal dialogs
│       ├── input/                 # Input handling
│       │   ├── mouse.go           # Mouse interactions
│       │   ├── keyboard.go        # Keyboard shortcuts
│       │   └── dragdrop.go        # Drag & drop logic
│       └── render/                # Rendering utilities
│           ├── grid.go            # Grid rendering
│           ├── icons.go           # Icon loading/caching
│           └── text.go            # Text rendering helpers
│
├── pkg/                            # Public reusable packages
│   └── paradox/                   # Paradox script utilities
│       ├── types.go               # Common Paradox types
│       └── utils.go               # Helper functions
│
├── assets/                         # Embedded assets
│   ├── fonts/                     # UI fonts
│   └── icons/                     # Built-in icons
│
├── test_data/                      # Test files
│   ├── focus_trees/
│   └── technologies/
│
├── go.mod
├── go.sum
├── README.md
├── _CONTEXT.md
├── _PLAN.md
└── .windsurfrules
```

**Key Design Decisions:**

1. **Layered Architecture:**
   - **Domain Layer:** Pure Go structs, no dependencies, business rules
   - **Infrastructure:** File I/O, parsing, external integrations
   - **Application:** Orchestration, validation, state management
   - **Presentation:** Ebitengine UI, user interactions

2. **Parser Design:**
   - Custom lexer/parser for Paradox scripting language
   - Two-phase: Lexer (tokens) → Parser (AST)
   - Separate parsers for focus/tech/gfx files
   - Preserve comments and formatting for round-trip editing

3. **UI Architecture (Ebitengine):**
   - **Scene-based:** Different screens (startup, focus editor, tech editor)
   - **Component-based:** Reusable UI components (canvas, nodes, panels)
   - **Event-driven:** Input handlers → Commands → State updates → Re-render
   - **Canvas:** Infinite scrollable grid with zoom

4. **State Management:**
   - Single source of truth in `app.State`
   - Immutable operations where possible
   - Command pattern for undo/redo support (future)

5. **Validation Strategy:**
   - **Domain validation:** In model methods (e.g., `Focus.Validate()`)
   - **Cross-entity validation:** In `app.Validator`
   - **Real-time feedback:** Validate on every change

6. **File Operations:**
   - Always create `.bak` before overwriting
   - UTF-8 encoding enforced
   - Atomic writes (write to temp → rename)
   - Error recovery with detailed messages

**Data Flow Example (User edits focus position):**

```
User drags node
    ↓
ui/input/dragdrop.go detects drag
    ↓
app/commands.go: MoveFocusCommand
    ↓
app/state.go: Update focus position
    ↓
app/validator.go: Validate new position
    ↓
ui/scenes/focus_editor.go: Re-render canvas
```

**Rationale:**
- **Go:** Type safety, excellent tooling, cross-platform, fast compilation
- **Ebitengine:** Lightweight, pure Go, cross-platform, good for 2D grid-based UIs
- **Layered architecture:** Clear separation enables testing, maintainability
- **Custom parser:** No existing Paradox parser in Go, full control over formatting
- **Scene-based UI:** Clean separation between different editing modes

**Alternatives Considered:**
- ❌ **Electron + React:** Too heavy, not native, large bundle size
- ❌ **Qt/GTK bindings:** Complex C dependencies, harder cross-platform builds
- ❌ **Fyne:** Less control over custom rendering for complex canvas
- ❌ **Using existing parser libraries:** Paradox format too specific

**Impact:**
- Clear structure for development
- Easy to test each layer independently
- Scalable for future features (undo/redo, multi-file editing)
- Cross-platform without platform-specific code

**Status:** Active

---

## 📚 Paradox Scripting Language Reference

### Syntax Examples

**Focus Tree Structure:**
```paradox
focus_tree = {
    id = brazil_focus
    country = {
        factor = 0
        modifier = {
            add = 10
            tag = BRA
        }
    }
    
    focus = {
        id = my_focus
        icon = GFX_focus_icon
        x = 5
        y = 3
        cost = 70
        
        prerequisite = { focus = parent_focus }
        mutually_exclusive = { focus = other_focus }
        
        completion_reward = {
            add_political_power = 50
        }
    }
}
```

**Technology Structure:**
```paradox
technologies = {
    tech_id = {
        allow = {
            has_country_flag = UNLOCK:folder_name
        }
        
        category_all_infantry = {
            soft_attack = 0.05
            defense = 0.03
        }
        
        path = {
            leads_to_tech = next_tech
            research_cost_coeff = 1.0
        }
        
        folder = {
            name = infantry_folder
            position = { x = 2 y = 5 }
        }
        
        research_cost = 1.5
        categories = { land_doctrine }
    }
}
```

**GFX File Structure:**
```paradox
spriteTypes = {
    spriteType = {
        name = GFX_focus_my_icon
        texturefile = "gfx/interface/goals/my_icon.dds"
    }
}
```

### Key Parsing Challenges
1. **Nested Blocks:** Arbitrary depth of `{ }` nesting
2. **Implicit Arrays:** Multiple blocks with same key = array
3. **Mixed Types:** Values can be strings, numbers, dates, or blocks
4. **Comments:** Lines starting with `#` should be preserved
5. **Whitespace:** Formatting should be preserved for round-trip editing
6. **Keywords vs Identifiers:** Context-dependent (e.g., `focus` is both keyword and value)

### Available Documentation Files
- **hoi4_focus_tree_documentation.md** - Complete focus tree structure reference (616 lines)
- **hoi4_tech_structure_documentation.md** - Technology file structure reference (469 lines)
- **hoi4_images_rules.md** - Icon and GFX file integration rules (73 lines)
- **test_tech.txt** - Real example technology file for testing

---

## 💻 Existing Code Reference

### Implemented Files (Ready to Use)

**Domain Layer (internal/domain/):**
- `position.go` - Position struct with X, Y coordinates and helper methods (Equals, Add)
- `focus.go` - Focus struct with 20+ properties, validation, prerequisite checking
- `technology.go` - Technology struct with effects, paths, validation
- `tree.go` - FocusTree and TechnologyTree with validation (circular deps, prerequisites, position conflicts)

**Application Layer (internal/app/):**
- `state.go` - Application state management (ModPath, CurrentMode, loaded trees, camera, zoom)

**UI Layer (internal/ui/):**
- `game.go` - Main Ebitengine Game struct with Update/Draw/Layout
- `scenes/scene.go` - Scene interface and SceneManager for scene switching
- `scenes/startup.go` - StartupScene placeholder (shows debug text)

**Infrastructure Layer (internal/parser/):**
- `lexer.go` - Lexer struct with Token types defined, methods stubbed (NextToken needs implementation)
- `parser.go` - Parser struct with basic structure, Parse() method stubbed

**Utilities (pkg/paradox/):**
- `types.go` - Block struct, IsKeyword() function with common Paradox keywords

**Entry Point:**
- `cmd/hoi4modder/main.go` - Application entry point, creates Game, runs Ebitengine

### What Needs Implementation

**Priority 1 - Native File Picker Integration:** ← **CURRENT PRIORITY**

> 📄 **Detailed guide:** [FILE_PICKER_PLAN.md](FILE_PICKER_PLAN.md) - complete implementation plan with code examples

**Step 1: Add File Dialog Library**
- Add dependency: `github.com/sqweek/dialog` (native file dialogs for Windows/Linux/Mac)
- Run: `go get github.com/sqweek/dialog`
- Update go.mod and go.sum

**Step 2: Create ModLoader (internal/app/mod_loader.go)**
- `DetectBasePath(filePath string) (string, error)` - extract Base_path from file path
  - Example: `E:/mods/my_mod/common/national_focus/brazil.txt` → `E:/mods/my_mod`
  - Look for `common/` directory in parent paths
- `ValidateModStructure(basePath string) error` - check if valid HOI4 mod
  - Verify `common/` directory exists
  - Check for `national_focus/` or `technologies/` subdirectories
- `DetectFileType(filePath string) (FileType, error)` - determine if focus or tech file
  - Check path contains `national_focus` or `technologies`

**Step 3: Update State (internal/app/state.go)**
- Add field: `BasePath string` - root directory of the mod
- Add field: `SelectedFilePath string` - full path to selected file
- Add field: `FileType FileType` - enum: Focus or Technology
- Add method: `SetBasePath(path string)` - store Base_path
- Add method: `LoadFile(filePath string) error` - load and validate file

**Step 4: Update StartupScene (internal/ui/scenes/startup.go)**
- Remove 'O' key handler
- Add "Open File..." button (centered on screen)
- On button click:
  1. Call `dialog.File().Filter("Text files", "txt").Load()`
  2. Get selected file path
  3. Call `ModLoader.DetectBasePath(filePath)`
  4. Call `ModLoader.ValidateModStructure(basePath)`
  5. Store in `state.SetBasePath(basePath)`
  6. Load file content
  7. Switch to FileViewer scene
- Display selected file info:
  - File name
  - Base_path
  - File type (Focus/Technology)

**Step 5: Error Handling**
- Show error dialog if:
  - File selection cancelled
  - Invalid mod structure
  - File read error
  - Not a .txt file
- Use `dialog.Message().Error()` for error messages

**Step 6: UI Improvements**
- Add visual button component (internal/ui/components/button.go)
- Button states: normal, hover, pressed
- Keyboard shortcut: Ctrl+O to open file picker

**Priority 2 - Lexer (internal/parser/lexer.go):**
- Implement `NextToken()` method to tokenize Paradox scripts
- Handle: identifiers, strings (quoted), numbers, dates (1939.1.1), comments (#)
- Recognize delimiters: `{`, `}`, `=`, `<`, `>`
- Skip whitespace while preserving it for formatting

**Priority 3 - Parser (internal/parser/):**
- Implement `Parse()` in parser.go to build AST from tokens
- Create `focus_parser.go` to parse focus_tree blocks into domain.FocusTree
- Create `tech_parser.go` to parse technologies blocks into domain.TechnologyTree
- Handle nested blocks, implicit arrays (multiple same keys)

**Priority 4 - Canvas Rendering:**
- Canvas component for grid rendering
- Node rendering (white squares with text)
- Connection lines between nodes

---

## 📝 Development Log

### [2025-01-02] - Architecture Design & Roadmap Planning

**Work Done:**
- Analyzed project requirements from documentation files
- Designed layered architecture with clear separation of concerns
- Created detailed project structure with all packages and files
- Defined data flow and component interactions
- Documented key design decisions and rationale
- Moved development phases from _CONTEXT.md to _PLAN.md
- Structured 3-phase roadmap with clear deliverables
- Expanded Current Focus with detailed subtasks

**Discoveries:**
- Paradox scripting language requires custom parser (no existing Go libraries)
- Ebitengine's game loop model fits well for interactive canvas editing
- Scene-based architecture provides clean separation between editing modes
- Three clear phases: MVP (visualization) → Extended (icons + editing) → Advanced (full features)

**Technical Decisions:**
- Layered architecture: Domain → Infrastructure → Application → Presentation
- Custom lexer/parser for Paradox scripts with formatting preservation
- Component-based UI with reusable elements (canvas, nodes, panels)
- Command pattern foundation for future undo/redo support
- Phase 1 focus: Read-only visualization with basic parser

**Changes:**
- Moved Phase 1/2/3 from _CONTEXT.md to _PLAN.md as Development Roadmap
- Updated _CONTEXT.md to reference _PLAN.md for roadmap details
- Expanded Active Tasks with detailed subtasks for Phase 1
- Structured roadmap with goals, features, and deliverables per phase

**Next Steps:**
- ~~Create initial project structure~~ ✅ Done
- ~~Implement domain models~~ ✅ Done
- Build lexer for Paradox scripting language
- Implement parser for focus/tech files

---

### [2025-01-02] - Project Structure & Domain Models Implementation

**Work Done:**
- ✅ Created complete project directory structure
  - cmd/hoi4modder/ - application entry point
  - internal/domain/ - core models (Focus, Technology, Position, Tree)
  - internal/parser/ - lexer and parser placeholders
  - internal/serializer/ - file writers placeholders
  - internal/app/ - application state management
  - internal/ui/ - Ebitengine UI with scene manager
  - pkg/paradox/ - Paradox script utilities
  - assets/ - fonts and icons directories
  - test_data/ - test files structure
- ✅ Initialized Go module with Ebitengine v2.9.3
- ✅ Implemented domain models:
  - Position: X-Y coordinates with helper methods
  - Focus: Complete national focus structure with validation
  - Technology: Complete tech structure with effects and paths
  - FocusTree: Tree management with circular dependency detection
  - TechnologyTree: Tech collection with path validation
- ✅ Created basic Ebitengine application:
  - Game struct with Update/Draw/Layout
  - SceneManager for scene switching
  - StartupScene placeholder
- ✅ Application builds and runs successfully (bin/modder.exe)
- ✅ Created .gitignore for Go projects

**Discoveries:**
- User already has test data (test_tech.txt in project root)
- Project compiles cleanly with Go 1.24.9
- Ebitengine window opens successfully with placeholder scene
- Domain validation includes circular dependencies, prerequisites, and position conflicts

**Technical Implementation:**
- Focus struct: 20+ properties including prerequisites (OR/AND logic), mutual exclusivity
- Technology struct: Effects as nested maps, paths with cost coefficients
- Tree validation: Recursive circular dependency detection, prerequisite existence checks
- Scene-based UI: Clean separation with OnEnter/OnExit lifecycle

**Next Steps:**
- ~~Implement Paradox script lexer (tokenization)~~ → Moved to Priority 2
- ~~Build parser for focus_tree and technologies blocks~~ → Moved to Priority 3
- ~~Test parser with test_tech.txt file~~ → After parser implementation
- Implement file browser UI for mod selection → **Changed to Priority 1**

---

### [2025-01-02] - Development Plan Adjustment

**Decision:** Reordered implementation priorities - GUI first, then parsing

**Rationale:**
- More logical workflow: user needs to select files before we can parse them
- Immediate visual feedback: user can see file list and content right away
- Incremental development: can test file loading without parser
- Better UX: user can explore mod structure before parsing is ready

**Changes to Plan:**
- **Priority 1:** File browser UI (directory picker, file scanner, file list, file viewer)
- **Priority 2:** Lexer implementation (moved from Priority 1)
- **Priority 3:** Parser implementation (moved from Priority 2)
- **Priority 4:** Canvas rendering (moved from Priority 3)

**New Implementation Order:**
1. Directory picker dialog in startup scene
2. File scanner to find .txt files in mod directories
3. File list UI component with mouse interaction
4. File viewer scene to display raw file content
5. Then proceed with lexer and parser

**Benefits:**
- User can immediately work with the application
- Can test file I/O independently from parsing
- Easier debugging (see raw file content before parsing)
- More natural development flow

---

### [2025-01-02] - File Browser Implementation

**Work Done:**
- ✅ Created `FileScanner` (internal/app/file_scanner.go)
  - Scans `common/national_focus/` and `common/technologies/`
  - Returns FileInfo with metadata (path, name, size, category)
  - Validates mod directory structure
- ✅ Extended State (internal/app/state.go)
  - Added: AvailableFiles, SelectedFile, FileContent fields
  - Methods: SetAvailableFiles(), SelectFile(), SetFileContent()
- ✅ Updated StartupScene (internal/ui/scenes/startup.go)
  - Press 'O' to scan current directory
  - Display scrollable file list
  - Mouse hover highlighting
  - Click to select file
  - Mouse wheel scrolling
- ✅ Created FileViewerScene (internal/ui/scenes/file_viewer.go)
  - Display raw file content
  - Scrolling with mouse wheel and arrow keys
  - Visual scrollbar indicator
  - ESC to return to file list
- ✅ Updated SceneManager to pass State to all scenes
- ✅ Created test structure: common/national_focus/, common/technologies/
- ✅ Application compiles and runs successfully

**Discoveries:**
- Basic file browser works but UX not ideal (keyboard shortcut not intuitive)
- Need proper file picker for real-world usage
- Must auto-detect Base_path from selected file for mod resource loading

**Next Steps:**
- Replace 'O' key with native file picker dialog
- Implement Base_path auto-detection
- Add proper error handling and validation

---

### [2025-01-02] - Native File Picker Planning

**Decision:** Implement native file picker with Base_path auto-detection

**Rationale:**
- **Better UX:** Native OS file dialog is familiar to users
- **Real-world ready:** Can work with any mod location, not just test directory
- **Base_path critical:** Need to know mod root for loading GFX and icon files later
- **Proper validation:** Can validate mod structure before loading

**Architecture:**
- **Library:** `github.com/sqweek/dialog` - cross-platform native dialogs
- **ModLoader:** New component to detect Base_path from file path
- **Workflow:** File selection → Base_path detection → Validation → Load

**Implementation Plan:**
1. Add dialog library dependency
2. Create ModLoader with path detection logic
3. Update State with BasePath and FileType fields
4. Replace keyboard shortcut with button UI
5. Add error handling with dialogs
6. Create reusable Button component

**Expected Benefits:**
- Professional file selection experience
- Works with mods in any location
- Automatic mod structure validation
- Foundation for Phase 2 (icon loading)

**📄 Detailed Implementation Guide:** See [FILE_PICKER_PLAN.md](FILE_PICKER_PLAN.md) for:
- Step-by-step implementation with code examples
- User workflow diagram
- Testing checklist (12 items)
- Files to create/modify
- Time estimates (~60 minutes)

---

### [2025-01-02] - Native File Picker Implementation

**Work Done:**
- ✅ Added `github.com/sqweek/dialog` library for native file dialogs
- ✅ Created `ModLoader` (internal/app/mod_loader.go)
  - `DetectBasePath()` - extracts mod root from file path
  - `ValidateModStructure()` - validates HOI4 mod directory structure
  - `DetectFileType()` - determines if file is Focus or Technology
  - `LoadModFile()` - complete file loading with all validations
- ✅ Extended State (internal/app/state.go)
  - Added: BasePath, SelectedFilePath, FileType fields
  - Methods: SetBasePath(), LoadFile()
- ✅ Created Button component (internal/ui/components/button.go)
  - Visual states: normal, hover, pressed
  - Mouse interaction handling
  - Reusable UI component
- ✅ Redesigned StartupScene (internal/ui/scenes/startup.go)
  - Removed keyboard shortcut ('O' key)
  - Added centered "Open File..." button
  - Native file picker dialog with .txt filter
  - Ctrl+O keyboard shortcut
  - Error handling with dialogs
  - Display selected file metadata
- ✅ Updated FileViewerScene (internal/ui/scenes/file_viewer.go)
  - Display file type (Focus/Technology)
  - Display Base_path
  - Better layout with metadata
- ✅ Application compiles and runs successfully

**Technical Implementation:**
- Base_path detection algorithm:
  ```
  E:/mods/my_mod/common/national_focus/brazil.txt
  → Find "common" in path
  → Extract everything before "common"
  → Result: E:/mods/my_mod
  ```
- Validation checks:
  - Base_path exists
  - Contains `common/` directory
  - Has `national_focus/` or `technologies/` subdirectories
- File type detection from path keywords

**User Experience:**
- Clean startup screen with centered button
- Native OS file picker (Windows/Linux/Mac)
- Automatic mod structure validation
- Clear error messages via dialogs
- File metadata display after selection

**Benefits Achieved:**
- ✅ Professional file selection UX
- ✅ Works with mods in any location
- ✅ Automatic Base_path detection
- ✅ Foundation ready for GFX/icon loading in Phase 2
- ✅ No hardcoded paths

**Next Steps:**
- Test with real HOI4 mod files
- Implement Paradox script lexer
- Build parser for focus/tech files

---

### [2025-01-02] - Bug Fix: Base_path Detection

**Problem:** 
- Windows paths were duplicated: `C:\C:Users\...` instead of `C:\Users\...`
- Error: `DNS_ERROR_INVALID_NAME (123)` when validating mod structure
- `os.Stat()` failed on malformed paths

**Root Cause:**
- In `DetectBasePath()`, `filepath.Join(parts[:commonIndex]...)` already included drive letter
- Then line 62 added drive letter again: `parts[0] + separator + basePath`
- Result: `C:\` + `C:\Users\...` = `C:\C:Users\...`

**Solution:**
- Changed logic to: `parts[0] + separator + filepath.Join(parts[1:commonIndex]...)`
- Now correctly builds: `C:\` + `Users\...\mod\name` = `C:\Users\...\mod\name`

**Testing:**
- ✅ Tested with real HOI4 mod: `C:\Users\krzor\Documents\Paradox Interactive\Hearts of Iron IV\mod\USSR_Class_Struggle`
- ✅ Base_path correctly detected
- ✅ Mod structure validation passes
- ✅ File loads and displays correctly

**Files Modified:**
- `internal/app/mod_loader.go` - Fixed `DetectBasePath()` function

