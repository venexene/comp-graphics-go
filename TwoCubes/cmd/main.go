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
	"github.com/venexene/comp-graphics-go/textures"
	"github.com/venexene/comp-graphics-go/utils"
)


// Размеры окна
const (
    width  = 1000
    height = 500
)


func main() {
	// Блокировка потока для OpenGL
	// Это нужно, тк контекст в OpenGL привязан к потоку ОС, в котором он создан
	// Однако планировщик Go может переключать горутины между потоками
	// LockOSThread прикрепляет текущую горутину к одному потоку ОС
	runtime.LockOSThread()

	window := initGlfw() // Инициализация окна GLFW
	defer glfw.Terminate() // Опрделение отичстики GLFW при выходе из программы

	program := initOpenGL() // Инициализация OpenGL и шейдеров

	// VAO (Vertex Array Object) - объект, который хранит все настройки того, как интерпретировать данные в VBO (пресет для чтения вершинных данных)
	// VBO (Vertex Buffer Object) - буфер в видеопамяти GPU, который хранит сырые данные вершин (позиции, цвета, нормали, текстурные координаты)
	// EBO (Element Buffer Object) - буфер, хранящий индексы вершин, определяющие порядок соединения вершин в треугольники
	vao, vbo, ebo := objects.CreateCube() // Создание геомерии куба

	// Настройка состояния OpenGL
	// gl.DEPTH_TEST включает z-буфер, те позволяет граням перекрывать друг друга
	gl.Enable(gl.DEPTH_TEST)
	
	// Создание view матрицы
	// Позиция камеры (eye) - (3.0, 2.0, 3.0)
	// Точка взгляда камеры (center) - (0.0, 0.0, 0.0)
	// Вектор направления вверх (up) - (0.0, 1.0, 0.0)
	view := mgl32.LookAt(5.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0)

	// Создание projection матрицы
	// Угол обзора (fov) - 45 градусов
	// Соотношение сторон (aspect) - width/height 
	// Ближняя плоскость отсечения (near) - 0.1
	// Дальняя плоскость отсечения (far) - 100.0
	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(width)/height, 0.1, 100.0)

    texture := textures.LoadTexture("textures/texture.png") // Загрузка текстуры

	// Основной цикл рендеринга
	for !window.ShouldClose() { // Проверка условия закрытия окна
		// Вызов отрисовки
		utils.DrawScene(window, program, vao, view, projection, texture)
	}

	// Очистка ресурсов
	gl.DeleteVertexArrays(1, &vao)
	gl.DeleteBuffers(1, &vbo)
	gl.DeleteBuffers(1, &ebo)
}


// Инициализация окна GLFW
func initGlfw() *glfw.Window {
	// Инициализация GLFW
    if err := glfw.Init(); err != nil {
        panic(err)
    }
    
	// Настройка подсказок окна
    glfw.WindowHint(glfw.Resizable, glfw.False) // Отключение возможности менять размер окна
    glfw.WindowHint(glfw.ContextVersionMajor, 4) // Установка мажорной версии OpenGL
    glfw.WindowHint(glfw.ContextVersionMinor, 1) // Установка минорной версии OpenG:
    glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile) // Установка профиля OpenGL
    glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True) // Установка совместимоси версий OpenGL

	// Создание окна
    window, err := glfw.CreateWindow(width, height, "Two Cubes", nil, nil)
    if err != nil {
        panic(err)
    }
    window.MakeContextCurrent() // Привязка контекста OpenGL

    return window
}


// Инициализация OpenGL
func initOpenGL() uint32 {
	// Инициализация OpenGL
    if err := gl.Init(); err != nil {
        panic(err)
    }

	// Загрузка шейдеров из файлов
    vertexShaderSource, err := shaders.LoadShaderFile("shaders/vertex.glsl")
    if err != nil {
        panic(err)
    }

    fragmentShaderSource, err := shaders.LoadShaderFile("shaders/fragment.glsl")
    if err != nil {
        panic(err)
    }

	// Получение и логирование версии OpenGL
	// gl.GetString возвращает C-строку
	// gl.GoStr преобразует ее в Go-строку
    version := gl.GoStr(gl.GetString(gl.VERSION))
    log.Println("OpenGL version", version)

	// Компиляция вершинного шейдера
	vertexShader, err := shaders.CompileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	// Компиляция фрагментного шейдера
	fragmentShader, err := shaders.CompileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	// Создание программы и прикрепление шейдеров
	// В контексте OpenGL, программа (program) - это объект, который объединяет скомпилированные шейдеры
	// в полностью готовый к использованию рендеринговый конвейер.
	// Это "исполняемый файл" для GPU.
	program := gl.CreateProgram() // Создаем программу и получаем ее ID
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program) // Линковка программы

	// Проверка успешности линковки
	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status) // Функция получения информации о программе
	if status == gl.FALSE { // Проверка статуса
		// Получение размера лога ошибок
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := make([]byte, logLength) // Выделение памяти под лог
		gl.GetProgramInfoLog(program, logLength, nil, &log[0]) // Получение текста ошибки
		panic(fmt.Sprintf("Failed to link program: %s", string(log))) // Паника с текстом ошибки
	}

	// Очистка
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

    return program
}
