// scene/object.go — Состояние объектов сцены и управление выбором.
//
// Назначение: определяет трансформации объектов (позиция, масштаб, вращение)
// и механизм переключения выбранного объекта (Tab).
//
// Ключевые типы:
//   ObjectState  — трансформация главного объекта сцены.
//   SceneObject  — дополнительный объект сцены (модель + трансформация + имя).
//   Selection    — состояние выбора: главный или один из дополнительных объектов.
//
// Ключевые функции:
//   DefaultObjectState()  — нейтральная трансформация (единичная).
//   ModelMatrix()         — вычисление модельной матрицы (T*Rz*Ry*Rx*S).
//   NewSelection()        — создание нового состояния выбора.
//   CycleForward()        — переход к следующему объекту.
//   SelectedName()        — имя текущего выбранного объекта.
//
// Зависимости: используется в scene.DrawScene(), input.ProcessInput(), utils.
package scene

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/venexene/comp-graphics-go/objects"
)

// ObjectState — состояние трансформации главного объекта сцены.
// Определяет, как объект расположен, повёрнут и масштабирован
// в мировом пространстве.
type ObjectState struct {
	// Position — положение объекта в мировом пространстве (world space).
	Position mgl32.Vec3
	// Scale — равномерный масштаб по всем трём осям.
	// Диапазон: [0.1, 3.0].
	Scale float32
	// RotationX — угол поворота вокруг оси X в радианах.
	RotationX float32
	// RotationY — угол поворота вокруг оси Y в радианах.
	RotationY float32
	// RotationZ — угол поворота вокруг оси Z в радианах.
	RotationZ float32
}

// DefaultObjectState — возвращает нейтральную трансформацию:
// позиция (0,0,0), масштаб 1, без поворота.
func DefaultObjectState() ObjectState {
	return ObjectState{
		Position:  mgl32.Vec3{0, 0, 0},
		Scale:     1.0,
		RotationX: 0.0,
		RotationY: 0.0,
		RotationZ: 0.0,
	}
}

// ModelMatrix — вычисляет модельную матрицу (Model Space → World Space).
// Порядок преобразований: масштаб → вращение X → вращение Y → вращение Z → перенос.
// Итоговая матрица: T * (Rz * Ry * Rx * S).
// Возвращает: матрицу 4×4.
func (o *ObjectState) ModelMatrix() mgl32.Mat4 {
	scale := mgl32.Scale3D(o.Scale, o.Scale, o.Scale)
	rotX := mgl32.HomogRotate3D(o.RotationX, mgl32.Vec3{1, 0, 0})
	rotY := mgl32.HomogRotate3D(o.RotationY, mgl32.Vec3{0, 1, 0})
	rotZ := mgl32.HomogRotate3D(o.RotationZ, mgl32.Vec3{0, 0, 1})
	trans := mgl32.Translate3D(o.Position.X(), o.Position.Y(), o.Position.Z())
	return trans.Mul4(rotZ.Mul4(rotY.Mul4(rotX.Mul4(scale))))
}

// SceneObject — дополнительный объект сцены со своей трансформацией.
// Используется для heart и default моделей, расположенных по бокам от снеговика.
type SceneObject struct {
	Model     *objects.Model // указатель на загруженную 3D-модель
	Position  mgl32.Vec3     // позиция в world space
	Scale     float32        // равномерный масштаб
	RotationX float32        // поворот вокруг X (радианы)
	RotationY float32        // поворот вокруг Y (радианы)
	RotationZ float32        // поворот вокруг Z (радианы)
	Name      string         // отображаемое имя объекта
}

// ModelMatrix — вычисляет модельную матрицу для SceneObject.
// Аналогично ObjectState.ModelMatrix(): T * Rz * Ry * Rx * S.
func (s *SceneObject) ModelMatrix() mgl32.Mat4 {
	scale := mgl32.Scale3D(s.Scale, s.Scale, s.Scale)
	rotX := mgl32.HomogRotate3D(s.RotationX, mgl32.Vec3{1, 0, 0})
	rotY := mgl32.HomogRotate3D(s.RotationY, mgl32.Vec3{0, 1, 0})
	rotZ := mgl32.HomogRotate3D(s.RotationZ, mgl32.Vec3{0, 0, 1})
	trans := mgl32.Translate3D(s.Position.X(), s.Position.Y(), s.Position.Z())
	return trans.Mul4(rotZ.Mul4(rotY.Mul4(rotX.Mul4(scale))))
}

// Selection — управление выбором между главным и дополнительными объектами.
// index=0 — главный объект (снеговик), index=1..n — объекты из слайса.
// Переключение по Tab в input.ProcessInput().
type Selection struct {
	objects  []*SceneObject // дополнительные объекты
	mainName string         // путь к файлу главного объекта
	index    int            // 0 = главный объект, 1..n = objects[0..n-1]
}

// NewSelection — создаёт новое состояние выбора.
// Принимает: mainName — путь к файлу главной модели (для отображения имени).
func NewSelection(mainName string) *Selection {
	return &Selection{mainName: mainName}
}

// RegisterObjects — сохраняет дополнительные объекты для циклического выбора.
// Вызывается: из main() после загрузки всех моделей.
func (s *Selection) RegisterObjects(objs ...*SceneObject) {
	s.objects = objs
}

// SetMainName — задаёт отображаемое имя главного объекта.
func (s *Selection) SetMainName(name string) {
	s.mainName = name
}

// CycleForward — переключает выбор на следующий объект по циклу.
// Главный → Объект1 → Объект2 → ... → Главный.
// Побочные эффекты: изменяет index.
func (s *Selection) CycleForward() {
	if len(s.objects) > 0 {
		s.index = (s.index + 1) % (1 + len(s.objects))
	}
}

// SelectedName — возвращает отображаемое имя текущего выбранного объекта.
// Для главного объекта извлекает имя из пути (без расширения).
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

// IsMain — возвращает true, если выбран главный объект (index == 0).
func (s *Selection) IsMain() bool {
	return s.index == 0
}

// SelectedSceneObject — возвращает выбранный дополнительный объект.
// Если выбран главный объект или индекс вне диапазона — возвращает nil.
func (s *Selection) SelectedSceneObject() *SceneObject {
	if s.index == 0 || s.index > len(s.objects) {
		return nil
	}
	return s.objects[s.index-1]
}
