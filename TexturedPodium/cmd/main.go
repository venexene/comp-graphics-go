// Пакет main — точка входа. Инициализирует GLFW, OpenGL, шейдеры,
// сцену (текстуры + пьедестал + сердце), запускает главный цикл рендера.
// Зависимости: go-gl/gl, go-gl/glfw, go-gl/mathgl — внешние.
// shaders, ui, utils — внутренние.
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/venexene/comp-graphics-go/shaders"
	"github.com/venexene/comp-graphics-go/ui"
	"github.com/venexene/comp-graphics-go/utils"
)

// Размеры окна: 1280×720 пикселей, не изменяемое.
const (
	width  = 1280
	height = 720
)

func main() {
	runtime.LockOSThread()

	
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

	
	if err := ui.InitializeUI(window); err != nil {
		panic(fmt.Errorf("failed to initialize UI: %w", err))
	}
	defer ui.Cleanup()

	
	if err := utils.InitScene(projectRoot); err != nil {
		panic(fmt.Errorf("failed to init scene: %w", err))
	}
	defer utils.Cleanup()

	
	gl.Enable(gl.DEPTH_TEST)

	
	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(width)/height, 0.1, 100.0)

	fmt.Println("Winner's podium with multi-texturing")
	fmt.Println("Lighting variants loaded:", shaders.GetLightingVariantCount())

	fmt.Println("\n" + ui.GetUIOverlayText(
		utils.GetLightingName,
		utils.GetShadingMode,
		utils.GetLinearCoef,
		utils.GetQuadraticCoef,
		utils.GetAmbientStrength,
		utils.GetLightPosition,
		utils.GetAttenuationMode,
		utils.GetMaterialWeight,
		utils.GetNumberWeight,
	))

	
	for !window.ShouldClose() {
		ui.BeginFrame()
		utils.DrawScene(window, projection)

		lx, ly, lz := utils.GetLightPosition()
		cx, cy, cz := utils.GetCubeColor()
		title := fmt.Sprintf(
			"Podium | Light:(%.1f,%.1f,%.1f) | Model: %s | Shading: %s | Atten: %s | MatW:%.2f NumW:%.2f | Color:(%.1f,%.1f,%.1f)",
			lx, ly, lz,
			utils.GetLightingName(),
			utils.GetShadingMode(),
			utils.GetAttenuationMode(),
			utils.GetMaterialWeight(),
			utils.GetNumberWeight(),
			cx, cy, cz,
		)
		window.SetTitle(title)

		ui.EndFrame()
	}
}

func initGlfw() *glfw.Window {
    if err := glfw.Init(); err != nil {
        panic(err)
    }
    
    glfw.WindowHint(glfw.Resizable, glfw.False)
    glfw.WindowHint(glfw.ContextVersionMajor, 4)
    glfw.WindowHint(glfw.ContextVersionMinor, 1)
    glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
    glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

    window, err := glfw.CreateWindow(width, height, "Winner's Podium", nil, nil)
    if err != nil {
        panic(err)
    }
    window.MakeContextCurrent()

    return window
}

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

func findProjectRoot() (string, error) {
	
	dir := "."
	if exe, err := os.Executable(); err == nil {
		dir = filepath.Dir(exe)
	}

	
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