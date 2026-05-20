package lighting

import "github.com/go-gl/mathgl/mgl32"

// MaterialConfig holds surface material properties passed to shaders.
type MaterialConfig struct {
	Ambient   mgl32.Vec3
	Diffuse   mgl32.Vec3
	Specular  mgl32.Vec3
	SheenCoef float32 // shininess exponent
}

// DefaultMaterial returns neutral white material with moderate shininess.
func DefaultMaterial() MaterialConfig {
	return MaterialConfig{
		Ambient:   mgl32.Vec3{0.2, 0.2, 0.2},
		Diffuse:   mgl32.Vec3{1.0, 1.0, 1.0},
		Specular:  mgl32.Vec3{1.0, 1.0, 1.0},
		SheenCoef: 32.0,
	}
}
