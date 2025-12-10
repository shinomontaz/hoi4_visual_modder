package components

import (
	"bytes"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lukegb/dds"
)

// IconLoader handles loading and caching of technology/focus icons
type IconLoader struct {
	basePath string // Mod folder path
	gamePath string // HOI4 game installation path (fallback)
	cache    map[string]*ebiten.Image
	mu       sync.RWMutex

	// Placeholder for missing icons
	placeholder *ebiten.Image
}

// NewIconLoader creates a new icon loader
func NewIconLoader(basePath string) *IconLoader {
	loader := &IconLoader{
		basePath: basePath,
		gamePath: detectGamePath(), // Auto-detect game installation
		cache:    make(map[string]*ebiten.Image),
	}

	// Create placeholder image (gray square with X)
	loader.placeholder = loader.createPlaceholder()

	return loader
}

// SetGamePath sets the HOI4 game installation path for fallback icon loading
func (il *IconLoader) SetGamePath(gamePath string) {
	il.mu.Lock()
	defer il.mu.Unlock()
	il.gamePath = gamePath
}

// detectGamePath tries to auto-detect HOI4 installation path
func detectGamePath() string {
	// Common installation paths
	possiblePaths := []string{
		`C:\Program Files (x86)\Steam\steamapps\common\Hearts of Iron IV`,
		`C:\Program Files\Steam\steamapps\common\Hearts of Iron IV`,
		`D:\Steam\steamapps\common\Hearts of Iron IV`,
		`E:\Steam\steamapps\common\Hearts of Iron IV`,
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(filepath.Join(path, "hoi4.exe")); err == nil {
			println("Auto-detected game path:", path)
			return path
		}
	}

	println("Game path not auto-detected, will use placeholder for missing icons")
	return ""
}

// LoadTechIcon loads a technology icon by name
func (il *IconLoader) LoadTechIcon(iconName string) *ebiten.Image {
	if iconName == "" {
		return il.placeholder
	}

	// Check cache first
	il.mu.RLock()
	if img, exists := il.cache[iconName]; exists {
		il.mu.RUnlock()
		return img
	}
	il.mu.RUnlock()

	// Try to load from file
	img := il.loadIconFromFile(iconName, "technologies")

	// Cache the result (even if it's placeholder)
	il.mu.Lock()
	il.cache[iconName] = img
	il.mu.Unlock()

	return img
}

// LoadFocusIcon loads a focus icon by name
func (il *IconLoader) LoadFocusIcon(iconName string) *ebiten.Image {
	if iconName == "" {
		return il.placeholder
	}

	// Check cache first
	il.mu.RLock()
	if img, exists := il.cache[iconName]; exists {
		il.mu.RUnlock()
		return img
	}
	il.mu.RUnlock()

	// Try to load from file
	img := il.loadIconFromFile(iconName, "goals")

	// Cache the result
	il.mu.Lock()
	il.cache[iconName] = img
	il.mu.Unlock()

	return img
}

// loadIconFromFile attempts to load icon from various possible locations
func (il *IconLoader) loadIconFromFile(iconName, category string) *ebiten.Image {
	// First try mod folder
	modPaths := []string{
		filepath.Join(il.basePath, "gfx", "interface", category, iconName+".dds"),
		filepath.Join(il.basePath, "gfx", "interface", category, iconName+".png"),
	}

	for _, path := range modPaths {
		if img := il.tryLoadImage(path); img != nil {
			return img
		}
	}

	// If not found in mod, try game folder (fallback)
	if il.gamePath != "" {
		gamePaths := []string{
			filepath.Join(il.gamePath, "gfx", "interface", category, iconName+".dds"),
			filepath.Join(il.gamePath, "gfx", "interface", category, iconName+".png"),
		}

		for _, path := range gamePaths {
			if img := il.tryLoadImage(path); img != nil {
				return img
			}
		}
	}

	// Icon not found, return placeholder
	return il.placeholder
}

// tryLoadImage attempts to load an image from a file
func (il *IconLoader) tryLoadImage(path string) *ebiten.Image {
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	// Try to decode based on extension
	ext := filepath.Ext(path)

	var img image.Image
	switch ext {
	case ".dds":
		// Decode DDS using lukegb/dds library
		img, err = dds.Decode(bytes.NewReader(data))
		if err != nil {
			// DDS decode failed, return nil
			return nil
		}

	case ".png", ".jpg", ".jpeg":
		img, _, err = image.Decode(bytes.NewReader(data))
		if err != nil {
			return nil
		}

	default:
		return nil
	}

	// Convert to ebiten image
	if img != nil {
		return ebiten.NewImageFromImage(img)
	}

	return nil
}

// createPlaceholder creates a placeholder image for missing icons
func (il *IconLoader) createPlaceholder() *ebiten.Image {
	// Create a 64x64 gray square
	size := 64
	img := ebiten.NewImage(size, size)

	// Fill with gray color
	grayColor := color.RGBA{100, 100, 100, 255}
	img.Fill(grayColor)

	// TODO: Draw an X or question mark
	// For now, just return gray square

	return img
}

// ClearCache clears the icon cache
func (il *IconLoader) ClearCache() {
	il.mu.Lock()
	defer il.mu.Unlock()

	il.cache = make(map[string]*ebiten.Image)
}

// GetCacheSize returns the number of cached icons
func (il *IconLoader) GetCacheSize() int {
	il.mu.RLock()
	defer il.mu.RUnlock()

	return len(il.cache)
}

// PreloadIcons preloads a list of icons
func (il *IconLoader) PreloadIcons(iconNames []string, category string) {
	for _, name := range iconNames {
		if category == "technologies" {
			il.LoadTechIcon(name)
		} else {
			il.LoadFocusIcon(name)
		}
	}
}

// SetBasePath updates the base path for icon loading
func (il *IconLoader) SetBasePath(basePath string) {
	il.mu.Lock()
	defer il.mu.Unlock()

	il.basePath = basePath
	// Clear cache when base path changes
	il.cache = make(map[string]*ebiten.Image)
}

// GetPlaceholder returns the placeholder image
func (il *IconLoader) GetPlaceholder() *ebiten.Image {
	return il.placeholder
}
