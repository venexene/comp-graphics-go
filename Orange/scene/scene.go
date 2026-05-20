// Файл: scene/scene.go
// Назначение: сцена с апельсином и сердцем, загрузка текстур, установка
// uniform-переменных и отрисовка объектов.
//
// Ключевые структуры:
//   OrangeScene — сцена апельсина: модель + текстуры (диффузная, normal map, AO, roughness).
//   HeartScene  — сцена сердца: модель + текстуры (диффузная, normal map, metalness, roughness).
//
// Ключевые функции:
//   NewOrangeScene — загружает модель апельсина и его текстуры.
//   NewHeartScene — загружает модель сердца и его текстуры.
//   setTransformUniforms — устанавливает uniform-матрицы (Model, View, Projection, нормалей).
//   setMaterialUniforms — устанавливает uniform материала.
//   setLightUniforms — устанавливает uniform источника света.
//   bindOrangeTextures — привязывает текстуры апельсина к текстурным блокам.
//   bindHeartTextures — привязывает текстуры сердца к текстурным блокам.
//   DrawSceneObjects — главная функция отрисовки: очистка буферов, установка uniform, рисование.
//   DrawOrangeScene — упрощённая обёртка для отрисовки только апельсина.
//   CreateWhiteTexture — создаёт текстуру 1x1 белого пикселя (заглушка).
//
// Зависимости:
//   Внутренние: lighting, objects, textures.
//   Внешние: github.com/go-gl/gl/v4.6-core/gl, github.com/go-gl/mathgl/mgl32.

package scene

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/venexene/comp-graphics-go/lighting"
	"github.com/venexene/comp-graphics-go/objects"
	"github.com/venexene/comp-graphics-go/textures"
)

// OrangeScene — сцена, содержащая модель апельсина и его текстуры.
// Поля:
//
//	Model        — 3D-модель апельсина (*objects.Model, загружен из orange.obj).
//	DiffuseTex   — диффузная текстура (GL_TEXTURE0), food_0022_color_1k.jpg.
//	NormalTex    — карта нормалей (GL_TEXTURE1), food_0022_normal_opengl_1k.png.
//	             Это tangent-space normal map (OpenGL-формат).
//	AOTex        — карта ambient occlusion (GL_TEXTURE2), food_0022_ao_1k.jpg.
//	RoughnessTex — карта шероховатости (GL_TEXTURE3), food_0022_roughness_1k.jpg.
type OrangeScene struct {
	Model        *objects.Model
	DiffuseTex   uint32
	NormalTex    uint32
	AOTex        uint32
	RoughnessTex uint32
}

// NewOrangeScene — загружает модель апельсина и все текстуры.
// Параметры:
//
//	basePath — корневой путь проекта (где находится go.mod).
//
// Пути к файлам формируются как basePath + "/models/orange.obj" и т.д.
// Если загрузка любого файла не удалась, возвращается ошибка.
func NewOrangeScene(basePath string) (*OrangeScene, error) {
	s := &OrangeScene{}
	var err error
	s.Model, err = objects.LoadOBJ(basePath + "/models/orange.obj")
	if err != nil {
		return nil, err
	}
	s.DiffuseTex, err = textures.LoadTexture(basePath + "/textures/orange/food_0022_color_1k.jpg")
	if err != nil {
		return nil, err
	}
	s.NormalTex, err = textures.LoadTexture(basePath + "/textures/orange/food_0022_normal_opengl_1k.png")
	if err != nil {
		return nil, err
	}
	s.AOTex, err = textures.LoadTexture(basePath + "/textures/orange/food_0022_ao_1k.jpg")
	if err != nil {
		return nil, err
	}
	s.RoughnessTex, err = textures.LoadTexture(basePath + "/textures/orange/food_0022_roughness_1k.jpg")
	if err != nil {
		return nil, err
	}
	return s, nil
}
// Cleanup — освобождает OpenGL-ресурсы сцены апельсина.
// Вызывается при завершении программы.
// Удаляет: модель (VAO/VBO), все текстуры (Diffuse, Normal, AO, Roughness).
func (s *OrangeScene) Cleanup() {
	if s.Model != nil {
		s.Model.Delete()
	}
	for _, tex := range []uint32{s.DiffuseTex, s.NormalTex, s.AOTex, s.RoughnessTex} {
		if tex != 0 {
			gl.DeleteTextures(1, &tex)
		}
	}
}

// CreateWhiteTexture — создаёт текстуру 1x1 белого пикселя.
// Используется как заглушка, если текстура не загружена.
// Параметры фильтрации: GL_LINEAR, GL_REPEAT.
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
// setTransformUniforms — устанавливает uniform-переменные трансформации в шейдере.
// Параметры:
//
//	modelMat — матрица Model (object space → world space).
//	view     — матрица View (world space → view space).
//	proj     — матрица Projection (view space → clip space).
//	u        — кеш location-ов uniform-переменных.
//
// Побочные эффекты: gl.UniformMatrix4fv для Model, View, Projection, Normal.
// Матрица нормалей (normal_mat) вычисляется как transpose(inverse(model)),
// затем приводится к mat3 для исключения компоненты сдвига.
// Это гарантирует корректное преобразование нормалей при неравномерном масштабе.
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
		m := modelMat.Inv().Transpose()
		normalMat := mgl32.Mat3{
			m[0], m[1], m[2],
			m[4], m[5], m[6],
			m[8], m[9], m[10],
		}
		gl.UniformMatrix3fv(u.Normal, 1, false, &normalMat[0])
	}
}
// setMaterialUniforms — устанавливает uniform-переменные материала в шейдере.
// Параметры:
//
//	mat — конфигурация материала (цвета, блеск, шероховатость).
//	u   — кеш location-ов uniform-переменных.
//
// Побочные эффекты: gl.Uniform3f/1f для material.ambient, diffuse, specular, sheen_coef, roughness.
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
	if u.Roughness != -1 {
		gl.Uniform1f(u.Roughness, mat.Roughness)
	}
}
// setLightUniforms — устанавливает uniform-переменные источника света в шейдере.
// Параметры:
//
//	l — конфигурация источника света (позиция, цвета, затухание).
//	u — кеш location-ов uniform-переменных.
//
// Побочные эффекты: gl.Uniform3f/1f/1i для всех полей light.*, linear_coef,
// quadratic_coef, attenuation_mode.
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
		gl.Uniform1i(u.AttenuationMode, int32(l.Mode))
	}
}
// drawObject — вспомогательная функция для отрисовки одного объекта.
// Устанавливает все uniform-переменные (трансформации, материал, свет)
// и вызывает model.Draw().
// Параметры:
//
//	model   — 3D-модель для отрисовки.
//	modelMat — матрица Model для данного объекта.
//	view, proj — матрицы View и Projection.
//	mat     — конфигурация материала.
//	l       — конфигурация источника света.
//	u       — кеш uniform-переменных.
func drawObject(model *objects.Model, modelMat mgl32.Mat4, view, proj mgl32.Mat4,
	mat *lighting.MaterialConfig, l *lighting.LightConfig, u *lighting.UniformCache) {
	setTransformUniforms(modelMat, view, proj, u)
	setMaterialUniforms(mat, u)
	setLightUniforms(l, u)
	if u.ViewPos != -1 {
	}
	model.Draw()
}
// bindOrangeTextures — привязывает текстуры апельсина к текстурным блокам OpenGL.
// Соответствие текстурных блоков:
//
//	GL_TEXTURE0 — DiffuseTex  → uniform u_diffuseMap  (sampler2D, значение 0)
//	GL_TEXTURE1 — NormalTex   → uniform u_normalMap   (sampler2D, значение 1)
//	GL_TEXTURE2 — AOTex       → uniform u_aoMap       (sampler2D, значение 2)
//	GL_TEXTURE3 — RoughnessTex → uniform u_roughnessMap (sampler2D, значение 3)
//
// Вызывается перед отрисовкой апельсина в главном цикле рендера.
func bindOrangeTextures(s *OrangeScene, u *lighting.UniformCache) {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, s.DiffuseTex)
	if u.DiffuseMap != -1 {
		gl.Uniform1i(u.DiffuseMap, 0)
	}
	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, s.NormalTex)
	if u.NormalMap != -1 {
		gl.Uniform1i(u.NormalMap, 1)
	}
	gl.ActiveTexture(gl.TEXTURE2)
	gl.BindTexture(gl.TEXTURE_2D, s.AOTex)
	if u.AOMap != -1 {
		gl.Uniform1i(u.AOMap, 2)
	}
	gl.ActiveTexture(gl.TEXTURE3)
	gl.BindTexture(gl.TEXTURE_2D, s.RoughnessTex)
	if u.RoughnessMap != -1 {
		gl.Uniform1i(u.RoughnessMap, 3)
	}
}
// DrawOrangeScene — упрощённая обёртка для отрисовки только апельсина.
// Вызывает DrawSceneObjects с heart=nil.
func DrawOrangeScene(
	program uint32,
	s *OrangeScene,
	cam *Camera,
	projection mgl32.Mat4,
	lightCfg *lighting.LightConfig,
	matCfg *lighting.MaterialConfig,
	objectState *ObjectState,
) {
	DrawSceneObjects(program, s, nil, cam, projection, lightCfg, matCfg, objectState, nil)
}

// HeartScene — сцена, содержащая модель сердца и его текстуры.
// Поля:
//
//	Model        — 3D-модель сердца (из heart.obj).
//	DiffuseTex   — диффузная текстура (GL_TEXTURE0), ..._Color.jpg.
//	NormalTex    — карта нормалей (GL_TEXTURE1), ..._NormalGL.jpg.
//	MetalnessTex — карта металличности (GL_TEXTURE2), ..._Metalness.jpg.
//	RoughnessTex — карта шероховатости (GL_TEXTURE3), ..._Roughness.jpg.
type HeartScene struct {
	Model        *objects.Model
	DiffuseTex   uint32
	NormalTex    uint32
	MetalnessTex uint32
	RoughnessTex uint32
}

// NewHeartScene — загружает модель сердца и его текстуры.
// Параметры:
//
//	basePath — корневой путь проекта.
//
// Загружает: heart.obj, карту цвета, карту нормалей (OpenGL-формат),
// карту металличности, карту шероховатости.
func NewHeartScene(basePath string) (*HeartScene, error) {
	s := &HeartScene{}
	var err error
	s.Model, err = objects.LoadOBJ(basePath + "/models/heart.obj")
	if err != nil {
		return nil, err
	}
	s.DiffuseTex, err = textures.LoadTexture(basePath + "/textures/ornament2/ChristmasTreeOrnament021_1K-JPG_Color.jpg")
	if err != nil {
		return nil, err
	}
	s.NormalTex, err = textures.LoadTexture(basePath + "/textures/ornament2/ChristmasTreeOrnament021_1K-JPG_NormalGL.jpg")
	if err != nil {
		return nil, err
	}
	s.MetalnessTex, err = textures.LoadTexture(basePath + "/textures/ornament2/ChristmasTreeOrnament021_1K-JPG_Metalness.jpg")
	if err != nil {
		return nil, err
	}
	s.RoughnessTex, err = textures.LoadTexture(basePath + "/textures/ornament2/ChristmasTreeOrnament021_1K-JPG_Roughness.jpg")
	if err != nil {
		return nil, err
	}
	return s, nil
}
// Cleanup — освобождает OpenGL-ресурсы сцены сердца.
// Удаляет модель и все текстуры.
func (s *HeartScene) Cleanup() {
	if s.Model != nil {
		s.Model.Delete()
	}
	for _, tex := range []uint32{s.DiffuseTex, s.NormalTex, s.MetalnessTex, s.RoughnessTex} {
		if tex != 0 {
			gl.DeleteTextures(1, &tex)
		}
	}
}

// bindHeartTextures — привязывает текстуры сердца к текстурным блокам OpenGL.
// Соответствие:
//
//	GL_TEXTURE0 — DiffuseTex   → u_diffuseMap   (0)
//	GL_TEXTURE1 — NormalTex    → u_normalMap    (1)
//	GL_TEXTURE2 — MetalnessTex → u_aoMap        (2) — используется как замена AO
//	GL_TEXTURE3 — RoughnessTex → u_roughnessMap (3)
func bindHeartTextures(s *HeartScene, u *lighting.UniformCache) {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, s.DiffuseTex)
	if u.DiffuseMap != -1 {
		gl.Uniform1i(u.DiffuseMap, 0)
	}
	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, s.NormalTex)
	if u.NormalMap != -1 {
		gl.Uniform1i(u.NormalMap, 1)
	}
	gl.ActiveTexture(gl.TEXTURE2)
	gl.BindTexture(gl.TEXTURE_2D, s.MetalnessTex)
	if u.AOMap != -1 {
		gl.Uniform1i(u.AOMap, 2)
	}
	gl.ActiveTexture(gl.TEXTURE3)
	gl.BindTexture(gl.TEXTURE_2D, s.RoughnessTex)
	if u.RoughnessMap != -1 {
		gl.Uniform1i(u.RoughnessMap, 3)
	}
}
// DrawSceneObjects — главная функция отрисовки кадра.
// Вызывается каждый кадр из utils.DrawScene.
//
// Порядок действий:
// 1. Очистка буфера цвета и глубины.
// 2. Активация шейдерной программы.
// 3. Обновление кеша uniform-переменных.
// 4. Вычисление матрицы View и позиции камеры.
// 5. Установка uniform камеры (view_pos).
// 6. Отрисовка апельсина (если задан): установка uniform, привязка текстур, Draw.
// 7. Отрисовка сердца (если задано): установка uniform, привязка текстур, Draw.
//
// Параметры:
//
//	program      — идентификатор активной шейдерной программы.
//	orange       — сцена апельсина (может быть nil).
//	heart        — сцена сердца (может быть nil).
//	cam          — камера (для ViewMatrix и EyePosition).
//	projection   — матрица проекции (перспективная).
//	lightCfg     — конфигурация источника света.
//	matCfg       — конфигурация материала.
//	orangeState  — состояние апельсина (позиция, вращение, масштаб).
//	heartState   — состояние сердца (позиция, вращение, масштаб).
func DrawSceneObjects(
	program uint32,
	orange *OrangeScene,
	heart *HeartScene,
	cam *Camera,
	projection mgl32.Mat4,
	lightCfg *lighting.LightConfig,
	matCfg *lighting.MaterialConfig,
	orangeState *ObjectState,
	heartState *ObjectState,
) {
	gl.ClearColor(0.15, 0.15, 0.18, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)
	var u lighting.UniformCache
	u.Refresh(program)
	view := cam.ViewMatrix()
	eye := cam.EyePosition()
	if u.ViewPos != -1 {
		gl.Uniform3f(u.ViewPos, eye.X(), eye.Y(), eye.Z())
	}
	if orange != nil && orangeState != nil {
		modelMat := orangeState.ModelMatrix()
		setTransformUniforms(modelMat, view, projection, &u)
		setMaterialUniforms(matCfg, &u)
		setLightUniforms(lightCfg, &u)
		bindOrangeTextures(orange, &u)
		orange.Model.Draw()
	}
	if heart != nil && heartState != nil {
		modelMat := heartState.ModelMatrix()
		setTransformUniforms(modelMat, view, projection, &u)
		setMaterialUniforms(matCfg, &u)
		setLightUniforms(lightCfg, &u)
		bindHeartTextures(heart, &u)
		heart.Model.Draw()
	}
}
