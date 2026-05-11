package x264

import (
	"image"
)

// ImageToYUV converts an image.Image to YUV420p planar data suitable for
// EncodeYUV. The image should have even dimensions.
func ImageToYUV(img image.Image) []byte {
	if img == nil {
		return nil
	}
	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	if w == 0 || h == 0 {
		return nil
	}

	size := w * h * 3 / 2
	yuv := make([]byte, size)

	sw := w / 2
	yOff := 0
	uOff := w * h
	vOff := w*h + sw*(h/2)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, _ := img.At(x+bounds.Min.X, y+bounds.Min.Y).RGBA()
			r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)

			yVal := (19595*uint32(r8) + 38470*uint32(g8) + 7471*uint32(b8) + 32768) >> 16
			yuv[yOff+y*w+x] = byte(yVal)

			if x%2 == 0 && y%2 == 0 {
				var tr, tg, tb uint64
				var cnt int
				for dy := 0; dy < 2; dy++ {
					for dx := 0; dx < 2; dx++ {
						px, py := x+dx, y+dy
						if px < w && py < h {
							rr, gg, bb, _ := img.At(px+bounds.Min.X, py+bounds.Min.Y).RGBA()
							tr += uint64(rr)
							tg += uint64(gg)
							tb += uint64(bb)
							cnt++
						}
					}
				}
				if cnt == 0 {
					cnt = 1
				}
				c := uint32(cnt)
				ar := uint8((tr / uint64(c)) >> 8)
				ag := uint8((tg / uint64(c)) >> 8)
				ab := uint8((tb / uint64(c)) >> 8)

				uVal := (8421376 - 11022*int32(ar) - 21709*int32(ag) + 32768*int32(ab)) >> 16
				vVal := (8421376 + 32768*int32(ar) - 27439*int32(ag) - 5329*int32(ab)) >> 16

				if uVal < 0 {
					uVal = 0
				} else if uVal > 255 {
					uVal = 255
				}
				if vVal < 0 {
					vVal = 0
				} else if vVal > 255 {
					vVal = 255
				}

				uvIdx := (y/2)*sw + (x / 2)
				yuv[uOff+uvIdx] = byte(uVal)
				yuv[vOff+uvIdx] = byte(vVal)
			}
		}
	}

	return yuv
}
