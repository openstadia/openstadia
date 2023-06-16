package rtc

import (
	"fmt"
	"github.com/pion/mediadevices/pkg/io/video"
	"golang.org/x/image/colornames"
	"image"
)

func Mark(show *bool) video.TransformFunc {
	return func(r video.Reader) video.Reader {
		return video.ReaderFunc(func() (image.Image, func(), error) {
			for {
				img, _, err := r.Read()
				if err != nil {
					return nil, func() {}, err
				}

				switch v := img.(type) {
				case *image.RGBA:
					for yi := 0; yi < 16; yi++ {
						for xi := 0; xi < 16; xi++ {
							if *show {
								v.Set(xi, yi, colornames.Red)
							} else {
								v.Set(xi, yi, colornames.White)
							}
						}
					}
				default:
					fmt.Printf("unexpected type %T\n", v)
				}

				if *show {

				}

				return img, func() {}, nil
			}
		})
	}
}
