// Package lighting provides light source configuration and attenuation control
// for the OpenGL scene.
package lighting

import "github.com/go-gl/mathgl/mgl32"

// AttenuationMode defines how light intensity falls off with distance from the source.
type AttenuationMode int

const (
	AttenuationBoth     AttenuationMode = iota // use both linear and quadratic falloff
	AttenuationLinear                          // linear falloff only: 1 / (c + k_l * d)
	AttenuationQuadratic                       // quadratic falloff only: 1 / (c + k_q * d²)
)

func (m AttenuationMode) String() string {
	switch m {
	case AttenuationLinear:
		return "Linear"
	case AttenuationQuadratic:
		return "Quadratic"
	default:
		return "Both"
	}
}

// LightConfig holds all parameters of the point light source.
type LightConfig struct {
	// Position in world space.
	Position mgl32.Vec3

	// Base light colors for each component.
	Ambient  mgl32.Vec3
	Diffuse  mgl32.Vec3
	Specular mgl32.Vec3

	// Base attenuation coefficients (before user multiplier).
	Constant  float32
	Linear    float32
	Quadratic float32

	// User-controlled attenuation multipliers.
	LinearCoef    float32
	QuadraticCoef float32

	// Ambient light intensity multiplier [0, 1].
	AmbientStrength float32

	// Current attenuation mode.
	Mode AttenuationMode
}

// DefaultLight returns the default light configuration.
func DefaultLight() LightConfig {
	return LightConfig{
		Position:        mgl32.Vec3{2.0, 4.0, 3.0},
		Ambient:         mgl32.Vec3{0.2, 0.2, 0.2},
		Diffuse:         mgl32.Vec3{1.0, 1.0, 1.0},
		Specular:        mgl32.Vec3{1.0, 1.0, 1.0},
		Constant:        1.0,
		Linear:          0.09,
		Quadratic:       0.032,
		LinearCoef:      0.5,
		QuadraticCoef:   0.5,
		AmbientStrength: 0.6,
		Mode:            AttenuationBoth,
	}
}

// CycleAttenuationMode advances the attenuation mode to the next state.
func (l *LightConfig) CycleAttenuationMode() {
	l.Mode = (l.Mode + 1) % 3
}
