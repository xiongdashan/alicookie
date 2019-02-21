package models

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
)

type CutImage struct {
	name string
}

func NewCutImg() *CutImage {
	return &CutImage{}
}

func (n *CutImage) Cut(buf []byte) *bytes.Buffer {

	m, err := png.Decode(bytes.NewReader(buf))

	if err != nil {
		fmt.Println(err)
		return nil
	}
	rgbImg := m.(*image.NRGBA)
	subImg := rgbImg.SubImage(image.Rect(893, 345, 811, 385)).(*image.NRGBA)

	// f, _ := os.Create("./111.png")
	// defer f.Close()

	// png.Encode(f, subImg)

	// return nil

	imgBuf := new(bytes.Buffer)

	png.Encode(imgBuf, subImg)

	return imgBuf

}
