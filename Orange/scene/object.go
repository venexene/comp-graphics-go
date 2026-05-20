package scene

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/venexene/comp-graphics-go/objects"
)

type ObjectState struct {
	Position  mgl32.Vec3
	Scale     float32
	RotationX float32
	RotationY float32
	RotationZ float32
}

func DefaultObjectState() ObjectState {
	return ObjectState{
		Position:  mgl32.Vec3{-3.36, -1.51, 0.91},
		Scale:     0.5,
		RotationX: 0.0,
		RotationY: 0.0,
		RotationZ: 0.0,
	}
}
func (o *ObjectState) ModelMatrix() mgl32.Mat4 {
	scale := mgl32.Scale3D(o.Scale, o.Scale, o.Scale)
	rotX := mgl32.HomogRotate3D(o.RotationX, mgl32.Vec3{1, 0, 0})
	rotY := mgl32.HomogRotate3D(o.RotationY, mgl32.Vec3{0, 1, 0})
	rotZ := mgl32.HomogRotate3D(o.RotationZ, mgl32.Vec3{0, 0, 1})
	trans := mgl32.Translate3D(o.Position.X(), o.Position.Y(), o.Position.Z())
	return trans.Mul4(rotZ.Mul4(rotY.Mul4(rotX.Mul4(scale))))
}

type SceneObject struct {
	Model     *objects.Model
	Position  mgl32.Vec3
	Scale     float32
	RotationX float32
	RotationY float32
	RotationZ float32
	Name      string
}

func (s *SceneObject) ModelMatrix() mgl32.Mat4 {
	scale := mgl32.Scale3D(s.Scale, s.Scale, s.Scale)
	rotX := mgl32.HomogRotate3D(s.RotationX, mgl32.Vec3{1, 0, 0})
	rotY := mgl32.HomogRotate3D(s.RotationY, mgl32.Vec3{0, 1, 0})
	rotZ := mgl32.HomogRotate3D(s.RotationZ, mgl32.Vec3{0, 0, 1})
	trans := mgl32.Translate3D(s.Position.X(), s.Position.Y(), s.Position.Z())
	return trans.Mul4(rotZ.Mul4(rotY.Mul4(rotX.Mul4(scale))))
}

type Selection struct {
	objects  []*SceneObject
	mainName string
	index    int
}

func NewSelection(mainName string) *Selection {
	return &Selection{mainName: mainName}
}
func (s *Selection) RegisterObjects(objs ...*SceneObject) {
	s.objects = objs
}
func (s *Selection) SetMainName(name string) {
	s.mainName = name
}
func (s *Selection) CycleForward() {
	if len(s.objects) > 0 {
		s.index = (s.index + 1) % (1 + len(s.objects))
	}
}
func (s *Selection) SelectedName() string {
	if s.index == 0 {
		base := filepath.Base(s.mainName)
		return strings.TrimSuffix(base, filepath.Ext(base))
	}
	idx := s.index - 1
	if idx >= 0 && idx < len(s.objects) && s.objects[idx] != nil {
		if s.objects[idx].Name != "" {
			base := filepath.Base(s.objects[idx].Name)
			return strings.TrimSuffix(base, filepath.Ext(base))
		}
		return fmt.Sprintf("Object%d", idx)
	}
	return "None"
}
func (s *Selection) IsMain() bool {
	return s.index == 0
}
func (s *Selection) SelectedSceneObject() *SceneObject {
	if s.index == 0 || s.index > len(s.objects) {
		return nil
	}
	return s.objects[s.index-1]
}
