package scene

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	Target   mgl32.Vec3
	Distance float32
	Yaw      float32
	Pitch    float32
}

func DefaultCamera() Camera {
	return Camera{
		Target:   mgl32.Vec3{0.0, 0.05, 0.3},
		Distance: 5.0,
		Yaw:      0.35,
		Pitch:    0.25,
	}
}
func (c *Camera) EyePosition() mgl32.Vec3 {
	x := c.Distance * float32(math.Cos(float64(c.Yaw))) * float32(math.Cos(float64(c.Pitch)))
	y := c.Distance * float32(math.Sin(float64(c.Pitch)))
	z := c.Distance * float32(math.Sin(float64(c.Yaw))) * float32(math.Cos(float64(c.Pitch)))
	return c.Target.Add(mgl32.Vec3{x, y, z})
}
func (c *Camera) ViewMatrix() mgl32.Mat4 {
	eye := c.EyePosition()
	return mgl32.LookAt(
		eye.X(), eye.Y(), eye.Z(),
		c.Target.X(), c.Target.Y(), c.Target.Z(),
		0.0, 1.0, 0.0,
	)
}
func (c *Camera) PanForward(amount float32) {
	offset := mgl32.Vec3{
		float32(math.Cos(float64(c.Yaw)) * math.Cos(float64(c.Pitch))),
		float32(math.Sin(float64(c.Pitch))),
		float32(math.Sin(float64(c.Yaw)) * math.Cos(float64(c.Pitch))),
	}
	c.Target = c.Target.Add(offset.Normalize().Mul(amount))
}
func (c *Camera) PanRight(amount float32) {
	right := mgl32.Vec3{
		float32(-math.Sin(float64(c.Yaw))),
		0,
		float32(math.Cos(float64(c.Yaw))),
	}
	c.Target = c.Target.Add(right.Normalize().Mul(amount))
}
func (c *Camera) PanUp(amount float32) {
	c.Target = c.Target.Add(mgl32.Vec3{0, amount, 0})
}
func (c *Camera) Zoom(delta, minDist, maxDist float32) {
	c.Distance += delta
	if c.Distance < minDist {
		c.Distance = minDist
	}
	if c.Distance > maxDist {
		c.Distance = maxDist
	}
}
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
