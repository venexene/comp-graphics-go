// Пакет lighting содержит структуры для работы с освещением и материалами:
// LightConfig — настройки точечного источника света (позиция, цвет, затухание),
// MaterialConfig — поверхностные свойства материала (ambient, diffuse, specular),
// UniformCache — кеш позиций uniform-переменных для шейдеров.
// Зависимости: go-gl/gl — получение uniform-локаций,
// go-gl/mathgl — векторные и матричные типы.
package lighting

import "github.com/go-gl/mathgl/mgl32"

// MaterialConfig хранит свойства поверхности материала, передаваемые в шейдер
// через uniform-структуру material.
// Поля:
//
//	Ambient   — коэффициент фонового отражения (RGB), диапазон [0.0, 1.0],
//	            умножается на light.ambient для получения фоновой компоненты.
//	Diffuse   — коэффициент диффузного отражения (RGB), диапазон [0.0, 1.0],
//	            умножается на light.diffuse и N·L.
//	Specular  — коэффициент зеркального отражения (RGB), диапазон [0.0, 1.0],
//	            умножается на light.specular и (R·V)^sheen.
//	SheenCoef — показатель степени для зеркального блика (ширина), чем больше,
//	            тем уже блик. Типичный диапазон [1, 128].
//	Roughness — параметр шероховатости для Oren-Nayar, диапазон [0.0, 1.0].
//	            0 = идеально гладкая поверхность (Ламберт), 1 = максимально шероховатая.
type MaterialConfig struct {
	Ambient   mgl32.Vec3
	Diffuse   mgl32.Vec3
	Specular  mgl32.Vec3
	SheenCoef float32 
	Roughness float32 
}

// DefaultMaterial возвращает MaterialConfig с нейтральными значениями.
// Используется по умолчанию для всех объектов сцены.
func DefaultMaterial() MaterialConfig {
	return MaterialConfig{
		Ambient:   mgl32.Vec3{0.2, 0.2, 0.2},
		Diffuse:   mgl32.Vec3{1.0, 1.0, 1.0},
		Specular:  mgl32.Vec3{1.0, 1.0, 1.0},
		SheenCoef: 32.0,
		Roughness: 0.5,
	}
}
