// scene/camera.go — Орбитальная камера (сферические координаты).
//
// Назначение: реализует камеру, вращающуюся вокруг целевой точки (орбита).
// Использует сферические координаты: расстояние, рыскание (yaw), тангаж (pitch).
//
// Ключевые типы:
//   Camera — камера с полями Target, Distance, Yaw, Pitch.
//
// Ключевые функции:
//   DefaultCamera()  — возвращает камеру в начальной позиции.
//   EyePosition()    — вычисляет позицию камеры в мировом пространстве.
//   ViewMatrix()     — возвращает LookAt-матрицу.
//   PanForward/Right/Up() — перемещение целевой точки.
//   Zoom()           — изменение расстояния до цели.
//   Rotate()         — вращение камеры (yaw/pitch).
//
// Зависимости: используется в scene.DrawScene() и input.ProcessInput().
package scene

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

// Camera — орбитальная камера.
// Хранит: целевую точку (Target), расстояние до неё (Distance),
// угол рыскания (Yaw) вокруг оси Y и угол тангажа (Pitch) относительно
// горизонтали. Все углы в радианах.
// Позиция камеры вычисляется как Target + сферические координаты.
// Начальное положение: Distance=5, Yaw=0, Pitch=0 — камера смотрит
// вдоль +Z на начало координат.
type Camera struct {
	// Target — точка, вокруг которой вращается камера (world space).
	Target mgl32.Vec3
	// Distance — расстояние от камеры до Target.
	// Диапазон: [1.0, 50.0], начальное значение 5.0.
	Distance float32
	// Yaw — горизонтальный угол (рыскание), радианы.
	// Вращение вокруг оси Y по часовой стрелке.
	Yaw float32
	// Pitch — вертикальный угол (тангаж), радианы.
	// Диапазон: от -π/2+0.1 до π/2-0.1 (ограничение gimbal lock).
	Pitch float32
}

// DefaultCamera — возвращает камеру по умолчанию.
// Камера находится на расстоянии 5 от начала координат, смотрит на (0,0,0).
func DefaultCamera() Camera {
	return Camera{
		Target:   mgl32.Vec3{0, 0, 0},
		Distance: 5.0,
		Yaw:      0.0,
		Pitch:    0.0,
	}
}

// EyePosition — вычисляет позицию камеры в мировом пространстве.
// Использует сферические координаты:
//   x = Distance * cos(Yaw) * cos(Pitch)
//   y = Distance * sin(Pitch)
//   z = Distance * sin(Yaw) * cos(Pitch)
// Возвращает: Target + (x,y,z) — позиция камеры в world space.
func (c *Camera) EyePosition() mgl32.Vec3 {
	x := c.Distance * float32(math.Cos(float64(c.Yaw))) * float32(math.Cos(float64(c.Pitch)))
	y := c.Distance * float32(math.Sin(float64(c.Pitch)))
	z := c.Distance * float32(math.Sin(float64(c.Yaw))) * float32(math.Cos(float64(c.Pitch)))
	return c.Target.Add(mgl32.Vec3{x, y, z})
}

// ViewMatrix — возвращает LookAt-матрицу вида.
// Вычисляет позицию глаза через EyePosition() и строит матрицу
// по правилу: LookAt(eye, target, up=(0,1,0)).
// Возвращает: матрицу 4×4 для преобразования World Space → View Space.
func (c *Camera) ViewMatrix() mgl32.Mat4 {
	eye := c.EyePosition()
	return mgl32.LookAt(
		eye.X(), eye.Y(), eye.Z(),
		c.Target.X(), c.Target.Y(), c.Target.Z(),
		0.0, 1.0, 0.0,
	)
}

// PanForward — перемещает целевую точку (Target) вдоль направления взгляда камеры.
// Принимает: amount — величина смещения (отрицательное = вперёд, положительное = назад).
// Побочные эффекты: изменяет Target.
func (c *Camera) PanForward(amount float32) {
	offset := mgl32.Vec3{
		float32(math.Cos(float64(c.Yaw)) * math.Cos(float64(c.Pitch))),
		float32(math.Sin(float64(c.Pitch))),
		float32(math.Sin(float64(c.Yaw)) * math.Cos(float64(c.Pitch))),
	}
	c.Target = c.Target.Add(offset.Normalize().Mul(amount))
}

// PanRight — перемещает целевую точку вправо относительно камеры.
// Принимает: amount — величина смещения (отрицательное = влево).
// Побочные эффекты: изменяет Target.
func (c *Camera) PanRight(amount float32) {
	right := mgl32.Vec3{
		float32(-math.Sin(float64(c.Yaw))),
		0,
		float32(math.Cos(float64(c.Yaw))),
	}
	c.Target = c.Target.Add(right.Normalize().Mul(amount))
}

// PanUp — перемещает целевую точку вверх по мировой оси Y.
// Принимает: amount — величина смещения (отрицательное = вниз).
// Побочные эффекты: изменяет Target.
func (c *Camera) PanUp(amount float32) {
	c.Target = c.Target.Add(mgl32.Vec3{0, amount, 0})
}

// Zoom — изменяет расстояние камеры до целевой точки.
// Принимает: delta — изменение расстояния; minDist, maxDist — ограничения.
// Побочные эффекты: изменяет Distance. Значение не выходит за [minDist, maxDist].
func (c *Camera) Zoom(delta, minDist, maxDist float32) {
	c.Distance += delta
	if c.Distance < minDist {
		c.Distance = minDist
	}
	if c.Distance > maxDist {
		c.Distance = maxDist
	}
}

// Rotate — поворачивает камеру (рыскание и тангаж).
// Принимает: dYaw — изменение угла рыскания; dPitch — изменение тангажа.
// Тангаж ограничен (-π/2+0.1, π/2-0.1) для предотвращения gimbal lock.
// Побочные эффекты: изменяет Yaw и Pitch.
func (c *Camera) Rotate(dYaw, dPitch float32) {
	c.Yaw += dYaw
	c.Pitch += dPitch
	if c.Pitch < -math.Pi/2+0.1 {
		c.Pitch = -math.Pi/2 + 0.1
	}
	if c.Pitch > math.Pi/2-0.1 {
		c.Pitch = math.Pi/2 - 0.1
	}
}
