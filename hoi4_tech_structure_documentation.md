# Hearts of Iron IV Technology File Structure Documentation

## Overview
This document describes the structure of Hearts of Iron IV technology files, specifically for the BlackIce mod's WWI land doctrine system. This documentation is optimized for LLM consumption in IDEs like Windsurf, VSCode, and Roocode to facilitate the creation of a visual technology tree editor.

## File Format
- **Type**: Plain text with custom Paradox scripting language
- **Extension**: `.txt`
- **Encoding**: UTF-8
- **Structure**: Nested key-value pairs with specific syntax

## Basic Syntax Rules

### Comments
```paradox
# Single line comment
### Multi-character comment marker for sections ###
```

### Variables
```paradox
@VARIABLE_NAME = value  # Define at file start
base = @VARIABLE_NAME   # Use in calculations
```

### Blocks
```paradox
key_name = {
    # Content inside curly braces
    nested_key = value
    another_block = {
        # Can be nested
    }
}
```

## Root Structure

```paradox
technologies = {
    # All technology definitions go here
    tech_id_1 = { ... }
    tech_id_2 = { ... }
    # ...
}
```

## Technology Definition Structure

### Complete Technology Block
```paradox
technology_id = {
    # 1. Access Control
    allow = { 
        # Conditions for technology to appear
        has_country_flag = UNLOCK:folder_name
        OR = {
            original_tag = GER
            has_tech = prerequisite_tech
        }
    }
    
    # 2. Effects/Bonuses
    category_all_infantry = {
        breakthrough = 0.05      # +5% breakthrough
        soft_attack = 0.03      # +3% soft attack
        defense = -0.02         # -2% defense (negative values)
        max_organisation = 10   # +10 organization (absolute)
    }
    
    # Global modifiers
    land_reinforce_rate = 0.01
    planning_speed = 0.05
    max_dig_in = 5
    attrition = -0.02
    
    # 3. Paths/Prerequisites
    path = {
        leads_to_tech = next_tech_id_1
        research_cost_coeff = 1  # Cost multiplier
    }
    path = {
        leads_to_tech = next_tech_id_2
        research_cost_coeff = 0.75  # 75% of normal cost
    }
    
    # 4. Research Properties
    research_cost = 1.5  # Base research cost
    xp_research_type = army  # army/navy/air
    xp_boost_cost = 15  # XP cost to boost
    xp_research_bonus = 1  # Bonus when boosted
    
    # 5. UI Position
    folder = {
        name = ww1_land_doctrine_folder
        position = { x = -3 y = 5 }  # Grid position
    }
    
    # 6. Completion Effects
    on_research_complete = {
        set_temp_variable = { tech = token:technology_id }
        is_tech_ai_valid = yes
        custom_effect_tooltip = tooltip_key
    }
    
    # 7. Categories/Tags
    categories = {
        land_doctrine
        ww1_doctrine
        german_doctrine  # Nation-specific
    }
    
    # 8. AI Weights
    ai_will_do = {
        base = 10000  # Base priority
        modifier = {
            factor = 2
            original_tag = GER
        }
    }
    
    ai_research_weights = {
        land_doctrine = 1.5
        ww1_doctrine = 1.5
    }
    
    # 9. Hidden Properties
    hidden_modifier = { ww1_doctrine_level = 1 }
    
    # 10. Special Properties
    enable_tactic = tactic_name
    enable_building = {
        building = bunker
        level = 3
    }
    
    # 11. Mutual Exclusion
    XOR = {
        other_exclusive_tech
    }
}
```

## Key Data Types

### Unit Categories
Common unit categories that can receive bonuses:
- `category_all_infantry` - All infantry types
- `category_light_infantry` - Light/special infantry
- `category_all_armor` - All armor units
- `category_artillery` - Artillery units
- `category_all_support` - Support companies
- `category_special_forces` - Elite units
- `category_all_DIV_HQ` - Division HQ elements
- `cavalry`, `infantry`, `guards_infantry`, `ss_infantry` - Specific unit types

### Modifier Types
Common modifiers that can be applied:

#### Combat Stats
- `breakthrough` - Offensive capability (percentage)
- `soft_attack` - Attack vs soft targets (percentage)
- `hard_attack` - Attack vs armored targets (percentage)
- `defense` - Defensive capability (percentage)
- `air_attack` - Anti-air capability (percentage)

#### Organization & Morale
- `max_organisation` - Maximum organization (absolute value)
- `default_morale` - Base morale (percentage)
- `army_morale_factor` - Army-wide morale modifier

#### Movement & Supply
- `maximum_speed` - Unit speed modifier
- `org_loss_when_moving` - Organization loss when moving
- `supply_consumption_factor` - Supply usage modifier
- `no_supply_grace` - Hours before supply penalties
- `out_of_supply_factor` - Penalty reduction when out of supply

#### Entrenchment & Planning
- `max_dig_in` - Maximum entrenchment level
- `dig_in_speed_factor` - Entrenchment speed
- `planning_speed` - Planning bonus accumulation
- `max_planning` - Maximum planning bonus

#### Special
- `partisan_effect` - Partisan activity modifier
- `resistance_damage_to_garrison` - Resistance damage reduction
- `special_forces_cap` - Special forces capacity
- `land_night_attack` - Night combat bonus

### Conditions (for `allow` blocks)
```paradox
allow = {
    has_country_flag = flag_name
    has_tech = technology_id
    NOT = { has_tech = other_tech }
    OR = {
        original_tag = GER
        original_tag = PRE
    }
    AND = {
        condition_1
        condition_2
    }
}
```

## Folder/Tree Organization

### Position Grid System
- **X-axis**: Horizontal position (negative = left, positive = right)
- **Y-axis**: Vertical position (higher values = lower on screen)
- **Grid spacing**: Usually 2 units apart for clear separation
- **Typical range**: x from -10 to 10, y from 0 to 20

### Visual Connection Rules
1. Technologies connected via `path` blocks
2. Prerequisites must be at lower Y values than descendants
3. Branching paths should maintain clear visual separation

## Special Mechanics

### Mutually Exclusive Technologies
```paradox
tech_option_a = {
    XOR = {
        tech_option_b
    }
    # Player must choose between option_a OR option_b
}
```

### Adoptable Doctrines Pattern
```paradox
# Nation-specific tech
german_exclusive_tech = {
    allow = {
        OR = {
            original_tag = GER
            has_tech = adopt_german_tactics  # Adoption mechanism
        }
    }
}

# Adoption enabler
adopt_german_tactics = {
    allow = {
        NOT = { 
            original_tag = GER
            has_tech = adopt_french_tactics  # Prevent multiple adoptions
        }
    }
}
```

### Dynamic Visibility
Technologies can be hidden/shown based on:
1. Country tags
2. Previously researched technologies
3. Country flags
4. Game state conditions

## Editor Requirements

### Core Features Needed
1. **Visual Grid Display**
   - Render technologies as nodes on X-Y grid
   - Show connections between technologies (paths)
   - Display technology icons/names

2. **Interactive Manipulation**
   - Drag & drop to change positions
   - Click to edit properties
   - Right-click context menus
   - Connect/disconnect paths visually

3. **Property Editor Panel**
   - Edit all technology properties
   - Add/remove modifiers
   - Manage categories and tags
   - Configure AI weights

4. **Validation System**
   - Check for circular dependencies
   - Validate position conflicts
   - Ensure path continuity
   - Verify syntax correctness

5. **Import/Export**
   - Parse existing .txt files
   - Generate properly formatted output
   - Maintain comments and structure
   - Handle special characters and encoding

### Data Model Suggestions

```typescript
interface Technology {
    id: string;
    allow?: Condition;
    effects: Effect[];
    paths: Path[];
    researchCost: number;
    position: Position;
    categories: string[];
    folder: string;
    // ... other properties
}

interface Effect {
    target: string;  // e.g., "category_all_infantry"
    modifiers: Record<string, number>;
}

interface Path {
    leadsTo: string;  // technology ID
    costCoeff: number;
}

interface Position {
    x: number;
    y: number;
}

interface Condition {
    type: 'AND' | 'OR' | 'NOT' | 'SIMPLE';
    conditions?: Condition[];
    value?: string;
}
```

### UI/UX Considerations

1. **Grid View**
   - Zoom in/out functionality
   - Pan across large trees
   - Snap-to-grid for positions
   - Visual indicators for exclusive choices

2. **Connection Visualization**
   - Arrows showing prerequisites
   - Different line styles for cost coefficients
   - Highlight paths on hover

3. **Search & Filter**
   - Search by technology ID
   - Filter by categories
   - Show/hide based on country
   - Highlight search results

4. **Editing Modes**
   - Move mode (reposition)
   - Connect mode (create paths)
   - Edit mode (modify properties)
   - View mode (read-only)

5. **Validation Feedback**
   - Real-time error highlighting
   - Warning indicators
   - Suggestion tooltips
   - Validation report panel

## Common Patterns

### Technology Tree Branches
```paradox
# Root technology
root_tech = {
    path = { leads_to_tech = branch_a_start }
    path = { leads_to_tech = branch_b_start }
    position = { x = 0 y = 0 }
}

# Branch A
branch_a_start = {
    path = { leads_to_tech = branch_a_tech_2 }
    position = { x = -4 y = 2 }
}

# Branch B
branch_b_start = {
    path = { leads_to_tech = branch_b_tech_2 }
    position = { x = 4 y = 2 }
}
```

### Progressive Upgrades
```paradox
basic_tech = {
    category_all_infantry = {
        soft_attack = 0.02
    }
}

improved_tech = {
    category_all_infantry = {
        soft_attack = 0.03  # Additional bonus
        breakthrough = 0.02  # New bonus
    }
}

advanced_tech = {
    category_all_infantry = {
        soft_attack = 0.05  # Further improvement
        breakthrough = 0.04
        defense = 0.03  # Another new bonus
    }
}
```

## File Structure Best Practices

1. **Organization**
   - Group related technologies together
   - Use clear section comments
   - Define variables at the top
   - Maintain consistent indentation

2. **Naming Conventions**
   - Use descriptive IDs: `german_stosstrupp_doctrine`
   - Avoid spaces (use underscores)
   - Be consistent with prefixes/suffixes
   - Use lowercase for IDs

3. **Balance Considerations**
   - Small incremental bonuses (0.02-0.05)
   - Higher costs for stronger effects
   - Consider cumulative effects
   - Test mutual exclusions

4. **Documentation**
   - Comment complex conditions
   - Explain unusual mechanics
   - Note dependencies
   - Document balance decisions

## Error Handling

Common issues to check for:
1. Missing closing braces
2. Duplicate technology IDs
3. Invalid reference IDs in paths
4. Circular dependencies
5. Position conflicts (same x,y)
6. Invalid modifier targets
7. Syntax errors in conditions

## Testing Checklist

- [ ] All technologies load without errors
- [ ] Paths connect correctly in-game
- [ ] Positions display properly
- [ ] Effects apply as expected
- [ ] AI researches appropriately
- [ ] Tooltips show correctly
- [ ] Exclusive choices work
- [ ] Adoption mechanics function

## Additional Resources

- Paradox modding wiki
- HOI4 defines files for valid modifiers
- BlackIce mod documentation
- Community modding forums

---

*This documentation is designed for LLM consumption to facilitate the development of a visual technology tree editor for Hearts of Iron IV mods.*
