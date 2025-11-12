# MindMap - Terminal Mind Mapping Tool

A fast, keyboard-driven mind mapping tool for the terminal built with Go and Bubble Tea.

## Features

- **Canvas-based layout** with pan and zoom
- **Keyboard-driven navigation** with spatial awareness
- **Hierarchical structure** with parent-child relationships
- **Automatic layout adjustment** prevents node overlapping
- **Color-coded branches** for visual organization
- **JSON persistence** for saving/loading maps
- **Border-to-border connections** for clean visuals

## Quick Start

```bash
# Build
go build -o mindmap

# Run
./mindmap
```

## Keyboard Controls

### Navigation
- **Arrow Keys** (←↑↓→): Select nearest node in that direction (spatial navigation)
- **WASD** or **hjkl**: Pan the camera view
- **[** / **]**: Cycle through nodes sequentially

### Node Creation
- **Tab**: Create child node (next level, positioned to the right)
- **Enter**: Create sibling node (same level, positioned below)
  - Note: At root node, both Tab and Enter create children

### Node Editing
- **e**: Edit selected node text
- **x** or **Delete**: Delete selected node (cannot delete root)

### View Controls
- **+** / **=**: Zoom in
- **-** / **_**: Zoom out
- **0**: Reset camera to origin
- **c**: Center camera on selected node

### Connections
- **L**: Create manual link between nodes (select source, then target)

### File Operations
- **Ctrl+S**: Save to `mindmap.json`
- **Ctrl+O**: Load from `mindmap.json`

### Help & Exit
- **?**: Show help message in status bar
- **q** or **Ctrl+C**: Quit application

## Visual Indicators

- **▶** arrow: Shows currently selected node
- **Rounded corners** (╭╮╰╯): Selected node borders
- **Square corners** (┌┐└┘): Unselected node borders
- **Colors**: Each root child gets a unique color; descendants inherit it

## Project Structure

```
.
├── main.go           # Entry point, initializes Bubble Tea program
├── model.go          # Core data model and state management
├── node.go           # Node and Edge data structures
├── camera.go         # Viewport and coordinate transformation
├── update.go         # Input handling and state updates
├── renderer.go       # Canvas rendering and visual output
├── persistence.go    # JSON save/load functionality
└── README.md         # This file
```

## Architecture Overview

### Data Model (`model.go`)

**Core Structures:**
- `Model`: Main Bubble Tea model containing all state
  - `Nodes`: Map of node ID → Node
  - `Edges`: Slice of connections between nodes
  - `Camera`: Viewport position and zoom
  - `Selected`: Currently selected node ID
  - `Mode`: Current interaction mode (Normal/Edit/Link)
  - `ColorPalette`: Colors for root children branches

**Key Functions:**
- `AddChildNode(text)`: Creates child to the right, inherits/assigns color
- `AddSiblingNode(text)`: Creates sibling below, same color as current
- `pushDownNodesBelow(y, amount)`: Shifts all nodes below Y downward
- `GetChildrenOf(parentID)`: Returns all direct children of a node
- `DeleteNode(id)`: Removes node and associated edges

### Node System (`node.go`)

**Node Structure:**
```go
type Node struct {
    ID       string   // Unique identifier
    Text     string   // Display text
    X, Y     float64  // World coordinates
    Width    int      // Calculated from text
    Height   int      // Calculated from text
    ParentID string   // Parent node ID for hierarchy
    Color    string   // Branch color (hex)
    Links    []string // Connected node IDs
}
```

**Edge Structure:**
```go
type Edge struct {
    FromID string
    ToID   string
}
```

### Camera System (`camera.go`)

**Camera Structure:**
```go
type Camera struct {
    X    float64 // Position in world space
    Y    float64
    Zoom float64 // 1.0 = normal, 0.5 = zoomed out, 2.0 = zoomed in
}
```

**Key Functions:**
- `WorldToScreen(wx, wy, screenW, screenH)`: Converts world coords → screen coords
- `ScreenToWorld(sx, sy, screenW, screenH)`: Converts screen coords → world coords
- `Pan(dx, dy)`: Moves camera
- `ZoomIn()` / `ZoomOut()`: Adjusts zoom level (clamped 0.25-4.0)

### Input Handling (`update.go`)

**Modes:**
- `ModeNormal`: Navigation and node manipulation
- `ModeEdit`: Text input for creating/editing nodes
- `ModeLink`: Creating connections between nodes

**Key Functions:**
- `handleNormalMode(msg)`: Processes navigation and commands
- `handleEditMode(msg)`: Processes text input
- `handleLinkMode(msg)`: Processes link creation
- `selectNodeInDirection(dx, dy)`: Smart spatial navigation with alignment priority

**Spatial Navigation Algorithm:**
- Prioritizes nodes that are visually aligned (2x weight)
- Score = `secondaryDistance * 2.0 + primaryDistance`
- Lower score = better match

### Rendering System (`renderer.go`)

**ColoredCell Structure:**
```go
type ColoredCell struct {
    Char  rune   // Character to display
    Color string // Hex color code
}
```

**Rendering Pipeline:**
1. Create 2D grid of ColoredCells
2. Draw edges first (behind nodes)
3. Draw nodes on top
4. Convert grid to colored string output
5. Add status bar

**Key Functions:**
- `drawNode(grid, node, isSelected)`: Renders node box with text
- `drawEdge(grid, from, to)`: Draws line connecting borders
- `drawLine(grid, x1, y1, x2, y2, color)`: Bresenham's line algorithm

**Border Connection Logic:**
- Horizontal: Right edge → Left edge
- Vertical: Bottom edge → Top edge
- Uses node centers for alignment

### Persistence (`persistence.go`)

**File Format (JSON):**
```json
{
  "nodes": [
    {
      "id": "0",
      "text": "Root Idea",
      "x": 0,
      "y": 0,
      "width": 15,
      "height": 3,
      "parent_id": "",
      "color": "",
      "links": []
    }
  ],
  "edges": [
    {"from": "0", "to": "1"}
  ],
  "camera": {
    "x": 0,
    "y": 0,
    "zoom": 1.0
  }
}
```

## Color System

**Palette (8 colors):**
1. `#FF6B6B` - Red
2. `#4ECDC4` - Cyan
3. `#45B7D1` - Blue
4. `#FFA07A` - Light Salmon
5. `#98D8C8` - Mint
6. `#F7DC6F` - Yellow
7. `#BB8FCE` - Purple
8. `#85C1E2` - Sky Blue

**Color Assignment:**
- Root node: No color (default)
- Root's children: Assigned next color from palette (cycles after 8)
- Descendants: Inherit parent's color
- Siblings: Share same color (same parent)

## Automatic Layout System

**Problem:** Adding nodes can cause overlaps with nodes below.

**Solution:** `pushDownNodesBelow(thresholdY, amount)`
- When adding sibling: Push all nodes with `Y >= newNodeY` down
- When adding child (with siblings): Push all nodes with `Y >= newNodeY` down
- Amount = new node height + vertical spacing

**Triggered By:**
- `AddSiblingNode()`: Always pushes down
- `AddChildNode()`: Only if parent already has children

## Node Positioning Rules

### Child Nodes (Tab)
- **X**: `parent.X + parent.Width + 5` (horizontal spacing)
- **Y**:
  - First child: `parent.Y` (aligned with parent)
  - Subsequent: `lowestChild.Y + lowestChild.Height + 3`

### Sibling Nodes (Enter)
- **X**: `sibling.X` (same horizontal position)
- **Y**: `sibling.Y + sibling.Height + 3` (vertical spacing)

## Development Notes

### Coordinate Systems
- **World coordinates**: Float64, infinite canvas
- **Screen coordinates**: Integer, terminal grid
- Camera transforms between them

### Node ID System
- Root node: Always `"0"`
- New nodes: Sequential integers as strings (`"1"`, `"2"`, etc.)
- `NextID` counter tracks next available ID

### Grid Rendering
- Grid is `Width × (Height-1)` (last line reserved for status bar)
- Origin is top-left
- Positive X = right, Positive Y = down

### Edge Cases Handled
- Root cannot be deleted
- Root cannot have siblings (both Enter/Tab create children)
- No selection: Create child at camera center
- Zoom limits: 0.25x to 4.0x
- Text truncation in small nodes
- Nodes outside viewport rendered as dots

## Future Enhancements (Not Implemented)

### Potential Features
- [ ] Outline/tree view toggle
- [ ] Node collapsing/expanding
- [ ] Search/filter nodes
- [ ] Export to various formats (PNG, SVG, Markdown)
- [ ] Themes and custom color palettes
- [ ] Undo/redo
- [ ] Multi-line text input
- [ ] Node tags and metadata
- [ ] Curved connection lines
- [ ] Mouse support
- [ ] Multiple files/tabs
- [ ] Node icons/emojis
- [ ] Auto-save
- [ ] Git integration

### Known Limitations
- Single file only (hardcoded `mindmap.json`)
- No undo/redo
- Text input is single-line only
- No node resizing (auto-calculated from text)
- No manual node positioning (auto-layout only)
- Color cycling after 8 root children (repeats colors)

## Dependencies

```
github.com/charmbracelet/bubbletea  # TUI framework
github.com/charmbracelet/lipgloss   # Styling and colors
```

## Building and Running

```bash
# Install dependencies
go mod download

# Build
go build -o mindmap

# Run
./mindmap

# Build and run
go run .
```

## License

Open source - do whatever you want with it!

## Tips & Tricks

1. **Quick Navigation**: Use arrow keys for intuitive spatial movement
2. **Camera Following**: Press `c` to center on selected node if lost
3. **Color Branches**: Create multiple root children first for distinct branches
4. **Fast Input**: Type quickly in edit mode - no need to wait
5. **Zoom for Overview**: Zoom out (`-`) to see entire map structure
6. **Save Often**: Ctrl+S frequently to avoid losing work

## Troubleshooting

**Problem**: Colors not showing
- Ensure terminal supports 24-bit color (most modern terminals do)
- Try a different terminal emulator

**Problem**: Weird characters in borders
- Terminal needs Unicode support
- Try modern terminal (kitty, alacritty, iTerm2, Windows Terminal)

**Problem**: Layout feels cramped
- Adjust spacing constants in `model.go`:
  - `spacing` (horizontal): Default 5.0
  - `verticalSpacing`: Default 3.0

**Problem**: Keyboard not responding
- Check if terminal is capturing keys (some multiplexers intercept)
- Restart application

---

**Built with ❤️ using Go and Bubble Tea**
