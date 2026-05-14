package utils

import (
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/venexene/comp-graphics-go/objects"
)

// RotationState - структура для хранения состояния поворотов
type RotationState struct {
    LocalRotation    float32 // Поворот каждого кубика вокруг локальной оси Y
    PiedestalRotation float32 // Поворот всего пьедестала вокруг его центра
    GlobalRotation    float32 // Поворот всего пьедестала вокруг глобальной оси Y
}

// Структура кубика пьедестала
type Cube struct {
	position mgl32.Vec3  // Позиция относительно центра пьедестала
	color    mgl32.Vec3  // Цвет кубика
}

// Создание пьедестала почета
func createPiedestal() []Cube {
	return []Cube{
		// 1 место - Золото (нижний кубик)
		{
			position: mgl32.Vec3{0.0, 0.0, 0.0},
			color:    mgl32.Vec3{1.0, 0.8, 0.0}, // Золотой
		},
		// 1 место - Золото (верхний кубик)
		{
			position: mgl32.Vec3{0.0, 1.0, 0.0},
			color:    mgl32.Vec3{1.0, 0.8, 0.0}, // Золотой
		},
		// 2 место - Серебро
		{
			position: mgl32.Vec3{-1.0, 0.0, 0.0},
			color:    mgl32.Vec3{0.8, 0.8, 0.9}, // Серебряный
		},
		// 3 место - Бронза
		{
			position: mgl32.Vec3{1.0, 0.0, 0.0},
			color:    mgl32.Vec3{0.8, 0.5, 0.2}, // Бронзовый
		},
	}
}

// Рендеринг
func DrawScene(window *glfw.Window, program uint32, vao uint32, view, projection mgl32.Mat4, rotationState *RotationState) {
	
	// Очистка экрана
	gl.ClearColor(0.2, 0.3, 0.3, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.UseProgram(program)

	// Получение uniform-переменных
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	viewUniform := gl.GetUniformLocation(program, gl.Str("view\x00"))
	projUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	colorUniform := gl.GetUniformLocation(program, gl.Str("cubeColor\x00"))

	// Передача view и projection матриц
	gl.UniformMatrix4fv(viewUniform, 1, false, &view[0])
	gl.UniformMatrix4fv(projUniform, 1, false, &projection[0])

	// Привязка VAO
	gl.BindVertexArray(vao)

	// Создание пьедестала
	cubes := createPiedestal()

	// Центр пьедестала в мировой системе координат
	piedestalCenter := mgl32.Vec3{2.0, 0.0, 0.0}

	// Глобальный поворот вокруг центра сцены (0,0,0)
	globalRotation := mgl32.HomogRotate3D(rotationState.GlobalRotation, mgl32.Vec3{0, 1, 0})

	// Поворот пьедестала вокруг его центра (локальная ось Y пьедестала)
	piedestalRotation := mgl32.HomogRotate3D(rotationState.PiedestalRotation, mgl32.Vec3{0, 1, 0})

	for _, cube := range cubes {
		// 1. Локальный поворот кубика вокруг своей оси
		localRotation := mgl32.HomogRotate3D(rotationState.LocalRotation, mgl32.Vec3{0, 1, 0})
		
		// 2. Матрица для позиции кубика относительно центра пьедестала
		translationToLocal := mgl32.Translate3D(cube.position.X(), cube.position.Y(), cube.position.Z())
		
		// 3. Сборка матрицы для кубика в локальных координатах пьедестала
		// Сначала локальный поворот, потом смещение к позиции на пьедестале
		cubeLocalTransform := translationToLocal.Mul4(localRotation)
		
		// 4. Применяем поворот пьедестала вокруг его центра
		piedestalTransform := piedestalRotation.Mul4(cubeLocalTransform)
		
		// 5. Смещаем пьедестал в нужную точку сцены
		translationToWorld := mgl32.Translate3D(piedestalCenter.X(), piedestalCenter.Y(), piedestalCenter.Z())
		
		// 6. Применяем глобальный поворот вокруг центра сцены
		finalTransform := globalRotation.Mul4(translationToWorld).Mul4(piedestalTransform)

		// Передаем матрицу модели в шейдер
		gl.UniformMatrix4fv(modelUniform, 1, false, &finalTransform[0])
		
		// Передаем цвет кубика
		gl.Uniform3f(colorUniform, cube.color.X(), cube.color.Y(), cube.color.Z())
		
		// Отрисовка кубика
		gl.DrawElements(gl.TRIANGLES, int32(len(objects.CubeIndices)), gl.UNSIGNED_INT, unsafe.Pointer(nil))
	}

	glfw.PollEvents()
	window.SwapBuffers()
}