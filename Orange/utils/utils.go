// Файл: utils/utils.go
// Назначение: центральный управляющий модуль — инициализация сцены, главный цикл рендера,
// обработка ввода, глобальное состояние.
//
// Ключевые переменные (глобальное состояние):
//   cam — камера сцены.
//   lightCfg — конфигурация источника света.
//   matCfg — конфигурация материала.
//   inputState — состояние ввода (для edge-триггеров).
//   orangeScene — сцена апельсина.
//   heartScene — сцена сердца.
//   orangeState — состояние апельсина (позиция, вращение, масштаб).
//   heartState — состояние сердца.
//
// Ключевые функции:
//   InitScene — загружает сцены (апельсин + сердце) и инициализирует состояния объектов.
//   DrawScene — обрабатывает ввод, рисует кадр, обновляет окно.
//   Cleanup — освобождает все OpenGL-ресурсы.
//   GetLightingName, GetShadingMode, ... — геттеры для состояния (используются UI).
//
// Зависимости:
//   Внутренние: input, lighting, scene, shaders.
//   Внешние: github.com/go-gl/glfw/v3.3/glfw, github.com/go-gl/mathgl/mgl32.

package utils

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/venexene/comp-graphics-go/input"
	"github.com/venexene/comp-graphics-go/lighting"
	"github.com/venexene/comp-graphics-go/scene"
	"github.com/venexene/comp-graphics-go/shaders"
)

// Глобальное состояние сцены.
var (
	cam         = scene.DefaultCamera()   // орбитальная камера
	lightCfg    = lighting.DefaultLight() // точечный источник света
	matCfg      = lighting.DefaultMaterial() // материал поверхности
	inputState  input.State               // состояние клавиш для edge-триггеров
	orangeScene *scene.OrangeScene        // сцена апельсина
	heartScene  *scene.HeartScene         // сцена сердца
	orangeState scene.ObjectState         // трансформация апельсина
	heartState  scene.ObjectState         // трансформация сердца
)

// initObjectStates — устанавливает начальные позиции, масштаб и вращение объектов.
// Апельсин: слева в центре, масштаб 0.5.
// Сердце: рядом с апельсином, масштаб 0.18, слегка повёрнуто.
func initObjectStates() {
	orangeState = scene.ObjectState{
		Position:  mgl32.Vec3{-1.56, -1.51, 0.91},
		Scale:     0.5,
		RotationX: 0.0,
		RotationY: 0.0,
		RotationZ: 0.0,
	}
	heartState = scene.ObjectState{
		Position:  mgl32.Vec3{-1.8, -2.495 - 0.25, 0.0},
		Scale:     0.18,
		RotationX: 0.0,
		RotationY: 0.25,
		RotationZ: 0.0,
	}
}

// InitScene — загружает сцены (апельсин + сердце) и инициализирует состояния.
// Вызывается однократно при запуске программы.
func InitScene(basePath string) error {
	var err error
	orangeScene, err = scene.NewOrangeScene(basePath)
	if err != nil {
		return err
	}
	heartScene, err = scene.NewHeartScene(basePath)
	if err != nil {
		return err
	}
	initObjectStates()
	return nil
}
// DrawScene — главная функция одного кадра рендера.
// Вызывается в цикле из main.go.
// Последовательность:
// 1. Обработка ввода (клавиатура) — input.ProcessInput.
// 2. Получение активной шейдерной программы.
// 3. Отрисовка всей сцены (апельсин + сердце) — scene.DrawSceneObjects.
// 4. Обработка событий GLFW (PollEvents).
// 5. Обмен буферов (SwapBuffers).
func DrawScene(window *glfw.Window, projection mgl32.Mat4) {
	input.ProcessInput(window, &cam, &orangeState, &lightCfg, nil, &inputState)
	program := shaders.GetCurrentLightingProgram()
	scene.DrawSceneObjects(program, orangeScene, heartScene, &cam, projection, &lightCfg, &matCfg, &orangeState, &heartState)
	glfw.PollEvents()
	window.SwapBuffers()
}

// Cleanup — освобождает все ресурсы сцены (модели, текстуры).
// Вызывается при завершении программы (defer).
func Cleanup() {
	if orangeScene != nil {
		orangeScene.Cleanup()
	}
	if heartScene != nil {
		heartScene.Cleanup()
	}
}

// GetLightingName — возвращает имя текущей модели освещения (для UI).
func GetLightingName() string { return shaders.GetCurrentLightingName() }

// GetShadingMode — возвращает режим затенения (для UI).
func GetShadingMode() string { return shaders.GetCurrentShadingMode().String() }

// GetLinearCoef — возвращает множитель линейного затухания (для UI).
func GetLinearCoef() float32 { return lightCfg.LinearCoef }

// GetQuadraticCoef — возвращает множитель квадратичного затухания (для UI).
func GetQuadraticCoef() float32 { return lightCfg.QuadraticCoef }

// GetAmbientStrength — возвращает силу ambient-освещения (для UI).
func GetAmbientStrength() float32 { return lightCfg.AmbientStrength }

// GetAttenuationMode — возвращает строку с режимом затухания (для UI).
func GetAttenuationMode() string { return lightCfg.Mode.String() }

// GetLightPosition — возвращает позицию источника света (для UI).
func GetLightPosition() (float32, float32, float32) {
	return lightCfg.Position.X(), lightCfg.Position.Y(), lightCfg.Position.Z()
}

// GetSelectedObjectName — возвращает имя выбранного объекта (для UI).
func GetSelectedObjectName() string { return "Orange + Heart" }
