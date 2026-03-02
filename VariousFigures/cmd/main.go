package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/venexene/comp-graphics-go/objects"
	"github.com/venexene/comp-graphics-go/shaders"
)

// Размеры окна
const (
	width  = 1200
	height = 800 // Увеличил высоту для размещения 3 кубиков внизу
)

func main() {
	// Блокировка потока для OpenGL
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()

	// Создаем шейдерные программы
	programDefault := initOpenGL("shaders/vertex.glsl", "shaders/fragment.glsl")
	programStriped := initOpenGL("shaders/vertex.glsl", "shaders/striped_fragment.glsl")
	programChecker := initOpenGL("shaders/vertex.glsl", "shaders/checker_fragment.glsl")
	programDiagonal := initOpenGL("shaders/vertex.glsl", "shaders/diagonal_fragment.glsl")
	programHorizontal := initOpenGL("shaders/vertex.glsl", "shaders/horizontal_fragment.glsl")

	// Создание геометрии пятиугольника
	pentagonVAO, pentagonVBO := objects.CreatePentagon()
	defer func() {
		gl.DeleteVertexArrays(1, &pentagonVAO)
		gl.DeleteBuffers(1, &pentagonVBO)
	}()

	// Создание геометрии куба (для всех кубиков используем одну и ту же геометрию)
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
	view := mgl32.LookAt(0.0, 0.0, 12.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0)

	// Создание projection матрицы
	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(width)/height, 0.1, 100.0)

	// Основной цикл рендеринга
	for !window.ShouldClose() {
		drawScene(window, 
			programDefault, programStriped, 
			programChecker, programDiagonal, programHorizontal,
			pentagonVAO, cubeVAO, squareVAO, 
			view, projection)
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

	window, err := glfw.CreateWindow(width, height, "Pentagon, Cube and Procedural Cubes", nil, nil)
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
func drawScene(window *glfw.Window, 
	programDefault, programStriped, 
	programChecker, programDiagonal, programHorizontal uint32, 
	pentagonVAO, cubeVAO, squareVAO uint32, 
	view, projection mgl32.Mat4) {
	
	// Очистка экрана
	gl.ClearColor(0.2, 0.3, 0.3, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// ========== ВЕРХНИЙ РЯД ==========
	
	// Рисуем пятиугольник (слева вверху)
	drawPentagon(programDefault, pentagonVAO, view, projection, -4.0, 2.0, 0.0)
	
	// Рисуем куб (в центре вверху)
	drawCube(programDefault, cubeVAO, view, projection, 0.0, 2.0, 0.0)
	
	// Рисуем полосатый квадрат (справа вверху)
	drawStripedSquare(programStriped, squareVAO, view, projection, 4.0, 2.0, 0.0)

	// ========== НИЖНИЙ РЯД - ТРИ КУБИКА С ПРОЦЕДУРНЫМИ ТЕКСТУРАМИ ==========
	
	// Кубик с квадратиками (слева внизу)
	drawCube(programChecker, cubeVAO, view, projection, -4.0, -2.0, 0.0)
	
	// Кубик с диагональной штриховкой (в центре внизу)
	drawCube(programDiagonal, cubeVAO, view, projection, 0.0, -2.0, 0.0)
	
	// Кубик с горизонтальной полоской (справа внизу)
	drawCube(programHorizontal, cubeVAO, view, projection, 4.0, -2.0, 0.0)

	glfw.PollEvents()
	window.SwapBuffers()
}

// Вспомогательная функция для рисования пятиугольника
func drawPentagon(program uint32, vao uint32, view, projection mgl32.Mat4, x, y, z float32) {
	gl.UseProgram(program)
	
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	viewUniform := gl.GetUniformLocation(program, gl.Str("view\x00"))
	projUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	useTextureUniform := gl.GetUniformLocation(program, gl.Str("useTexture\x00"))
	
	gl.UniformMatrix4fv(viewUniform, 1, false, &view[0])
	gl.UniformMatrix4fv(projUniform, 1, false, &projection[0])
	gl.Uniform1i(useTextureUniform, 0)
	
	model := mgl32.Translate3D(x, y, z)
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
	
	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.TRIANGLES, 0, 15)
}

// Вспомогательная функция для рисования куба
func drawCube(program uint32, vao uint32, view, projection mgl32.Mat4, x, y, z float32) {
	gl.UseProgram(program)
	
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	viewUniform := gl.GetUniformLocation(program, gl.Str("view\x00"))
	projUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	useTextureUniform := gl.GetUniformLocation(program, gl.Str("useTexture\x00"))
	
	gl.UniformMatrix4fv(viewUniform, 1, false, &view[0])
	gl.UniformMatrix4fv(projUniform, 1, false, &projection[0])
	gl.Uniform1i(useTextureUniform, 0)
	
	model := mgl32.Translate3D(x, y, z)
	// Добавляем небольшой поворот для лучшей видимости граней
	rotation := mgl32.HomogRotate3D(mgl32.DegToRad(30.0), mgl32.Vec3{1, 1, 0}.Normalize())
	model = model.Mul4(rotation)
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
	
	gl.BindVertexArray(vao)
	gl.DrawElements(gl.TRIANGLES, 36, gl.UNSIGNED_INT, nil)
}

// Вспомогательная функция для рисования полосатого квадрата
func drawStripedSquare(program uint32, vao uint32, view, projection mgl32.Mat4, x, y, z float32) {
	gl.UseProgram(program)
	
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	viewUniform := gl.GetUniformLocation(program, gl.Str("view\x00"))
	projUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	useTextureUniform := gl.GetUniformLocation(program, gl.Str("useTexture\x00"))
	
	gl.UniformMatrix4fv(viewUniform, 1, false, &view[0])
	gl.UniformMatrix4fv(projUniform, 1, false, &projection[0])
	gl.Uniform1i(useTextureUniform, 0)
	
	model := mgl32.Translate3D(x, y, z)
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
	
	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
}