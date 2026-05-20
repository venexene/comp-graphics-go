// Файл: lighting/uniforms.go
// Назначение: кеширование location-ов uniform-переменных шейдерной программы.
//
// Ключевые типы:
//   UniformCache — хранит location-ы всех uniform-переменных для быстрой установки.
//
// Ключевые функции:
//   Refresh — получает location-ы из указанной шейдерной программы.
//
// Зависимости:
//   Внутренние: — (используется scene/scene.go).
//   Внешние: github.com/go-gl/gl/v4.6-core/gl.

package lighting

import (
	"github.com/go-gl/gl/v4.6-core/gl"
)

// UniformCache — кеш location-ов uniform-переменных.
// Каждое поле хранит результат gl.GetUniformLocation.
// Значение -1 означает, что uniform не найдена в программе (не используется).
// Refresh() должен вызываться каждый раз при смене шейдерной программы.
//
// Список uniform-переменных:
//
//	transform.model       — матрица Model (mat4)
//	transform.view        — матрица View (mat4)
//	transform.projection  — матрица Projection (mat4)
//	transform.normal_mat  — матрица нормалей (mat3) = transpose(inverse(model))
//	transform.view_pos    — позиция камеры в world space (vec3)
//	material.ambient      — ambient цвет материала (vec3)
//	material.diffuse      — diffuse цвет материала (vec3)
//	material.specular     — specular цвет материала (vec3)
//	material.sheen_coef   — коэффициент блеска (float)
//	roughness             — шероховатость (float)
//	light.ambient         — ambient цвет источника (vec3)
//	light.diffuse         — diffuse цвет источника (vec3)
//	light.specular        — specular цвет источника (vec3)
//	light.position        — позиция источника в world space (vec3)
//	light.constant        — константа затухания (float)
//	light.linear          — линейный коэффициент затухания (float)
//	light.quadratic       — квадратичный коэффициент затухания (float)
//	light.ambient_strength — сила ambient-освещения (float)
//	linear_coef           — множитель линейного затухания (float)
//	quadratic_coef        — множитель квадратичного затухания (float)
//	attenuation_mode      — режим затухания: 0=оба, 1=линейное, 2=квадратичное (int)
//	u_diffuseMap          — текстурный блок диффузной карты (sampler2D, GL_TEXTURE0)
//	u_normalMap           — текстурный блок карты нормалей (sampler2D, GL_TEXTURE1)
//	u_aoMap               — текстурный блок карты ambient occlusion (sampler2D, GL_TEXTURE2)
//	u_roughnessMap        — текстурный блок карты шероховатости (sampler2D, GL_TEXTURE3)
type UniformCache struct {
	Model            int32
	View             int32
	Projection       int32
	Normal           int32
	ViewPos          int32
	MaterialAmbient  int32
	MaterialDiffuse  int32
	MaterialSpecular int32
	MaterialSheen    int32
	Roughness        int32
	LightAmbient     int32
	LightDiffuse     int32
	LightSpecular    int32
	LightPosition    int32
	LightConstant    int32
	LightLinear      int32
	LightQuadratic   int32
	AmbientStrength  int32
	LinearCoef       int32
	QuadraticCoef    int32
	AttenuationMode  int32
	DiffuseMap       int32
	NormalMap        int32
	AOMap            int32
	RoughnessMap     int32
}

// Refresh — получает location-ы uniform-переменных из шейдерной программы.
// Параметры:
//
//	program — идентификатор шейдерной программы (gl.CreateProgram).
//
// Должен вызываться при каждой смене программы (например, при переключении
// варианта освещения), т.к. разные программы могут иметь разные location-ы.
func (c *UniformCache) Refresh(program uint32) {
	c.Model = gl.GetUniformLocation(program, gl.Str("transform.model\x00"))
	c.View = gl.GetUniformLocation(program, gl.Str("transform.view\x00"))
	c.Projection = gl.GetUniformLocation(program, gl.Str("transform.projection\x00"))
	c.Normal = gl.GetUniformLocation(program, gl.Str("transform.normal_mat\x00"))
	c.ViewPos = gl.GetUniformLocation(program, gl.Str("transform.view_pos\x00"))
	c.MaterialAmbient = gl.GetUniformLocation(program, gl.Str("material.ambient\x00"))
	c.MaterialDiffuse = gl.GetUniformLocation(program, gl.Str("material.diffuse\x00"))
	c.MaterialSpecular = gl.GetUniformLocation(program, gl.Str("material.specular\x00"))
	c.MaterialSheen = gl.GetUniformLocation(program, gl.Str("material.sheen_coef\x00"))
	c.Roughness = gl.GetUniformLocation(program, gl.Str("roughness\x00"))
	c.LightAmbient = gl.GetUniformLocation(program, gl.Str("light.ambient\x00"))
	c.LightDiffuse = gl.GetUniformLocation(program, gl.Str("light.diffuse\x00"))
	c.LightSpecular = gl.GetUniformLocation(program, gl.Str("light.specular\x00"))
	c.LightPosition = gl.GetUniformLocation(program, gl.Str("light.position\x00"))
	c.LightConstant = gl.GetUniformLocation(program, gl.Str("light.constant\x00"))
	c.LightLinear = gl.GetUniformLocation(program, gl.Str("light.linear\x00"))
	c.LightQuadratic = gl.GetUniformLocation(program, gl.Str("light.quadratic\x00"))
	c.AmbientStrength = gl.GetUniformLocation(program, gl.Str("light.ambient_strength\x00"))
	c.LinearCoef = gl.GetUniformLocation(program, gl.Str("linear_coef\x00"))
	c.QuadraticCoef = gl.GetUniformLocation(program, gl.Str("quadratic_coef\x00"))
	c.AttenuationMode = gl.GetUniformLocation(program, gl.Str("attenuation_mode\x00"))
	c.DiffuseMap = gl.GetUniformLocation(program, gl.Str("u_diffuseMap\x00"))
	c.NormalMap = gl.GetUniformLocation(program, gl.Str("u_normalMap\x00"))
	c.AOMap = gl.GetUniformLocation(program, gl.Str("u_aoMap\x00"))
	c.RoughnessMap = gl.GetUniformLocation(program, gl.Str("u_roughnessMap\x00"))
}
