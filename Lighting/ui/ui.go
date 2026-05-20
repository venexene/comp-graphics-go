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
	getLightPosition func() (float32, float32, float32),
	getSelectedObjectName func() string,
	getAttenuationMode func() string,
) string {
	lx, ly, lz := getLightPosition()
	text := fmt.Sprintf(`
╔═══════════════════════════════════╗
║     LIGHTING CONTROLS UI          ║
╠═══════════════════════════════════╣
║ Model: %-27s ║
║ Shading: %-25s ║
║ Selected: %-23s ║
╠═══════════════════════════════════╣
║ LIGHT POSITION                    ║
║   X: %6.2f  Y: %6.2f  Z: %6.2f   ║
╠═══════════════════════════════════╣
║ ATTENUATION                       ║
║   Mode (M):      %-15s ║
║   Linear Coef (Z/X):  %.2f   ║
║   Quad. Coef (C/V):   %.2f   ║
╠═══════════════════════════════════╣
║ LIGHT POWER                       ║
║   Ambient Strength (B/N):  %.2f   ║
╠═══════════════════════════════════╣
║ KEYBOARD SHORTCUTS                ║
║   T/G     - Switch Lighting Model ║
║   Y       - Toggle Gouraud/Phong  ║
║   M       - Cycle Attenuation     ║
║   Tab     - Select Next Object    ║
║   Arrows  - Rotate Camera         ║
║   WASD    - Pan Camera            ║
║   +/-     - Zoom In/Out           ║
║   Q/E     - Scale Object          ║
║   R/F     - Rotate Object Z       ║
║   1/2     - Rotate Object X       ║
║   3/4     - Rotate Object Y       ║
║   IJKL    - Move Object XY        ║
║   U/O     - Move Object Z         ║
║   Alt+IJKL- Move Light XY         ║
║   Alt+U/O - Move Light Z          ║
║   Ctrl+R  - Reset All             ║
╚═══════════════════════════════════╝
`,
		getLightingName(),
		getShadingMode(),
		getSelectedObjectName(),
		lx, ly, lz,
		getAttenuationMode(),
		getLinearCoef(),
		getQuadraticCoef(),
		getAmbientStrength(),
	)
	return text
}
