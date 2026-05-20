// lighting/material.go — Конфигурация материала поверхности.
//
// Назначение: определяет оптические свойства поверхности объекта (цвета
// фонового, диффузного, зеркального отражения и коэффициент блеска).
//
// Ключевые типы:
//   MaterialConfig — параметры материала, передаваемые в шейдеры как uniform
//   структуры material (см. uniforms.go:Refresh()).
//
// Ключевые функции:
//   DefaultMaterial() — возвращает нейтральный белый материал.
//
// Зависимости: используется в scene.DrawScene() при установке uniform-переменных
//   материала; MaterialConfig хранится в глобальной переменной utils.matCfg.
package lighting

import "github.com/go-gl/mathgl/mgl32"

// MaterialConfig — параметры материала поверхности.
// Моделирует упрощённую модель Фонга: фоновая, диффузная и зеркальная
// составляющие + экспонента блеска.
// Все поля передаются в шейдеры как uniform-структура material.
type MaterialConfig struct {
	// Ambient — коэффициент фонового отражения (k_a).
	// Компоненты R,G,B в диапазоне [0, 1].
	// В шейдере: material.ambient.
	Ambient mgl32.Vec3
	// Diffuse — коэффициент диффузного отражения (k_d).
	// Компоненты R,G,B в диапазоне [0, 1].
	// В шейдере: material.diffuse.
	Diffuse mgl32.Vec3
	// Specular — коэффициент зеркального отражения (k_s).
	// Компоненты R,G,B в диапазоне [0, 1].
	// В шейдере: material.specular.
	Specular mgl32.Vec3
	// SheenCoef — экспонента блеска (shininess). Чем выше, тем меньше
	// и острее блик. Диапазон: [1, 256], начальное значение 32.
	// В шейдере: material.sheen_coef.
	SheenCoef float32
}

// DefaultMaterial — возвращает материал по умолчанию.
// Белый диффузный цвет, умеренный блеск (32), слабое фоновое отражение.
func DefaultMaterial() MaterialConfig {
	return MaterialConfig{
		Ambient:   mgl32.Vec3{0.2, 0.2, 0.2},
		Diffuse:   mgl32.Vec3{1.0, 1.0, 1.0},
		Specular:  mgl32.Vec3{1.0, 1.0, 1.0},
		SheenCoef: 32.0,
	}
}
