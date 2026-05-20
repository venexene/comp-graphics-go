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
