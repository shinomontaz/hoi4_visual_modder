# Hearts of Iron IV National Focus Tree Structure Documentation

## Overview
This document describes the structure of Hearts of Iron IV National Focus trees, using Brazil's focus tree as a reference example. This documentation is optimized for LLM consumption in IDEs like Windsurf, VSCode, and Roocode to facilitate the creation of a visual focus tree editor.

## File Format
- **Type**: Plain text with Paradox scripting language
- **Extension**: `.txt`
- **Encoding**: UTF-8
- **Structure**: Hierarchical tree of interconnected focus nodes

## Root Structure

```paradox
focus_tree = {
    id = tree_identifier        # Unique ID for this focus tree
    
    country = {
        factor = 0              # Base weight
        modifier = {
            add = 10           # Add weight conditionally
            tag = BRA          # For specific country tag
        }
    }
    
    continuous_focus_position = { x = 1250 y = 1300 }  # Position for continuous focuses
    
    default = no               # Is this the default tree?
    reset_on_civilwar = no     # Reset tree during civil war?
    
    # Focus definitions
    focus = { ... }
    focus = { ... }
    # ... more focuses
}
```

## Focus Definition Structure

### Complete Focus Block
```paradox
focus = {
    # 1. Identification
    id = focus_unique_id          # Unique identifier
    icon = GFX_icon_name          # Icon graphic reference
    
    # 2. Position
    x = 3                         # Absolute X position
    y = 0                         # Absolute Y position
    # OR relative positioning
    relative_position_id = parent_focus_id
    
    # 3. Dependencies
    prerequisite = { 
        focus = required_focus_1 
    }
    prerequisite = { 
        focus = required_focus_2 
        focus = required_focus_3  # OR relationship within block
    }
    mutually_exclusive = { 
        focus = exclusive_focus 
    }
    
    # 4. Availability
    available = {
        date > 1937.1.1           # Date condition
        has_government = fascism  # Government type
        has_country_flag = flag_name
        custom_trigger_tooltip = {
            tooltip = tooltip_key
            # Complex conditions
        }
    }
    
    # 5. Cost and Time
    cost = 70                     # Days to complete (base)
    
    # 6. Bypass Conditions
    bypass = {
        # Conditions that auto-complete the focus
        has_war_with = USA
    }
    
    # 7. Cancellation
    cancel_if_invalid = yes       # Cancel if becomes unavailable
    continue_if_invalid = no      # Continue if becomes invalid
    available_if_capitulated = no # Available if capitulated
    
    # 8. Search Filters (for UI)
    search_filters = { FOCUS_FILTER_POLITICAL FOCUS_FILTER_INDUSTRY }
    
    # 9. Completion Rewards
    completion_reward = {
        # Political Power
        add_political_power = 50
        
        # Ideas/National Spirits
        add_ideas = idea_name
        remove_ideas = old_idea
        swap_ideas = {
            remove_idea = old_idea
            add_idea = new_idea
        }
        
        # Resources
        add_resource = {
            type = oil
            amount = 12
            state = 499
        }
        
        # Buildings
        random_owned_controlled_state = {
            add_extra_state_shared_building_slots = 1
            add_building_construction = {
                type = industrial_complex
                level = 1
                instant_build = yes
            }
        }
        
        # Research
        add_research_slot = 1
        add_tech_bonus = {
            name = bonus_name
            bonus = 0.25
            uses = 1
            category = land_doctrine
        }
        
        # Military
        army_experience = 10
        navy_experience = 25
        air_experience = 15
        
        # Territory
        300 = { add_claim_by = BRA }
        
        # Threat
        add_named_threat = { 
            threat = 3 
            name = threat_name 
        }
        
        # Events
        country_event = { 
            id = event.1 
            days = 5 
        }
        
        # Leaders
        retire_country_leader = yes
        kill_country_leader = yes
        
        # Diplomacy
        ENG = {
            country_event = { id = brazil.1 }
        }
    }
    
    # 10. AI Weights
    ai_will_do = {
        factor = 10               # Base AI priority
        modifier = {
            factor = 0
            is_historical_focus_on = yes
        }
    }
    
    # 11. Completion Tooltip
    complete_tooltip = {
        # What shows when hovering over completed focus
        add_extra_state_shared_building_slots = 1
    }
}
```

## Position System

### Absolute Positioning
```paradox
focus = {
    x = 5    # Horizontal position (grid units)
    y = 3    # Vertical position (grid units)
}
```

### Relative Positioning
```paradox
focus = {
    relative_position_id = parent_focus
    x = 2    # Offset from parent
    y = 1    # Offset from parent
}
```

### Grid Rules
- **Grid spacing**: Usually 2 units between focuses
- **X-axis range**: Typically -20 to 40
- **Y-axis range**: Typically 0 to 20
- **Branch spacing**: 2-4 units horizontally between parallel branches

## Dependency Types

### Prerequisites
```paradox
# Single prerequisite (AND)
prerequisite = { focus = focus_a }

# Multiple prerequisites (AND)
prerequisite = { focus = focus_a }
prerequisite = { focus = focus_b }

# OR prerequisites (any one within block)
prerequisite = { 
    focus = option_a 
    focus = option_b 
}
```

### Mutual Exclusivity
```paradox
mutually_exclusive = { 
    focus = alternative_a 
}
mutually_exclusive = { 
    focus = alternative_b 
}
```

## Building Types
Common building types used in completion rewards:
- `industrial_complex` - Civilian factories
- `arms_factory` - Military factories
- `dockyard` - Naval dockyards
- `infrastructure` - Infrastructure
- `air_base` - Air bases
- `naval_base` - Naval bases
- `bunker` - Land fortifications
- `coastal_bunker` - Coastal fortifications
- `steel_refinery` - Steel production
- `aluminium_refinery` - Aluminium production
- `artillery_assembly` - Artillery production

## Resource Types
- `oil`
- `coal`
- `iron`
- `steel`
- `aluminium`
- `rubber`
- `tungsten`
- `chromium`
- `bauxite`

## Common Completion Reward Patterns

### Industrial Development
```paradox
completion_reward = {
    random_owned_controlled_state = {
        add_extra_state_shared_building_slots = 1
        add_building_construction = {
            type = industrial_complex
            level = 1
            instant_build = yes
        }
    }
}
```

### Specific State Development
```paradox
completion_reward = {
    500 = {  # State ID
        if = {
            limit = { is_controlled_by = BRA }
            add_building_construction = {
                type = infrastructure
                level = 1
                instant_build = yes
            }
        }
    }
}
```

### Research Bonuses
```paradox
completion_reward = {
    add_tech_bonus = {
        name = infantry_weapons_bonus
        bonus = 0.25
        uses = 1
        category = infantry_weapons
    }
}
```

### Political Changes
```paradox
completion_reward = {
    set_politics = {
        ruling_party = fascism
        elections_allowed = no
    }
    add_political_power = 150
}
```

## Visual Tree Organization Patterns

### Main Branch Structure
```
        Root Focus (y=0)
             |
    ┌────────┼────────┐
Branch A  Branch B  Branch C
  (x=-4)    (x=0)    (x=4)
```

### Mutually Exclusive Paths
```
    Common Prerequisite
           |
    ┌──────┴──────┐
Option A      Option B
(exclusive)   (exclusive)
```

### Complex Prerequisites
```
Focus A     Focus B
    └───┬───┘
      Focus C
        |
      Focus D
```

## Editor Requirements

### Core Features Needed

1. **Visual Tree Display**
   - Render focuses as nodes on grid
   - Show prerequisite lines/arrows
   - Display mutual exclusivity indicators
   - Show focus icons and names
   - Status indicators (available/completed/current)

2. **Interactive Manipulation**
   - Drag & drop focus positioning
   - Click to edit properties
   - Connect/disconnect prerequisites visually
   - Create mutual exclusivity links
   - Copy/paste focus branches

3. **Property Editor Panel**
   - Edit all focus properties
   - Manage completion rewards
   - Configure availability conditions
   - Set AI weights
   - Icon selector

4. **Validation System**
   - Check for circular dependencies
   - Validate position conflicts
   - Ensure prerequisite validity
   - Check mutual exclusivity logic
   - Validate date conditions

5. **Tree Management**
   - Branch organization tools
   - Auto-align features
   - Spacing adjustment
   - Batch operations

### Data Model Suggestions

```typescript
interface FocusTree {
    id: string;
    country: CountryFilter;
    continuousFocusPosition?: Position;
    default: boolean;
    resetOnCivilWar: boolean;
    focuses: Focus[];
}

interface Focus {
    id: string;
    icon: string;
    position: Position;
    relativePositionId?: string;
    prerequisites: Prerequisite[];
    mutuallyExclusive: string[];
    available?: Condition;
    bypass?: Condition;
    cost: number;
    completionReward: Reward[];
    aiWillDo?: AIWeight;
    searchFilters?: string[];
}

interface Position {
    x: number;
    y: number;
}

interface Prerequisite {
    focuses: string[];  // OR relationship between these
}

interface Reward {
    type: 'political_power' | 'idea' | 'building' | 'resource' | 'tech_bonus' | 'experience';
    data: any;  // Type-specific data
}

interface Condition {
    type: 'AND' | 'OR' | 'NOT' | 'SIMPLE';
    conditions?: Condition[];
    check?: string;
    value?: any;
}
```

### UI/UX Considerations

1. **Grid View**
   - Zoom in/out functionality
   - Pan across large trees
   - Grid snap for positioning
   - Visual grid overlay
   - Minimap for navigation

2. **Connection Visualization**
   - Different line styles for prerequisites vs mutual exclusivity
   - Highlight connections on hover
   - Show connection direction clearly
   - Color coding for different dependency types

3. **Focus Node Display**
   - Icon display
   - Name label
   - Cost indicator
   - Status coloring
   - Completion reward preview

4. **Editing Modes**
   - Move mode (reposition focuses)
   - Connect mode (create dependencies)
   - Edit mode (modify properties)
   - Select mode (multi-select operations)
   - Preview mode (test flow)

5. **Branch Management**
   - Collapse/expand branches
   - Move entire branches
   - Duplicate branches
   - Delete branches with confirmation

## Common Patterns

### Industrial Branch Pattern
```paradox
# Base industrial focus
focus = {
    id = industrial_effort
    x = 5
    y = 0
    cost = 70
    completion_reward = {
        add_ideas = industrial_bonus
    }
}

# Civilian industry path
focus = {
    id = civilian_industry_1
    prerequisite = { focus = industrial_effort }
    x = -2
    y = 1
    relative_position_id = industrial_effort
    # ...
}

# Military industry path
focus = {
    id = military_industry_1
    prerequisite = { focus = industrial_effort }
    x = 2
    y = 1
    relative_position_id = industrial_effort
    # ...
}
```

### Political Branch Pattern
```paradox
# Political root
focus = {
    id = political_effort
    x = 15
    y = 0
    # ...
}

# Ideology choices (mutually exclusive)
focus = {
    id = democratic_path
    prerequisite = { focus = political_effort }
    mutually_exclusive = { focus = fascist_path focus = communist_path }
    # ...
}

focus = {
    id = fascist_path
    prerequisite = { focus = political_effort }
    mutually_exclusive = { focus = democratic_path focus = communist_path }
    # ...
}
```

### Date-Locked Progression
```paradox
focus = {
    id = early_focus
    available = {
        date > 1937.1.1
    }
    # ...
}

focus = {
    id = mid_focus
    prerequisite = { focus = early_focus }
    available = {
        date > 1939.1.1
    }
    # ...
}

focus = {
    id = late_focus
    prerequisite = { focus = mid_focus }
    available = {
        date > 1941.1.1
    }
    # ...
}
```

## Validation Rules

1. **Unique IDs**: Each focus must have a unique id
2. **Valid Prerequisites**: Referenced focuses must exist
3. **No Circular Dependencies**: A→B→C→A is invalid
4. **Position Conflicts**: No two focuses at same x,y
5. **Mutual Exclusivity Logic**: Must be symmetric
6. **Date Consistency**: Child focuses shouldn't unlock before parents
7. **Cost Values**: Must be positive numbers
8. **State IDs**: Must be valid game state IDs

## Best Practices

1. **Naming Conventions**
   - Use descriptive IDs: `develop_civilian_industry_1`
   - Prefix with country tag for unique trees: `BRA_industrial_effort`
   - Use consistent suffixes: `_1`, `_2` for sequences

2. **Tree Organization**
   - Group related focuses into branches
   - Maintain consistent spacing (2 units)
   - Align focuses vertically in chains
   - Use relative positioning for branch consistency

3. **Cost Balancing**
   - Early focuses: 35-70 days
   - Mid-game focuses: 70-140 days
   - Late-game focuses: 140+ days
   - Consider player agency vs historical pacing

4. **Reward Scaling**
   - Early rewards: Small bonuses
   - Progressive scaling through branches
   - Avoid overwhelming early game
   - Balance between paths

## Common Issues to Check

1. Missing prerequisites breaking chains
2. Overlapping focus positions
3. Unreachable focuses (impossible conditions)
4. Missing icons causing display issues
5. Invalid state IDs in rewards
6. Typos in focus ID references
7. Asymmetric mutual exclusivity
8. Date conditions preventing progression

## Testing Checklist

- [ ] All focuses display correctly
- [ ] Prerequisites connect properly
- [ ] Mutual exclusivities work
- [ ] Date conditions trigger correctly
- [ ] Rewards apply as expected
- [ ] AI selects focuses appropriately
- [ ] Tree resets properly on civil war (if applicable)
- [ ] Bypass conditions work
- [ ] Cancellation conditions work

---

*This documentation is designed for LLM consumption to facilitate the development of a visual national focus tree editor for Hearts of Iron IV.*
