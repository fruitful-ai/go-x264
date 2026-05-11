#include "x264_wrapper.h"
#include <x264.h>
#include <stdlib.h>
#include <string.h>

struct X264Encoder {
    x264_t* enc;
    x264_picture_t pic;
    x264_picture_t pic_out;
    int width;
    int height;
    int64_t pts;
};

static void set_params(x264_param_t* param, int width, int height, int fps) {
    x264_param_default_preset(param, "veryfast", "zerolatency");

    param->i_width = width;
    param->i_height = height;
    param->i_fps_num = fps;
    param->i_fps_den = 1;

    param->i_csp = X264_CSP_I420;

    param->i_bframe = 0;
    param->i_bframe_adaptive = X264_B_ADAPT_NONE;

    param->i_level_idc = 30;

    param->analyse.i_me_range = 16;
    param->analyse.i_subpel_refine = 1;

    param->b_repeat_headers = 1;
    param->b_annexb = 1;

    x264_param_apply_profile(param, "baseline");
}

X264Encoder* x264_encoder_create(int width, int height, int fps) {
    x264_param_t param;
    set_params(&param, width, height, fps);

    x264_t* enc = x264_encoder_open(&param);
    if (enc == NULL) {
        return NULL;
    }

    X264Encoder* wrapper = (X264Encoder*)malloc(sizeof(X264Encoder));
    if (wrapper == NULL) {
        x264_encoder_close(enc);
        return NULL;
    }

    memset(wrapper, 0, sizeof(X264Encoder));
    wrapper->enc = enc;
    wrapper->width = width;
    wrapper->height = height;
    wrapper->pts = 0;

    x264_picture_init(&wrapper->pic);
    wrapper->pic.img.i_csp = X264_CSP_I420;
    wrapper->pic.img.i_plane = 3;
    wrapper->pic.img.plane[0] = NULL;
    wrapper->pic.img.plane[1] = NULL;
    wrapper->pic.img.plane[2] = NULL;
    wrapper->pic.img.plane[3] = NULL;
    wrapper->pic.img.i_stride[0] = width;
    wrapper->pic.img.i_stride[1] = width / 2;
    wrapper->pic.img.i_stride[2] = width / 2;
    wrapper->pic.img.i_stride[3] = 0;
    wrapper->pic.i_pts = 0;
    wrapper->pic.i_type = X264_TYPE_AUTO;

    return wrapper;
}

int x264_encode_frame(X264Encoder* enc, uint8_t* yuv, uint8_t** out, int* out_size) {
    if (enc == NULL || enc->enc == NULL || yuv == NULL || out == NULL || out_size == NULL) {
        return -1;
    }

    enc->pic.img.plane[0] = yuv;
    enc->pic.img.plane[1] = yuv + enc->width * enc->height;
    enc->pic.img.plane[2] = yuv + enc->width * enc->height * 5 / 4;
    enc->pic.img.plane[3] = NULL;

    enc->pic.i_pts = enc->pts;

    x264_nal_t* nals = NULL;
    int i_nals = 0;

    int ret = x264_encoder_encode(enc->enc, &nals, &i_nals, &enc->pic, &enc->pic_out);
    if (ret < 0) {
        *out = NULL;
        *out_size = 0;
        return ret;
    }

    if (i_nals == 0) {
        *out = NULL;
        *out_size = 0;
        return 0;
    }

    size_t total_size = 0;
    for (int i = 0; i < i_nals; i++) {
        total_size += nals[i].i_payload;
    }

    if (total_size == 0) {
        *out = NULL;
        *out_size = 0;
        return 0;
    }

    uint8_t* output = (uint8_t*)malloc(total_size);
    if (output == NULL) {
        *out = NULL;
        *out_size = 0;
        return -1;
    }

    size_t offset = 0;
    for (int i = 0; i < i_nals; i++) {
        memcpy(output + offset, nals[i].p_payload, nals[i].i_payload);
        offset += nals[i].i_payload;
    }

    *out = output;
    *out_size = (int)total_size;
    enc->pts++;

    x264_encoder_delayed_frames(enc->enc);

    return 0;
}

void x264_encoder_destroy(X264Encoder* enc) {
    if (enc == NULL) {
        return;
    }

    if (enc->enc != NULL) {
        x264_encoder_close(enc->enc);
        enc->enc = NULL;
    }

    free(enc);
}