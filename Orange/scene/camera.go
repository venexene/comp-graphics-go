// Файл: scene/camera.go
// Назначение: реализация орбитальной камеры (arcball) для обзора сцены.
//
// Ключевые типы:
//   Camera — орбитальная камера с target, distance, yaw и pitch.
//
// Ключевые функции:
//   DefaultCamera — возвращает камеру с начальной позицией.
//   EyePosition — вычисляет позицию камеры в world space.
//   ViewMatrix — строит матрицу View через mgl32.LookAt.
//   PanForward/PanRight/PanUp — перемещение точки взгляда (target).
//   Zoom — изменение расстояния до target.
//   Rotate — изменение углов yaw/pitch (орбитальное вращение).
//
// Зависимости:
//   Внутренние: — (используется scene/scene.go, input/input.go, utils/utils.go).
//   Внешние: github.com/go-gl/mathgl/mgl32.

package scene

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

// Camera — орбитальная камера (target-distance-yaw-pitch).
// Вращение камеры — это вращение вокруг точки Target на сфере
// радиусом Distance с углами Yaw (горизонтальный) и Pitch (вертикальный).
//
// Поля:
//
//	Target   — точка, на которую смотрит камера (world space).
//	         Единицы: единицы модели. Начальное значение: (0.0, 0.05, 0.3).
//	Distance — расстояние от камеры до Target, диапазон [1.0, 50.0].
//	Yaw      — горизонтальный угол поворота вокруг Target (радианы).
//	         Значение 0 — камера по оси +Z.
//	Pitch    — вертикальный угол поворота (радианы).
//	         Ограничен (-π/2+0.1, π/2-0.1) для избежания переворота.
type Camera struct {
	Target   mgl32.Vec3
	Distance float32
	Yaw      float32
	Pitch    float32
}

// DefaultCamera — возвращает камеру с начальными параметрами.
// Камера расположена на расстоянии 5 единиц от центра апельсина,
// немного сверху и справа.
func DefaultCamera() Camera {
	return Camera{
		Target:   mgl32.Vec3{0.0, 0.05, 0.3},
		Distance: 5.0,
		Yaw:      0.35,
		Pitch:    0.25,
	}
}

// EyePosition — вычисляет позицию камеры в мировом пространстве.
// Использует сферические координаты относительно Target:
//
//	x = Distance * cos(Yaw) * cos(Pitch)
//	y = Distance * sin(Pitch)
//	z = Distance * sin(Yaw) * cos(Pitch)
//
// Возвращает: Target + (x, y, z).
func (c *Camera) EyePosition() mgl32.Vec3 {
	x := c.Distance * float32(math.Cos(float64(c.Yaw))) * float32(math.Cos(float64(c.Pitch)))
	y := c.Distance * float32(math.Sin(float64(c.Pitch)))
	z := c.Distance * float32(math.Sin(float64(c.Yaw))) * float32(math.Cos(float64(c.Pitch)))
	return c.Target.Add(mgl32.Vec3{x, y, z})
}

// ViewMatrix — строит матрицу вида (View) через mgl32.LookAt.
// Используется в цикле рендера для преобразования world space → view space.
// Uniform-переменная в шейдере: transform.view.
func (c *Camera) ViewMatrix() mgl32.Mat4 {
	eye := c.EyePosition()
	return mgl32.LookAt(
		eye.X(), eye.Y(), eye.Z(),
		c.Target.X(), c.Target.Y(), c.Target.Z(),
		0.0, 1.0, 0.0,
	)
}

// PanForward — смещает target камеры в направлении взгляда.
// amount > 0 — отодвигает target (камера "отлетает" назад).
func (c *Camera) PanForward(amount float32) {
	offset := mgl32.Vec3{
		float32(math.Cos(float64(c.Yaw)) * math.Cos(float64(c.Pitch))),
		float32(math.Sin(float64(c.Pitch))),
		float32(math.Sin(float64(c.Yaw)) * math.Cos(float64(c.Pitch))),
	}
	c.Target = c.Target.Add(offset.Normalize().Mul(amount))
}

// PanRight — смещает target камеры вправо относительно направления взгляда.
func (c *Camera) PanRight(amount float32) {
	right := mgl32.Vec3{
		float32(-math.Sin(float64(c.Yaw))),
		0,
		float32(math.Cos(float64(c.Yaw))),
	}
	c.Target = c.Target.Add(right.Normalize().Mul(amount))
}

// PanUp — смещает target камеры вертикально вверх (по оси Y мирового пространства).
func (c *Camera) PanUp(amount float32) {
	c.Target = c.Target.Add(mgl32.Vec3{0, amount, 0})
}

// Zoom — изменяет расстояние от камеры до target.
// Параметры:
//
//	delta   — приращение расстояния (положительное — отдаление).
//	minDist — минимальное расстояние (1.0).
//	maxDist — максимальное расстояние (50.0).
func (c *Camera) Zoom(delta, minDist, maxDist float32) {
	c.Distance += delta
	if c.Distance < minDist {
		c.Distance = minDist
	}
	if c.Distance > maxDist {
		c.Distance = maxDist
	}
}

// Rotate — вращает камеру вокруг target (орбитальное вращение).
// Параметры:
//
//	dYaw   — приращение горизонтального угла (клавиши Left/Right).
//	dPitch — приращение вертикального угла (клавиши Up/Down).
//
// Pitch ограничен, чтобы камера не переворачивалась.
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
