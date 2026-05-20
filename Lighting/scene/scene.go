// scene/scene.go — Отрисовка сцены и управление uniform-переменными.
//
// Назначение: предоставляет функции для отрисовки всей сцены: очистка
// буферов, активация шейдерной программы, установка всех uniform-переменных
// (трансформации, материал, свет, текстура) для каждого объекта.
//
// Ключевые типы: нет собственных; использует UniformCache, LightConfig,
//   MaterialConfig (lighting), Model (objects), ObjectState, SceneObject, Camera.
//
// Ключевые функции:
//   CreateWhiteTexture()     — создаёт белую текстуру 1×1 для моделей без текстур.
//   setTransformUniforms()   — загрузка матриц Model/View/Projection/Normal в шейдер.
//   setMaterialUniforms()    — загрузка параметров материала в шейдер.
//   setLightUniforms()       — загрузка параметров света в шейдер.
//   drawObject()             — отрисовка одного объекта со всеми uniform.
//   DrawScene()              — полная отрисовка сцены (все объекты).
//
// Зависимости: вызывается из utils.DrawScene().
package scene

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/venexene/comp-graphics-go/lighting"
	"github.com/venexene/comp-graphics-go/objects"
)

// CreateWhiteTexture — создаёт текстуру 1×1 белого цвета.
// Используется для моделей, у которых нет собственных текстурных координат
// или файла текстуры. Позволяет шейдерам всегда читать из texture(diffuse_map, uv),
// не беспокоясь о bind-е текстуры.
// Возвращает: ID текстуры OpenGL.
// Побочные эффекты: создаёт объект текстуры в OpenGL.
func CreateWhiteTexture() uint32 {
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	whitePixel := []uint8{255, 255, 255, 255}
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, 1, 1, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(whitePixel))

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	gl.BindTexture(gl.TEXTURE_2D, 0)
	return texture
}

// setTransformUniforms — загружает матрицы трансформации в шейдерную программу.
// Принимает: modelMat — модельная матрица (Model Space → World Space);
//   view — видовая матрица (World Space → View Space);
//   proj — матрица проекции (View Space → Clip Space);
//   u — кэш расположений uniform-переменных.
// Вычисляет матрицу нормалей как (M⁻¹)ᵀ, извлекает подматрицу 3×3.
// Побочные эффекты: вызывает glUniformMatrix4fv/3fv (изменяет состояние OpenGL).
func setTransformUniforms(modelMat mgl32.Mat4, view, proj mgl32.Mat4, u *lighting.UniformCache) {
	if u.View != -1 {
		gl.UniformMatrix4fv(u.View, 1, false, &view[0])
	}
	if u.Projection != -1 {
		gl.UniformMatrix4fv(u.Projection, 1, false, &proj[0])
	}
	if u.Model != -1 {
		gl.UniformMatrix4fv(u.Model, 1, false, &modelMat[0])
	}
	if u.Normal != -1 {
		// Матрица нормалей = (M⁻¹)ᵀ, где M — верхняя левая 3×3 модельной матрицы.
		m := modelMat.Inv().Transpose()
		normalMat := mgl32.Mat3{
			m[0], m[1], m[2],
			m[4], m[5], m[6],
			m[8], m[9], m[10],
		}
		gl.UniformMatrix3fv(u.Normal, 1, false, &normalMat[0])
	}
}

// setMaterialUniforms — загружает параметры материала в шейдер.
// Принимает: mat — конфигурация материала (цвета ambient/diffuse/specular, sheen).
// Побочные эффекты: вызывает glUniform3f/1f для каждого поля.
func setMaterialUniforms(mat *lighting.MaterialConfig, u *lighting.UniformCache) {
	if u.MaterialAmbient != -1 {
		gl.Uniform3f(u.MaterialAmbient, mat.Ambient.X(), mat.Ambient.Y(), mat.Ambient.Z())
	}
	if u.MaterialDiffuse != -1 {
		gl.Uniform3f(u.MaterialDiffuse, mat.Diffuse.X(), mat.Diffuse.Y(), mat.Diffuse.Z())
	}
	if u.MaterialSpecular != -1 {
		gl.Uniform3f(u.MaterialSpecular, mat.Specular.X(), mat.Specular.Y(), mat.Specular.Z())
	}
	if u.MaterialSheen != -1 {
		gl.Uniform1f(u.MaterialSheen, mat.SheenCoef)
	}
}

// setLightUniforms — загружает параметры точечного источника света в шейдер.
// Принимает: l — конфигурацию света (позиция, цвета, коэффициенты затухания).
// Побочные эффекты: вызывает glUniform3f/1f/1i для каждого поля.
func setLightUniforms(l *lighting.LightConfig, u *lighting.UniformCache) {
	if u.LightAmbient != -1 {
		gl.Uniform3f(u.LightAmbient, l.Ambient.X(), l.Ambient.Y(), l.Ambient.Z())
	}
	if u.LightDiffuse != -1 {
		gl.Uniform3f(u.LightDiffuse, l.Diffuse.X(), l.Diffuse.Y(), l.Diffuse.Z())
	}
	if u.LightSpecular != -1 {
		gl.Uniform3f(u.LightSpecular, l.Specular.X(), l.Specular.Y(), l.Specular.Z())
	}
	if u.LightPosition != -1 {
		gl.Uniform3f(u.LightPosition, l.Position.X(), l.Position.Y(), l.Position.Z())
	}
	if u.LightConstant != -1 {
		gl.Uniform1f(u.LightConstant, l.Constant)
	}
	if u.LightLinear != -1 {
		gl.Uniform1f(u.LightLinear, l.Linear)
	}
	if u.LightQuadratic != -1 {
		gl.Uniform1f(u.LightQuadratic, l.Quadratic)
	}
	if u.AmbientStrength != -1 {
		gl.Uniform1f(u.AmbientStrength, l.AmbientStrength)
	}
	if u.LinearCoef != -1 {
		gl.Uniform1f(u.LinearCoef, l.LinearCoef)
	}
	if u.QuadraticCoef != -1 {
		gl.Uniform1f(u.QuadraticCoef, l.QuadraticCoef)
	}
	if u.AttenuationMode != -1 {
		// 0 = Both, 1 = Linear, 2 = Quadratic
		gl.Uniform1i(u.AttenuationMode, int32(l.Mode))
	}
}

// drawObject — отрисовывает один 3D-объект.
// Принимает: model — модель OpenGL (VAO/VBO); modelMat — её матрица трансформации;
//   view/proj — видовую и проекционную матрицы; mat — материал; l — свет;
//   u — кэш uniform-переменных.
// Последовательность: трансформации → материал → свет → текстура → Draw.
// Побочные эффекты: изменяет uniform-переменные текущей шейдерной программы,
//   вызывает glDrawArrays.
func drawObject(model *objects.Model, modelMat mgl32.Mat4, view, proj mgl32.Mat4,
	mat *lighting.MaterialConfig, l *lighting.LightConfig, u *lighting.UniformCache) {

	setTransformUniforms(modelMat, view, proj, u)
	setMaterialUniforms(mat, u)
	setLightUniforms(l, u)

	if u.DiffuseMap != -1 {
		// Привязка сэмплера к текстурному юниту 0.
		gl.Uniform1i(u.DiffuseMap, 0)
	}

	model.Draw()
}

// DrawScene — полная отрисовка одного кадра сцены.
// Последовательность:
// 1. Очистка цветового буфера и z-буфера (цвет фона — тёмный сине-зелёный).
// 2. Активация шейдерной программы (glUseProgram).
// 3. Кэширование uniform-переменных для этой программы.
// 4. Вычисление view-матрицы камеры.
// 5. Установка позиции камеры (transform.view_pos).
// 6. Бинд белой текстуры по умолчанию.
// 7. Отрисовка главной модели (снеговик).
// 8. Отрисовка дополнительных объектов (сердце, default).
func DrawScene(
	program uint32,
	mainModel *objects.Model,
	mainState *ObjectState,
	extras []*SceneObject,
	cam *Camera,
	projection mgl32.Mat4,
	lightCfg *lighting.LightConfig,
	matCfg *lighting.MaterialConfig,
	defaultTex uint32,
) {
	gl.ClearColor(0.2, 0.3, 0.3, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.UseProgram(program)

	// Кэширование uniform-переменных при каждой смене шейдера.
	var u lighting.UniformCache
	u.Refresh(program)

	view := cam.ViewMatrix()
	eye := cam.EyePosition()

	if u.ViewPos != -1 {
		// Позиция камеры нужна для вычисления направления взгляда в зеркальных моделях.
		gl.Uniform3f(u.ViewPos, eye.X(), eye.Y(), eye.Z())
	}

	// Активация белой текстуры-заглушки на текстурном юните 0.
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, defaultTex)

	// Отрисовка главной модели (снеговик).
	if mainModel != nil {
		drawObject(mainModel, mainState.ModelMatrix(), view, projection, matCfg, lightCfg, &u)
	}

	// Отрисовка дополнительных сценических объектов.
	for _, obj := range extras {
		if obj == nil || obj.Model == nil {
			continue
		}
		drawObject(obj.Model, obj.ModelMatrix(), view, projection, matCfg, lightCfg, &u)
	}
}
