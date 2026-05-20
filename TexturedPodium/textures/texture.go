// Пакет textures отвечает за загрузку изображений с диска и создание
// OpenGL-текстур (GL_TEXTURE_2D). Поддерживаются форматы PNG и JPEG.
// Ключевые функции: LoadTexture — загружает файл, декодирует, создаёт
// текстуру с LINEAR-фильтрацией и REPEAT-обёрткой.
// Зависимости: Go standard library (image, image/png, image/jpeg),
// go-gl/gl — загрузка данных в видеопамять.
package textures

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/go-gl/gl/v4.6-core/gl"
)

// LoadTexture открывает файл изображения (PNG/JPEG), декодирует его в RGBA
// и загружает в видеопамять OpenGL как GL_TEXTURE_2D.
// Вход: filepath — путь к файлу относительно projectRoot или абсолютный.
// Возвращает: uint32 — идентификатор OpenGL-текстуры, error — ошибка
// открытия/декодирования/загрузки.
// Побочные эффекты: генерирует новый GL-объект текстуры, настраивает
// фильтрацию (LINEAR) и обёртку (REPEAT).
// Вызывается на этапе инициализации сцены (InitScene).
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
