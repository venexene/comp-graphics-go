// cmd/main.go — Точка входа в программу.
//
// Назначение: инициализация окна GLFW, контекста OpenGL, загрузка всех
// 3D-моделей, настройка сцены и главный цикл рендеринга.
//
// Ключевые структуры: нет собственных структур; использует ObjectState,
// SceneObject, LightConfig, Camera из пакетов scene, lighting.
//
// Ключевые функции:
//   main()          — точка входа, инициализация + цикл рендера.
//   initGlfw()      — создание окна GLFW (Core Profile 4.1).
//   initOpenGL()    — инициализация OpenGL и компиляция шейдеров.
//   findProjectRoot() — поиск корня проекта (файла go.mod) вверх по дереву.
//
// Зависимости: objects (загрузка OBJ), shaders (шейдерные программы),
//   scene (камера, свет, отрисовка), ui (текстовый оверлей), utils (мост).
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/venexene/comp-graphics-go/objects"
	"github.com/venexene/comp-graphics-go/shaders"
	"github.com/venexene/comp-graphics-go/ui"
	"github.com/venexene/comp-graphics-go/utils"
)

// width, height — размеры окна в пикселях (1280×720).
const (
	width  = 1280
	height = 720
)

// main — точка входа. Последовательность:
// 1. Блокировка горутины в OS-потоке (требование OpenGL).
// 2. Поиск корня проекта (go.mod) для разрешения путей к шейдерам/моделям.
// 3. Создание окна GLFW, инициализация OpenGL, компиляция всех шейдерных программ.
// 4. Инициализация UI (текстовый оверлей), создание текстуры-заглушки.
// 5. Загрузка трёх OBJ-моделей: снеговик (центр), сердце (+X), дефолт (-X).
// 6. Включение теста глубины, настройка камеры LookAt, проекции Perspective.
// 7. Регистрация объектов для циклического переключения (Tab).
// 8. Главный цикл: обработка ввода → отрисовка сцены → обновление заголовка.
func main() {
	runtime.LockOSThread()

	// Поиск go.mod вверх от CWD — чтобы шейдеры и модели находились
	// независимо от того, откуда запущена программа.
	projectRoot, err := findProjectRoot()
	if err != nil {
		panic(fmt.Errorf("cannot find project root: %w", err))
	}
	shaders.SetBasePath(projectRoot)

	window := initGlfw()
	defer glfw.Terminate()

	err = initOpenGL()
	if err != nil {
		panic(err)
	}
	defer shaders.CleanupLightingVariants()

	// Инициализация текстового UI.
	if err := ui.InitializeUI(window); err != nil {
		panic(fmt.Errorf("failed to initialize UI: %w", err))
	}
	defer ui.Cleanup()

	// Создание белой текстуры 1×1 для моделей без собственных текстур.
	utils.InitScene()

	// Путь к основному объекту сцены — снеговику.
	objPath := filepath.Join(projectRoot, "models", "snowman.obj")
	if len(os.Args) > 1 {
		objPath = os.Args[1]
	}

	model, err := objects.LoadOBJ(objPath)
	if err != nil {
		panic(fmt.Errorf("failed to load OBJ %s: %w", objPath, err))
	}
	defer model.Delete()

	// Дополнительные модели: сердце справа от снеговика, default слева.
	heartModel, err := objects.LoadOBJ(filepath.Join(projectRoot, "models", "heart.obj"))
	if err != nil {
		panic(fmt.Errorf("failed to load heart OBJ: %w", err))
	}
	defer heartModel.Delete()

	defaultModel, err := objects.LoadOBJ(filepath.Join(projectRoot, "models", "default.obj"))
	if err != nil {
		panic(fmt.Errorf("failed to load default OBJ: %w", err))
	}
	defer defaultModel.Delete()

	// Размещение моделей в ряд вдоль оси X с шагом 4.0.
	spacing := float32(4.0)
	// Извлечение имени из пути (без расширения).
	heartName := strings.TrimSuffix(filepath.Base("models/heart.obj"), filepath.Ext("models/heart.obj"))
	defaultName := strings.TrimSuffix(filepath.Base("models/default.obj"), filepath.Ext("models/default.obj"))

	heartObj := &utils.SceneObject{Model: heartModel, Position: mgl32.Vec3{spacing, 0.0, 0.0}, Scale: 0.6, RotationZ: 0.0, Name: heartName}
	defaultObj := &utils.SceneObject{Model: defaultModel, Position: mgl32.Vec3{-spacing, 0.0, 0.0}, Scale: 0.6, RotationZ: 0.0, Name: defaultName}

	// Включение z-буфера для корректной сортировки по глубине.
	gl.Enable(gl.DEPTH_TEST)

	// Начальная камера: смотрит на (0,0,0) из точки (5,2,5).
	view := mgl32.LookAt(5.0, 2.0, 5.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0)

	// Матрица перспективной проекции: FOV 45°, соотношение 16:9,
	// ближняя плоскость 0.1, дальняя 100.0.
	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(width)/height, 0.1, 100.0)

	fmt.Println("OBJ model loaded:", objPath)
	fmt.Println("Lighting variants loaded:", shaders.GetLightingVariantCount())

	// Регистрация дополнительных объектов для переключения по Tab.
	utils.RegisterSceneObjects(heartObj, defaultObj)
	utils.SetMainObjectName(objPath)

	// Вывод UI-оверлея в терминал при старте.
	fmt.Println("\n" + ui.GetUIOverlayText(
		utils.GetLightingName,
		utils.GetShadingMode,
		utils.GetLinearCoef,
		utils.GetQuadraticCoef,
		utils.GetAmbientStrength,
		utils.GetLightPosition,
		utils.GetSelectedObjectName,
		utils.GetAttenuationMode,
	))

	// Главный цикл рендеринга.
	for !window.ShouldClose() {
		ui.BeginFrame()
		// Обработка ввода + очистка буферов + отрисовка всех объектов.
		utils.DrawScene(window, model, view, projection, heartObj, defaultObj)

		// Обновление заголовка окна — отображение текущего состояния сцены.
		lx, ly, lz := utils.GetLightPosition()
		selected := utils.GetSelectedObjectName()
		title := fmt.Sprintf("Selected: %s | Light:(%.2f,%.2f,%.2f) | Model: %s | Shading: %s | Atten: %s | Linear: %.2f | Quad: %.2f | Ambient: %.2f",
			selected,
			lx, ly, lz,
			utils.GetLightingName(),
			utils.GetShadingMode(),
			utils.GetAttenuationMode(),
			utils.GetLinearCoef(),
			utils.GetQuadraticCoef(),
			utils.GetAmbientStrength(),
		)
		window.SetTitle(title)

		ui.EndFrame()
	}
}

// initGlfw — создание окна GLFW.
// Возвращает: указатель на окно GLFW (размер 1280×720, Core Profile 4.1).
// Побочные эффекты: инициализация GLFW, создание контекста OpenGL,
//   контекст делается текущим для вызывающего потока.
// Вызывается: однократно при запуске.
func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "OBJ Viewer with Camera Controls", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

// initOpenGL — инициализация OpenGL и компиляция шейдеров.
// Возвращает: ошибку, если gl.Init() или компиляция шейдеров провалились.
// Побочные эффекты: загружает функции OpenGL через gl.Init(),
//   компилирует все 8 шейдерных программ (см. lightingVariants в shaders/lighting.go).
// Вызывается: однократно после initGlfw().
func initOpenGL() error {
	if err := gl.Init(); err != nil {
		return err
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	if err := shaders.InitLightingVariants(); err != nil {
		return err
	}

	return nil
}

// findProjectRoot locates the project root directory by walking up from the
// executable location (or CWD) until a go.mod file is found.
func findProjectRoot() (string, error) {
	// Prefer executable path, fall back to CWD
	dir := "."
	if exe, err := os.Executable(); err == nil {
		dir = filepath.Dir(exe)
	}

	// Walk up until we find go.mod
	for {
		modPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(modPath); err == nil {
			return filepath.Abs(dir)
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached filesystem root without finding go.mod; try CWD
			break
		}
		dir = parent
	}

	// Fallback: try current working directory
	if cwd, err := os.Getwd(); err == nil {
		dir = cwd
		for {
			modPath := filepath.Join(dir, "go.mod")
			if _, err := os.Stat(modPath); err == nil {
				return filepath.Abs(dir)
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}

	return "", fmt.Errorf("go.mod not found in any parent directory")
}