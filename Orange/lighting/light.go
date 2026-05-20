// Файл: lighting/light.go
// Назначение: конфигурация точечного источника света (PointLight).
//
// Ключевые типы:
//   AttenuationMode — режим затухания света (оба, линейное, квадратичное).
//   LightConfig — структура источника: позиция, цвета, коэффициенты затухания.
//
// Ключевые функции:
//   DefaultLight — возвращает LightConfig с начальными параметрами.
//   CycleAttenuationMode — циклически переключает режим затухания.
//
// Зависимости:
//   Внутренние: — (используется scene/scene.go, utils/utils.go, input/input.go).
//   Внешние: github.com/go-gl/mathgl/mgl32 — математические типы (Vec3).

package lighting

import "github.com/go-gl/mathgl/mgl32"

// AttenuationMode — режим затухания точечного источника.
// Значения:
//
//	AttenuationBoth     — линейное + квадратичное затухание (0).
//	AttenuationLinear    — только линейное (1).
//	AttenuationQuadratic — только квадратичное (2).
//
// Uniform-переменная в шейдере: attenuation_mode (int).
// Переключается клавишей M или через GUI.
type AttenuationMode int

const (
	AttenuationBoth AttenuationMode = iota
	AttenuationLinear
	AttenuationQuadratic
)

// String — возвращает текстовое название режима затухания.
// Используется в заголовке окна и в отладочном выводе.
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
// Поля:
//
//	Position        — позиция источника в мировом пространстве (world space).
//	                Единицы измерения: единицы модели апельсина.
//	                Начальное значение: (3.0, 5.0, 2.5) — справа сверху.
//	Ambient         — цвет фонового (ambient) освещения, RGB [0,1].
//	Diffuse         — цвет диффузного (diffuse) освещения, RGB [0,1].
//	Specular        — цвет зеркального (specular) освещения, RGB [0,1].
//	Constant        — константа в знаменателе затухания (обычно 1.0).
//	Linear          — коэффициент линейного затухания (0.09 по умолчанию).
//	Quadratic       — коэффициент квадратичного затухания (0.032).
//	LinearCoef      — множитель линейного затухания, регулируется GUI [0, ∞).
//	QuadraticCoef   — множитель квадратичного затухания, регулируется GUI [0, ∞).
//	AmbientStrength — сила фонового освещения [0, 1], регулируется слайдером.
//	Mode            — режим затухания (AttenuationMode).
//
// Uniform-переменные в шейдере (см. uniforms.go):
//   light.position, light.ambient, light.diffuse, light.specular
//   light.constant, light.linear, light.quadratic
//   light.ambient_strength
//   linear_coef, quadratic_coef, attenuation_mode
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

// DefaultLight — возвращает LightConfig с параметрами по умолчанию.
// Источник расположен справа и сверху от апельсина.
// Цвета подобраны для реалистичного освещения оранжевой текстуры.
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

// CycleAttenuationMode — переключает режим затухания по циклу:
// Both → Linear → Quadratic → Both.
// Вызывается по нажатию клавиши M.
func (l *LightConfig) CycleAttenuationMode() {
	l.Mode = (l.Mode + 1) % 3
}
