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
)

// Размеры окна
const (
    width  = 1200
    height = 600
)

func main() {
	// Блокировка потока для OpenGL
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()

	// Создаем две шейдерные программы
	programDefault := initOpenGL("shaders/vertex.glsl", "shaders/fragment.glsl")
	programStriped := initOpenGL("shaders/vertex.glsl", "shaders/striped_fragment.glsl")

	// Создание геометрии пятиугольника
	pentagonVAO, pentagonVBO := objects.CreatePentagon()
	defer func() {
		gl.DeleteVertexArrays(1, &pentagonVAO)
		gl.DeleteBuffers(1, &pentagonVBO)
	}()

	// Создание геометрии куба
	cubeVAO, cubeVBO, cubeEBO := objects.CreateColoredCube()
	defer func() {
		gl.DeleteVertexArrays(1, &cubeVAO)
		gl.DeleteBuffers(1, &cubeVBO)
		gl.DeleteBuffers(1, &cubeEBO)
	}()

	// Создание геометрии полосатого квадрата
	squareVAO, squareVBO := objects.CreateStripedSquare()
	defer func() {
		gl.DeleteVertexArrays(1, &squareVAO)
		gl.DeleteBuffers(1, &squareVBO)
	}()

	// Настройка состояния OpenGL
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	
	// Создание view матрицы
	view := mgl32.LookAt(0.0, 0.0, 8.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0)

	// Создание projection матрицы
	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(width)/height, 0.1, 100.0)

	// Основной цикл рендеринга
	for !window.ShouldClose() {
		drawScene(window, programDefault, programStriped, pentagonVAO, cubeVAO, squareVAO, view, projection)
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

	window, err := glfw.CreateWindow(width, height, "Pentagon, Cube and Striped Square", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

// Инициализация OpenGL с указанными шейдерами
func initOpenGL(vertexPath, fragmentPath string) uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}

	vertexShaderSource, err := shaders.LoadShaderFile(vertexPath)
	if err != nil {
		panic(err)
	}

	fragmentShaderSource, err := shaders.LoadShaderFile(fragmentPath)
	if err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)
	log.Printf("Loading shaders: %s, %s", vertexPath, fragmentPath)

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

// Отрисовка сцены
func drawScene(window *glfw.Window, programDefault, programStriped uint32, pentagonVAO, cubeVAO, squareVAO uint32, view, projection mgl32.Mat4) {
	// Очистка экрана
	gl.ClearColor(0.2, 0.3, 0.3, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	
	// ПЯТИУГОЛЬНИК
	gl.UseProgram(programDefault)
	
	// Получаем и устанавливаем uniform-переменные для programDefault
	modelUniform := gl.GetUniformLocation(programDefault, gl.Str("model\x00"))
	viewUniform := gl.GetUniformLocation(programDefault, gl.Str("view\x00"))
	projUniform := gl.GetUniformLocation(programDefault, gl.Str("projection\x00"))
	useTextureUniform := gl.GetUniformLocation(programDefault, gl.Str("useTexture\x00"))
	
	// Устанавливаем view и projection для этой программы
	gl.UniformMatrix4fv(viewUniform, 1, false, &view[0])
	gl.UniformMatrix4fv(projUniform, 1, false, &projection[0])
	gl.Uniform1i(useTextureUniform, 0) // Не используем текстуру
	
	// Рисуем пятиугольник
	pentagonModel := mgl32.Translate3D(-4.0, 0.0, 0.0)
	gl.UniformMatrix4fv(modelUniform, 1, false, &pentagonModel[0])
	
	gl.BindVertexArray(pentagonVAO)
	gl.DrawArrays(gl.TRIANGLES, 0, 15)
	
	// КУБ
	cubeModel := mgl32.Translate3D(0.0, 0.0, 0.0)
	rotation := mgl32.HomogRotate3D(mgl32.DegToRad(45.0), mgl32.Vec3{1, 1, 0}.Normalize())
	cubeModel = cubeModel.Mul4(rotation)
	gl.UniformMatrix4fv(modelUniform, 1, false, &cubeModel[0])
	
	gl.BindVertexArray(cubeVAO)
	gl.DrawElements(gl.TRIANGLES, 36, gl.UNSIGNED_INT, nil)
	
	// КВАДРАТ
	gl.UseProgram(programStriped)
	
	// Получаем и устанавливаем uniform-переменные для programStriped
	modelUniformStriped := gl.GetUniformLocation(programStriped, gl.Str("model\x00"))
	viewUniformStriped := gl.GetUniformLocation(programStriped, gl.Str("view\x00"))
	projUniformStriped := gl.GetUniformLocation(programStriped, gl.Str("projection\x00"))
	useTextureUniformStriped := gl.GetUniformLocation(programStriped, gl.Str("useTexture\x00"))
	
	// Устанавливаем view и projection для этой программы
	gl.UniformMatrix4fv(viewUniformStriped, 1, false, &view[0])
	gl.UniformMatrix4fv(projUniformStriped, 1, false, &projection[0])
	gl.Uniform1i(useTextureUniformStriped, 0) // Не используем текстуру
	
	// Рисуем полосатый квадрат
	squareModel := mgl32.Translate3D(4.0, 0.0, 0.0)
	// Небольшой поворот, чтобы показать, что полоски привязаны к геометрии
	gl.UniformMatrix4fv(modelUniformStriped, 1, false, &squareModel[0])
	
	gl.BindVertexArray(squareVAO)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)

	glfw.PollEvents()
	window.SwapBuffers()
}