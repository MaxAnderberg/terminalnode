package main

import "math"

// Camera represents the viewport into the world space
type Camera struct {
	X    float64 `json:"x"`    // Camera position in world space
	Y    float64 `json:"y"`
	Zoom float64 `json:"zoom"` // 1.0 = normal, 0.5 = zoomed out, 2.0 = zoomed in

	// Target values for smooth interpolation
	TargetX    float64 `json:"-"` // Target X position (not serialized)
	TargetY    float64 `json:"-"` // Target Y position (not serialized)
	TargetZoom float64 `json:"-"` // Target zoom level (not serialized)
}

// NewCamera creates a new camera at the origin
func NewCamera() Camera {
	return Camera{
		X:          0,
		Y:          0,
		Zoom:       1.0,
		TargetX:    0,
		TargetY:    0,
		TargetZoom: 1.0,
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

// Pan moves the camera by the given offset (sets target for smooth movement)
func (c *Camera) Pan(dx, dy float64) {
	c.TargetX += dx
	c.TargetY += dy
}

// ZoomIn increases the zoom level (sets target for smooth movement)
func (c *Camera) ZoomIn() {
	c.TargetZoom *= 1.2
	if c.TargetZoom > 4.0 {
		c.TargetZoom = 4.0
	}
}

// ZoomOut decreases the zoom level (sets target for smooth movement)
func (c *Camera) ZoomOut() {
	c.TargetZoom *= 0.8
	if c.TargetZoom < 0.25 {
		c.TargetZoom = 0.25
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

// Update smoothly interpolates the camera towards its target position and zoom
// smoothness controls how smooth the movement is (0.0-1.0, where higher = smoother but slower)
// Returns true if the camera is still moving
func (c *Camera) Update(smoothness float64) bool {
	const threshold = 0.001 // Stop interpolating when close enough

	isMoving := false

	// Interpolate X position
	if math.Abs(c.X-c.TargetX) > threshold {
		c.X += (c.TargetX - c.X) * smoothness
		isMoving = true
	} else {
		c.X = c.TargetX
	}

	// Interpolate Y position
	if math.Abs(c.Y-c.TargetY) > threshold {
		c.Y += (c.TargetY - c.Y) * smoothness
		isMoving = true
	} else {
		c.Y = c.TargetY
	}

	// Interpolate Zoom
	if math.Abs(c.Zoom-c.TargetZoom) > threshold {
		c.Zoom += (c.TargetZoom - c.Zoom) * smoothness
		isMoving = true
	} else {
		c.Zoom = c.TargetZoom
	}

	return isMoving
}
