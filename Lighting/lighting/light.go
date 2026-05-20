// lighting/light.go — Конфигурация точечного источника света.
//
// Назначение: определяет структуру данных точечного света (позиция, цвета,
// коэффициенты затухания) и enum режимов затухания.
//
// Ключевые типы:
//   AttenuationMode — перечисление режимов затухания (Both / Linear / Quadratic).
//   LightConfig — все параметры точечного источника света.
//
// Ключевые функции:
//   DefaultLight()          — возвращает конфигурацию по умолчанию.
//   CycleAttenuationMode()  — циклическое переключение режима затухания.
//
// Зависимости: используется scene.DrawScene() при передаче uniform-переменных
//   в шейдер; input.ProcessInput() изменяет поля LightConfig по горячим клавишам.
package lighting

import "github.com/go-gl/mathgl/mgl32"

// AttenuationMode — режим затухания освещения с расстоянием.
// Определяет, какие коэффициенты (линейный и/или квадратичный)
// участвуют в знаменателе формулы затухания.
type AttenuationMode int

const (
	// AttenuationBoth — используются оба коэффициента (по умолчанию):
	//   1 / (c + k_l * d + k_q * d²)
	AttenuationBoth AttenuationMode = iota
	// AttenuationLinear — только линейное затухание:
	//   1 / (c + k_l * d)
	AttenuationLinear
	// AttenuationQuadratic — только квадратичное затухание:
	//   1 / (c + k_q * d²)
	AttenuationQuadratic
)

// String — строковое представление режима затухания для вывода в UI.
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

// LightConfig — конфигурация точечного источника света.
// Моделирует физический точечный источник (point light) в мировом пространстве.
// Поля Position, Ambient, Diffuse, Specular передаются в шейдеры как uniform
// структуры light (см. uniforms.go:Refresh()).
type LightConfig struct {
	// Position — положение источника в мировом пространстве (world space).
	// Единицы: те же, что у координат модели (единицы OBJ-файла).
	// Диапазон: произвольный.
	Position mgl32.Vec3

	// Ambient — цвет фонового (ambient) освещения источника.
	// Компоненты R,G,B в диапазоне [0, 1].
	Ambient mgl32.Vec3
	// Diffuse — цвет диффузного освещения источника.
	// Компоненты R,G,B в диапазоне [0, 1].
	Diffuse mgl32.Vec3
	// Specular — цвет зеркального освещения источника.
	// Компоненты R,G,B в диапазоне [0, 1].
	Specular mgl32.Vec3

	// Constant — постоянный коэффициент затухания (c в формуле).
	// Обычно 1.0; гарантирует, что знаменатель не равен нулю.
	Constant float32
	// Linear — линейный коэффициент затухания (k_l).
	// Единицы: 1/единица_расстояния. Типичное значение: 0.09.
	Linear float32
	// Quadratic — квадратичный коэффициент затухания (k_q).
	// Единицы: 1/единица_расстояния². Типичное значение: 0.032.
	Quadratic float32

	// LinearCoef — множитель линейного затухания, управляемый пользователем (Z/X).
	// Умножается на Linear в формуле затухания.
	// Диапазон: [0, ∞), начальное значение 0.5.
	LinearCoef float32
	// QuadraticCoef — множитель квадратичного затухания, управляемый пользователем (C/V).
	// Умножается на Quadratic в формуле затухания.
	// Диапазон: [0, ∞), начальное значение 0.5.
	QuadraticCoef float32

	// AmbientStrength — глобальный множитель фонового освещения [0, 1].
	// Управляется клавишами B/N.
	AmbientStrength float32

	// Mode — выбранный режим затухания (Both / Linear / Quadratic).
	// Переключается клавишей M.
	Mode AttenuationMode
}

// DefaultLight — возвращает конфигурацию света по умолчанию.
// Свет расположен справа-сверху-спереди от начала координат,
// имеет белый диффузный и зеркальный цвет, слабый фоновый.
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

// CycleAttenuationMode — циклически переключает режим затухания:
// Both → Linear → Quadratic → Both.
// Вызывается из input.ProcessInput() при нажатии M.
func (l *LightConfig) CycleAttenuationMode() {
	l.Mode = (l.Mode + 1) % 3
}
