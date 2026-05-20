// Пакет utils — связующее звено между компонентами: инициализирует сцену
// (загрузка текстур, создание пьедестала, загрузка сердца), запускает
// рендер каждого кадра (DrawScene), предоставляет геттеры для UI и заголовка.
package utils

import (
	"fmt"
	"path/filepath"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/venexene/comp-graphics-go/input"
	"github.com/venexene/comp-graphics-go/lighting"
	"github.com/venexene/comp-graphics-go/objects"
	"github.com/venexene/comp-graphics-go/scene"
	"github.com/venexene/comp-graphics-go/shaders"
	"github.com/venexene/comp-graphics-go/textures"
)

// Глобальные объекты сцены: камера, конфигурация света и материала,
// состояние ввода, пьедестал и текстуры.
var (
	cam        = scene.DefaultCamera()
	lightCfg   = lighting.DefaultLight()
	matCfg     = lighting.DefaultMaterial()
	inputState input.State

	podium       *scene.Podium
	numberTexIDs map[int]uint32 // ID текстур номеров: 1→1.jpg, 2→2.png, 3→3.png
	matTexIDs    map[int]uint32 // ID текстур материалов: 1→metal, 2→marble, 3→wood
)

// InitScene загружает все текстуры (3 номера + 3 материала + onyx),
// создаёт геометрию куба (CreateCube), конструирует пьедестал (NewPodium),
// загружает модель сердца (LoadOBJ) и прикрепляет его к пьедесталу.
// Вызывается один раз при запуске программы.
func InitScene(basePath string) error {
	// Загрузка текстур номеров 1, 2, 3.
	numberTexIDs = make(map[int]uint32)
	for _, n := range []int{1, 2, 3} {
		var ext string
		switch n {
		case 1:
			ext = ".jpg"
		default:
			ext = ".png"
		}
		path := filepath.Join(basePath, "textures", "imgs", fmt.Sprintf("%d%s", n, ext))
		tex, err := textures.LoadTexture(path)
		if err != nil {
			return fmt.Errorf("failed to load number texture %d: %w", n, err)
		}
		numberTexIDs[n] = tex
	}

	// Загрузка текстур материалов: ① metal, ② marble, ③ wood.
	matTexIDs = make(map[int]uint32)
	metalTex, err := textures.LoadTexture(filepath.Join(basePath, "textures", "materials", "metal.jpg"))
	if err != nil {
		return fmt.Errorf("failed to load metal texture: %w", err)
	}
	marbleTex, err := textures.LoadTexture(filepath.Join(basePath, "textures", "materials", "marble.jpg"))
	if err != nil {
		return fmt.Errorf("failed to load marble texture: %w", err)
	}
	woodTex, err := textures.LoadTexture(filepath.Join(basePath, "textures", "materials", "wood.jpg"))
	if err != nil {
		return fmt.Errorf("failed to load wood texture: %w", err)
	}
	matTexIDs[1] = metalTex
	matTexIDs[2] = marbleTex
	matTexIDs[3] = woodTex

	// Цвета кубиков: ① жёлтый, ② серый, ③ оранжевый.
	cubeColors := map[int]mgl32.Vec3{
		1: {1.0, 1.0, 0.0},
		2: {0.5, 0.5, 0.5},
		3: {1.0, 0.5, 0.0},
	}

	// Создание пьедестала (куб 0.8, spacing 0.8 — вплотную).
	podium = scene.NewPodium(0.8, 0.8, numberTexIDs, matTexIDs, cubeColors)

	// Загрузка сердца и текстуры onyx.
	heartModel, err := objects.LoadOBJ(filepath.Join(basePath, "models", "heart.obj"))
	if err != nil {
		return fmt.Errorf("failed to load heart model: %w", err)
	}
	onyxTex, err := textures.LoadTexture(filepath.Join(basePath, "textures", "materials", "onyx.jpg"))
	if err != nil {
		return fmt.Errorf("failed to load onyx texture: %w", err)
	}
	podium.SetHeart(heartModel, onyxTex, mgl32.Vec3{1.0, 0.0, 0.0})

	return nil
}

// DrawScene обрабатывает ввод, выбирает шейдерную программу и рендерит кадр.
// Вызывается в главном цикле (cmd/main.go).
func DrawScene(window *glfw.Window, projection mgl32.Mat4) {
	input.ProcessInput(window, &cam, &scene.ObjectState{}, &lightCfg, nil, &inputState)
	program := shaders.GetCurrentLightingProgram()
	scene.DrawPodium(program, podium, &cam, projection, &lightCfg, &matCfg,
		lightCfg.MaterialWeight, lightCfg.NumberWeight)
	glfw.PollEvents()
	window.SwapBuffers()
}

// Cleanup освобождает GPU-ресурсы (VAO кубов и сердца).
func Cleanup() {
	if podium != nil {
		podium.Delete()
	}
}

// --- Геттеры для UI и заголовка окна ---

func GetLightingName() string       { return shaders.GetCurrentLightingName() }
func GetShadingMode() string        { return shaders.GetCurrentShadingMode().String() }
func GetLinearCoef() float32        { return lightCfg.LinearCoef }
func GetQuadraticCoef() float32     { return lightCfg.QuadraticCoef }
func GetAmbientStrength() float32   { return lightCfg.AmbientStrength }
func GetSelectedObjectName() string { return "Podium" }
func GetAttenuationMode() string    { return lightCfg.Mode.String() }
func GetMaterialWeight() float32    { return lightCfg.MaterialWeight }
func GetNumberWeight() float32      { return lightCfg.NumberWeight }

func GetLightPosition() (float32, float32, float32) {
	return lightCfg.Position.X(), lightCfg.Position.Y(), lightCfg.Position.Z()
}

func GetCubeColor() (float32, float32, float32) {
	return lightCfg.CubeColor.X(), lightCfg.CubeColor.Y(), lightCfg.CubeColor.Z()
}
