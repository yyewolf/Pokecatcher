package main

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"strings"
)

var decodedImages map[string]image.Image

func decodeKnown() {
	decodedImages = make(map[string]image.Image)
	List := box.List()
	for i := range List {
		if strings.Contains(List[i], "img") {
			c, _ := box.FindString(List[i])
			Name := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(List[i], "img/", ""), "img\\", ""), ".png", "")
			img, err := loadImg(c)
			if err != nil {
				logDebug("[ERROR] ", err)
			}
			decodedImages[Name] = img
		}
	}
}

func loadImg(Image string) (image.Image, error) {
	i := strings.NewReader(Image)
	img, _, err := image.Decode(i)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func diff(a, b uint32) int64 {
	if a > b {
		return int64(a - b)
	}
	return int64(b - a)
}

func compareIMG(f1 image.Image, f2 image.Image) float64 {
	//f1 => resized
	//f2 => pokécord

	b := f1.Bounds()
	var sum int64
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r1, g1, b1, _ := f1.At(x, y).RGBA()
			r2, g2, b2, _ := f2.At(x, y).RGBA()
			sum += diff(r1, r2)
			sum += diff(g1, g2)
			sum += diff(b1, b2)
		}
	}

	nPixels := (b.Max.X - b.Min.X) * (b.Max.Y - b.Min.Y)

	return float64(sum*100) / (float64(nPixels) * 0xffff * 3)
}
