package main

import (
	"fmt"
	"image/color"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/google/uuid"
)

func main() {
	files, err := ioutil.ReadDir("image")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	for _, file := range files {
		fmt.Println(file.Name())
		if err := genLgtmImage("./image/" + file.Name()); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	}
}

func genLgtmImage(filePath string) error {
	loadImage, err := gg.LoadImage(filePath)
	if err != nil {
		return err
	}

	rct := loadImage.Bounds()
	// Width > Height
	if rct.Dx() > rct.Dy() {
		// Resize the cropped image to width = 400px preserving the aspect ratio.
		loadImage = imaging.Resize(loadImage, 400, 0, imaging.Lanczos)
	} else {
		// Resize the cropped image to Height = 400px preserving the aspect ratio.
		loadImage = imaging.Resize(loadImage, 0, 400, imaging.Lanczos)
	}

	resizedRct := loadImage.Bounds()
	dc := gg.NewContext(resizedRct.Dx(), resizedRct.Dy())

	dc.DrawImage(loadImage, 0, 0)

	fontPath := filepath.Join("fonts", "MPLUSRounded1c-Medium.ttf")
	if err := dc.LoadFontFace(fontPath, 80); err != nil {
		return err
	}
	dc.SetColor(color.White)
	s := "LGTM"
	textWidth, textHeight := dc.MeasureString(s)

	x := (float64(dc.Width()) - textWidth) / 2
	y := (float64(dc.Height()) / 2) + (textHeight / 2)
	dc.DrawString(s, x, y)

	var outputFilename = "./lgtm_image/" + uuid.NewString() + ".png"
	if err := dc.SavePNG(outputFilename); err != nil {
		return err
	}
	return nil
}
