package utils

import (
	"math"
	"time"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/venexene/comp-graphics-go/objects"
	"github.com/venexene/comp-graphics-go/shaders"
)

var startTime = time.Now()

var (
	lastTKeyState = false
	lastGKeyState = false
	lastYKeyState = false
)

var (
	defaultTexture    uint32
	lightPosition     = mgl32.Vec3{2.0, 4.0, 3.0}
	lightConstant     float32 = 1.0
	lightLinear       float32 = 0.09
	lightQuadratic    float32 = 0.032
	lightLinearCoef   float32 = 0.5
	lightQuadraticCoef float32 = 0.5
	ambientStrength   float32 = 0.6
)

// CreateWhiteTexture creates a simple 1x1 white texture for the default material
func CreateWhiteTexture() uint32 {
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	// Create a simple 1x1 white pixel
	whitePixel := []uint8{255, 255, 255, 255}
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, 1, 1, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(whitePixel))

	// Set texture parameters
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	gl.BindTexture(gl.TEXTURE_2D, 0)
	return texture
}

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
	camera      = Camera{Target: mgl32.Vec3{0, 0, 0}, Distance: 5.0, Yaw: 0.0, Pitch: 0.0}
	objectState = ObjectState{Position: mgl32.Vec3{0, 0, 0}, Scale: 1.0, RotationX: 0.0, RotationY: 0.0, RotationZ: 0.0}
)

// InitScene initializes the rendering scene (textures, etc)
func InitScene() {
	defaultTexture = CreateWhiteTexture()
}

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

	// Переключение модели освещения и режима затенения
	currentTKeyState := window.GetKey(glfw.KeyT) == glfw.Press
	currentGKeyState := window.GetKey(glfw.KeyG) == glfw.Press
	currentYKeyState := window.GetKey(glfw.KeyY) == glfw.Press

	if currentTKeyState && !lastTKeyState {
		shaders.CycleLightingVariant(true)
	}
	if currentGKeyState && !lastGKeyState {
		shaders.CycleLightingVariant(false)
	}
	if currentYKeyState && !lastYKeyState {
		shaders.ToggleShadingMode()
	}

	lastTKeyState = currentTKeyState
	lastGKeyState = currentGKeyState
	lastYKeyState = currentYKeyState

	// Уровни затухания света
	if window.GetKey(glfw.KeyZ) == glfw.Press {
		lightLinearCoef -= 0.01
		if lightLinearCoef < 0.0 {
			lightLinearCoef = 0.0
		}
	}
	if window.GetKey(glfw.KeyX) == glfw.Press {
		lightLinearCoef += 0.01
	}
	if window.GetKey(glfw.KeyC) == glfw.Press {
		lightQuadraticCoef -= 0.01
		if lightQuadraticCoef < 0.0 {
			lightQuadraticCoef = 0.0
		}
	}
	if window.GetKey(glfw.KeyV) == glfw.Press {
		lightQuadraticCoef += 0.01
	}

	// Фоновый свет
	if window.GetKey(glfw.KeyB) == glfw.Press {
		ambientStrength -= 0.01
		if ambientStrength < 0.0 {
			ambientStrength = 0.0
		}
	}
	if window.GetKey(glfw.KeyN) == glfw.Press {
		ambientStrength += 0.01
		if ambientStrength > 1.0 {
			ambientStrength = 1.0
		}
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
func DrawScene(window *glfw.Window, model *objects.Model, view, projection mgl32.Mat4) {
	// Обработка ввода
	ProcessInput(window)

	// Очистка экрана
	gl.ClearColor(0.2, 0.3, 0.3, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// Получить текущую программу освещения
	program := shaders.GetCurrentLightingProgram()
	gl.UseProgram(program)

	// Получение uniform-переменных
	modelUniform := gl.GetUniformLocation(program, gl.Str("transform.model\x00"))
	viewUniform := gl.GetUniformLocation(program, gl.Str("transform.view\x00"))
	projUniform := gl.GetUniformLocation(program, gl.Str("transform.projection\x00"))
	normalUniform := gl.GetUniformLocation(program, gl.Str("transform.normal_mat\x00"))
	viewPosUniform := gl.GetUniformLocation(program, gl.Str("transform.view_pos\x00"))
	materialAmbientUniform := gl.GetUniformLocation(program, gl.Str("material.ambient\x00"))
	materialDiffuseUniform := gl.GetUniformLocation(program, gl.Str("material.diffuse\x00"))
	materialSpecularUniform := gl.GetUniformLocation(program, gl.Str("material.specular\x00"))
	materialSheenUniform := gl.GetUniformLocation(program, gl.Str("material.sheen_coef\x00"))
	lightAmbientUniform := gl.GetUniformLocation(program, gl.Str("light.ambient\x00"))
	lightDiffuseUniform := gl.GetUniformLocation(program, gl.Str("light.diffuse\x00"))
	lightSpecularUniform := gl.GetUniformLocation(program, gl.Str("light.specular\x00"))
	lightPositionUniform := gl.GetUniformLocation(program, gl.Str("light.position\x00"))
	lightConstantUniform := gl.GetUniformLocation(program, gl.Str("light.constant\x00"))
	lightLinearUniform := gl.GetUniformLocation(program, gl.Str("light.linear\x00"))
	lightQuadraticUniform := gl.GetUniformLocation(program, gl.Str("light.quadratic\x00"))
	lightAmbientStrengthUniform := gl.GetUniformLocation(program, gl.Str("light.ambient_strength\x00"))
	linearCoefUniform := gl.GetUniformLocation(program, gl.Str("linear_coef\x00"))
	quadraticCoefUniform := gl.GetUniformLocation(program, gl.Str("quadratic_coef\x00"))
	diffuseMapUniform := gl.GetUniformLocation(program, gl.Str("diffuse_map\x00"))

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
	if viewPosUniform != -1 {
		gl.Uniform3f(viewPosUniform, cameraPosition.X(), cameraPosition.Y(), cameraPosition.Z())
	}

	// СОЗДАНИЕ МАТРИЦЫ МОДЕЛИ
	scale := mgl32.Scale3D(objectState.Scale, objectState.Scale, objectState.Scale)
	
	// Вращение только вокруг Z
	rotationZ := mgl32.HomogRotate3D(objectState.RotationZ, mgl32.Vec3{0, 0, 1})
	
	// Перемещение
	translation := mgl32.Translate3D(objectState.Position.X(), objectState.Position.Y(), objectState.Position.Z())
	
	// Комбинированная матрица модели: translation * rotation * scale
	modelMat := translation.Mul4(rotationZ.Mul4(scale))

	// Используем только цвет
	if modelUniform != -1 {
		gl.UniformMatrix4fv(modelUniform, 1, false, &modelMat[0])
	}
	if normalUniform != -1 {
		// Extract the upper-left 3x3 and compute the normal matrix
		m := modelMat.Inv().Transpose()
		normalMat := mgl32.Mat3{
			m[0], m[1], m[2],
			m[4], m[5], m[6],
			m[8], m[9], m[10],
		}
		gl.UniformMatrix3fv(normalUniform, 1, false, &normalMat[0])
	}

	if materialAmbientUniform != -1 {
		gl.Uniform3f(materialAmbientUniform, 0.2, 0.2, 0.2)
	}
	if materialDiffuseUniform != -1 {
		gl.Uniform3f(materialDiffuseUniform, 1.0, 1.0, 1.0)
	}
	if materialSpecularUniform != -1 {
		gl.Uniform3f(materialSpecularUniform, 1.0, 1.0, 1.0)
	}
	if materialSheenUniform != -1 {
		gl.Uniform1f(materialSheenUniform, 32.0)
	}
	if lightAmbientUniform != -1 {
		gl.Uniform3f(lightAmbientUniform, 0.2, 0.2, 0.2)
	}
	if lightDiffuseUniform != -1 {
		gl.Uniform3f(lightDiffuseUniform, 1.0, 1.0, 1.0)
	}
	if lightSpecularUniform != -1 {
		gl.Uniform3f(lightSpecularUniform, 1.0, 1.0, 1.0)
	}
	if lightPositionUniform != -1 {
		gl.Uniform3f(lightPositionUniform, lightPosition.X(), lightPosition.Y(), lightPosition.Z())
	}
	if lightConstantUniform != -1 {
		gl.Uniform1f(lightConstantUniform, lightConstant)
	}
	if lightLinearUniform != -1 {
		gl.Uniform1f(lightLinearUniform, lightLinear)
	}
	if lightQuadraticUniform != -1 {
		gl.Uniform1f(lightQuadraticUniform, lightQuadratic)
	}
	if lightAmbientStrengthUniform != -1 {
		gl.Uniform1f(lightAmbientStrengthUniform, ambientStrength)
	}
	if linearCoefUniform != -1 {
		gl.Uniform1f(linearCoefUniform, lightLinearCoef)
	}
	if quadraticCoefUniform != -1 {
		gl.Uniform1f(quadraticCoefUniform, lightQuadraticCoef)
	}
	if diffuseMapUniform != -1 {
		gl.Uniform1i(diffuseMapUniform, 0)
	}

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, defaultTexture)

	// Отрисовка загруженной модели
	model.Draw()

	glfw.PollEvents()
	window.SwapBuffers()
}