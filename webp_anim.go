// Copyright 2026 github.com/coalaura. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webp

/*
#cgo CFLAGS: -I./internal/libwebp-1.6.0/
#cgo CFLAGS: -I./internal/libwebp-1.6.0/src/
#cgo CFLAGS: -I./internal/include/
#cgo CFLAGS: -Wno-pointer-sign -DWEBP_USE_THREAD -O3 -ffast-math
#cgo !windows LDFLAGS: -lm

#include "webp.h"
#include <webp/demux.h>
#include <webp/mux.h>
#include <webp/encode.h>
#include <stdlib.h>
#include <string.h>
*/
import "C"

import (
	"errors"
	"image"
	"image/color"
	"io"
	"runtime"
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

// AnimationInfo describes the canvas and playback settings supplied to DecodeFrames.
type AnimationInfo struct {
	Width      int
	Height     int
	FrameCount int
	LoopCount  int
	Background color.RGBA
}

// DecodeAll reads a WEBP image from r and returns all frames.
// Delay values are in milliseconds. If opt is nil, defaults are used.
func DecodeAll(r io.Reader, opt *DecodeOptions) (*Animation, error) {
	useThreads := useDecodeThreads(opt)

	data, release, err := readAllPooled(r)
	if err != nil {
		return nil, err
	}
	if release != nil {
		defer release()
	}
	return decodeAll(data, useThreads)
}

func decodeAll(data []byte, useThreads bool) (*Animation, error) {
	frames := make([]image.Image, 0)
	delays := make([]int, 0)
	info, err := decodeFrames(data, useThreads, func(frame *image.RGBA, delay int) error {
		pix := make([]byte, len(frame.Pix))
		copy(pix, frame.Pix)
		frames = append(frames, &image.RGBA{Pix: pix, Stride: frame.Stride, Rect: frame.Rect})
		delays = append(delays, delay)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &Animation{Image: frames, Delay: delays, LoopCount: info.LoopCount, Background: info.Background}, nil
}

// DecodeFrames calls fn for each composited RGBA animation frame. The frame
// buffer is owned by the decoder and is only valid until fn returns.
func DecodeFrames(r io.Reader, opt *DecodeOptions, fn func(frame *image.RGBA, delay int) error) (AnimationInfo, error) {
	if fn == nil {
		return AnimationInfo{}, errors.New("webp: DecodeFrames, callback is nil")
	}
	data, release, err := readAllPooled(r)
	if err != nil {
		return AnimationInfo{}, err
	}
	defer release()
	return decodeFrames(data, useDecodeThreads(opt), fn)
}

func decodeFrames(data []byte, useThreads bool, fn func(frame *image.RGBA, delay int) error) (AnimationInfo, error) {
	if len(data) == 0 {
		return AnimationInfo{}, errors.New("webp: DecodeAll, empty data")
	}

	var decOptions C.WebPAnimDecoderOptions
	if C.WebPAnimDecoderOptionsInit(&decOptions) == 0 {
		return AnimationInfo{}, errors.New("webp: DecodeAll, init decoder options failed")
	}
	decOptions.color_mode = C.MODE_RGBA
	decOptions.use_threads = C.int(boolToInt(useThreads))

	cptr := C.malloc(C.size_t(len(data)))
	if cptr == nil {
		return AnimationInfo{}, errors.New("webp: DecodeAll, alloc input failed")
	}
	defer C.free(cptr)
	C.memcpy(cptr, unsafe.Pointer(&data[0]), C.size_t(len(data)))
	runtime.KeepAlive(data)
	webpData := C.WebPData{bytes: (*C.uint8_t)(cptr), size: C.size_t(len(data))}
	dec := C.WebPAnimDecoderNew(&webpData, &decOptions)
	if dec == nil {
		return AnimationInfo{}, errors.New("webp: DecodeAll, create decoder failed")
	}
	defer C.WebPAnimDecoderDelete(dec)

	var info C.WebPAnimInfo
	if C.WebPAnimDecoderGetInfo(dec, &info) == 0 {
		return AnimationInfo{}, errors.New("webp: DecodeAll, get info failed")
	}

	frameCount := int(info.frame_count)
	width := int(info.canvas_width)
	height := int(info.canvas_height)
	stride, frameSize, sizeErr := decodeBufferSize(width, height, 4)
	if sizeErr != nil || frameCount < 0 {
		return AnimationInfo{}, errors.New("webp: DecodeAll, invalid animation dimensions")
	}

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

	index := 0
	for C.WebPAnimDecoderHasMoreFrames(dec) != 0 {
		var buf *C.uint8_t
		var timestamp C.int
		if C.WebPAnimDecoderGetNext(dec, &buf, &timestamp) == 0 {
			return AnimationInfo{}, errors.New("webp: DecodeAll, decode frame failed")
		}
		delay := 0
		if index < len(delays) {
			delay = delays[index]
		}
		frame := &image.RGBA{Pix: unsafe.Slice((*byte)(unsafe.Pointer(buf)), frameSize), Stride: stride, Rect: image.Rect(0, 0, width, height)}
		if err := fn(frame, delay); err != nil {
			return AnimationInfo{}, err
		}
		index++
	}
	return AnimationInfo{Width: width, Height: height, FrameCount: index, LoopCount: int(info.loop_count), Background: rgbaFromWebPColor(uint32(info.bgcolor))}, nil
}

// EncodeAll writes all frames in anim to w as an animated WEBP.
// Delay values are in milliseconds. If opt is nil, defaults are used.
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
	useThreads := false
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
		useThreads = opt.UseThreads
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
	if useThreads {
		config.thread_level = 1
	} else {
		config.thread_level = 0
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
	size, ok := checkedSizeToInt(uint64(out.size))
	if !ok {
		return errors.New("webp: EncodeAll, output is too large")
	}
	return writeAll(w, unsafe.Slice((*byte)(unsafe.Pointer(out.bytes)), size))
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
