package textures

import (
	"os"
	"image"
    "image/draw"
	_ "image/png"

	"github.com/go-gl/gl/v4.6-core/gl"
)

// Загрузка текстуры 
func LoadTexture(file string) uint32 {
    // Открытие файла
    imgFile, err := os.Open(file)
    if err != nil {
        panic(err)
    }
    defer imgFile.Close()

    // Декодирование изображения
    img, _, err := image.Decode(imgFile)
    if err != nil {
        panic(err)
    }

    // Преобразования в RGBA
    rgba := image.NewRGBA(img.Bounds())
    draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

    // Создание текстуры
    var texture uint32
    gl.GenTextures(1, &texture)
    gl.BindTexture(gl.TEXTURE_2D, texture)

    // Загрузка изображения в текстуру
    gl.TexImage2D(
        gl.TEXTURE_2D, 0, gl.RGBA,
        int32(rgba.Rect.Size().X), int32(rgba.Rect.Size().Y),
        0, gl.RGBA, gl.UNSIGNED_BYTE,
        gl.Ptr(rgba.Pix),
    )

    // Настройка параметров текстуры
    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

    return texture
}