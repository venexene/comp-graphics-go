package utils

import (
	"time"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/venexene/comp-graphics-go/objects"
)

var startTime = time.Now()

// Структура для хранения состояния камеры
type Camera struct {
	Position  mgl32.Vec3
	Target    mgl32.Vec3
	Up        mgl32.Vec3
	Distance  float32
	RotationX float32
	RotationY float32
}

// Структура для хранения состояния куба
type CubeState struct {
	Position  mgl32.Vec3
	Scale     float32
	RotationX float32
	RotationY float32
	RotationZ float32
}

var (
	camera = Camera{
		Position:  mgl32.Vec3{0.0, 0.0, 5.0},
		Target:    mgl32.Vec3{0.0, 0.0, 0.0},
		Up:        mgl32.Vec3{0.0, 1.0, 0.0},
		Distance:  5.0,
		RotationX: 0.0,
		RotationY: 0.0,
	}


	cube = CubeState{
		Position:  mgl32.Vec3{0.0, 0.0, 0.0},
		Scale:     1.0,
		RotationX: 0.0,
		RotationY: 0.0,
		RotationZ: 0.0,
	}
)

// Обработка нажатий клавиш
func ProcessInput(window *glfw.Window) {
	// Управление кубом (WASD для перемещения в плоскости XZ)
	if window.GetKey(glfw.KeyW) == glfw.Press {
		cube.Position[2] -= 0.001 // Вперед
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		cube.Position[2] += 0.001 // Назад
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
		cube.Position[0] -= 0.001 // Влево
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		cube.Position[0] += 0.001 // Вправо
	}
	
	// Вертикальное перемещение куба
	if window.GetKey(glfw.KeySpace) == glfw.Press {
		cube.Position[1] += 0.001 // Вверх
	}
	if window.GetKey(glfw.KeyLeftShift) == glfw.Press || window.GetKey(glfw.KeyRightShift) == glfw.Press {
		cube.Position[1] -= 0.001 // Вниз
	}

	// Масштабирование куба (Q - уменьшить, E - увеличить)
	if window.GetKey(glfw.KeyQ) == glfw.Press {
		cube.Scale -= 0.001
		if cube.Scale < 0.1 {
			cube.Scale = 0.1
		}
	}
	if window.GetKey(glfw.KeyE) == glfw.Press {
		cube.Scale += 0.001
		if cube.Scale > 3.0 {
			cube.Scale = 3.0
		}
	}

	// Вращение куба (стрелки)
	if window.GetKey(glfw.KeyUp) == glfw.Press {
		cube.RotationX += 0.001
	}
	if window.GetKey(glfw.KeyDown) == glfw.Press {
		cube.RotationX -= 0.001
	}
	if window.GetKey(glfw.KeyLeft) == glfw.Press {
		cube.RotationY -= 0.001
	}
	if window.GetKey(glfw.KeyRight) == glfw.Press {
		cube.RotationY += 0.001
	}

	// Вращение вокруг оси Z (R - по часовой, F - против часовой)
	if window.GetKey(glfw.KeyR) == glfw.Press {
		cube.RotationZ += 0.001
	}
	if window.GetKey(glfw.KeyF) == glfw.Press {
		cube.RotationZ -= 0.001
	}


	// Приближение/отдаление камеры
	if window.GetKey(glfw.KeyKPAdd) == glfw.Press || window.GetKey(glfw.KeyEqual) == glfw.Press {
		camera.Distance -= 0.1
		if camera.Distance < 2.0 {
			camera.Distance = 2.0
		}
	}
	if window.GetKey(glfw.KeyKPSubtract) == glfw.Press || window.GetKey(glfw.KeyMinus) == glfw.Press {
		camera.Distance += 0.1
		if camera.Distance > 20.0 {
			camera.Distance = 20.0
		}
	}

	// Сброс позиции куба (R)
	if window.GetKey(glfw.KeyR) == glfw.Press && 
	   window.GetKey(glfw.KeyLeftControl) == glfw.Press {
		cube.Position = mgl32.Vec3{0.0, 0.0, 0.0}
		cube.Scale = 1.0
		cube.RotationX = 0.0
		cube.RotationY = 0.0
		cube.RotationZ = 0.0
	}
}

// Рендеринг
func DrawScene(window *glfw.Window, program uint32, vao uint32, view, projection mgl32.Mat4) {
	// Обработка ввода
	ProcessInput(window)

	// Очистка экрана
	gl.ClearColor(0.2, 0.3, 0.3, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.UseProgram(program)

	// Получение uniform-переменных
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	viewUniform := gl.GetUniformLocation(program, gl.Str("view\x00"))
	projUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	useTextureUniform := gl.GetUniformLocation(program, gl.Str("useTexture\x00"))

	// Обновление view матрицы на основе позиции камеры
	view = mgl32.LookAt(
		camera.Position[0], camera.Position[1], camera.Position[2],
		camera.Target[0], camera.Target[1], camera.Target[2],
		camera.Up[0], camera.Up[1], camera.Up[2],
	)

	// Передача матриц в шейдеры
	gl.UniformMatrix4fv(viewUniform, 1, false, &view[0])
	gl.UniformMatrix4fv(projUniform, 1, false, &projection[0])

	// Привязка VAO
	gl.BindVertexArray(vao)

	// СОЗДАНИЕ МАТРИЦЫ МОДЕЛИ ДЛЯ КУБА
	// Масштабирование
	scale := mgl32.Scale3D(cube.Scale, cube.Scale, cube.Scale)
	
	// Вращение
	rotationX := mgl32.HomogRotate3D(cube.RotationX, mgl32.Vec3{1, 0, 0})
	rotationY := mgl32.HomogRotate3D(cube.RotationY, mgl32.Vec3{0, 1, 0})
	rotationZ := mgl32.HomogRotate3D(cube.RotationZ, mgl32.Vec3{0, 0, 1})
	rotation := rotationZ.Mul4(rotationY.Mul4(rotationX))
	
	// Перемещение
	translation := mgl32.Translate3D(cube.Position.X(), cube.Position.Y(), cube.Position.Z())
	
	// Комбинированная матрица модели: translation * rotation * scale
	model := translation.Mul4(rotation.Mul4(scale))

	// Используем только цвет, без текстуры
	gl.Uniform1i(useTextureUniform, 0)
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
	
	// Отрисовка куба
	gl.DrawElements(gl.TRIANGLES, int32(len(objects.CubeIndices)), gl.UNSIGNED_INT, unsafe.Pointer(nil))

	glfw.PollEvents()
	window.SwapBuffers()
}