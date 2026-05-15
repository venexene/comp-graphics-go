package ui

import (
	"fmt"
)

// LightingControlState holds UI state
type LightingControlState struct {
	ShowPanel bool
}

var state = LightingControlState{ShowPanel: true}

// InitializeUI initializes the UI system
func InitializeUI(window interface{}) error {
	// UI is rendered as text overlay via window title
	return nil
}

// BeginFrame starts a new UI frame
func BeginFrame() {
	// No-op for text-based UI
}

// EndFrame ends the UI frame
func EndFrame() {
	// No-op for text-based UI
}

// ShowLightingControlPanel displays the lighting control UI
func ShowLightingControlPanel(
	getLightingName func() string,
	cycleLighting func(forward bool),
	toggleShading func(),
	getShadingMode func() string,
	getLinearCoef func() float32,
	setLinearCoef func(float32),
	getQuadraticCoef func() float32,
	setQuadraticCoef func(float32),
	getAmbientStrength func() float32,
	setAmbientStrength func(float32),
) {
	// Store references for rendering in overlay
	_ = getLightingName
	_ = cycleLighting
	_ = toggleShading
	_ = getShadingMode
	_ = getLinearCoef
	_ = setLinearCoef
	_ = getQuadraticCoef
	_ = setQuadraticCoef
	_ = getAmbientStrength
	_ = setAmbientStrength
}

// Cleanup releases UI resources
func Cleanup() {
	// No-op for text-based UI
}

// GetUIOverlayText returns the text to render as UI overlay
func GetUIOverlayText(
	getLightingName func() string,
	getShadingMode func() string,
	getLinearCoef func() float32,
	getQuadraticCoef func() float32,
	getAmbientStrength func() float32,
) string {
	text := fmt.Sprintf(`
╔═══════════════════════════════════╗
║     LIGHTING CONTROLS UI          ║
╠═══════════════════════════════════╣
║ Model: %-27s ║
║ Shading: %-25s ║
╠═══════════════════════════════════╣
║ ATTENUATION PARAMETERS            ║
║   Linear Falloff (Z/X):    %.2f   ║
║   Quadratic Falloff (C/V): %.2f   ║
╠═══════════════════════════════════╣
║ LIGHT POWER                       ║
║   Ambient Strength (B/N):  %.2f   ║
╠═══════════════════════════════════╣
║ KEYBOARD SHORTCUTS                ║
║   T/G     - Switch Lighting Model ║
║   Y       - Toggle Gouraud/Phong  ║
║   Arrows  - Rotate Camera         ║
║   WASD    - Pan Camera            ║
║   +/-     - Zoom In/Out           ║
║   Q/E     - Scale Object          ║
║   IJKL    - Move Object XY        ║
║   U/O     - Move Object Z         ║
║   Ctrl+R  - Reset All             ║
╚═══════════════════════════════════╝
`,
		getLightingName(),
		getShadingMode(),
		getLinearCoef(),
		getQuadraticCoef(),
		getAmbientStrength(),
	)
	return text
}
