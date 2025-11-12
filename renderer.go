package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ColoredCell holds a character and its color
type ColoredCell struct {
	Char  rune
	Color string
}

// View renders the mind map
func (m Model) View() string {
	if m.Width == 0 || m.Height == 0 {
		return ""
	}

	// Create a 2D grid for rendering with color information
	grid := make([][]ColoredCell, m.Height-1) // -1 for status bar
	for i := range grid {
		grid[i] = make([]ColoredCell, m.Width)
		for j := range grid[i] {
			grid[i][j] = ColoredCell{Char: ' ', Color: ""}
		}
	}

	// Draw edges first (so they appear behind nodes)
	m.drawEdges(grid)

	// Draw nodes
	m.drawNodes(grid)

	// Convert grid to string with colors
	var sb strings.Builder
	for _, row := range grid {
		for _, cell := range row {
			if cell.Color != "" {
				// Apply color using lipgloss
				style := lipgloss.NewStyle().Foreground(lipgloss.Color(cell.Color))
				sb.WriteString(style.Render(string(cell.Char)))
			} else {
				sb.WriteRune(cell.Char)
			}
		}
		sb.WriteRune('\n')
	}

	// Add status bar
	statusBar := m.renderStatusBar()
	sb.WriteString(statusBar)

	return sb.String()
}

// drawNodes renders all nodes onto the grid
func (m Model) drawNodes(grid [][]ColoredCell) {
	for id, node := range m.Nodes {
		m.drawNode(grid, node, id == m.Selected)
	}
}

// drawNode renders a single node onto the grid
func (m Model) drawNode(grid [][]ColoredCell, node *Node, isSelected bool) {
	// Convert world coordinates to screen coordinates
	sx, sy := m.Camera.WorldToScreen(node.X, node.Y, m.Width, m.Height-1)

	// Check if node is visible
	if sy >= len(grid) || sy < 0 {
		return
	}

	// Apply zoom to size
	width := int(float64(node.Width) * m.Camera.Zoom)
	height := int(float64(node.Height) * m.Camera.Zoom)

	// Don't render if too small
	if width < 3 || height < 2 {
		// Just draw a point
		if sy >= 0 && sy < len(grid) && sx >= 0 && sx < len(grid[0]) {
			grid[sy][sx] = ColoredCell{Char: '●', Color: node.Color}
		}
		return
	}

	// Get border runes based on selection
	// Selected nodes use rounded double-line borders for emphasis
	// Unselected nodes use single-line rounded corners for clean look
	var top, bottom, left, right, topLeft, topRight, bottomLeft, bottomRight rune
	if isSelected {
		top, bottom, left, right = '━', '━', '┃', '┃'
		topLeft, topRight, bottomLeft, bottomRight = '┏', '┓', '┗', '┛'
	} else {
		top, bottom, left, right = '─', '─', '│', '│'
		topLeft, topRight, bottomLeft, bottomRight = '╭', '╮', '╰', '╯'
	}

	// Add selection indicator
	if isSelected && sy >= 0 && sy < len(grid) && sx-2 >= 0 && sx-2 < len(grid[0]) {
		grid[sy][sx-2] = ColoredCell{Char: '▶', Color: node.Color}
	}

	// Draw top border
	if sy >= 0 && sy < len(grid) {
		if sx >= 0 && sx < len(grid[0]) {
			grid[sy][sx] = ColoredCell{Char: topLeft, Color: node.Color}
		}
		for x := sx + 1; x < sx+width-1 && x < len(grid[0]); x++ {
			if x >= 0 {
				grid[sy][x] = ColoredCell{Char: top, Color: node.Color}
			}
		}
		if sx+width-1 >= 0 && sx+width-1 < len(grid[0]) {
			grid[sy][sx+width-1] = ColoredCell{Char: topRight, Color: node.Color}
		}
	}

	// Draw middle (text with improved padding)
	// Use the same wrapping logic as calculateNodeSize
	const maxTextWidth = 22
	lines := wrapText(node.Text, maxTextWidth)
	for i := 1; i < height-1; i++ {
		y := sy + i
		if y < 0 || y >= len(grid) {
			continue
		}

		// Left border
		if sx >= 0 && sx < len(grid[0]) {
			grid[y][sx] = ColoredCell{Char: left, Color: node.Color}
		}

		// Add left padding space
		if sx+1 >= 0 && sx+1 < len(grid[0]) {
			grid[y][sx+1] = ColoredCell{Char: ' ', Color: ""}
		}

		// Text content
		lineIdx := i - 1
		if lineIdx < len(lines) {
			text := lines[lineIdx]
			maxRenderWidth := width - 4 // Account for borders and padding (2 spaces)
			if len(text) > maxRenderWidth {
				text = text[:maxRenderWidth]
			}

			for j, ch := range text {
				x := sx + j + 2 // +2 for border and left padding
				if x >= 0 && x < len(grid[0]) {
					grid[y][x] = ColoredCell{Char: ch, Color: node.Color}
				}
			}
		}

		// Add right padding space
		if sx+width-2 >= 0 && sx+width-2 < len(grid[0]) {
			grid[y][sx+width-2] = ColoredCell{Char: ' ', Color: ""}
		}

		// Right border
		if sx+width-1 >= 0 && sx+width-1 < len(grid[0]) {
			grid[y][sx+width-1] = ColoredCell{Char: right, Color: node.Color}
		}
	}

	// Draw bottom border
	if sy+height-1 >= 0 && sy+height-1 < len(grid) {
		if sx >= 0 && sx < len(grid[0]) {
			grid[sy+height-1][sx] = ColoredCell{Char: bottomLeft, Color: node.Color}
		}
		for x := sx + 1; x < sx+width-1 && x < len(grid[0]); x++ {
			if x >= 0 {
				grid[sy+height-1][x] = ColoredCell{Char: bottom, Color: node.Color}
			}
		}
		if sx+width-1 >= 0 && sx+width-1 < len(grid[0]) {
			grid[sy+height-1][sx+width-1] = ColoredCell{Char: bottomRight, Color: node.Color}
		}
	}
}

// drawEdges renders all edges onto the grid
func (m Model) drawEdges(grid [][]ColoredCell) {
	for _, edge := range m.Edges {
		fromNode := m.Nodes[edge.FromID]
		toNode := m.Nodes[edge.ToID]
		if fromNode != nil && toNode != nil {
			m.drawEdge(grid, fromNode, toNode)
		}
	}
}

// drawEdge draws a line between two nodes, connecting at their borders
func (m Model) drawEdge(grid [][]ColoredCell, from, to *Node) {
	// Get center points to determine direction
	fromCX, fromCY := from.GetCenter()
	toCX, toCY := to.GetCenter()

	var fx, fy, tx, ty float64

	// Determine connection points based on relative positions
	// Horizontal connections (most common)
	if toCX > fromCX { // "to" is to the right of "from"
		// Connect from right edge of "from" to left edge of "to"
		fx = from.X + float64(from.Width)
		fy = fromCY
		tx = to.X
		ty = toCY
	} else if toCX < fromCX { // "to" is to the left of "from"
		// Connect from left edge of "from" to right edge of "to"
		fx = from.X
		fy = fromCY
		tx = to.X + float64(to.Width)
		ty = toCY
	} else { // Vertically aligned
		if toCY > fromCY { // "to" is below "from"
			// Connect from bottom of "from" to top of "to"
			fx = fromCX
			fy = from.Y + float64(from.Height)
			tx = toCX
			ty = to.Y
		} else { // "to" is above "from"
			// Connect from top of "from" to bottom of "to"
			fx = fromCX
			fy = from.Y
			tx = toCX
			ty = to.Y + float64(to.Height)
		}
	}

	// Convert to screen coordinates
	sx1, sy1 := m.Camera.WorldToScreen(fx, fy, m.Width, m.Height-1)
	sx2, sy2 := m.Camera.WorldToScreen(tx, ty, m.Width, m.Height-1)

	// Draw line using Bresenham's algorithm with the "to" node's color
	m.drawLine(grid, sx1, sy1, sx2, sy2, to.Color)
}

// drawLine draws a smooth Bezier curve between two points
func (m Model) drawLine(grid [][]ColoredCell, x1, y1, x2, y2 int, color string) {
	// Calculate control points for cubic Bezier curve
	// Place control points horizontally offset for smooth horizontal connections
	dx := float64(x2 - x1)
	dy := float64(y2 - y1)

	// Adjust control point distance based on the distance between nodes
	dist := math.Sqrt(dx*dx + dy*dy)
	cpOffset := math.Min(dist*0.4, 30.0) // 40% of distance, max 30 units

	// Control points for horizontal flow
	cp1x := float64(x1) + cpOffset
	cp1y := float64(y1)
	cp2x := float64(x2) - cpOffset
	cp2y := float64(y2)

	// If connection is more vertical than horizontal, adjust control points vertically
	if math.Abs(dy) > math.Abs(dx) {
		cp1x = float64(x1)
		cp1y = float64(y1) + cpOffset*math.Copysign(1, dy)
		cp2x = float64(x2)
		cp2y = float64(y2) - cpOffset*math.Copysign(1, dy)
	}

	// Draw the Bezier curve using parametric equation
	// Sample enough points for smooth rendering
	steps := int(dist * 2) // Ensure we have enough resolution
	if steps < 10 {
		steps = 10
	}

	prevX, prevY := x1, y1
	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)

		// Cubic Bezier formula: B(t) = (1-t)³P0 + 3(1-t)²tP1 + 3(1-t)t²P2 + t³P3
		omt := 1 - t
		omt2 := omt * omt
		omt3 := omt2 * omt
		t2 := t * t
		t3 := t2 * t

		x := omt3*float64(x1) + 3*omt2*t*cp1x + 3*omt*t2*cp2x + t3*float64(x2)
		y := omt3*float64(y1) + 3*omt2*t*cp1y + 3*omt*t2*cp2y + t3*float64(y2)

		curX, curY := int(math.Round(x)), int(math.Round(y))

		// Draw line segment from previous point to current point
		m.drawLineSegment(grid, prevX, prevY, curX, curY, color)

		prevX, prevY = curX, curY
	}
}

// drawLineSegment draws a small line segment and picks the best character for direction
func (m Model) drawLineSegment(grid [][]ColoredCell, x1, y1, x2, y2 int, color string) {
	dx := x2 - x1
	dy := y2 - y1

	// Plot start point
	if y1 >= 0 && y1 < len(grid) && x1 >= 0 && x1 < len(grid[0]) {
		if grid[y1][x1].Char == ' ' {
			lineChar := m.getLineChar(dx, dy)
			grid[y1][x1] = ColoredCell{Char: lineChar, Color: color}
		}
	}

	// If points are the same, we're done
	if x1 == x2 && y1 == y2 {
		return
	}

	// Use Bresenham to fill in the segment
	absDx := abs(dx)
	absDy := abs(dy)

	sx := -1
	if x1 < x2 {
		sx = 1
	}
	sy := -1
	if y1 < y2 {
		sy = 1
	}

	err := absDx - absDy

	for {
		if x1 == x2 && y1 == y2 {
			break
		}

		e2 := 2 * err
		if e2 > -absDy {
			err -= absDy
			x1 += sx
		}
		if e2 < absDx {
			err += absDx
			y1 += sy
		}

		// Plot point if within bounds
		if y1 >= 0 && y1 < len(grid) && x1 >= 0 && x1 < len(grid[0]) {
			if grid[y1][x1].Char == ' ' {
				lineChar := m.getLineChar(dx, dy)
				grid[y1][x1] = ColoredCell{Char: lineChar, Color: color}
			}
		}
	}
}

// getLineChar returns the best Unicode box-drawing character for a given direction
func (m Model) getLineChar(dx, dy int) rune {
	// Determine angle and pick appropriate character
	if dx == 0 && dy == 0 {
		return '·'
	}

	// Calculate approximate angle
	absDx := abs(dx)
	absDy := abs(dy)

	// Mostly horizontal
	if absDx > absDy*2 {
		return '─'
	}
	// Mostly vertical
	if absDy > absDx*2 {
		return '│'
	}

	// Diagonal
	if (dx > 0 && dy < 0) || (dx < 0 && dy > 0) {
		return '╱'
	}
	return '╲'
}

// renderStatusBar creates the status bar at the bottom
func (m Model) renderStatusBar() string {
	var modeStr string
	switch m.Mode {
	case ModeNormal:
		modeStr = "NORMAL"
	case ModeEdit:
		modeStr = fmt.Sprintf("EDIT: %s_", m.EditBuffer)
	case ModeLink:
		modeStr = fmt.Sprintf("LINK: %s → ?", m.LinkSourceID)
	}

	left := fmt.Sprintf(" %s ", modeStr)
	middle := m.StatusMsg

	// Debug: show selected node color
	debugColor := ""
	if selectedNode := m.GetSelectedNode(); selectedNode != nil {
		debugColor = fmt.Sprintf(" | Color: %s", selectedNode.Color)
	}

	right := fmt.Sprintf(" Nodes: %d | Zoom: %.1fx | Pos: (%.0f, %.0f)%s | ?: help ",
		len(m.Nodes), m.Camera.Zoom, m.Camera.X, m.Camera.Y, debugColor)

	// Calculate spacing
	totalWidth := m.Width
	usedWidth := lipgloss.Width(left) + lipgloss.Width(middle) + lipgloss.Width(right)
	spacing := ""
	if usedWidth < totalWidth {
		spacing = strings.Repeat(" ", totalWidth-usedWidth)
	}

	// Style the status bar with improved visual hierarchy
	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#E0E0E0")).
		Background(lipgloss.Color("#2A2A2A"))

	modeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#000000")).
		Background(lipgloss.Color("#00D787")).
		Bold(true).
		Padding(0, 1)

	if m.Mode == ModeEdit {
		modeStyle = modeStyle.
			Background(lipgloss.Color("#FFB86C")).
			Foreground(lipgloss.Color("#000000"))
	} else if m.Mode == ModeLink {
		modeStyle = modeStyle.
			Background(lipgloss.Color("#FF79C6")).
			Foreground(lipgloss.Color("#000000"))
	}

	// Enhanced visual separation
	leftPart := modeStyle.Render(modeStr)
	middlePart := statusStyle.Render(" " + middle)
	rightPart := statusStyle.Render(right)

	return leftPart + statusStyle.Render(spacing) + middlePart + rightPart
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// distance calculates distance between two points
func distance(x1, y1, x2, y2 float64) float64 {
	dx := x2 - x1
	dy := y2 - y1
	return math.Sqrt(dx*dx + dy*dy)
}
