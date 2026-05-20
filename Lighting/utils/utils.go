// utils/utils.go — Слой совместимости и глобальное состояние.
//
// Назначение: предоставляет слой совместимости, который делегирует вызовы
// специализированным пакетам (lighting, scene, input, shaders).
// Хранит глобальные экземпляры камеры, состояния объекта, конфигурации
// света и материала, а также выбранного объекта. Новый код должен
// импортировать соответствующие пакеты напрямую.
//
// Ключевые функции:
//   InitScene()             — создание белой текстуры по умолчанию.
//   RegisterSceneObjects()  — регистрация дополнительных объектов для выбора.
//   SetMainObjectName()     — задание имени главного объекта.
//   DrawScene()             — полный цикл: ввод → отрисовка → swap.
//   GetLightingName() и др. — геттеры для UI/заголовка окна.
//
// Зависимости: объединяет все пакеты проекта; вызывается из main().
package utils

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/venexene/comp-graphics-go/input"
	"github.com/venexene/comp-graphics-go/lighting"
	"github.com/venexene/comp-graphics-go/objects"
	"github.com/venexene/comp-graphics-go/scene"
	"github.com/venexene/comp-graphics-go/shaders"
)

// SceneObject — реэкспорт типа для обратной совместимости.
type SceneObject = scene.SceneObject

// Глобальное состояние сцены (синглтоны).
var (
	defaultTex uint32 // ID белой текстуры-заглушки

	cam        = scene.DefaultCamera()    // камера
	mainState  = scene.DefaultObjectState() // трансформация главного объекта
	lightCfg   = lighting.DefaultLight()  // конфигурация света
	matCfg     = lighting.DefaultMaterial() // конфигурация материала
	sel        = scene.NewSelection("Main") // состояние выбора объекта
	inputState input.State                 // состояние клавиш (edge detection)
)

// InitScene — инициализация сцены: создание белой текстуры 1×1.
// Вызывается: из main() после инициализации шейдеров.
func InitScene() {
	defaultTex = scene.CreateWhiteTexture()
}

// RegisterSceneObjects — регистрирует дополнительные объекты для циклического
// переключения по Tab. Вызывается: из main() после загрузки моделей.
func RegisterSceneObjects(objs ...*SceneObject) {
	sel.RegisterObjects(objs...)
}

// SetMainObjectName — задаёт отображаемое имя главной модели (снеговика).
// Вызывается: из main().
func SetMainObjectName(name string) {
	sel.SetMainName(name)
}

// DrawScene — полный кадр рендеринга.
// Последовательность:
// 1. Обработка ввода (input.ProcessInput).
// 2. Получение текущей шейдерной программы.
// 3. Отрисовка сцены (scene.DrawScene).
// 4. Обработка событий GLFW (glfw.PollEvents).
// 5. Смена буферов (SwapBuffers).
// Вызывается: каждый кадр в главном цикле main().
func DrawScene(window *glfw.Window, model *objects.Model, view, projection mgl32.Mat4, extras ...*SceneObject) {
	input.ProcessInput(window, &cam, &mainState, &lightCfg, sel, &inputState)
	program := shaders.GetCurrentLightingProgram()
	scene.DrawScene(program, model, &mainState, extras, &cam, projection, &lightCfg, &matCfg, defaultTex)
	glfw.PollEvents()
	window.SwapBuffers()
}

// ===== Геттеры для UI / заголовка окна =====

func GetLightingName() string       { return shaders.GetCurrentLightingName() }
func GetShadingMode() string        { return shaders.GetCurrentShadingMode().String() }
func GetLinearCoef() float32        { return lightCfg.LinearCoef }
func GetQuadraticCoef() float32     { return lightCfg.QuadraticCoef }
func GetAmbientStrength() float32   { return lightCfg.AmbientStrength }
func GetSelectedObjectName() string { return sel.SelectedName() }
func GetAttenuationMode() string    { return lightCfg.Mode.String() }

func GetLightPosition() (float32, float32, float32) {
	return lightCfg.Position.X(), lightCfg.Position.Y(), lightCfg.Position.Z()
}
