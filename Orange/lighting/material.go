// Файл: lighting/material.go
// Назначение: конфигурация материала поверхности (Material).
//
// Ключевые типы:
//   MaterialConfig — параметры материала: ambient, diffuse, specular цвета, блеск, шероховатость.
//
// Ключевые функции:
//   DefaultMaterial — возвращает MaterialConfig с параметрами по умолчанию.
//
// Зависимости:
//   Внутренние: — (используется scene/scene.go, utils/utils.go).
//   Внешние: github.com/go-gl/mathgl/mgl32.

package lighting

import "github.com/go-gl/mathgl/mgl32"

// MaterialConfig — параметры материала поверхности.
// Поля:
//
//	Ambient   — цвет фонового отражения, RGB [0,1].
//	Diffuse   — цвет диффузного отражения, RGB [0,1].
//	Specular  — цвет зеркального отражения, RGB [0,1].
//	SheenCoef — коэффициент блеска (shininess) для specular.
//	           Диапазон [1.0, 256.0]. Чем выше, тем меньше и ярче пятно блика.
//	           В шейдере: material.sheen_coef.
//	Roughness — шероховатость поверхности [0, 1], используется для карты шероховатости.
//
// Uniform-переменные в шейдере (см. uniforms.go):
//   material.ambient, material.diffuse, material.specular, material.sheen_coef
//   roughness
type MaterialConfig struct {
	Ambient   mgl32.Vec3
	Diffuse   mgl32.Vec3
	Specular  mgl32.Vec3
	SheenCoef float32
	Roughness float32
}

// DefaultMaterial — возвращает MaterialConfig с нейтральными параметрами.
// Цвета подобраны так, чтобы итоговый цвет определялся текстурой,
// а не материалом (диффузный белый, specular светло-серый).
func DefaultMaterial() MaterialConfig {
	return MaterialConfig{
		Ambient:   mgl32.Vec3{0.15, 0.15, 0.15},
		Diffuse:   mgl32.Vec3{1.0, 1.0, 1.0},
		Specular:  mgl32.Vec3{0.8, 0.8, 0.8},
		SheenCoef: 48.0,
		Roughness: 0.4,
	}
}
