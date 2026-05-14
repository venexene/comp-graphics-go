package main

import (
	"log"
	"runtime"
	"fmt"

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
	// Блокировка потока для OpenGL
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()

	program := initOpenGL()

	// Создание геометрии куба
	vao, vbo, ebo := objects.CreateCube()

	// Настройка состояния OpenGL
	gl.Enable(gl.DEPTH_TEST)
	
	// Создание view матрицы
	view := mgl32.LookAt(5.0, 4.0, 8.0, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0)

	// Создание projection
	projection := mgl32.Perspective(mgl32.DegToRad(40.0), float32(width)/height, 0.1, 100.0)

	// Инициализация состояния поворотов
	rotationState := &utils.RotationState{}

	// Установка callback для клавиатуры
	window.SetKeyCallback(keyCallback(rotationState))

	// Основной цикл рендеринга
	for !window.ShouldClose() {
		utils.DrawScene(window, program, vao, view, projection, rotationState)
	}

	// Очистка ресурсов
	gl.DeleteVertexArrays(1, &vao)
	gl.DeleteBuffers(1, &vbo)
	gl.DeleteBuffers(1, &ebo)
}

// Callback для обработки нажатий клавиш
func keyCallback(rotationState *utils.RotationState) glfw.KeyCallback {
	return func(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action == glfw.Press || action == glfw.Repeat {
			switch key {
			// Управление поворотом кубиков вокруг локальной оси Y
			case glfw.Key1:
				rotationState.LocalRotation += 0.1
			case glfw.Key2:
				rotationState.LocalRotation -= 0.1
			
			// Управление поворотом пьедестала вокруг его центра
			case glfw.Key3:
				rotationState.PiedestalRotation += 0.1
			case glfw.Key4:
				rotationState.PiedestalRotation -= 0.1
			
			// Управление поворотом пьедестала вокруг глобальной оси Y
			case glfw.Key5:
				rotationState.GlobalRotation += 0.1
			case glfw.Key6:
				rotationState.GlobalRotation -= 0.1
			
			// Сброс всех поворотов
			case glfw.Key0:
				rotationState.LocalRotation = 0
				rotationState.PiedestalRotation = 0
				rotationState.GlobalRotation = 0
			}
		}
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

	window, err := glfw.CreateWindow(width, height, "Piedestal of Honor", nil, nil)
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