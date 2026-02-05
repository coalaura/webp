package webp

/*
#cgo CFLAGS: -I./internal/libwebp-1.6.0/
#cgo CFLAGS: -I./internal/libwebp-1.6.0/src/
#cgo CFLAGS: -I./internal/include/
#cgo CFLAGS: -Wno-pointer-sign -DWEBP_USE_THREAD
#cgo !windows LDFLAGS: -lm

#include "webp.h"
#include <webp/demux.h>
#include <webp/mux.h>
#include <webp/encode.h>
#include <stdlib.h>
*/
import "C"

import (
	"errors"
	"image"
	"image/color"
	"io"
	"unsafe"
)

// Animation represents an animated WebP image.
// Delay values are in milliseconds.
type Animation struct {
	Image      []image.Image
	Delay      []int
	LoopCount  int
	Background color.RGBA
}

// DecodeAll reads a WEBP image from r and returns all frames.
// Delay values are in milliseconds.
func DecodeAll(r io.Reader) (*Animation, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return decodeAll(data)
}

func decodeAll(data []byte) (*Animation, error) {
	if len(data) == 0 {
		return nil, errors.New("webp: DecodeAll, empty data")
	}

	var decOptions C.WebPAnimDecoderOptions
	if C.WebPAnimDecoderOptionsInit(&decOptions) == 0 {
		return nil, errors.New("webp: DecodeAll, init decoder options failed")
	}
	decOptions.color_mode = C.MODE_RGBA
	decOptions.use_threads = 1

	cdata := C.CBytes(data)
	if cdata == nil {
		return nil, errors.New("webp: DecodeAll, alloc failed")
	}
	defer C.free(cdata)
	webpData := C.WebPData{bytes: (*C.uint8_t)(cdata), size: C.size_t(len(data))}
	dec := C.WebPAnimDecoderNew(&webpData, &decOptions)
	if dec == nil {
		return nil, errors.New("webp: DecodeAll, create decoder failed")
	}
	defer C.WebPAnimDecoderDelete(dec)

	var info C.WebPAnimInfo
	if C.WebPAnimDecoderGetInfo(dec, &info) == 0 {
		return nil, errors.New("webp: DecodeAll, get info failed")
	}

	frameCount := int(info.frame_count)
	width := int(info.canvas_width)
	height := int(info.canvas_height)

	delays := make([]int, 0, frameCount)
	demux := C.WebPAnimDecoderGetDemuxer(dec)
	if demux != nil {
		var iter C.WebPIterator
		if C.WebPDemuxGetFrame(demux, 1, &iter) != 0 {
			for {
				delays = append(delays, int(iter.duration))
				if C.WebPDemuxNextFrame(&iter) == 0 {
					break
				}
			}
			C.WebPDemuxReleaseIterator(&iter)
		}
	}

	frames := make([]image.Image, 0, frameCount)
	timestamps := make([]int, 0, frameCount)
	for C.WebPAnimDecoderHasMoreFrames(dec) != 0 {
		var buf *C.uint8_t
		var timestamp C.int
		if C.WebPAnimDecoderGetNext(dec, &buf, &timestamp) == 0 {
			return nil, errors.New("webp: DecodeAll, decode frame failed")
		}
		timestamps = append(timestamps, int(timestamp))
		size := width * height * 4
		pix := make([]byte, size)
		copy(pix, ((*[1 << 30]byte)(unsafe.Pointer(buf)))[0:size:size])
		frames = append(frames, &image.RGBA{
			Pix:    pix,
			Stride: 4 * width,
			Rect:   image.Rect(0, 0, width, height),
		})
	}

	if len(delays) != len(frames) {
		delays = make([]int, len(frames))
		for i := range frames {
			if i+1 < len(timestamps) {
				delays[i] = timestamps[i+1] - timestamps[i]
			} else if i > 0 {
				delays[i] = timestamps[i] - timestamps[i-1]
			}
		}
	}

	return &Animation{
		Image:      frames,
		Delay:      delays,
		LoopCount:  int(info.loop_count),
		Background: rgbaFromWebPColor(uint32(info.bgcolor)),
	}, nil
}

// EncodeAll writes all frames in anim to w as an animated WEBP.
// Delay values are in milliseconds.
func EncodeAll(w io.Writer, anim *Animation, opt *Options) error {
	if anim == nil {
		return errors.New("webp: EncodeAll, animation is nil")
	}
	if len(anim.Image) == 0 {
		return errors.New("webp: EncodeAll, no frames")
	}
	if len(anim.Delay) != len(anim.Image) {
		return errors.New("webp: EncodeAll, delay count mismatch")
	}

	first := anim.Image[0]
	if first == nil {
		return errors.New("webp: EncodeAll, nil frame")
	}
	bounds := first.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	if width <= 0 || height <= 0 {
		return errors.New("webp: EncodeAll, invalid frame size")
	}
	for i, img := range anim.Image {
		if img == nil {
			return errors.New("webp: EncodeAll, nil frame")
		}
		if img.Bounds().Dx() != width || img.Bounds().Dy() != height {
			return errors.New("webp: EncodeAll, inconsistent frame sizes")
		}
		if anim.Delay[i] < 0 {
			return errors.New("webp: EncodeAll, negative delay")
		}
	}

	var encOptions C.WebPAnimEncoderOptions
	if C.WebPAnimEncoderOptionsInit(&encOptions) == 0 {
		return errors.New("webp: EncodeAll, init encoder options failed")
	}
	encOptions.anim_params.loop_count = C.int(anim.LoopCount)
	encOptions.anim_params.bgcolor = C.uint32_t(webpColorFromRGBA(anim.Background))

	enc := C.WebPAnimEncoderNew(C.int(width), C.int(height), &encOptions)
	if enc == nil {
		return errors.New("webp: EncodeAll, create encoder failed")
	}
	defer C.WebPAnimEncoderDelete(enc)

	method := DefaultMethod
	targetSize := 0
	alphaQuality := DefaultAlphaQuality
	autoFilter := false
	exact := false
	lossless := false
	quality := float32(DefaultQuality)
	if opt != nil {
		method = opt.Method
		targetSize = opt.TargetSize
		alphaQuality = opt.AlphaQuality
		if alphaQuality == 0 {
			alphaQuality = DefaultAlphaQuality
		}
		autoFilter = opt.AutoFilter
		exact = opt.Exact
		lossless = opt.Lossless
		if targetSize == 0 {
			quality = opt.Quality
		}
	}

	var config C.WebPConfig
	if lossless {
		if C.WebPConfigInit(&config) == 0 {
			return errors.New("webp: EncodeAll, init config failed")
		}
		config.lossless = 1
		config.quality = C.float(quality)
	} else {
		if C.WebPConfigPreset(&config, C.WEBP_PRESET_DEFAULT, C.float(quality)) == 0 {
			return errors.New("webp: EncodeAll, init config failed")
		}
	}
	config.method = C.int(method)
	config.target_size = C.int(targetSize)
	config.alpha_quality = C.int(alphaQuality)
	if autoFilter {
		config.autofilter = 1
	}
	if exact {
		config.exact = 1
	}
	if C.WebPValidateConfig(&config) == 0 {
		return errors.New("webp: EncodeAll, invalid config")
	}

	timestamp := 0
	for i, img := range anim.Image {
		rgba := toRGBAImage(img)
		var pic C.WebPPicture
		if C.WebPPictureInit(&pic) == 0 {
			return errors.New("webp: EncodeAll, init picture failed")
		}
		pic.use_argb = 1
		pic.width = C.int(width)
		pic.height = C.int(height)
		if len(rgba.Pix) == 0 {
			C.WebPPictureFree(&pic)
			return errors.New("webp: EncodeAll, empty frame data")
		}
		if C.WebPPictureImportRGBA(&pic, (*C.uint8_t)(unsafe.Pointer(&rgba.Pix[0])), C.int(rgba.Stride)) == 0 {
			C.WebPPictureFree(&pic)
			return errors.New("webp: EncodeAll, import RGBA failed")
		}
		if C.WebPAnimEncoderAdd(enc, &pic, C.int(timestamp), &config) == 0 {
			errStr := C.WebPAnimEncoderGetError(enc)
			C.WebPPictureFree(&pic)
			if errStr != nil && *errStr != 0 {
				return errors.New("webp: EncodeAll, " + C.GoString(errStr))
			}
			return errors.New("webp: EncodeAll, add frame failed")
		}
		C.WebPPictureFree(&pic)
		timestamp += anim.Delay[i]
	}

	if C.WebPAnimEncoderAdd(enc, nil, C.int(timestamp), nil) == 0 {
		errStr := C.WebPAnimEncoderGetError(enc)
		if errStr != nil && *errStr != 0 {
			return errors.New("webp: EncodeAll, " + C.GoString(errStr))
		}
		return errors.New("webp: EncodeAll, finalize failed")
	}

	var out C.WebPData
	if C.WebPAnimEncoderAssemble(enc, &out) == 0 {
		errStr := C.WebPAnimEncoderGetError(enc)
		if errStr != nil && *errStr != 0 {
			return errors.New("webp: EncodeAll, " + C.GoString(errStr))
		}
		return errors.New("webp: EncodeAll, assemble failed")
	}
	defer C.WebPDataClear(&out)

	if out.bytes == nil || out.size == 0 {
		return errors.New("webp: EncodeAll, empty output")
	}
	data := C.GoBytes(unsafe.Pointer(out.bytes), C.int(out.size))
	_, err := w.Write(data)
	return err
}

func rgbaFromWebPColor(argb uint32) color.RGBA {
	return color.RGBA{
		A: uint8(argb & 0xff),
		R: uint8((argb >> 8) & 0xff),
		G: uint8((argb >> 16) & 0xff),
		B: uint8((argb >> 24) & 0xff),
	}
}

func webpColorFromRGBA(c color.RGBA) uint32 {
	return uint32(c.A) | uint32(c.R)<<8 | uint32(c.G)<<16 | uint32(c.B)<<24
}
