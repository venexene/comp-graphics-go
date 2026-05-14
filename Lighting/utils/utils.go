package utils

import (
	"math"
	"time"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/venexene/comp-graphics-go/objects"
)

var startTime = time.Now()

// Структура для хранения состояния камеры
type Camera struct {
	Target   mgl32.Vec3
	Distance float32
	Yaw      float32 // горизонтальное вращение (вокруг Y)
	Pitch    float32 // вертикальное вращение
}

// Структура для хранения состояния объекта
type ObjectState struct {
	Position  mgl32.Vec3
	Scale     float32
	RotationX float32
	RotationY float32
	RotationZ float32
}

var (
	camera = Camera{
		Target:   mgl32.Vec3{0.0, 0.0, 0.0},
		Distance: 5.0,
		Yaw:      0.0,
		Pitch:    0.0,
	}

	objectState = ObjectState{
		Position:  mgl32.Vec3{0.0, 0.0, 0.0},
		Scale:     1.0,
		RotationX: 0.0,
		RotationY: 0.0,
		RotationZ: 0.0,
	}
)

// Обработка нажатий клавиш
func ProcessInput(window *glfw.Window) {
	// Управление камерой (WASD для панорамирования цели)
	offset := mgl32.Vec3{
		float32(math.Cos(float64(camera.Yaw)) * math.Cos(float64(camera.Pitch))),
		float32(math.Sin(float64(camera.Pitch))),
		float32(math.Sin(float64(camera.Yaw)) * math.Cos(float64(camera.Pitch))),
	}
	right := mgl32.Vec3{
		float32(-math.Sin(float64(camera.Yaw))),
		0,
		float32(math.Cos(float64(camera.Yaw))),
	}

	if window.GetKey(glfw.KeyW) == glfw.Press {
		// Панорамирование вперед
		dir := offset.Normalize().Mul(-0.01)
		camera.Target = camera.Target.Add(dir)
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		// Панорамирование назад
		dir := offset.Normalize().Mul(0.01)
		camera.Target = camera.Target.Add(dir)
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
		// Панорамирование влево
		dir := right.Normalize().Mul(-0.01)
		camera.Target = camera.Target.Add(dir)
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		// Панорамирование вправо
		dir := right.Normalize().Mul(0.01)
		camera.Target = camera.Target.Add(dir)
	}
	
	// Вертикальное панорамирование камеры
	if window.GetKey(glfw.KeySpace) == glfw.Press {
		camera.Target = camera.Target.Add(mgl32.Vec3{0, 0.01, 0})
	}
	if window.GetKey(glfw.KeyLeftShift) == glfw.Press || window.GetKey(glfw.KeyRightShift) == glfw.Press {
		camera.Target = camera.Target.Add(mgl32.Vec3{0, -0.01, 0})
	}

	// Масштабирование объекта (Q - уменьшить, E - увеличить)
	if window.GetKey(glfw.KeyQ) == glfw.Press {
		objectState.Scale -= 0.001
		if objectState.Scale < 0.1 {
			objectState.Scale = 0.1
		}
	}
	if window.GetKey(glfw.KeyE) == glfw.Press {
		objectState.Scale += 0.001
		if objectState.Scale > 3.0 {
			objectState.Scale = 3.0
		}
	}

	// Перемещение объекта (IJKL для XZ, UO для Y)
	if window.GetKey(glfw.KeyI) == glfw.Press {
		objectState.Position[2] -= 0.01 // Вперед
	}
	if window.GetKey(glfw.KeyK) == glfw.Press {
		objectState.Position[2] += 0.01 // Назад
	}
	if window.GetKey(glfw.KeyJ) == glfw.Press {
		objectState.Position[0] -= 0.01 // Влево
	}
	if window.GetKey(glfw.KeyL) == glfw.Press {
		objectState.Position[0] += 0.01 // Вправо
	}
	if window.GetKey(glfw.KeyU) == glfw.Press {
		objectState.Position[1] += 0.01 // Вверх
	}
	if window.GetKey(glfw.KeyO) == glfw.Press {
		objectState.Position[1] -= 0.01 // Вниз
	}

	// Вращение камеры (стрелки)
	if window.GetKey(glfw.KeyUp) == glfw.Press {
		camera.Pitch -= 0.01
		if camera.Pitch < -math.Pi/2 + 0.1 {
			camera.Pitch = -math.Pi/2 + 0.1
		}
	}
	if window.GetKey(glfw.KeyDown) == glfw.Press {
		camera.Pitch += 0.01
		if camera.Pitch > math.Pi/2 - 0.1 {
			camera.Pitch = math.Pi/2 - 0.1
		}
	}
	if window.GetKey(glfw.KeyLeft) == glfw.Press {
		camera.Yaw -= 0.01
	}
	if window.GetKey(glfw.KeyRight) == glfw.Press {
		camera.Yaw += 0.01
	}

	// Вращение объекта вокруг оси Z (R - по часовой, F - против часовой)
	if window.GetKey(glfw.KeyR) == glfw.Press {
		objectState.RotationZ += 0.001
	}
	if window.GetKey(glfw.KeyF) == glfw.Press {
		objectState.RotationZ -= 0.001
	}


	// Приближение/отдаление камеры
	if window.GetKey(glfw.KeyKPAdd) == glfw.Press || window.GetKey(glfw.KeyEqual) == glfw.Press {
		camera.Distance -= 0.1
		if camera.Distance < 1.0 {
			camera.Distance = 1.0
		}
	}
	if window.GetKey(glfw.KeyKPSubtract) == glfw.Press || window.GetKey(glfw.KeyMinus) == glfw.Press {
		camera.Distance += 0.1
		if camera.Distance > 50.0 {
			camera.Distance = 50.0
		}
	}

	// Сброс позиции объекта и камеры (Ctrl+R)
	if window.GetKey(glfw.KeyR) == glfw.Press && 
	   window.GetKey(glfw.KeyLeftControl) == glfw.Press {
		objectState.Position = mgl32.Vec3{0.0, 0.0, 0.0}
		objectState.Scale = 1.0
		objectState.RotationX = 0.0
		objectState.RotationY = 0.0
		objectState.RotationZ = 0.0
		camera.Target = mgl32.Vec3{0.0, 0.0, 0.0}
		camera.Distance = 5.0
		camera.Yaw = 0.0
		camera.Pitch = 0.0
	}
}

// Рендеринг
func DrawScene(window *glfw.Window, program uint32, model *objects.Model, view, projection mgl32.Mat4) {
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

	// Вычисление позиции камеры на основе сферических координат
	x := camera.Distance * float32(math.Cos(float64(camera.Yaw))) * float32(math.Cos(float64(camera.Pitch)))
	y := camera.Distance * float32(math.Sin(float64(camera.Pitch)))
	z := camera.Distance * float32(math.Sin(float64(camera.Yaw))) * float32(math.Cos(float64(camera.Pitch)))
	
	cameraPosition := camera.Target.Add(mgl32.Vec3{x, y, z})

	// Обновление view матрицы на основе позиции камеры
	view = mgl32.LookAt(
		cameraPosition.X(), cameraPosition.Y(), cameraPosition.Z(),
		camera.Target.X(), camera.Target.Y(), camera.Target.Z(),
		0.0, 1.0, 0.0, // up vector
	)

	// Передача матриц в шейдеры
	gl.UniformMatrix4fv(viewUniform, 1, false, &view[0])
	gl.UniformMatrix4fv(projUniform, 1, false, &projection[0])

	// СОЗДАНИЕ МАТРИЦЫ МОДЕЛИ
	scale := mgl32.Scale3D(objectState.Scale, objectState.Scale, objectState.Scale)
	
	// Вращение только вокруг Z
	rotationZ := mgl32.HomogRotate3D(objectState.RotationZ, mgl32.Vec3{0, 0, 1})
	
	// Перемещение
	translation := mgl32.Translate3D(objectState.Position.X(), objectState.Position.Y(), objectState.Position.Z())
	
	// Комбинированная матрица модели: translation * rotation * scale
	modelMat := translation.Mul4(rotationZ.Mul4(scale))

	// Используем только цвет
	gl.UniformMatrix4fv(modelUniform, 1, false, &modelMat[0])

	// Отрисовка загруженной модели
	model.Draw()

	glfw.PollEvents()
	window.SwapBuffers()
}