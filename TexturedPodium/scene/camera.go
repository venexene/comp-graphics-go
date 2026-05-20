// Пакет scene содержит типы для управления сценой:
// Camera — орбитальная камера (сферические координаты),
// Podium — пьедестал почёта из четырёх кубиков,
// PodiumCube — один кубик с материалами и цветом.
// Зависимости: go-gl/mathgl — матрицы и векторы.
package scene

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

// Camera — орбитальная камера, вращается вокруг точки Target.
// Управляется: стрелки (вращение), WASD (панорама), +/- (зум).
// Система координат: сферическая (Yaw, Pitch, Distance) с центром в Target.
type Camera struct {
	Target   mgl32.Vec3 // точка, вокруг которой вращается камера (world space)
	Distance float32    // расстояние от камеры до Target [1.0, 50.0]
	Yaw      float32    // горизонтальный угол (радианы), 0 = взгляд вдоль +Z
	Pitch    float32    // вертикальный угол (радианы), 0 = горизонтально, Pi/2 = сверху
}

// DefaultCamera возвращает камеру с видом на пьедестал спереди-справа-сверху.
func DefaultCamera() Camera {
	return Camera{
		Target:   mgl32.Vec3{0, 0.5, 0},
		Distance: 6.0,
		Yaw:      0.7,
		Pitch:    0.35,
	}
}

// EyePosition вычисляет позицию камеры из сферических координат.
// Формула: Eye = Target + Distance * (cos(Yaw)*cos(Pitch), sin(Pitch), sin(Yaw)*cos(Pitch))
func (c *Camera) EyePosition() mgl32.Vec3 {
	x := c.Distance * float32(math.Cos(float64(c.Yaw))) * float32(math.Cos(float64(c.Pitch)))
	y := c.Distance * float32(math.Sin(float64(c.Pitch)))
	z := c.Distance * float32(math.Sin(float64(c.Yaw))) * float32(math.Cos(float64(c.Pitch)))
	return c.Target.Add(mgl32.Vec3{x, y, z})
}

// ViewMatrix возвращает матрицу View = LookAt(eye, target, up).
func (c *Camera) ViewMatrix() mgl32.Mat4 {
	eye := c.EyePosition()
	return mgl32.LookAt(
		eye.X(), eye.Y(), eye.Z(),
		c.Target.X(), c.Target.Y(), c.Target.Z(),
		0.0, 1.0, 0.0,
	)
}

// PanForward смещает Target вдоль направления взгляда камеры.
func (c *Camera) PanForward(amount float32) {
	offset := mgl32.Vec3{
		float32(math.Cos(float64(c.Yaw)) * math.Cos(float64(c.Pitch))),
		float32(math.Sin(float64(c.Pitch))),
		float32(math.Sin(float64(c.Yaw)) * math.Cos(float64(c.Pitch))),
	}
	c.Target = c.Target.Add(offset.Normalize().Mul(amount))
}

// PanRight смещает Target вправо от направления взгляда.
func (c *Camera) PanRight(amount float32) {
	right := mgl32.Vec3{
		float32(-math.Sin(float64(c.Yaw))),
		0,
		float32(math.Cos(float64(c.Yaw))),
	}
	c.Target = c.Target.Add(right.Normalize().Mul(amount))
}

// PanUp смещает Target вдоль мировой оси Y.
func (c *Camera) PanUp(amount float32) {
	c.Target = c.Target.Add(mgl32.Vec3{0, amount, 0})
}

// Zoom изменяет расстояние камеры до цели, ограниченное [minDist, maxDist].
func (c *Camera) Zoom(delta, minDist, maxDist float32) {
	c.Distance += delta
	if c.Distance < minDist {
		c.Distance = minDist
	}
	if c.Distance > maxDist {
		c.Distance = maxDist
	}
}

// Rotate поворачивает камеру вокруг цели. Pitch ограничен ±(Pi/2 - 0.1)
// во избежание gimbal lock.
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
