// Copyright 2026 github.com/coalaura. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webp

/*
#cgo CFLAGS: -I./internal/libwebp-1.6.0/
#cgo CFLAGS: -I./internal/libwebp-1.6.0/src/
#cgo CFLAGS: -I./internal/include/
#include "webp.h"
*/
import "C"

import (
	"errors"
	"image"
	"io"
	"runtime/cgo"
	"unsafe"
)

type encodeWriterState struct {
	w   io.Writer
	err error
}

//export goWebPWrite
func goWebPWrite(data *C.uint8_t, size C.size_t, customPtr unsafe.Pointer) C.int {
	state := cgo.Handle(uintptr(customPtr)).Value().(*encodeWriterState)
	n, ok := checkedSizeToInt(uint64(size))
	if !ok || (n > 0 && data == nil) {
		state.err = errors.New("webp: invalid native output chunk")
		return 0
	}
	if err := writeAll(state.w, unsafe.Slice((*byte)(unsafe.Pointer(data)), n)); err != nil {
		state.err = err
		return 0
	}
	return 1
}

// EncodeTo streams a WebP encoding directly to w without materializing the
// final encoded byte slice in Go memory. A write error can leave partial data
// in w and is returned unchanged.
func EncodeTo(w io.Writer, m image.Image, opt *Options) error {
	if w == nil || m == nil {
		return errors.New("webp: EncodeTo, nil writer or image")
	}
	method, targetSize, alphaQuality, autoFilter, useThreads, quality, lossless, exact, level, err := encodeSettings(opt)
	if err != nil {
		return err
	}
	var pix []byte
	var width, height, stride, mode int
	switch m := adjustImage(m).(type) {
	case *image.Gray:
		pix, width, height, stride, mode = m.Pix, m.Rect.Dx(), m.Rect.Dy(), m.Stride, 1
	case *RGBImage:
		pix, width, height, stride, mode = m.XPix, m.XRect.Dx(), m.XRect.Dy(), m.XStride, 3
	case *image.RGBA:
		pix, width, height, stride, mode = m.Pix, m.Rect.Dx(), m.Rect.Dy(), m.Stride, 4
	default:
		return errors.New("webp: EncodeTo, unsupported image type")
	}
	if err := validatePackedPixels(pix, width, height, stride, mode); err != nil {
		return err
	}
	state := &encodeWriterState{w: w}
	handle := cgo.NewHandle(state)
	defer handle.Delete()
	ok := C.webpEncodeToWriter(
		(*C.uint8_t)(unsafe.Pointer(&pix[0])), C.int(width), C.int(height), C.int(stride), C.int(mode),
		C.int(boolToInt(lossless)), C.int(boolToInt(exact)), C.int(level), C.float(quality), C.int(method),
		C.int(targetSize), C.int(alphaQuality), C.int(boolToInt(autoFilter)), C.int(boolToInt(useThreads)), unsafe.Pointer(handle),
	)
	if state.err != nil {
		return state.err
	}
	if ok == 0 {
		return errors.New("webp: EncodeTo, encode failed")
	}
	return nil
}
