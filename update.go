package main

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// tickMsg is sent on each animation frame
type tickMsg time.Time

// doTick returns a command that sends a tick message
func doTick() tea.Cmd {
	return tea.Tick(time.Second/60, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tickMsg:
		// Update camera smoothly towards target
		// smoothness: 0.2 = smooth, 0.5 = fast, adjust to preference
		m.Camera.Update(0.25)
		return m, doTick()
	}

	return m, nil
}

// handleKeyPress processes keyboard input based on current mode
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.Mode {
	case ModeNormal:
		return m.handleNormalMode(msg)
	case ModeEdit:
		return m.handleEditMode(msg)
	case ModeLink:
		return m.handleLinkMode(msg)
	}
	return m, nil
}

// handleNormalMode handles input in normal navigation mode
func (m Model) handleNormalMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	panSpeed := 5.0 / m.Camera.Zoom // Pan faster when zoomed out (increased from 2.0)

	switch msg.String() {
	// Quit
	case "ctrl+c", "q":
		return m, tea.Quit

	// Arrow keys: spatial node selection
	case "up":
		m.selectNodeInDirection(0, -1)
	case "down":
		m.selectNodeInDirection(0, 1)
	case "left":
		m.selectNodeInDirection(-1, 0)
	case "right":
		m.selectNodeInDirection(1, 0)

	// WASD/vim keys: pan camera
	case "w", "k":
		m.Camera.Pan(0, -panSpeed)
		m.StatusMsg = ""
	case "s", "j":
		m.Camera.Pan(0, panSpeed)
		m.StatusMsg = ""
	case "a", "h":
		m.Camera.Pan(-panSpeed, 0)
		m.StatusMsg = ""
	case "d", "l":
		m.Camera.Pan(panSpeed, 0)
		m.StatusMsg = ""

	// Zoom
	case "+", "=":
		m.Camera.ZoomIn()
		m.StatusMsg = ""
	case "-", "_":
		m.Camera.ZoomOut()
		m.StatusMsg = ""

	// Reset camera
	case "0":
		m.Camera = NewCamera()
		m.StatusMsg = "Camera reset"

	// Node creation - Enter for sibling, Tab for child
	case "enter":
		m.Mode = ModeEdit
		m.EditBuffer = ""
		m.IsCreatingNode = true
		m.IsCreatingChild = false
		m.StatusMsg = "New sibling: type text and press Enter"

	case "tab":
		m.Mode = ModeEdit
		m.EditBuffer = ""
		m.IsCreatingNode = true
		m.IsCreatingChild = true
		m.StatusMsg = "New child: type text and press Enter"

	// Edit selected node
	case "e":
		if node := m.GetSelectedNode(); node != nil {
			m.Mode = ModeEdit
			m.EditBuffer = node.Text
			m.IsCreatingNode = false
			m.StatusMsg = "Edit node text (ESC to cancel, Enter to save)"
		}

	// Delete selected node
	case "x", "delete", "backspace":
		if m.Selected != "" {
			m.DeleteNode(m.Selected)
		}

	// Create link
	case "L":
		if m.Selected != "" {
			m.Mode = ModeLink
			m.LinkSourceID = m.Selected
			m.StatusMsg = "Select target node (ESC to cancel)"
		}

	// Select nodes
	case "]":
		m.selectNextNode()
	case "[":
		m.selectPrevNode()

	// Center camera on selected node
	case "c":
		if node := m.GetSelectedNode(); node != nil {
			cx, cy := node.GetCenter()
			m.Camera.TargetX = cx
			m.Camera.TargetY = cy
			m.StatusMsg = "Centered on node"
		}

	// Save/Load
	case "ctrl+s":
		if err := m.SaveToFile("mindmap.json"); err != nil {
			m.StatusMsg = fmt.Sprintf("Error saving: %v", err)
		} else {
			m.StatusMsg = "Saved to mindmap.json"
		}
	case "ctrl+o":
		if err := m.LoadFromFile("mindmap.json"); err != nil {
			m.StatusMsg = fmt.Sprintf("Error loading: %v", err)
		} else {
			m.StatusMsg = "Loaded from mindmap.json"
		}

	// Help
	case "?":
		m.StatusMsg = "arrows:select wasd:pan +/-:zoom Enter:sibling Tab:child e:edit x:delete L:link c:center Ctrl+S:save q:quit"
	}

	return m, nil
}

// handleEditMode handles input when editing a node
func (m Model) handleEditMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.Mode = ModeNormal
		m.EditBuffer = ""
		m.IsCreatingNode = false
		m.StatusMsg = "Cancelled"
		return m, nil

	case "enter":
		if m.EditBuffer != "" {
			if m.IsCreatingNode {
				// Creating new node - check if child or sibling
				if m.IsCreatingChild {
					m.AddChildNode(m.EditBuffer)
				} else {
					m.AddSiblingNode(m.EditBuffer)
				}
			} else {
				// Editing existing node
				if node := m.GetSelectedNode(); node != nil {
					node.Text = m.EditBuffer
					node.UpdateSize()
					m.StatusMsg = "Node updated"
				}
			}
		}
		m.Mode = ModeNormal
		m.EditBuffer = ""
		m.IsCreatingNode = false
		m.IsCreatingChild = false
		return m, nil

	case "backspace":
		if len(m.EditBuffer) > 0 {
			m.EditBuffer = m.EditBuffer[:len(m.EditBuffer)-1]
		}

	default:
		// Add character to buffer
		if len(msg.String()) == 1 {
			m.EditBuffer += msg.String()
		}
	}

	return m, nil
}

// handleLinkMode handles input when creating a link
func (m Model) handleLinkMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.Mode = ModeNormal
		m.LinkSourceID = ""
		m.StatusMsg = "Link cancelled"
		return m, nil

	case "tab":
		m.selectNextNode()
	case "shift+tab":
		m.selectPrevNode()

	case "enter":
		if m.Selected != "" && m.LinkSourceID != "" && m.Selected != m.LinkSourceID {
			m.AddEdge(m.LinkSourceID, m.Selected)
		}
		m.Mode = ModeNormal
		m.LinkSourceID = ""
		return m, nil
	}

	return m, nil
}

// selectNextNode cycles to the next node
func (m *Model) selectNextNode() {
	if len(m.Nodes) == 0 {
		return
	}

	ids := make([]string, 0, len(m.Nodes))
	for id := range m.Nodes {
		ids = append(ids, id)
	}

	// Find current index
	currentIdx := -1
	for i, id := range ids {
		if id == m.Selected {
			currentIdx = i
			break
		}
	}

	// Select next
	nextIdx := (currentIdx + 1) % len(ids)
	m.Selected = ids[nextIdx]
	m.StatusMsg = ""
}

// selectPrevNode cycles to the previous node
func (m *Model) selectPrevNode() {
	if len(m.Nodes) == 0 {
		return
	}

	ids := make([]string, 0, len(m.Nodes))
	for id := range m.Nodes {
		ids = append(ids, id)
	}

	// Find current index
	currentIdx := -1
	for i, id := range ids {
		if id == m.Selected {
			currentIdx = i
			break
		}
	}

	// Select previous
	prevIdx := currentIdx - 1
	if prevIdx < 0 {
		prevIdx = len(ids) - 1
	}
	m.Selected = ids[prevIdx]
	m.StatusMsg = ""
}

// selectNodeInDirection selects the nearest node in the given direction using smart scoring
func (m *Model) selectNodeInDirection(dx, dy float64) {
	selectedNode := m.GetSelectedNode()
	if selectedNode == nil {
		return
	}

	// Get center of current node
	currentX, currentY := selectedNode.GetCenter()

	var bestNode *Node
	var bestScore float64 = -1

	// Find the best node in the given direction
	for _, node := range m.Nodes {
		if node.ID == m.Selected {
			continue // Skip current node
		}

		// Get center of candidate node
		nodeX, nodeY := node.GetCenter()

		// Calculate relative position
		relX := nodeX - currentX
		relY := nodeY - currentY

		// Check if node is in the right direction
		inDirection := false
		var primaryDist, secondaryDist float64

		if dx != 0 { // Horizontal movement (left/right)
			if (dx > 0 && relX > 0) || (dx < 0 && relX < 0) {
				inDirection = true
				primaryDist = absFloat(relX)     // Distance in direction of movement
				secondaryDist = absFloat(relY)   // Distance perpendicular to movement
			}
		} else if dy != 0 { // Vertical movement (up/down)
			if (dy > 0 && relY > 0) || (dy < 0 && relY < 0) {
				inDirection = true
				primaryDist = absFloat(relY)     // Distance in direction of movement
				secondaryDist = absFloat(relX)   // Distance perpendicular to movement
			}
		}

		if !inDirection {
			continue
		}

		// Score: prioritize nodes that are:
		// 1. Aligned (low secondary distance)
		// 2. Close (low primary distance)
		// Lower score is better
		score := secondaryDist*2.0 + primaryDist

		// Update best node if this has a better score
		if bestScore < 0 || score < bestScore {
			bestScore = score
			bestNode = node
		}
	}

	// Select the best node found
	if bestNode != nil {
		m.Selected = bestNode.ID
		m.StatusMsg = ""
	}
}

// absFloat returns absolute value of float64
func absFloat(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// trimString trims a string to max length with ellipsis
func trimString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// ellipsis adds ellipsis if string is too long
func ellipsis(s string, maxLen int) string {
	lines := strings.Split(s, "\n")
	if len(lines) > 0 {
		return trimString(lines[0], maxLen)
	}
	return trimString(s, maxLen)
}
