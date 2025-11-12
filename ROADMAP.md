# TerminalNode Roadmap

## ✅ Completed

### v0.1.0 - Initial Release
- [x] Canvas-based layout with pan and zoom
- [x] Keyboard-driven navigation with spatial awareness
- [x] Hierarchical node structure (children and siblings)
- [x] Automatic layout adjustment to prevent overlaps
- [x] Color-coded branches for visual organization
- [x] JSON persistence (save/load)
- [x] Border-to-border edge connections

### v0.1.1 - Camera & Color Fixes
- [x] Increased camera pan speed (2.5x faster)
- [x] Fixed root children color assignment bug
- [x] Added debug color display in status bar
- [x] Verified 'c' key centers camera on selected node

### v0.1.2 - Text Wrapping & UX Polish
- [x] Multi-line text support with automatic word wrapping
- [x] Smart word-boundary wrapping (22 chars ~4-5 words per line)
- [x] Dynamic node height adjustment based on wrapped content

### v0.1.3 - Visual Polish
- [x] Smooth Bezier curve connections between nodes
  - Cubic Bezier interpolation with intelligent control points
  - Diagonal line characters (╱╲) for smooth curve rendering
  - Maintains border-to-border connection points
- [x] Enhanced node box design
  - Selected nodes: bold double-line borders (━┃┏┓┗┛)
  - Unselected nodes: rounded corners (╭╮╰╯) for clean look
  - Improved internal padding (1 space on each side)
- [x] Polished status bar styling
  - Modern hex color scheme
  - Better visual hierarchy and contrast
  - Color-coded mode indicators (green/orange/pink)

---

## 🚧 In Progress

Nothing currently in progress.

---

## 📋 High Priority (Next Up)

### UX/UI Improvements
- [ ] **Search/Filter** - Quick win for navigation
  - Fuzzy search through node text
  - Highlight matching nodes
  - Jump to search results

- [ ] **Undo/Redo System** - Essential for production use
  - Track state history
  - Ctrl+Z / Ctrl+Y keybinds
  - Limited history (last 50 actions)

---

## 📌 Medium Priority

### Core Features
- [ ] **Node Collapsing/Expanding**
  - Toggle subtree visibility
  - Show collapse indicator (▶/▼)
  - Adjust layout when collapsing

- [ ] **Export Functionality**
  - Markdown export (easiest first)
  - JSON export (already have)
  - PNG/SVG export (future)

- [ ] **Auto-save**
  - Periodic save timer (every 30s)
  - Save on quit
  - Backup/recovery

### UX Polish
- [ ] **Better Color System**
  - Verify colors display correctly on all terminals
  - Custom palette support
  - Color picker for branches
  - Option to disable colors

- [ ] **Enhanced Status Bar Features**
  - Show more context (current path/breadcrumbs)
  - Shortcuts reminder

- [ ] **Visual Feedback**
  - Flash on node creation
  - Smooth transitions
  - Animation on delete

- [ ] **Welcome Screen**
  - Tutorial on first launch
  - Quick reference card
  - Sample mind map

---

## 🔮 Low Priority (Nice to Have)

### Advanced Features
- [ ] **Multiple Files/Tabs**
  - Switch between different mind maps
  - Tab bar UI
  - Recent files list

- [ ] **Themes**
  - Dark/light mode toggle
  - Custom color schemes
  - Preset themes

- [ ] **Node Icons/Emojis**
  - Visual categorization
  - Emoji picker
  - Custom icons

- [ ] **Tags & Metadata**
  - Add tags to nodes
  - Filter by tags
  - Metadata panel

- [ ] **Git Integration**
  - Track changes
  - Commit mind maps
  - Diff visualization

---

## 💭 Someday/Maybe (Parked Ideas)

### Mouse Support
- Click to select nodes
- Drag to move nodes
- Scroll to pan
- **Note**: Parked to focus on keyboard-first design

### Manual Node Positioning
- Drag nodes with keyboard (Shift+arrows)
- Grid snapping
- Alignment tools
- **Note**: Auto-layout is working well, may not be needed

### Advanced Connections
- Curved connection lines (more sophisticated)
- Connection labels
- Different line styles (dotted, dashed)
- Bidirectional arrows

### Collaboration
- Multi-user editing
- Share via URL
- Real-time sync
- **Note**: Far future, complex infrastructure

---

## 🐛 Known Issues

### To Fix
- None currently tracked

### Limitations
- Single file only (hardcoded `mindmap.json`)
- Text input is single-line only
- No undo/redo yet
- No node resizing (auto-calculated)
- Colors may not work on all terminals

---

## 🎯 Next Milestone: v0.2.0

**Focus**: Core UX improvements for production use

**Target Features**:
1. Search/Filter functionality
2. Undo/Redo system
3. Node collapsing/expanding
4. Export functionality (Markdown)
5. Auto-save

**Goal**: Make the tool production-ready for daily use with essential workflow features.

---

## 📝 Notes

- Keep keyboard-first philosophy
- Prioritize speed and responsiveness
- Maintain minimal/clean aesthetic
- Focus on core mind mapping workflow
- Git-friendly file format (JSON)

---

**Last Updated**: 2025-11-12 (v0.1.3 - Visual polish implemented)
