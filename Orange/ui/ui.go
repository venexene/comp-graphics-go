package ui

import (
	"fmt"
)

type LightingControlState struct {
	ShowPanel bool
}

var state = LightingControlState{ShowPanel: true}

func InitializeUI(window interface{}) error {
	return nil
}
func BeginFrame() {
}
func EndFrame() {
}
func Cleanup() {
}
func GetUIOverlayText(
	getLightingName func() string,
	getShadingMode func() string,
	getLinearCoef func() float32,
	getQuadraticCoef func() float32,
	getAmbientStrength func() float32,
	getLightPosition func() (float32, float32, float32),
	getAttenuationMode func() string,
) string {
	lx, ly, lz := getLightPosition()
	text := fmt.Sprintf(`
╔═══════════════════════════════════╗
║  ORANGE + HEART — BUMP MAPPING    ║
╠═══════════════════════════════════╣
║ Model: %-27s ║
║ Shading: %-25s ║
╠═══════════════════════════════════╣
║ LIGHT                             ║
║   Pos:  %5.2f %5.2f %5.2f          ║
║   Atten (M):   %-15s ║
║   Linear (Z/X):    %.2f         ║
║   Quad   (C/V):    %.2f         ║
║   Ambient (B/N):   %.2f         ║
╠═══════════════════════════════════╣
║ CAMERA                            ║
║   Arrows  - Rotate               ║
║   WASD    - Pan                  ║
║   Space/Shift - Up/Down          ║
║   +/-     - Zoom                 ║
╠═══════════════════════════════════╣
║ LIGHT & SHADER                    ║
║   Alt+IJKLUO - Move Light        ║
║   T/G     - Switch Model         ║
║   Y       - Toggle Shading       ║
║   M       - Cycle Attenuation    ║
║   Ctrl+R  - Reset All            ║
╚═══════════════════════════════════╝
`,
		getLightingName(),
		getShadingMode(),
		lx, ly, lz,
		getAttenuationMode(),
		getLinearCoef(),
		getQuadraticCoef(),
		getAmbientStrength(),
	)
	return text
}
