// Пакет lighting — UniformCache кеширует позиции uniform-переменных
// шейдерной программы, чтобы избежать дорогих вызовов glGetUniformLocation
// в каждом кадре. Обновляется при переключении шейдерной программы.
package lighting

import (
	"github.com/go-gl/gl/v4.6-core/gl"
)

// UniformCache хранит позиции всех uniform-переменных, используемых
// шейдерами освещения и мультитекстурирования.
// Значения: -1 = uniform не найден в текущей шейдерной программе.
type UniformCache struct {
	// Трансформации (struct Transform)
	Model      int32 // transform.model — матрица Model (локальная → мировая)
	View       int32 // transform.view — матрица View (мировая → камера)
	Projection int32 // transform.projection — матрица Projection (камера → экран)
	Normal     int32 // transform.normal_mat — матрица нормалей (3x3)
	ViewPos    int32 // transform.view_pos — позиция камеры в world space

	// Материал (struct Material)
	MaterialAmbient  int32 // material.ambient
	MaterialDiffuse  int32 // material.diffuse
	MaterialSpecular int32 // material.specular
	MaterialSheen    int32 // material.sheen_coef — показатель степени блика
	Roughness        int32 // roughness — шероховатость для Oren-Nayar

	// Источник света (struct PointLight)
	LightAmbient   int32 // light.ambient
	LightDiffuse   int32 // light.diffuse
	LightSpecular  int32 // light.specular
	LightPosition  int32 // light.position
	LightConstant  int32 // light.constant
	LightLinear    int32 // light.linear
	LightQuadratic int32 // light.quadratic

	// Пользовательские коэффициенты (отдельные uniform)
	AmbientStrength  int32 // light.ambient_strength — множитель фонового света
	LinearCoef       int32 // linear_coef — множитель линейного затухания
	QuadraticCoef    int32 // quadratic_coef — множитель квадратичного затухания
	AttenuationMode  int32 // attenuation_mode — тип затухания (0=Both, 1=Linear, 2=Quad)

	// Мультитекстурирование
	MaterialMap    int32 // u_materialTexture — сэмплер материала (GL_TEXTURE0)
	NumberMap      int32 // u_numberTexture — сэмплер номера (GL_TEXTURE1)
	CubeColor      int32 // u_cubeColor — базовый цвет кубика (vec3)
	MaterialWeight int32 // u_materialWeight — вес текстуры материала
	NumberWeight   int32 // u_numberWeight — вес текстуры номера
}

// Refresh запрашивает позиции uniform-переменных у текущей активной
// шейдерной программы. Вызывается каждый раз при смене программы
// (переключении модели освещения или типа шейдинга).
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
	c.MaterialMap = gl.GetUniformLocation(program, gl.Str("u_materialTexture\x00"))
	c.NumberMap = gl.GetUniformLocation(program, gl.Str("u_numberTexture\x00"))
	c.CubeColor = gl.GetUniformLocation(program, gl.Str("u_cubeColor\x00"))
	c.MaterialWeight = gl.GetUniformLocation(program, gl.Str("u_materialWeight\x00"))
	c.NumberWeight = gl.GetUniformLocation(program, gl.Str("u_numberWeight\x00"))
}
