package x264

/*
#cgo CFLAGS: -I.
#cgo pkg-config: x264

#include "x264_wrapper.h"
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"runtime"
	"unsafe"
)

type Encoder struct {
	enc    *C.X264Encoder
	fps    int
	width  int
	height int
}

func New(width, height, fps int) (*Encoder, error) {
	if width <= 0 || height <= 0 || fps <= 0 {
		return nil, errors.New("invalid dimensions or fps")
	}

	enc := C.x264_encoder_create(C.int(width), C.int(height), C.int(fps))
	if enc == nil {
		return nil, errors.New("failed to create x264 encoder")
	}

	e := &Encoder{
		enc:    enc,
		width:  width,
		height: height,
		fps:    fps,
	}
	runtime.SetFinalizer(e, (*Encoder).Close)
	return e, nil
}

func (e *Encoder) EncodeYUV(data []byte) ([]byte, error) {
	if e.enc == nil {
		return nil, errors.New("encoder is closed")
	}

	expectedSize := e.width * e.height * 3 / 2
	if len(data) < expectedSize {
		return nil, errors.New("YUV data too small")
	}

	var out *C.uint8_t
	var outSize C.int

	ret := C.x264_encode_frame(e.enc, (*C.uint8_t)(unsafe.Pointer(&data[0])), &out, &outSize)
	if ret != 0 {
		return nil, errors.New("encoding failed")
	}

	if outSize == 0 || out == nil {
		return []byte{}, nil
	}

	result := C.GoBytes(unsafe.Pointer(out), outSize)
	C.free(unsafe.Pointer(out))

	return result, nil
}

func (e *Encoder) Close() {
	if e.enc != nil {
		C.x264_encoder_destroy(e.enc)
		e.enc = nil
	}
}

// GetFps return the number of frames per second the encododer expects
func (e *Encoder) GetFps() int {
	return e.fps
}
