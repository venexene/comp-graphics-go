package utils

import (
	"time"
	"unsafe"
	"math"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/venexene/comp-graphics-go/objects"
)

var startTime = time.Now() // Время начала программы для поворотов

// Получение угла поворота на основе времени начала прогрмамы
func getRotationAngle(rotationSpeed float64) float32 {
    elapsedSeconds := time.Since(startTime).Seconds()
    // Просто линейно растущий угол без сброса
    angle := elapsedSeconds * rotationSpeed * 2 * math.Pi
    return float32(angle)
}

// Рендеринг
func DrawScene(window *glfw.Window, program uint32, vao uint32, view, projection mgl32.Mat4, texture uint32) {
	// Очистка экрана
	gl.ClearColor(0.2, 0.3, 0.3, 1.0) // Цвет очистки
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT) // Очистка буфера

	gl.UseProgram(program) // Активация программы шейдера

	// Получение uniform-переменных
	// gl.Str("model\x00") - создает C-строку с нулевым байтом в конце
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	viewUniform := gl.GetUniformLocation(program, gl.Str("view\x00"))
	projUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	useTextureUniform := gl.GetUniformLocation(program, gl.Str("useTexture\x00"))

	// Передача матриц в шейдеры
	// modelUniform - ID uniform-переменной
	// 1 - количество матриц
	// false - нужно ли транспонировать матрицу (для OpenGL обычно false)
	// &model[0] - указатель на первый элемент матрицы (данные передаются в GPU)
	gl.UniformMatrix4fv(viewUniform, 1, false, &view[0])
	gl.UniformMatrix4fv(projUniform, 1, false, &projection[0])

	// Привязка VAO и отрисовка
	gl.BindVertexArray(vao) // Активация настроек (VAO) вершин


	//ПЕРВЫЙ КУБ

	// Создание матрицы model для поворотов
	// Получение угла поворота
	angle := getRotationAngle(0.5)
	rotationY := mgl32.HomogRotate3D(angle, mgl32.Vec3{0, 1, 0}) // Поворот вокруг оси Y 
	rotationX := mgl32.HomogRotate3D(angle*0.7, mgl32.Vec3{1, 0, 0}) // Поворот вокруг оси X
	rotation := rotationY.Mul4(rotationX) // Комбинация поворотов
	translation := mgl32.Translate3D(0.0, 0.0, 2.0) // Смещение влево
	model1 := translation.Mul4(rotation) // Комбинация поворота и смещения
	gl.Uniform1i(useTextureUniform, 0) // Оповещещение о том, что не используется текстура
	gl.UniformMatrix4fv(modelUniform, 1, false, &model1[0])
	gl.DrawElements(gl.TRIANGLES, int32(len(objects.CubeIndices)), gl.UNSIGNED_INT, unsafe.Pointer(nil))


	//ВТОРОЙ КУБ

	// Создание матрицы model для поворотов
	// Получение угла поворота
	angle = getRotationAngle(0.5)
	rotationY = mgl32.HomogRotate3D(angle, mgl32.Vec3{0, 1, 0}) // Поворот вокруг оси Y 
	rotationX = mgl32.HomogRotate3D(angle*0.7, mgl32.Vec3{1, 0, 0}) // Поворот вокруг оси X
	rotation = rotationY.Mul4(rotationX) // Комбинация поворотов
	translation = mgl32.Translate3D(0.0, 0.0, -2.0) // Смещение влево
	model2 := translation.Mul4(rotation) // Комбинация поворота и смещения
	gl.ActiveTexture(gl.TEXTURE0) // Активация текстуры
    gl.BindTexture(gl.TEXTURE_2D, texture) // Привязка текстуры
	gl.Uniform1i(useTextureUniform, 1) // Оповещещение о том, что текстура используется
	gl.UniformMatrix4fv(modelUniform, 1, false, &model2[0])
	gl.DrawElements(gl.TRIANGLES, int32(len(objects.CubeIndices)), gl.UNSIGNED_INT, unsafe.Pointer(nil))

	glfw.PollEvents() // Проверка и обработка событий окна
	window.SwapBuffers() // Смена заднего и переднего буфера (двойная буферизация)
}