package x264

import (
	"image"
	"image/color"
	"testing"
)

func TestImageToYUV_ZeroSized(t *testing.T) {
	tests := []struct {
		name string
		img  image.Image
	}{
		{"nil", nil},
		{"zero width", image.NewRGBA(image.Rect(0, 0, 0, 100))},
		{"zero height", image.NewRGBA(image.Rect(0, 0, 100, 0))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			yuv := ImageToYUV(tt.img)
			if yuv != nil {
				t.Errorf("expected nil, got len=%d", len(yuv))
			}
		})
	}
}

func TestImageToYUV_Size(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 4, 4))
	yuv := ImageToYUV(img)
	expected := 4 * 4 * 3 / 2
	if len(yuv) != expected {
		t.Errorf("expected size %d, got %d", expected, len(yuv))
	}
}

func TestImageToYUV_SolidColors(t *testing.T) {
	tests := []struct {
		name string
		r, g, b uint8
		ey, eu, ev byte // expected Y, Cb, Cr
	}{
		{"red",   255, 0,   0,   76, 85, 255},
		{"green", 0,   255, 0,   150, 44, 21},
		{"blue",  0,   0,   255, 29, 255, 107},
		{"gray",  128, 128, 128, 128, 128, 128},
		{"black", 0,   0,   0,   0, 128, 128},
		{"white", 255, 255, 255, 255, 128, 128},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img := image.NewRGBA(image.Rect(0, 0, 2, 2))
			for y := 0; y < 2; y++ {
				for x := 0; x < 2; x++ {
					img.Set(x, y, color.RGBA{tt.r, tt.g, tt.b, 255})
				}
			}

			yuv := ImageToYUV(img)

			// Check first Y pixel
			if yuv[0] != tt.ey {
				t.Errorf("Y: expected %d, got %d", tt.ey, yuv[0])
			}

			// Check U (first chroma sample)
			uOff := 2 * 2
			if yuv[uOff] != tt.eu {
				t.Errorf("U: expected %d, got %d", tt.eu, yuv[uOff])
			}

			// Check V
			vOff := 2*2 + 1*1
			if yuv[vOff] != tt.ev {
				t.Errorf("V: expected %d, got %d", tt.ev, yuv[vOff])
			}
		})
	}
}

func TestImageToYUV_SubImage(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 4, 4))
	// Fill with blue
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			src.Set(x, y, color.RGBA{0, 0, 255, 255})
		}
	}
	// Extract 2x2 sub-image from the center
	sub := src.SubImage(image.Rect(1, 1, 3, 3))
	yuv := ImageToYUV(sub)

	if len(yuv) != 2*2*3/2 {
		t.Errorf("expected size %d, got %d", 2*2*3/2, len(yuv))
	}

	// All pixels should be blue
	if yuv[0] != 29 {
		t.Errorf("expected Y=29 for blue, got %d", yuv[0])
	}
	uOff := 2 * 2
	if yuv[uOff] != 255 {
		t.Errorf("expected U=255 for blue, got %d", yuv[uOff])
	}
}

func TestImageToYUV_NRGBA(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 2, 2))
	for y := 0; y < 2; y++ {
		for x := 0; x < 2; x++ {
			img.Set(x, y, color.NRGBA{255, 0, 0, 255})
		}
	}
	yuv := ImageToYUV(img)
	if len(yuv) != 2*2*3/2 {
		t.Errorf("expected size %d, got %d", 2*2*3/2, len(yuv))
	}
	if yuv[0] != 76 {
		t.Errorf("expected Y=76 for red, got %d", yuv[0])
	}
}

func TestImageToYUV_Gray(t *testing.T) {
	img := image.NewGray(image.Rect(0, 0, 2, 2))
	for y := 0; y < 2; y++ {
		for x := 0; x < 2; x++ {
			img.Set(x, y, color.Gray{128})
		}
	}
	yuv := ImageToYUV(img)
	if len(yuv) != 2*2*3/2 {
		t.Errorf("expected size %d, got %d", 2*2*3/2, len(yuv))
	}
	if yuv[0] != 128 {
		t.Errorf("expected Y=128 for gray, got %d", yuv[0])
	}
	uOff := 2 * 2
	if yuv[uOff] != 128 {
		t.Errorf("expected U=128 for gray, got %d", yuv[uOff])
	}
}

func TestImageToYUV_NonZeroOrigin(t *testing.T) {
	img := image.NewRGBA(image.Rect(10, 10, 12, 12))
	for y := 10; y < 12; y++ {
		for x := 10; x < 12; x++ {
			img.Set(x, y, color.RGBA{0, 255, 0, 255})
		}
	}
	yuv := ImageToYUV(img)
	if len(yuv) != 2*2*3/2 {
		t.Errorf("expected size %d, got %d", 2*2*3/2, len(yuv))
	}
	if yuv[0] != 150 {
		t.Errorf("expected Y=150 for green, got %d", yuv[0])
	}
}

func TestImageToYUV_ChromaAveraging(t *testing.T) {
	// Create a 2x2 image with 4 different colors to verify chroma averaging
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.RGBA{255, 0, 0, 255})   // red
	img.Set(1, 0, color.RGBA{0, 255, 0, 255})   // green
	img.Set(0, 1, color.RGBA{0, 0, 255, 255})   // blue
	img.Set(1, 1, color.RGBA{255, 255, 255, 255}) // white

	yuv := ImageToYUV(img)

	// Average R = (255 + 0 + 0 + 255) / 4 = 127.5 -> 127
	// Average G = (0 + 255 + 0 + 255) / 4 = 127.5 -> 127
	// Average B = (0 + 0 + 255 + 255) / 4 = 127.5 -> 127
	// U for gray = 128
	uOff := 2 * 2
	if yuv[uOff] != 128 {
		t.Errorf("expected U=128 for averaged gray, got %d", yuv[uOff])
	}
	if yuv[uOff+1] != 128 {
		t.Errorf("expected V=128 for averaged gray, got %d", yuv[uOff+1])
	}
}

func TestImageToYUV_LargerImage(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 64, 64))
	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			img.Set(x, y, color.RGBA{uint8(x * 4), uint8(y * 4), 128, 255})
		}
	}
	yuv := ImageToYUV(img)
	expected := 64 * 64 * 3 / 2
	if len(yuv) != expected {
		t.Errorf("expected size %d, got %d", expected, len(yuv))
	}
}

func TestCropEvenI420(t *testing.T) {
	tests := []struct {
		name   string
		w, h   int
	}{
		{"even 4x4", 4, 4},
		{"odd width 5x4", 5, 4},
		{"odd height 4x5", 4, 5},
		{"both odd 5x5", 5, 5},
		{"odd width 3x4", 3, 4},
		{"odd height 4x3", 4, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tw, th := tt.w & ^1, tt.h & ^1

			ySrc := tt.w * tt.h
			uvSrc := (tt.w / 2) * (tt.h / 2)
			src := make([]byte, ySrc+uvSrc*2)

			for i := range ySrc {
				src[i] = 0x11
			}
			for i := ySrc; i < ySrc+uvSrc; i++ {
				src[i] = 0x33
			}
			for i := ySrc + uvSrc; i < len(src); i++ {
				src[i] = 0x55
			}

			yDst := tw * th
			uvDst := (tw / 2) * (th / 2)
			dst := make([]byte, yDst+uvDst*2)

			CropEvenI420(src, tt.w, tt.h, dst)

			for i := range yDst {
				if dst[i] != 0x11 {
					t.Errorf("Y[%d] = 0x%02x, want 0x11", i, dst[i])
				}
			}
			for i := yDst; i < yDst+uvDst; i++ {
				if dst[i] != 0x33 {
					t.Errorf("U[%d] = 0x%02x, want 0x33", i-yDst, dst[i])
				}
			}
			for i := yDst + uvDst; i < len(dst); i++ {
				if dst[i] != 0x55 {
					t.Errorf("V[%d] = 0x%02x, want 0x55", i-yDst-uvDst, dst[i])
				}
			}
		})
	}
}

func TestCropEvenI420_ZeroOutput(t *testing.T) {
	// 1x1 crops to 0x0: no output pixels, should not panic.
	CropEvenI420([]byte{42}, 1, 1, nil)
}

func TestImageToYUV_WithEncoder(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 64, 64))
	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			img.Set(x, y, color.RGBA{uint8(x * 4), uint8(y * 4), 128, 255})
		}
	}

	enc, err := New(64, 64, 30)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer enc.Close()

	yuv := ImageToYUV(img)
	data, err := enc.EncodeYUV(yuv)
	if err != nil {
		t.Fatalf("EncodeYUV failed: %v", err)
	}
	if len(data) == 0 {
		t.Error("encoded data should not be empty")
	}
}
