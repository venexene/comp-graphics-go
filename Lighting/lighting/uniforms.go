// lighting/uniforms.go — Кэш расположений uniform-переменных шейдера.
//
// Назначение: хранит идентификаторы (locations) всех uniform-переменных,
// используемых в шейдерных программах освещения. Позволяет избежать
// дорогих вызовов glGetUniformLocation в каждом кадре.
//
// Ключевые типы:
//   UniformCache — структура с полями-слотами для каждого uniform.
//
// Ключевые функции:
//   Refresh() — запрашивает расположения uniform-переменных для данной программы.
//
// Зависимости: используется в scene.DrawScene() для кэширования uniform;
//   шейдеры ожидают uniform-структуры transform, material, light
//   и отдельные переменные linear_coef, quadratic_coef, attenuation_mode, diffuse_map.
package lighting

import (
	"github.com/go-gl/gl/v4.6-core/gl"
)

// UniformCache — кэш расположений uniform-переменных шейдерной программы.
// Заполняется вызовом Refresh(program) при переключении шейдера.
// Значение -1 означает, что данная uniform-переменная не найдена в программе
// (для вариаций, где она не используется, например diffuse_map в Ламберте).
type UniformCache struct {
	// Трансформации (struct transform в шейдере GLSL)
	Model      int32 // transform.model — модельная матрица (mat4)
	View       int32 // transform.view — видовая матрица (mat4)
	Projection int32 // transform.projection — матрица проекции (mat4)
	Normal     int32 // transform.normal_mat — матрица нормалей (mat3)
	ViewPos    int32 // transform.view_pos — позиция камеры в world space (vec3)

	// Материал (struct material в шейдере)
	MaterialAmbient  int32 // material.ambient — коэффициент фонового отражения (vec3)
	MaterialDiffuse  int32 // material.diffuse — коэффициент диффузного отражения (vec3)
	MaterialSpecular int32 // material.specular — коэффициент зеркального отражения (vec3)
	MaterialSheen    int32 // material.sheen_coef — экспонента блеска (float)

	// Источник света (struct light в шейдере)
	LightAmbient   int32 // light.ambient — цвет фоновой составляющей света (vec3)
	LightDiffuse   int32 // light.diffuse — цвет диффузной составляющей света (vec3)
	LightSpecular  int32 // light.specular — цвет зеркальной составляющей света (vec3)
	LightPosition  int32 // light.position — позиция источника в world space (vec3)
	LightConstant  int32 // light.constant — постоянный коэффициент затухания (float)
	LightLinear    int32 // light.linear — линейный коэффициент затухания (float)
	LightQuadratic int32 // light.quadratic — квадратичный коэффициент затухания (float)

	// Отдельные uniform-переменные
	AmbientStrength int32 // light.ambient_strength — множитель фонового света (float)
	LinearCoef      int32 // linear_coef — пользовательский множитель линейного затухания (float)
	QuadraticCoef   int32 // quadratic_coef — пользовательский множитель квадратичного затухания (float)
	AttenuationMode int32 // attenuation_mode — режим затухания: 0=Both,1=Linear,2=Quadratic (int)
	DiffuseMap      int32 // diffuse_map — сэмплер диффузной текстуры (sampler2D)
}

// Refresh — запрашивает у OpenGL расположения всех uniform-переменных
// для указанной шейдерной программы. Должен вызываться при каждой смене
// шейдерной программы (1 раз за кадр или при переключении варианта освещения).
// Принимает: program — идентификатор шейдерной программы OpenGL.
// Побочные эффекты: изменяет поля кэша; не вызывает OpenGL-ошибок,
//   если uniform не используется — его location будет -1.
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
	c.DiffuseMap = gl.GetUniformLocation(program, gl.Str("diffuse_map\x00"))
}
