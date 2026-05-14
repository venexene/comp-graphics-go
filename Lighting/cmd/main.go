package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/venexene/comp-graphics-go/objects"
	"github.com/venexene/comp-graphics-go/shaders"
	"github.com/venexene/comp-graphics-go/utils"
)

// Размеры окна
const (
    width  = 1000
    height = 500
)

func main() {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()

	err := initOpenGL()
	if err != nil {
		panic(err)
	}
	defer shaders.CleanupLightingVariants()

	// Initialize scene (textures, etc)
	utils.InitScene()

	objPath := "models/snowman.obj"
	if len(os.Args) > 1 {
		objPath = os.Args[1]
	}

	model, err := objects.LoadOBJ(objPath)
	if err != nil {
		panic(fmt.Errorf("failed to load OBJ %s: %w", objPath, err))
	}
	defer model.Delete()

	// Настройка состояния OpenGL
	gl.Enable(gl.DEPTH_TEST)
	
	// Начальная позиция камеры
	view := mgl32.LookAt(5.0, 2.0, 5.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0)

	// Создание projection матрицы
	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(width)/height, 0.1, 100.0)

	fmt.Println("OBJ model loaded:", objPath)
	fmt.Println("Lighting variants loaded:", shaders.GetLightingVariantCount())
	fmt.Println("Controls:")
	fmt.Println("  T/G: Next/Previous lighting variant")
	fmt.Println("  Arrow keys: Rotate camera")
	fmt.Println("  WASD: Pan camera")
	fmt.Println("  +/-: Zoom")
	fmt.Println("  Q/E: Scale object")
	fmt.Println("  IJKL/U/O: Move object")
	fmt.Println("  Ctrl+R: Reset")

	// Основной цикл рендеринга
	for !window.ShouldClose() {
		utils.DrawScene(window, model, view, projection)
	}
}

// Инициализация окна GLFW
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

// Инициализация OpenGL
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