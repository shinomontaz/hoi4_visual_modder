package parser

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/shinomontaz/hoi4_visual_modder/internal/domain"
)

// BookmarkParser parses bookmark files to extract country lists
type BookmarkParser struct {
	modPath  string
	gamePath string
}

// NewBookmarkParser creates a new bookmark parser
func NewBookmarkParser(modPath, gamePath string) *BookmarkParser {
	return &BookmarkParser{
		modPath:  modPath,
		gamePath: gamePath,
	}
}

// ParseBookmarks loads all bookmarks from mod and game folders
func (bp *BookmarkParser) ParseBookmarks() ([]*domain.Bookmark, error) {
	bookmarks := make([]*domain.Bookmark, 0)

	// Try mod path first
	if bp.modPath != "" {
		modBookmarks, err := bp.parseBookmarksFromPath(bp.modPath)
		if err == nil && len(modBookmarks) > 0 {
			bookmarks = append(bookmarks, modBookmarks...)
		}
	}

	// If no bookmarks in mod, try game path
	if len(bookmarks) == 0 && bp.gamePath != "" {
		gameBookmarks, err := bp.parseBookmarksFromPath(bp.gamePath)
		if err == nil {
			bookmarks = append(bookmarks, gameBookmarks...)
		}
	}

	if len(bookmarks) == 0 {
		return nil, fmt.Errorf("no bookmarks found in mod or game folders")
	}

	return bookmarks, nil
}

// parseBookmarksFromPath parses all bookmark files in a given base path
func (bp *BookmarkParser) parseBookmarksFromPath(basePath string) ([]*domain.Bookmark, error) {
	bookmarksDir := filepath.Join(basePath, "common", "bookmarks")

	// Check if directory exists
	if _, err := os.Stat(bookmarksDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("bookmarks directory not found: %s", bookmarksDir)
	}

	// Find all .txt files
	files, err := filepath.Glob(filepath.Join(bookmarksDir, "*.txt"))
	if err != nil {
		return nil, fmt.Errorf("failed to list bookmark files: %w", err)
	}

	bookmarks := make([]*domain.Bookmark, 0)

	// Parse each file
	for _, file := range files {
		fileBookmarks, err := bp.parseBookmarkFile(file)
		if err != nil {
			// Log error but continue with other files
			fmt.Printf("Warning: failed to parse %s: %v\n", file, err)
			continue
		}
		bookmarks = append(bookmarks, fileBookmarks...)
	}

	return bookmarks, nil
}

// parseBookmarkFile parses a single bookmark file
func (bp *BookmarkParser) parseBookmarkFile(filePath string) ([]*domain.Bookmark, error) {
	// Read file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Parse using existing parser
	p := NewParser(string(content))
	program, err := p.Parse()
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	bookmarks := make([]*domain.Bookmark, 0)

	// Find "bookmarks" block
	for _, stmt := range program.Statements {
		if assign, ok := stmt.(*AssignmentStatement); ok {
			if assign.Name.Value == "bookmarks" {
				if block, ok := assign.Value.(*BlockStatement); ok {
					bookmarks = bp.parseBookmarksBlock(block)
				}
			}
		}
	}

	return bookmarks, nil
}

// parseBookmarksBlock parses the bookmarks { ... } block
func (bp *BookmarkParser) parseBookmarksBlock(block *BlockStatement) []*domain.Bookmark {
	bookmarks := make([]*domain.Bookmark, 0)

	for _, stmt := range block.Statements {
		if assign, ok := stmt.(*AssignmentStatement); ok {
			if assign.Name.Value == "bookmark" {
				if bookmarkBlock, ok := assign.Value.(*BlockStatement); ok {
					bookmark := bp.parseBookmark(bookmarkBlock)
					if bookmark != nil {
						bookmarks = append(bookmarks, bookmark)
					}
				}
			}
		}
	}

	return bookmarks
}

// parseBookmark parses a single bookmark { ... } block
func (bp *BookmarkParser) parseBookmark(block *BlockStatement) *domain.Bookmark {
	bookmark := &domain.Bookmark{
		Countries: make([]*domain.BookmarkCountry, 0),
	}

	for _, stmt := range block.Statements {
		if assign, ok := stmt.(*AssignmentStatement); ok {
			key := assign.Name.Value

			switch key {
			case "name":
				bookmark.Name = bp.extractString(assign.Value)
			case "desc":
				bookmark.Description = bp.extractString(assign.Value)
			case "date":
				bookmark.Date = bp.extractString(assign.Value)
			case "default_country":
				bookmark.DefaultCountry = bp.extractString(assign.Value)
			case "effect", "picture", "default", "available":
				// Skip these fields
				continue
			default:
				// Check if this is a country block
				// Countries can be 2-3 letter tags (USA, GER, ENG, etc.) or special cases like "---"
				if countryBlock, ok := assign.Value.(*BlockStatement); ok {
					// If it's a block and looks like a country tag, parse it
					if len(key) >= 2 && len(key) <= 3 {
						country := bp.parseCountry(key, countryBlock)
						if country != nil {
							bookmark.Countries = append(bookmark.Countries, country)
						}
					}
				}
			}
		}
	}

	return bookmark
}

// parseCountry parses a country block
func (bp *BookmarkParser) parseCountry(tag string, block *BlockStatement) *domain.BookmarkCountry {
	country := &domain.BookmarkCountry{
		Tag:     tag,
		Name:    tag,  // Default to tag
		IsMajor: true, // Default to major (minor = yes means it's minor)
		Ideas:   make([]string, 0),
		Focuses: make([]string, 0),
	}

	for _, stmt := range block.Statements {
		if assign, ok := stmt.(*AssignmentStatement); ok {
			key := assign.Name.Value

			switch key {
			case "history":
				country.History = bp.extractString(assign.Value)
			case "ideology":
				country.Ideology = bp.extractString(assign.Value)
			case "minor":
				// "minor = yes" means it's a minor power
				val := bp.extractString(assign.Value)
				country.IsMajor = val != "yes"
			case "ideas":
				country.Ideas = bp.extractArray(assign.Value)
			case "focuses":
				country.Focuses = bp.extractArray(assign.Value)
			}
		}
	}

	return country
}

// extractString extracts string value from expression
func (bp *BookmarkParser) extractString(expr Expression) string {
	switch v := expr.(type) {
	case *StringLiteral:
		return v.Value
	case *Identifier:
		return v.Value
	default:
		return ""
	}
}

// extractArray extracts array of strings from block
func (bp *BookmarkParser) extractArray(expr Expression) []string {
	result := make([]string, 0)

	if block, ok := expr.(*BlockStatement); ok {
		for _, stmt := range block.Statements {
			if assign, ok := stmt.(*AssignmentStatement); ok {
				result = append(result, assign.Name.Value)
			}
		}
	}

	return result
}
