// Файл: scene/object.go
// Назначение: состояния объектов сцены, их матрицы трансформации и механизм выбора (selection).
//
// Ключевые типы:
//   ObjectState — состояние основного объекта: позиция, масштаб, вращение.
//   SceneObject — объект сцены со своей моделью, позицией, масштабом, вращением.
//   Selection — механизм выбора объекта (для управления через клавиатуру).
//
// Ключевые функции:
//   DefaultObjectState — возвращает начальное состояние для основного объекта.
//   ObjectState.ModelMatrix — строит матрицу Model (трансформация объекта).
//   SceneObject.ModelMatrix — строит матрицу Model для объекта сцены.
//   NewSelection — создаёт новый Selection.
//   Selection.CycleForward — переключает выбранный объект по циклу.
//
// Зависимости:
//   Внутренние: objects (objects/model.go).
//   Внешние: github.com/go-gl/mathgl/mgl32.

package scene

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/venexene/comp-graphics-go/objects"
)

// ObjectState — состояние основного объекта сцены (например, апельсина).
// Используется для построения матрицы Model.
// Поля:
//
//	Position  — позиция объекта в world space (x, y, z).
//	Scale     — равномерный масштаб, диапазон [0.1, 3.0].
//	RotationX — угол вращения вокруг оси X (радианы).
//	RotationY — угол вращения вокруг оси Y (радианы).
//	RotationZ — угол вращения вокруг оси Z (радианы).
type ObjectState struct {
	Position  mgl32.Vec3
	Scale     float32
	RotationX float32
	RotationY float32
	RotationZ float32
}

// DefaultObjectState — возвращает состояние объекта по умолчанию.
// Апельсин расположен слева от центра сцены, сердце — рядом.
func DefaultObjectState() ObjectState {
	return ObjectState{
		Position:  mgl32.Vec3{-3.36, -1.51, 0.91},
		Scale:     0.5,
		RotationX: 0.0,
		RotationY: 0.0,
		RotationZ: 0.0,
	}
}

// ModelMatrix — строит матрицу Model для объекта.
// Порядок трансформаций (справа налево):
//
//	Model = T * Rz * Ry * Rx * S
//
// где T — перенос, Rz/Ry/Rx — вращения, S — масштабирование.
// Этот порядок означает: сначала масштаб, потом вращения, потом перенос.
// Матрица используется в шейдере как transform.model.
func (o *ObjectState) ModelMatrix() mgl32.Mat4 {
	scale := mgl32.Scale3D(o.Scale, o.Scale, o.Scale)
	rotX := mgl32.HomogRotate3D(o.RotationX, mgl32.Vec3{1, 0, 0})
	rotY := mgl32.HomogRotate3D(o.RotationY, mgl32.Vec3{0, 1, 0})
	rotZ := mgl32.HomogRotate3D(o.RotationZ, mgl32.Vec3{0, 0, 1})
	trans := mgl32.Translate3D(o.Position.X(), o.Position.Y(), o.Position.Z())
	return trans.Mul4(rotZ.Mul4(rotY.Mul4(rotX.Mul4(scale))))
}

// SceneObject — объект сцены с собственной моделью и трансформацией.
// Используется для дополнительных объектов (например, сердца).
// Поля:
//
//	Model     — указатель на загруженную 3D-модель (*objects.Model).
//	Position  — позиция объекта в world space.
//	Scale     — равномерный масштаб.
//	RotationX — вращение вокруг оси X.
//	RotationY — вращение вокруг оси Y.
//	RotationZ — вращение вокруг оси Z.
//	Name      — имя объекта (для отображения в интерфейсе).
type SceneObject struct {
	Model     *objects.Model
	Position  mgl32.Vec3
	Scale     float32
	RotationX float32
	RotationY float32
	RotationZ float32
	Name      string
}

// ModelMatrix — строит матрицу Model для SceneObject.
// Аналогична ObjectState.ModelMatrix, порядок: T * Rz * Ry * Rx * S.
func (s *SceneObject) ModelMatrix() mgl32.Mat4 {
	scale := mgl32.Scale3D(s.Scale, s.Scale, s.Scale)
	rotX := mgl32.HomogRotate3D(s.RotationX, mgl32.Vec3{1, 0, 0})
	rotY := mgl32.HomogRotate3D(s.RotationY, mgl32.Vec3{0, 1, 0})
	rotZ := mgl32.HomogRotate3D(s.RotationZ, mgl32.Vec3{0, 0, 1})
	trans := mgl32.Translate3D(s.Position.X(), s.Position.Y(), s.Position.Z())
	return trans.Mul4(rotZ.Mul4(rotY.Mul4(rotX.Mul4(scale))))
}

// Selection — механизм выбора объекта для управления с клавиатуры.
// Позволяет циклически переключаться между основным объектом и
// дополнительными объектами сцены.
// index == 0 — выбран основной объект (main).
// index >= 1 — выбран objects[index-1].
type Selection struct {
	objects  []*SceneObject
	mainName string
	index    int
}

// NewSelection — создаёт Selection для объекта с именем mainName.
func NewSelection(mainName string) *Selection {
	return &Selection{mainName: mainName}
}

// RegisterObjects — регистрирует массив дополнительных объектов сцены.
func (s *Selection) RegisterObjects(objs ...*SceneObject) {
	s.objects = objs
}

// SetMainName — устанавливает имя основного объекта.
func (s *Selection) SetMainName(name string) {
	s.mainName = name
}

// CycleForward — переключает выбор на следующий объект (по циклу).
// Последовательность: main → obj[0] → obj[1] → ... → main.
func (s *Selection) CycleForward() {
	if len(s.objects) > 0 {
		s.index = (s.index + 1) % (1 + len(s.objects))
	}
}

// SelectedName — возвращает имя выбранного объекта для отображения.
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

// IsMain — возвращает true, если выбран основной объект (index == 0).
func (s *Selection) IsMain() bool {
	return s.index == 0
}

// SelectedSceneObject — возвращает указатель на выбранный SceneObject,
// или nil если выбран основной объект.
func (s *Selection) SelectedSceneObject() *SceneObject {
	if s.index == 0 || s.index > len(s.objects) {
		return nil
	}
	return s.objects[s.index-1]
}
