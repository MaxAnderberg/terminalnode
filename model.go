package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Mode represents the current interaction mode
type Mode int

const (
	ModeNormal Mode = iota // Navigation mode
	ModeEdit               // Editing node text
	ModeLink               // Creating links between nodes
)

// Model is the Bubble Tea model for the mind map
type Model struct {
	// Mind map data
	Nodes    map[string]*Node
	Edges    []Edge
	Camera   Camera
	Selected string // Currently selected node ID

	// UI state
	Mode            Mode
	EditBuffer      string
	IsCreatingNode  bool // True when creating new node, false when editing
	IsCreatingChild bool // True for child (Tab), false for sibling (Enter)
	Width           int
	Height          int
	NextID          int
	StatusMsg       string
	LinkSourceID    string // When in link mode, the source node

	// Colors
	ColorPalette   []string
	NextColorIndex int

	// Styles
	normalStyle   lipgloss.Style
	selectedStyle lipgloss.Style
	statusStyle   lipgloss.Style
}

// NewModel creates a new mind map model
func NewModel() Model {
	nodes := make(map[string]*Node)

	// Create initial node at center
	initialNode := NewNode("0", "Root Idea", 0, 0)
	nodes["0"] = initialNode

	return Model{
		Nodes:    nodes,
		Edges:    make([]Edge, 0),
		Camera:   NewCamera(),
		Selected: "0",
		Mode:     ModeNormal,
		NextID:   1,
		Width:    80,
		Height:   24,

		// Color palette for root children branches
		ColorPalette: []string{
			"#FF6B6B", // Red
			"#4ECDC4", // Cyan
			"#45B7D1", // Blue
			"#FFA07A", // Light Salmon
			"#98D8C8", // Mint
			"#F7DC6F", // Yellow
			"#BB8FCE", // Purple
			"#85C1E2", // Sky Blue
		},
		NextColorIndex: 0,

		normalStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1),

		selectedStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("86")).
			Padding(0, 1).
			Bold(true),

		statusStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")),
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// GetSelectedNode returns the currently selected node
func (m *Model) GetSelectedNode() *Node {
	if m.Selected != "" {
		return m.Nodes[m.Selected]
	}
	return nil
}

// GetChildrenOf returns all children of a given parent node
func (m *Model) GetChildrenOf(parentID string) []*Node {
	children := make([]*Node, 0)
	for _, node := range m.Nodes {
		if node.ParentID == parentID {
			children = append(children, node)
		}
	}
	return children
}

// AddChildNode creates a new child node to the right of the selected node
func (m *Model) AddChildNode(text string) {
	id := fmt.Sprintf("%d", m.NextID)
	m.NextID++

	var x, y float64
	var parentID string

	// Position new node to the right of selected node
	if selectedNode := m.GetSelectedNode(); selectedNode != nil {
		spacing := 5.0          // Horizontal spacing
		verticalSpacing := 3.0  // Vertical spacing between children

		x = selectedNode.X + float64(selectedNode.Width) + spacing
		parentID = selectedNode.ID

		// Find existing children of this parent and position below them
		existingChildren := m.GetChildrenOf(selectedNode.ID)
		if len(existingChildren) > 0 {
			// Find the lowest child and position below it
			lowestY := selectedNode.Y
			lowestHeight := selectedNode.Height
			for _, child := range existingChildren {
				childBottom := child.Y + float64(child.Height)
				if childBottom > lowestY + float64(lowestHeight) {
					lowestY = child.Y
					lowestHeight = child.Height
				}
			}
			y = lowestY + float64(lowestHeight) + verticalSpacing

			// Push down nodes below this position
			newNodeHeight := 3 // Default height
			spaceNeeded := float64(newNodeHeight) + verticalSpacing
			m.pushDownNodesBelow(y, spaceNeeded)
		} else {
			// First child, align with parent
			y = selectedNode.Y
		}
	} else {
		// Fallback to camera center if no selected node
		x, y = m.Camera.GetViewportCenter()
	}

	node := NewNode(id, text, x, y)
	node.ParentID = parentID

	// Assign color based on parent
	if parentID == "0" {
		// Child of root: assign next color from palette
		node.Color = m.ColorPalette[m.NextColorIndex%len(m.ColorPalette)]
		m.NextColorIndex++
	} else if parentID != "" {
		// Inherit parent's color
		if parent := m.Nodes[parentID]; parent != nil {
			node.Color = parent.Color
		}
	}

	m.Nodes[id] = node

	// Automatically create edge from parent to new node
	if parentID != "" {
		m.AddEdge(parentID, id)
	}

	m.Selected = id
	m.StatusMsg = fmt.Sprintf("Created child node %s", id)
}

// AddSiblingNode creates a new sibling node below the selected node
func (m *Model) AddSiblingNode(text string) {
	selectedNode := m.GetSelectedNode()
	if selectedNode == nil {
		// If no selection, create a child instead
		m.AddChildNode(text)
		return
	}

	// Root node can't have siblings - create child instead
	if selectedNode.ID == "0" {
		m.AddChildNode(text)
		return
	}

	id := fmt.Sprintf("%d", m.NextID)
	m.NextID++

	// Position at same X as selected node, but below it
	verticalSpacing := 3.0
	x := selectedNode.X
	newNodeHeight := 3 // Default height for new node
	y := selectedNode.Y + float64(selectedNode.Height) + verticalSpacing

	// Calculate how much space the new node will take
	spaceNeeded := float64(newNodeHeight) + verticalSpacing

	// Push down all nodes that are below this Y position
	m.pushDownNodesBelow(y, spaceNeeded)

	node := NewNode(id, text, x, y)
	node.ParentID = selectedNode.ParentID // Same parent as sibling

	// Assign color based on parent
	if selectedNode.ParentID == "0" {
		// Sibling of root's child: assign NEW color from palette
		node.Color = m.ColorPalette[m.NextColorIndex%len(m.ColorPalette)]
		m.NextColorIndex++
	} else {
		// Regular sibling: inherit color
		node.Color = selectedNode.Color
	}

	m.Nodes[id] = node

	// Connect to same parent as the selected node
	if selectedNode.ParentID != "" {
		m.AddEdge(selectedNode.ParentID, id)
	}

	m.Selected = id
	m.StatusMsg = fmt.Sprintf("Created sibling node %s", id)
}

// pushDownNodesBelow moves all nodes below a certain Y position downward
func (m *Model) pushDownNodesBelow(thresholdY, amount float64) {
	for _, node := range m.Nodes {
		if node.Y >= thresholdY {
			node.Y += amount
		}
	}
}

// DeleteNode removes a node and its associated edges
func (m *Model) DeleteNode(id string) {
	if id == "0" {
		m.StatusMsg = "Cannot delete root node"
		return
	}

	delete(m.Nodes, id)

	// Remove associated edges
	newEdges := make([]Edge, 0)
	for _, edge := range m.Edges {
		if edge.FromID != id && edge.ToID != id {
			newEdges = append(newEdges, edge)
		}
	}
	m.Edges = newEdges

	// Deselect if this was selected
	if m.Selected == id {
		m.Selected = ""
		// Select first available node
		for nodeID := range m.Nodes {
			m.Selected = nodeID
			break
		}
	}

	m.StatusMsg = fmt.Sprintf("Deleted node %s", id)
}

// AddEdge creates a link between two nodes
func (m *Model) AddEdge(fromID, toID string) {
	// Check if edge already exists
	for _, edge := range m.Edges {
		if edge.FromID == fromID && edge.ToID == toID {
			m.StatusMsg = "Edge already exists"
			return
		}
	}

	m.Edges = append(m.Edges, Edge{FromID: fromID, ToID: toID})

	// Also add to node's links
	if node := m.Nodes[fromID]; node != nil {
		node.Links = append(node.Links, toID)
	}

	m.StatusMsg = fmt.Sprintf("Created link %s → %s", fromID, toID)
}

// GetNodeAt returns the node at the given screen coordinates (if any)
func (m *Model) GetNodeAt(screenX, screenY int) *Node {
	wx, wy := m.Camera.ScreenToWorld(screenX, screenY, m.Width, m.Height)

	for _, node := range m.Nodes {
		if wx >= node.X && wx <= node.X+float64(node.Width) &&
			wy >= node.Y && wy <= node.Y+float64(node.Height) {
			return node
		}
	}
	return nil
}
