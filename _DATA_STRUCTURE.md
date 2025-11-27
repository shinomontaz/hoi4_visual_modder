# HOI4 Data Structure Documentation

## 1. Mod Selection and Validation

### Mod Descriptor File (*.mod)
**Location:** `<Documents>/Paradox Interactive/Hearts of Iron IV/mod/`

**Structure:**
```
version="8.1.0"
tags={ "Historical" "Gameplay" }
replace_path="common/..."
name="BlackICE Historical Immersion Mod"
supported_version="1.14.2"
path="mod/BlackICE_Historical_Immersion_Mod/"  # ← KEY FIELD
```

**Validation Rules:**
1. User selects `.mod` file
2. Parse `path` variable to get mod folder name
3. Verify mod folder exists next to `.mod` file
4. No deep validation needed - if `.mod` + folder exist = valid

**Example:**
```
Selected: C:\Users\...\Hearts of Iron IV\mod\BlackICE.mod
Parsed path: "mod/BlackICE_Historical_Immersion_Mod/"
Expected folder: C:\Users\...\Hearts of Iron IV\mod\BlackICE_Historical_Immersion_Mod\
```

---

## 2. Game Installation Selection

### Game Folder Validation
**Common Paths:**
- `C:\Program Files (x86)\Steam\steamapps\common\Hearts of Iron IV`
- `C:\Program Files\Steam\steamapps\common\Hearts of Iron IV`
- `D:\Steam\steamapps\common\Hearts of Iron IV`

**Validation Rules:**
1. User selects folder OR clicks "Auto-detect"
2. Check for `hoi4.exe` (or platform-specific executable)
3. Check for standard folders: `common/`, `gfx/`, `history/`
4. If user creates fake structure - their problem

---

## 3. Country Selection

### Bookmark Files
**Location (priority order):**
1. `<mod_path>/common/bookmarks/*.txt` (check first)
2. `<game_path>/common/bookmarks/*.txt` (fallback)

**Structure:**
```
bookmarks = {
    bookmark = {
        name = "GATHERING_STORM_NAME"
        date = 1936.1.1.12
        default_country = "GER"
        
        "USA" = {
            history = "USA_GATHERING_STORM_DESC"
            ideology = liberalism
            minor = no  # ← implicit if not specified
        }
        
        "CAN" = {
            minor = yes  # ← minor nation flag
            history = "CAN_GATHERING_STORM_DESC"
        }
    }
}
```

**Parsing Logic:**
1. Find all `.txt` files in `common/bookmarks/`
2. Parse each `bookmark` block
3. Extract country tags (e.g., "USA", "GER", "SOV")
4. Group by major/minor (check `minor = yes` flag)
5. Display in UI with filtering

**Country Data:**
- **Tag:** 3-letter code (GER, SOV, USA)
- **Name:** Localized from `history` field or tag
- **Type:** Major/Minor
- **Ideology:** Starting ideology
- **Ideas:** Starting national spirits

---

## 4. Technology Trees

### Technology Tags and Folders
**Location:** `common/technology_tags/*.txt`

**Structure:**
```
technology_categories = {
    artillery
    light_fighter
    naval_equipment
    # ... hundreds of categories
}

technology_folders = {
    infantry_folder = {
        ledger = army
    }
    
    # Country-specific overrides
    luftwaffe_folder = {
        ledger = air
        available = { 
            has_country_flag = GER_air
        }
    }
    
    sovietair_folder = {
        ledger = air
        available = {
            has_country_flag = SOV_air
        }
    }
}
```

### Technology File Resolution

**Generic Folders → Generic Files:**
```
infantry_folder → common/technologies/infantry.txt
air_techs_folder → common/technologies/air_techs.txt
naval_folder → common/technologies/naval.txt
```

**Country-Specific Folders → Country Files:**
```
luftwaffe_folder (GER) → common/technologies/GER_air.txt
sovietair_folder (SOV) → common/technologies/SOV_air.txt
usair_folder (USA) → common/technologies/USA_air.txt
```

**Resolution Logic:**
1. Parse `technology_folders` from `common/technology_tags/*.txt`
2. For selected country, check `available` conditions
3. Match folder name to file name:
   - `luftwaffe_folder` → `GER_air.txt`
   - `sovietair_folder` → `SOV_air.txt`
   - `infantry_folder` → `infantry.txt`
4. Load technology tree from resolved file

**Example for USSR:**
```
Available folders:
- infantry_folder → infantry.txt (generic)
- support_folder → support.txt (generic)
- trm_armour_sov_folder → SOV_armor.txt (country-specific)
- sovietair_folder → SOV_air.txt (country-specific)
- industry_folder → industry.txt (generic)
```

---

## 5. National Focus Trees

### Focus File Location
**Path:** `common/national_focus/<country_tag>.txt`

**Examples:**
- `common/national_focus/germany.txt`
- `common/national_focus/soviet.txt`
- `common/national_focus/usa.txt`

**Generic Focus:**
- `common/national_focus/generic.txt` (fallback)

---

## Data Flow Summary

```
1. User selects .mod file
   ↓
2. Parse path, validate mod folder
   ↓
3. User selects/auto-detects game folder
   ↓
4. Validate game installation
   ↓
5. Load bookmarks (mod → game fallback)
   ↓
6. Display country list
   ↓
7. User selects country (e.g., SOV)
   ↓
8. Country Scene: Show options
   ├─ National Focus → load common/national_focus/soviet.txt
   └─ Technologies → load technology_tags, resolve folders
      ├─ Infantry → infantry.txt
      ├─ Air → SOV_air.txt (country-specific)
      └─ Armor → SOV_armor.txt (country-specific)
```

---

## Implementation Notes

### Mod Path vs Game Path
- **Mod path:** Always check first for overrides
- **Game path:** Fallback for vanilla content
- **Icons:** Check mod → game → placeholder

### Country Flags
Technology folders use country flags like:
- `has_country_flag = GER_air`
- `has_country_flag = SOV_armor`

These flags are set in country history files or focus trees.

### Technology Folder Naming Convention
Pattern: `<prefix>_<country>_<category>_folder`
- `trm_armour_ger_folder` → `GER_armor.txt`
- `luftwaffe_folder` → `GER_air.txt`
- `sovietair_folder` → `SOV_air.txt`

File resolution requires mapping folder names to file names (may need manual mapping table).
