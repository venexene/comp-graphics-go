package scene

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/venexene/comp-graphics-go/objects"
)

// ObjectState holds the transform state for the main scene object.
type ObjectState struct {
	Position  mgl32.Vec3
	Scale     float32
	RotationX float32
	RotationY float32
	RotationZ float32
}

// DefaultObjectState returns a neutral object transform.
func DefaultObjectState() ObjectState {
	return ObjectState{
		Position:  mgl32.Vec3{0, 0, 0},
		Scale:     1.0,
		RotationX: 0.0,
		RotationY: 0.0,
		RotationZ: 0.0,
	}
}

// ModelMatrix computes the model matrix: translation * rotZ * rotY * rotX * scale.
func (o *ObjectState) ModelMatrix() mgl32.Mat4 {
	scale := mgl32.Scale3D(o.Scale, o.Scale, o.Scale)
	rotX := mgl32.HomogRotate3D(o.RotationX, mgl32.Vec3{1, 0, 0})
	rotY := mgl32.HomogRotate3D(o.RotationY, mgl32.Vec3{0, 1, 0})
	rotZ := mgl32.HomogRotate3D(o.RotationZ, mgl32.Vec3{0, 0, 1})
	trans := mgl32.Translate3D(o.Position.X(), o.Position.Y(), o.Position.Z())
	return trans.Mul4(rotZ.Mul4(rotY.Mul4(rotX.Mul4(scale))))
}

// SceneObject represents a secondary object in the scene with its own transform.
type SceneObject struct {
	Model     *objects.Model
	Position  mgl32.Vec3
	Scale     float32
	RotationX float32
	RotationY float32
	RotationZ float32
	Name      string
}

// ModelMatrix computes the model matrix for this scene object.
func (s *SceneObject) ModelMatrix() mgl32.Mat4 {
	scale := mgl32.Scale3D(s.Scale, s.Scale, s.Scale)
	rotX := mgl32.HomogRotate3D(s.RotationX, mgl32.Vec3{1, 0, 0})
	rotY := mgl32.HomogRotate3D(s.RotationY, mgl32.Vec3{0, 1, 0})
	rotZ := mgl32.HomogRotate3D(s.RotationZ, mgl32.Vec3{0, 0, 1})
	trans := mgl32.Translate3D(s.Position.X(), s.Position.Y(), s.Position.Z())
	return trans.Mul4(rotZ.Mul4(rotY.Mul4(rotX.Mul4(scale))))
}

// Selection manages which object is currently selected for manipulation.
type Selection struct {
	objects    []*SceneObject
	mainName   string
	index      int // 0 = main object, 1..n = objects[0..n-1]
}

// NewSelection creates a new selection state.
func NewSelection(mainName string) *Selection {
	return &Selection{mainName: mainName}
}

// RegisterObjects stores the extra scene objects for selection cycling.
func (s *Selection) RegisterObjects(objs ...*SceneObject) {
	s.objects = objs
}

// SetMainName sets the display name for the main (primary) object.
func (s *Selection) SetMainName(name string) {
	s.mainName = name
}

// CycleForward advances the selection to the next object.
func (s *Selection) CycleForward() {
	if len(s.objects) > 0 {
		s.index = (s.index + 1) % (1 + len(s.objects))
	}
}

// SelectedName returns the display name of the currently selected object.
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

// IsMain returns true if the main object is selected.
func (s *Selection) IsMain() bool {
	return s.index == 0
}

// SelectedSceneObject returns the selected extra scene object, or nil if main is selected.
func (s *Selection) SelectedSceneObject() *SceneObject {
	if s.index == 0 || s.index > len(s.objects) {
		return nil
	}
	return s.objects[s.index-1]
}
