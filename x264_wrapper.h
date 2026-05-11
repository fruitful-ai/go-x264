#ifndef X264_WRAPPER_H
#define X264_WRAPPER_H

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

typedef struct X264Encoder X264Encoder;

X264Encoder* x264_encoder_create(int width, int height, int fps);
int x264_encode_frame(X264Encoder* enc, uint8_t* yuv, uint8_t** out, int* out_size);
void x264_encoder_destroy(X264Encoder* enc);

#ifdef __cplusplus
}
#endif

#endif
