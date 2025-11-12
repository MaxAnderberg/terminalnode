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

// wrapText wraps text to fit within maxWidth, breaking on word boundaries
func wrapText(text string, maxWidth int) []string {
	if maxWidth < 5 {
		maxWidth = 5 // Minimum sensible width
	}

	// First split by explicit newlines
	paragraphs := strings.Split(text, "\n")
	var wrappedLines []string

	for _, paragraph := range paragraphs {
		if len(paragraph) == 0 {
			wrappedLines = append(wrappedLines, "")
			continue
		}

		// Split paragraph into words
		words := strings.Fields(paragraph)
		if len(words) == 0 {
			wrappedLines = append(wrappedLines, "")
			continue
		}

		var currentLine string
		for _, word := range words {
			// If adding this word would exceed maxWidth
			if len(currentLine) > 0 && len(currentLine)+1+len(word) > maxWidth {
				// If the word itself is longer than maxWidth, we need to break it
				if len(word) > maxWidth {
					// Add current line if not empty
					if len(currentLine) > 0 {
						wrappedLines = append(wrappedLines, currentLine)
						currentLine = ""
					}
					// Break the long word into chunks
					for len(word) > maxWidth {
						wrappedLines = append(wrappedLines, word[:maxWidth])
						word = word[maxWidth:]
					}
					currentLine = word
				} else {
					// Save current line and start new one
					wrappedLines = append(wrappedLines, currentLine)
					currentLine = word
				}
			} else {
				// Add word to current line
				if len(currentLine) > 0 {
					currentLine += " " + word
				} else {
					currentLine = word
				}
			}
		}

		// Add the last line if not empty
		if len(currentLine) > 0 {
			wrappedLines = append(wrappedLines, currentLine)
		}
	}

	if len(wrappedLines) == 0 {
		wrappedLines = append(wrappedLines, "")
	}

	return wrappedLines
}

// calculateNodeSize returns the width and height needed for a node's text
func calculateNodeSize(text string) (int, int) {
	const maxTextWidth = 22 // Roughly 4-5 words, similar to MindNode

	lines := wrapText(text, maxTextWidth)
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
