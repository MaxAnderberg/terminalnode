package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// MindMapData represents the serializable mind map data
type MindMapData struct {
	Nodes  map[string]*Node `json:"nodes"`
	Edges  []Edge           `json:"edges"`
	Camera Camera           `json:"camera"`
}

// SaveToFile saves the mind map to a JSON file
func (m *Model) SaveToFile(filename string) error {
	data := MindMapData{
		Nodes:  m.Nodes,
		Edges:  m.Edges,
		Camera: m.Camera,
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, jsonData, 0644)
}

// LoadFromFile loads the mind map from a JSON file
func (m *Model) LoadFromFile(filename string) error {
	jsonData, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var data MindMapData
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return err
	}

	m.Nodes = data.Nodes
	m.Edges = data.Edges
	m.Camera = data.Camera

	// Initialize camera targets (not serialized, so set them to current values)
	m.Camera.TargetX = m.Camera.X
	m.Camera.TargetY = m.Camera.Y
	m.Camera.TargetZoom = m.Camera.Zoom

	// Select first node if none selected
	if m.Selected == "" && len(m.Nodes) > 0 {
		for id := range m.Nodes {
			m.Selected = id
			break
		}
	}

	// Update NextID to be higher than any existing ID
	maxID := 0
	for id := range m.Nodes {
		var numID int
		if _, err := fmt.Sscanf(id, "%d", &numID); err == nil {
			if numID > maxID {
				maxID = numID
			}
		}
	}
	m.NextID = maxID + 1

	return nil
}
