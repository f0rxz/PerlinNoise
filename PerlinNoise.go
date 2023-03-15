package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
)

type NoiseImage struct {
	Image         *image.Gray
	Width, Height int
}

func NewNoiseImage(width, height int) *NoiseImage {
	var noiseImage = &NoiseImage{
		Image:  image.NewGray(image.Rect(0, 0, width, height)),
		Width:  width,
		Height: height,
	}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			noiseImage.Image.SetGray(x, y, color.Gray{Y: byte(rand.Intn(256))})
		}
	}
	return noiseImage
}

func (noiseImage *NoiseImage) Resize(width, height int) {
	var (
		newImage                  = image.NewGray(image.Rect(0, 0, width, height))
		widthSource, heightSource = noiseImage.Width, noiseImage.Height
		xs, ys                    = float32(widthSource) / float32(width), float32(heightSource) / float32(height)
	)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var sx = (float32(x)+0.5)*xs - 0.5
			if sx < 0 {
				sx = 0
			}

			var sy = (float32(y)+0.5)*ys - 0.5
			if sy < 0 {
				sy = 0
			}

			var (
				x0     = int(sx)
				y0     = int(sy)
				fracx  = sx - float32(x0)
				fracy  = sy - float32(y0)
				ifracx = 1.0 - fracx
				ifracy = 1.0 - fracy
			)

			var x1 = x0 + 1
			if x1 >= widthSource {
				x1 = x0
			}

			var y1 = y0 + 1
			if y1 >= heightSource {
				y1 = y0
			}

			var (
				c1 = float32(noiseImage.Image.GrayAt(x0, y0).Y)
				c2 = float32(noiseImage.Image.GrayAt(x1, y0).Y)
				c3 = float32(noiseImage.Image.GrayAt(x0, y1).Y)
				c4 = float32(noiseImage.Image.GrayAt(x1, y1).Y)

				l0 = ifracx*c1 + fracx*c2
				l1 = ifracx*c3 + fracx*c4
				rf = ifracy*l0 + fracy*l1
			)

			newImage.SetGray(x, y, color.Gray{Y: byte(rf)})
		}
	}

	noiseImage.Image = newImage
	noiseImage.Width = width
	noiseImage.Height = height
}

func SaveAsPng(outputName string, img *image.Gray) error {
	if file, err := os.Create(outputName); err == nil {
		defer file.Close()
		return png.Encode(file, img)
	} else {
		return err
	}
}

func NoiseImagesMiddle(noiseImages []*NoiseImage) *image.Gray {
	var maxWidth, maxHeight = 0, 0
	for _, noiseImage := range noiseImages {
		if noiseImage.Width > maxWidth {
			maxWidth = noiseImage.Width
		}
		if noiseImage.Height > maxHeight {
			maxHeight = noiseImage.Height
		}
	}

	var pixels = make([]int, maxWidth*maxHeight)
	for _, noiseImage := range noiseImages {
		for y := 0; y < noiseImage.Height; y++ {
			for x := 0; x < noiseImage.Width; x++ {
				pixels[x+y*maxWidth] += int(noiseImage.Image.GrayAt(x, y).Y)
			}
		}
	}

	var (
		newImage         = image.NewGray(image.Rect(0, 0, maxWidth, maxHeight))
		noiseImagesCount = len(noiseImages)
	)
	for y := 0; y < maxHeight; y++ {
		for x := 0; x < maxWidth; x++ {
			newImage.SetGray(x, y, color.Gray{Y: byte(math.Round(float64(pixels[x+y*maxWidth]) / float64(noiseImagesCount)))})
		}
	}

	return newImage
}

func main() {
	const (
		octaves = 12
		outSize = 1 << octaves
	)
	var (
		size        = outSize
		noiseImages = make([]*NoiseImage, octaves)
	)
	noiseImages[0] = NewNoiseImage(size, size)
	for i := 1; i < octaves; i++ {
		size /= 2
		var noiseImage = NewNoiseImage(size, size)
		noiseImage.Resize(outSize, outSize)
		noiseImages[i] = noiseImage
	}
	SaveAsPng("out.png", NoiseImagesMiddle(noiseImages))
}
