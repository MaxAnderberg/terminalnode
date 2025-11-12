package main

import (
	"fmt"
	"strings"
)

// Node represents a single node in the mind map
type Node struct {
	ID       string   `json:"id"`
	Text     string   `json:"text"`
	X        float64  `json:"x"`
	Y        float64  `json:"y"`
	Width    int      `json:"width"`
	Height   int      `json:"height"`
	ParentID string   `json:"parent_id"` // ID of parent node
	Color    string   `json:"color"`     // Color for this branch
	Links    []string `json:"links"`     // IDs of connected nodes
}

// NewNode creates a new node at the given position
func NewNode(id, text string, x, y float64) *Node {
	width, height := calculateNodeSize(text)
	return &Node{
		ID:     id,
		Text:   text,
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
		Links:  make([]string, 0),
	}
}

// calculateNodeSize returns the width and height needed for a node's text
func calculateNodeSize(text string) (int, int) {
	lines := strings.Split(text, "\n")
	height := len(lines) + 2 // +2 for borders
	width := 0
	for _, line := range lines {
		if len(line) > width {
			width = len(line)
		}
	}
	width += 4 // +4 for borders and padding
	if width < 10 {
		width = 10 // Minimum width
	}
	return width, height
}

// Edge represents a connection between two nodes
type Edge struct {
	FromID string `json:"from"`
	ToID   string `json:"to"`
}

// GetCenter returns the center point of the node
func (n *Node) GetCenter() (float64, float64) {
	return n.X + float64(n.Width)/2, n.Y + float64(n.Height)/2
}

// UpdateSize recalculates the node's size based on its text
func (n *Node) UpdateSize() {
	n.Width, n.Height = calculateNodeSize(n.Text)
}

// String returns a string representation of the node
func (n *Node) String() string {
	return fmt.Sprintf("Node[%s: '%s' at (%.1f, %.1f)]", n.ID, n.Text, n.X, n.Y)
}
