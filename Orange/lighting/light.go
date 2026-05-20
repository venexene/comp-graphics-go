package lighting

import "github.com/go-gl/mathgl/mgl32"

type AttenuationMode int

const (
	AttenuationBoth AttenuationMode = iota
	AttenuationLinear
	AttenuationQuadratic
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

type LightConfig struct {
	Position        mgl32.Vec3
	Ambient         mgl32.Vec3
	Diffuse         mgl32.Vec3
	Specular        mgl32.Vec3
	Constant        float32
	Linear          float32
	Quadratic       float32
	LinearCoef      float32
	QuadraticCoef   float32
	AmbientStrength float32
	Mode            AttenuationMode
}

func DefaultLight() LightConfig {
	return LightConfig{
		Position:        mgl32.Vec3{3.0, 5.0, 2.5},
		Ambient:         mgl32.Vec3{0.2, 0.2, 0.22},
		Diffuse:         mgl32.Vec3{1.0, 0.95, 0.85},
		Specular:        mgl32.Vec3{1.0, 1.0, 1.0},
		Constant:        1.0,
		Linear:          0.09,
		Quadratic:       0.032,
		LinearCoef:      0.4,
		QuadraticCoef:   0.3,
		AmbientStrength: 0.35,
		Mode:            AttenuationBoth,
	}
}
func (l *LightConfig) CycleAttenuationMode() {
	l.Mode = (l.Mode + 1) % 3
}
