package main

import "math"

// Camera represents the viewport into the world space
type Camera struct {
	X    float64 `json:"x"`    // Camera position in world space
	Y    float64 `json:"y"`
	Zoom float64 `json:"zoom"` // 1.0 = normal, 0.5 = zoomed out, 2.0 = zoomed in
}

// NewCamera creates a new camera at the origin
func NewCamera() Camera {
	return Camera{
		X:    0,
		Y:    0,
		Zoom: 1.0,
	}
}

// WorldToScreen converts world coordinates to screen coordinates
func (c *Camera) WorldToScreen(wx, wy float64, screenWidth, screenHeight int) (int, int) {
	// Center the camera
	centerX := float64(screenWidth) / 2
	centerY := float64(screenHeight) / 2

	// Apply zoom and camera offset
	sx := (wx-c.X)*c.Zoom + centerX
	sy := (wy-c.Y)*c.Zoom + centerY

	return int(math.Round(sx)), int(math.Round(sy))
}

// ScreenToWorld converts screen coordinates to world coordinates
func (c *Camera) ScreenToWorld(sx, sy, screenWidth, screenHeight int) (float64, float64) {
	centerX := float64(screenWidth) / 2
	centerY := float64(screenHeight) / 2

	// Reverse the transformation
	wx := (float64(sx)-centerX)/c.Zoom + c.X
	wy := (float64(sy)-centerY)/c.Zoom + c.Y

	return wx, wy
}

// Pan moves the camera by the given offset
func (c *Camera) Pan(dx, dy float64) {
	c.X += dx
	c.Y += dy
}

// ZoomIn increases the zoom level
func (c *Camera) ZoomIn() {
	c.Zoom *= 1.2
	if c.Zoom > 4.0 {
		c.Zoom = 4.0
	}
}

// ZoomOut decreases the zoom level
func (c *Camera) ZoomOut() {
	c.Zoom *= 0.8
	if c.Zoom < 0.25 {
		c.Zoom = 0.25
	}
}

// GetViewportCenter returns the world coordinates of the viewport center
func (c *Camera) GetViewportCenter() (float64, float64) {
	return c.X, c.Y
}

// IsVisible checks if a point is visible in the viewport
func (c *Camera) IsVisible(wx, wy float64, screenWidth, screenHeight int) bool {
	sx, sy := c.WorldToScreen(wx, wy, screenWidth, screenHeight)
	return sx >= 0 && sx < screenWidth && sy >= 0 && sy < screenHeight
}
