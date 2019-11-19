package util

import (
	"errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/nfnt/resize"
	"gorgonia.org/tensor"
)

const (
	height = 224
	width  = 224
)

var (
	mean = [3]float32{0.485, 0.456, 0.406}
	std  = [3]float32{0.229, 0.224, 0.225}
)

func ReadImage(filename string) (tensor.Tensor, error) {
	// Read input image
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	img = resize.Resize(width, height, img, resize.Bilinear)

	input := tensor.New(tensor.WithShape(1, 3, height, width), tensor.Of(tensor.Float32))
	err = imageToBCHW(img, input)
	if err != nil {
		return nil, err
	}
	err = normalize(input)
	if err != nil {
		return nil, err
	}

	return input, nil
}

func imageToBCHW(img image.Image, dst tensor.Tensor) error {
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r, g, b, a := img.At(x, y).RGBA()
			if a != 65535 {
				return errors.New("transparency not handled")
			}
			err := dst.SetAt(float32(float64(r)/0xffff), 0, 0, y, x)
			if err != nil {
				return err
			}
			err = dst.SetAt(float32(float64(g)/0xffff), 0, 1, y, x)
			if err != nil {
				return err
			}
			err = dst.SetAt(float32(float64(b)/0xffff), 0, 2, y, x)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func normalize(input tensor.Tensor) (err error) {
	for channel := 0; channel < 3; channel++ {
		m := mean[channel]
		s := std[channel]
		f := func(c float32) float32 { return (c - m) / s }

		cchannel, _ := input.Slice(nil, ss(channel), nil, nil)
		_, err := cchannel.Apply(f, tensor.WithReuse(cchannel))
		if err != nil {
			return err
		}
	}
	return nil
}

type ss int

func (s ss) Start() int { return int(s) }
func (s ss) End() int   { return int(s) + 1 }
func (s ss) Step() int  { return 0 }
