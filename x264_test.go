package x264

import (
	"testing"
)

func generateTestYUV(width, height int) []byte {
	size := width * height * 3 / 2
	yuv := make([]byte, size)

	for i := 0; i < width*height; i++ {
		yuv[i] = 128
	}
	for i := width * height; i < width*height*5/4; i++ {
		yuv[i] = 128
	}
	for i := width*height*5/4 - width/2*height/2; i < size; i++ {
		yuv[i] = 128
	}

	return yuv
}

func TestNewInvalidParams(t *testing.T) {
	tests := []struct {
		w, h, fps int
	}{
		{0, 480, 30},
		{640, 0, 30},
		{640, 480, 0},
		{-1, 480, 30},
	}

	for _, tt := range tests {
		enc, err := New(tt.w, tt.h, tt.fps)
		if enc != nil || err == nil {
			t.Errorf("New(%d, %d, %d) expected error", tt.w, tt.h, tt.fps)
		}
	}
}

func TestNewAndClose(t *testing.T) {
	enc, err := New(640, 480, 30)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	if enc == nil {
		t.Fatal("New() returned nil encoder")
	}

	enc.Close()

	_, err = enc.EncodeYUV([]byte{0})
	if err == nil {
		t.Error("EncodeYUV after Close should fail")
	}
}

func TestEncodeYUVInvalidSize(t *testing.T) {
	enc, err := New(640, 480, 30)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer enc.Close()

	_, err = enc.EncodeYUV([]byte{1, 2, 3})
	if err == nil {
		t.Error("EncodeYUV with small buffer should fail")
	}
}

func TestEncodeYUVValidFrame(t *testing.T) {
	enc, err := New(640, 480, 30)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer enc.Close()

	yuv := generateTestYUV(640, 480)

	data, err := enc.EncodeYUV(yuv)
	if err != nil {
		t.Fatalf("EncodeYUV failed: %v", err)
	}

	if len(data) == 0 {
		t.Error("Encoded data should not be empty")
	}

	if data[0] != 0x00 || data[1] != 0x00 || data[2] != 0x00 || data[3] != 0x01 {
		t.Error("H.264 NAL should start with start code")
	}
}

func TestMultipleEncodes(t *testing.T) {
	enc, err := New(640, 480, 30)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer enc.Close()

	for i := 0; i < 10; i++ {
		yuv := generateTestYUV(640, 480)
		data, err := enc.EncodeYUV(yuv)
		if err != nil {
			t.Errorf("EncodeYUV frame %d failed: %v", i, err)
		}
		if len(data) == 0 {
			t.Errorf("Frame %d: encoded data should not be empty", i)
		}
	}
}

func TestEncodeDifferentSizes(t *testing.T) {
	sizes := []struct {
		w, h int
	}{
		{320, 240},
		{640, 480},
		{1280, 720},
	}

	for _, s := range sizes {
		enc, err := New(s.w, s.h, 30)
		if err != nil {
			t.Errorf("New(%d, %d) failed: %v", s.w, s.h, err)
			continue
		}

		yuv := generateTestYUV(s.w, s.h)
		data, err := enc.EncodeYUV(yuv)
		if err != nil {
			t.Errorf("EncodeYUV(%dx%d) failed: %v", s.w, s.h, err)
		}

		enc.Close()

		if len(data) == 0 {
			t.Errorf("EncodeYUV(%dx%d) returned empty data", s.w, s.h)
		}
	}
}

func TestGetFps(t *testing.T) {
	enc, err := New(640, 640, 30)
	if err != nil {
		t.Fatal(err)
	}
	if fps := enc.GetFps(); fps != 30 {
		t.Fatalf("Wanted %d got %d", 30, fps)
	}
}
