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

	program := initOpenGL()

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

	// Основной цикл рендеринга
	for !window.ShouldClose() {
		utils.DrawScene(window, program, model, view, projection)
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
func initOpenGL() uint32 {
    if err := gl.Init(); err != nil {
        panic(err)
    }

    vertexShaderSource, err := shaders.LoadShaderFile("shaders/vertex.glsl")
    if err != nil {
        panic(err)
    }

    fragmentShaderSource, err := shaders.LoadShaderFile("shaders/fragment.glsl")
    if err != nil {
        panic(err)
    }

    version := gl.GoStr(gl.GetString(gl.VERSION))
    log.Println("OpenGL version", version)

    vertexShader, err := shaders.CompileShader(vertexShaderSource, gl.VERTEX_SHADER)
    if err != nil {
        panic(err)
    }

    fragmentShader, err := shaders.CompileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
    if err != nil {
        panic(err)
    }

    program := gl.CreateProgram()
    gl.AttachShader(program, vertexShader)
    gl.AttachShader(program, fragmentShader)
    gl.LinkProgram(program)

    var status int32
    gl.GetProgramiv(program, gl.LINK_STATUS, &status)
    if status == gl.FALSE {
        var logLength int32
        gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

        log := make([]byte, logLength)
        gl.GetProgramInfoLog(program, logLength, nil, &log[0])
        panic(fmt.Sprintf("Failed to link program: %s", string(log)))
    }

    gl.DeleteShader(vertexShader)
    gl.DeleteShader(fragmentShader)

    return program
}