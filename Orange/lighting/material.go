package lighting

import "github.com/go-gl/mathgl/mgl32"

type MaterialConfig struct {
	Ambient   mgl32.Vec3
	Diffuse   mgl32.Vec3
	Specular  mgl32.Vec3
	SheenCoef float32
	Roughness float32
}

func DefaultMaterial() MaterialConfig {
	return MaterialConfig{
		Ambient:   mgl32.Vec3{0.15, 0.15, 0.15},
		Diffuse:   mgl32.Vec3{1.0, 1.0, 1.0},
		Specular:  mgl32.Vec3{0.8, 0.8, 0.8},
		SheenCoef: 48.0,
		Roughness: 0.4,
	}
}
