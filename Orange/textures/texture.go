// Файл: textures/texture.go
// Назначение: загрузка текстур из файлов изображений (JPEG, PNG) в OpenGL.
//
// Ключевые функции:
//   LoadTexture — загружает изображение из файла, создаёт OpenGL-текстуру (TEXTURE_2D).
//
// Зависимости:
//   Внутренние: — (используется scene/scene.go).
//   Внешние: github.com/go-gl/gl/v4.6-core/gl,
//            image, image/draw, image/jpeg, image/png (стандартная библиотека Go).

package textures

import (
	"fmt"
	"github.com/go-gl/gl/v4.6-core/gl"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

// LoadTexture — загружает изображение из файла и создаёт OpenGL-текстуру.
// Параметры:
//
//	filepath — путь к файлу изображения (.jpg или .png).
//
// Возвращает:
//
//	uint32 — идентификатор OpenGL-текстуры (GL_TEXTURE_2D).
//
// Алгоритм:
// 1. Открывает файл, декодирует изображение (image.Decode поддерживает JPEG и PNG).
// 2. Преобразует в RGBA (единый формат для OpenGL).
// 3. Генерирует текстуру, загружает пиксели через gl.TexImage2D.
// 4. Устанавливает параметры фильтрации (GL_LINEAR) и повторения (GL_REPEAT).
// 5. Отвязывает текстуру.
//
// Примечание: карта нормалей загружается как обычная RGB-текстура (не sRGB),
// т.к. данные нормали — это векторы, а не цвета.
func LoadTexture(filepath string) (uint32, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return 0, fmt.Errorf("cannot open texture %s: %w", filepath, err)
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return 0, fmt.Errorf("cannot decode texture %s: %w", filepath, err)
	}
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, img.Bounds().Min, draw.Src)
	var tex uint32
	gl.GenTextures(1, &tex)
	gl.BindTexture(gl.TEXTURE_2D, tex)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Bounds().Dx()),
		int32(rgba.Bounds().Dy()),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix),
	)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	return tex, nil
}
