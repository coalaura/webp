// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Copyright 2026 github.com/coalaura. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
// cgo pointer:
//
// Go1.3: Changes to the garbage collector
// http://golang.org/doc/go1.3#garbage_collector
//
// Go1.6:
// https://github.com/golang/proposal/blob/master/design/12416-cgo-pointers.md
//

package webp

/*
#cgo CFLAGS: -I./internal/libwebp-1.6.0/
#cgo CFLAGS: -I./internal/libwebp-1.6.0/src/
#cgo CFLAGS: -I./internal/include/
#cgo CFLAGS: -Wno-pointer-sign -DWEBP_USE_THREAD -O3 -ffast-math
#cgo !windows LDFLAGS: -lm

#include "webp.h"

#include <webp/decode.h>

#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"image"
	"unsafe"
)

func copyWebPBytes(ptr unsafe.Pointer, size C.size_t) ([]byte, error) {
	n, ok := checkedSizeToInt(uint64(size))
	if !ok || n <= 0 || ptr == nil {
		return nil, errors.New("webp: invalid native output size")
	}
	output := make([]byte, n)
	copy(output, unsafe.Slice((*byte)(ptr), n))
	return output, nil
}

func webpGetInfo(data []byte) (width, height int, hasAlpha bool, hasAnimation bool, format int, err error) {
	if len(data) == 0 {
		err = errors.New("webpGetInfo: bad arguments, data is empty")
		return
	}
	if len(data) > maxWebpHeaderSize {
		data = data[:maxWebpHeaderSize]
	}

	var features C.WebPBitstreamFeatures
	if C.WebPGetFeatures((*C.uint8_t)(unsafe.Pointer(&data[0])), C.size_t(len(data)), &features) != C.VP8_STATUS_OK {
		err = errors.New("C.WebPGetFeatures: failed")
		return
	}
	width, height = int(features.width), int(features.height)
	hasAlpha = (features.has_alpha != 0)
	hasAnimation = (features.has_animation != 0)
	format = int(features.format)
	return
}

func webpDecodeGray(data []byte) (pix []byte, width, height int, err error) {
	if len(data) == 0 {
		err = errors.New("webpDecodeGray: bad arguments")
		return
	}

	var cw, ch C.int
	var cptr = C.webpDecodeGray((*C.uint8_t)(unsafe.Pointer(&data[0])), C.size_t(len(data)), &cw, &ch)
	if cptr == nil {
		err = errors.New("webpDecodeGray: failed")
		return
	}
	defer C.WebPFree(unsafe.Pointer(cptr))

	_, size, sizeErr := decodeBufferSize(int(cw), int(ch), 1)
	if sizeErr != nil {
		return nil, 0, 0, sizeErr
	}
	pix = make([]byte, size)
	copy(pix, unsafe.Slice((*byte)(unsafe.Pointer(cptr)), size))
	width, height = int(cw), int(ch)
	return
}

func webpDecodeRGB(data []byte, opt *DecodeOptions) (pix []byte, width, height int, err error) {
	if len(data) == 0 {
		err = errors.New("webpDecodeRGB: bad arguments")
		return
	}

	width, height, _, _, _, err = webpGetInfo(data)
	if err != nil {
		return
	}
	if width <= 0 || height <= 0 {
		err = errors.New("webpDecodeRGB: invalid image dimensions")
		return
	}

	stride, size, sizeErr := decodeBufferSize(width, height, 3)
	if sizeErr != nil {
		err = sizeErr
		return
	}
	var res C.int
	if opt == nil || (opt.Crop == (image.Rectangle{}) && opt.Width == 0 && opt.Height == 0) {
		pix = make([]byte, size)
		res = C.webpDecodeRGBIntoDefault(
			(*C.uint8_t)(unsafe.Pointer(&data[0])), C.size_t(len(data)),
			C.int(width), C.int(height), C.int(stride),
			(*C.uint8_t)(unsafe.Pointer(&pix[0])), C.int(boolToInt(useDecodeThreads(opt))),
		)
	} else {
		transform, transformErr := decodeTransformFor(width, height, opt, 0, 0)
		if transformErr != nil {
			return nil, 0, 0, transformErr
		}
		width, height = transform.width, transform.height
		stride, size, sizeErr = decodeBufferSize(width, height, 3)
		if sizeErr != nil {
			return nil, 0, 0, sizeErr
		}
		pix = make([]byte, size)
		scaledWidth, scaledHeight := 0, 0
		if opt.Width != 0 {
			scaledWidth, scaledHeight = width, height
		}
		res = C.webpDecodeRGBInto(
			(*C.uint8_t)(unsafe.Pointer(&data[0])), C.size_t(len(data)),
			C.int(width), C.int(height), C.int(stride),
			(*C.uint8_t)(unsafe.Pointer(&pix[0])), C.int(boolToInt(useDecodeThreads(opt))),
			C.int(transform.cropLeft), C.int(transform.cropTop), C.int(transform.cropWidth), C.int(transform.cropHeight),
			C.int(scaledWidth), C.int(scaledHeight),
		)
	}
	if res != C.VP8_STATUS_OK {
		pix = nil
		err = errors.New("webpDecodeRGB: failed")
		return
	}
	return
}

func webpDecodeRGBA(data []byte, opt *DecodeOptions) (pix []byte, width, height int, err error) {
	if len(data) == 0 {
		err = errors.New("webpDecodeRGBA: bad arguments")
		return
	}

	width, height, _, _, _, err = webpGetInfo(data)
	if err != nil {
		return
	}
	if width <= 0 || height <= 0 {
		err = errors.New("webpDecodeRGBA: invalid image dimensions")
		return
	}

	stride, size, sizeErr := decodeBufferSize(width, height, 4)
	if sizeErr != nil {
		err = sizeErr
		return
	}
	var res C.int
	if opt == nil || (opt.Crop == (image.Rectangle{}) && opt.Width == 0 && opt.Height == 0) {
		pix = make([]byte, size)
		res = C.webpDecodeRGBAIntoDefault(
			(*C.uint8_t)(unsafe.Pointer(&data[0])), C.size_t(len(data)),
			C.int(width), C.int(height), C.int(stride),
			(*C.uint8_t)(unsafe.Pointer(&pix[0])), C.int(boolToInt(useDecodeThreads(opt))),
		)
	} else {
		transform, transformErr := decodeTransformFor(width, height, opt, 0, 0)
		if transformErr != nil {
			return nil, 0, 0, transformErr
		}
		width, height = transform.width, transform.height
		stride, size, sizeErr = decodeBufferSize(width, height, 4)
		if sizeErr != nil {
			return nil, 0, 0, sizeErr
		}
		pix = make([]byte, size)
		scaledWidth, scaledHeight := 0, 0
		if opt.Width != 0 {
			scaledWidth, scaledHeight = width, height
		}
		res = C.webpDecodeRGBAInto(
			(*C.uint8_t)(unsafe.Pointer(&data[0])), C.size_t(len(data)),
			C.int(width), C.int(height), C.int(stride),
			(*C.uint8_t)(unsafe.Pointer(&pix[0])), C.int(boolToInt(useDecodeThreads(opt))),
			C.int(transform.cropLeft), C.int(transform.cropTop), C.int(transform.cropWidth), C.int(transform.cropHeight),
			C.int(scaledWidth), C.int(scaledHeight),
		)
	}
	if res != C.VP8_STATUS_OK {
		pix = nil
		err = errors.New("webpDecodeRGBA: failed")
		return
	}
	return
}

func webpDecodeGrayToSize(data []byte, width, height int, opt *DecodeOptions) (pix []byte, err error) {
	if len(data) == 0 {
		return nil, errors.New("webpDecodeGrayToSize: bad arguments")
	}
	sourceWidth, sourceHeight, _, _, _, infoErr := webpGetInfo(data)
	if infoErr != nil {
		return nil, infoErr
	}
	transform, transformErr := decodeTransformFor(sourceWidth, sourceHeight, opt, width, height)
	if transformErr != nil {
		return nil, transformErr
	}
	width, height = transform.width, transform.height
	stride, size, err := decodeBufferSize(width, height, 1)
	if err != nil {
		return nil, err
	}
	pix = make([]byte, size)
	res := C.webpDecodeGrayToSize(
		(*C.uint8_t)(unsafe.Pointer(&data[0])), C.size_t(len(data)),
		C.int(width), C.int(height), C.int(stride),
		(*C.uint8_t)(unsafe.Pointer(&pix[0])), C.int(boolToInt(useDecodeThreads(opt))),
		C.int(transform.cropLeft), C.int(transform.cropTop), C.int(transform.cropWidth), C.int(transform.cropHeight),
	)
	if res != C.VP8_STATUS_OK {
		pix = nil
		err = errors.New("webpDecodeGrayToSize: failed")
	}
	return
}

func webpDecodeRGBToSize(data []byte, width, height int, opt *DecodeOptions) (pix []byte, err error) {
	if len(data) == 0 {
		return nil, errors.New("webpDecodeRGBToSize: bad arguments")
	}
	sourceWidth, sourceHeight, _, _, _, infoErr := webpGetInfo(data)
	if infoErr != nil {
		return nil, infoErr
	}
	transform, transformErr := decodeTransformFor(sourceWidth, sourceHeight, opt, width, height)
	if transformErr != nil {
		return nil, transformErr
	}
	width, height = transform.width, transform.height
	stride, size, err := decodeBufferSize(width, height, 3)
	if err != nil {
		return nil, err
	}
	pix = make([]byte, size)
	res := C.webpDecodeRGBToSize(
		(*C.uint8_t)(unsafe.Pointer(&data[0])), C.size_t(len(data)),
		C.int(width), C.int(height), C.int(stride),
		(*C.uint8_t)(unsafe.Pointer(&pix[0])), C.int(boolToInt(useDecodeThreads(opt))),
		C.int(transform.cropLeft), C.int(transform.cropTop), C.int(transform.cropWidth), C.int(transform.cropHeight),
	)
	if res != C.VP8_STATUS_OK {
		pix = nil
		err = errors.New("webpDecodeRGBToSize: failed")
	}
	return
}

func webpDecodeRGBAToSize(data []byte, width, height int, opt *DecodeOptions) (pix []byte, err error) {
	if len(data) == 0 {
		return nil, errors.New("webpDecodeRGBAToSize: bad arguments")
	}
	sourceWidth, sourceHeight, _, _, _, infoErr := webpGetInfo(data)
	if infoErr != nil {
		return nil, infoErr
	}
	transform, transformErr := decodeTransformFor(sourceWidth, sourceHeight, opt, width, height)
	if transformErr != nil {
		return nil, transformErr
	}
	width, height = transform.width, transform.height
	stride, size, err := decodeBufferSize(width, height, 4)
	if err != nil {
		return nil, err
	}
	pix = make([]byte, size)
	res := C.webpDecodeRGBAToSize(
		(*C.uint8_t)(unsafe.Pointer(&data[0])), C.size_t(len(data)),
		C.int(width), C.int(height), C.int(stride),
		(*C.uint8_t)(unsafe.Pointer(&pix[0])), C.int(boolToInt(useDecodeThreads(opt))),
		C.int(transform.cropLeft), C.int(transform.cropTop), C.int(transform.cropWidth), C.int(transform.cropHeight),
	)
	if res != C.VP8_STATUS_OK {
		pix = nil
		err = errors.New("webpDecodeRGBAToSize: failed")
	}
	return
}

func webpEncodeGray(pix []byte, width, height, stride int, quality float32, method int, targetSize int, alphaQuality int, autoFilter int, threadLevel int) (output []byte, err error) {
	if len(pix) == 0 || width <= 0 || height <= 0 || stride <= 0 || quality < 0.0 {
		err = errors.New("webpEncodeGray: bad arguments")
		return
	}
	if targetSize < 0 || alphaQuality < 0 || alphaQuality > 100 {
		err = errors.New("webpEncodeGray: bad arguments")
		return
	}
	if !validCInts(method, targetSize, alphaQuality, autoFilter, threadLevel) {
		err = errors.New("webpEncodeGray: bad arguments")
		return
	}
	if err = validatePackedPixels(pix, width, height, stride, 1); err != nil {
		return
	}

	var cptr_size C.size_t
	var cptr = C.webpEncodeGray(
		(*C.uint8_t)(unsafe.Pointer(&pix[0])), C.int(width), C.int(height),
		C.int(stride), C.float(quality), C.int(method),
		C.int(targetSize), C.int(alphaQuality), C.int(autoFilter), C.int(threadLevel),
		&cptr_size,
	)
	if cptr == nil || cptr_size == 0 {
		err = errors.New("webpEncodeGray: failed")
		return
	}
	defer C.WebPFree(unsafe.Pointer(cptr))
	output, err = copyWebPBytes(unsafe.Pointer(cptr), cptr_size)
	return
}

func webpEncodeRGB(pix []byte, width, height, stride int, quality float32, method int, targetSize int, alphaQuality int, autoFilter int, threadLevel int) (output []byte, err error) {
	if len(pix) == 0 || width <= 0 || height <= 0 || stride <= 0 || quality < 0.0 {
		err = errors.New("webpEncodeRGB: bad arguments")
		return
	}
	if targetSize < 0 || alphaQuality < 0 || alphaQuality > 100 {
		err = errors.New("webpEncodeRGB: bad arguments")
		return
	}
	if !validCInts(method, targetSize, alphaQuality, autoFilter, threadLevel) {
		err = errors.New("webpEncodeRGB: bad arguments")
		return
	}
	if err = validatePackedPixels(pix, width, height, stride, 3); err != nil {
		return
	}

	var cptr_size C.size_t
	var cptr = C.webpEncodeRGB(
		(*C.uint8_t)(unsafe.Pointer(&pix[0])), C.int(width), C.int(height),
		C.int(stride), C.float(quality), C.int(method),
		C.int(targetSize), C.int(alphaQuality), C.int(autoFilter), C.int(threadLevel),
		&cptr_size,
	)
	if cptr == nil || cptr_size == 0 {
		err = errors.New("webpEncodeRGB: failed")
		return
	}
	defer C.WebPFree(unsafe.Pointer(cptr))
	output, err = copyWebPBytes(unsafe.Pointer(cptr), cptr_size)
	return
}

func webpEncodeRGBA(pix []byte, width, height, stride int, quality float32, method int, targetSize int, alphaQuality int, autoFilter int, threadLevel int) (output []byte, err error) {
	if len(pix) == 0 || width <= 0 || height <= 0 || stride <= 0 || quality < 0.0 {
		err = errors.New("webpEncodeRGBA: bad arguments")
		return
	}
	if targetSize < 0 || alphaQuality < 0 || alphaQuality > 100 {
		err = errors.New("webpEncodeRGBA: bad arguments")
		return
	}
	if !validCInts(method, targetSize, alphaQuality, autoFilter, threadLevel) {
		err = errors.New("webpEncodeRGBA: bad arguments")
		return
	}
	if err = validatePackedPixels(pix, width, height, stride, 4); err != nil {
		return
	}

	var cptr_size C.size_t
	var cptr = C.webpEncodeRGBA(
		(*C.uint8_t)(unsafe.Pointer(&pix[0])), C.int(width), C.int(height),
		C.int(stride), C.float(quality), C.int(method),
		C.int(targetSize), C.int(alphaQuality), C.int(autoFilter), C.int(threadLevel),
		&cptr_size,
	)
	if cptr == nil || cptr_size == 0 {
		err = errors.New("webpEncodeRGBA: failed")
		return
	}
	defer C.WebPFree(unsafe.Pointer(cptr))
	output, err = copyWebPBytes(unsafe.Pointer(cptr), cptr_size)
	return
}

func webpEncodeLosslessGray(pix []byte, width, height, stride int, method int, targetSize int, alphaQuality int, autoFilter int, threadLevel int) (output []byte, err error) {
	if len(pix) == 0 || width <= 0 || height <= 0 || stride <= 0 {
		err = errors.New("webpEncodeLosslessGray: bad arguments")
		return
	}
	if targetSize < 0 || alphaQuality < 0 || alphaQuality > 100 {
		err = errors.New("webpEncodeLosslessGray: bad arguments")
		return
	}
	if !validCInts(method, targetSize, alphaQuality, autoFilter, threadLevel) {
		err = errors.New("webpEncodeLosslessGray: bad arguments")
		return
	}
	if err = validatePackedPixels(pix, width, height, stride, 1); err != nil {
		return
	}

	var cptr_size C.size_t
	var cptr = C.webpEncodeLosslessGray(
		(*C.uint8_t)(unsafe.Pointer(&pix[0])), C.int(width), C.int(height),
		C.int(stride), C.int(method),
		C.int(targetSize), C.int(alphaQuality), C.int(autoFilter), C.int(threadLevel),
		&cptr_size,
	)
	if cptr == nil || cptr_size == 0 {
		err = errors.New("webpEncodeLosslessGray: failed")
		return
	}
	defer C.WebPFree(unsafe.Pointer(cptr))
	output, err = copyWebPBytes(unsafe.Pointer(cptr), cptr_size)
	return
}

func webpEncodeLosslessRGB(pix []byte, width, height, stride int, method int, targetSize int, alphaQuality int, autoFilter int, threadLevel int) (output []byte, err error) {
	if len(pix) == 0 || width <= 0 || height <= 0 || stride <= 0 {
		err = errors.New("webpEncodeLosslessRGB: bad arguments")
		return
	}
	if targetSize < 0 || alphaQuality < 0 || alphaQuality > 100 {
		err = errors.New("webpEncodeLosslessRGB: bad arguments")
		return
	}
	if !validCInts(method, targetSize, alphaQuality, autoFilter, threadLevel) {
		err = errors.New("webpEncodeLosslessRGB: bad arguments")
		return
	}
	if err = validatePackedPixels(pix, width, height, stride, 3); err != nil {
		return
	}

	var cptr_size C.size_t
	var cptr = C.webpEncodeLosslessRGB(
		(*C.uint8_t)(unsafe.Pointer(&pix[0])), C.int(width), C.int(height),
		C.int(stride), C.int(method),
		C.int(targetSize), C.int(alphaQuality), C.int(autoFilter), C.int(threadLevel),
		&cptr_size,
	)
	if cptr == nil || cptr_size == 0 {
		err = errors.New("webpEncodeLosslessRGB: failed")
		return
	}
	defer C.WebPFree(unsafe.Pointer(cptr))
	output, err = copyWebPBytes(unsafe.Pointer(cptr), cptr_size)
	return
}

func webpEncodeLosslessRGBA(exact int, pix []byte, width, height, stride int, method int, targetSize int, alphaQuality int, autoFilter int, threadLevel int) (output []byte, err error) {
	if len(pix) == 0 || width <= 0 || height <= 0 || stride <= 0 {
		err = errors.New("webpEncodeLosslessRGBA: bad arguments")
		return
	}
	if targetSize < 0 || alphaQuality < 0 || alphaQuality > 100 {
		err = errors.New("webpEncodeLosslessRGBA: bad arguments")
		return
	}
	if !validCInts(exact, method, targetSize, alphaQuality, autoFilter, threadLevel) {
		err = errors.New("webpEncodeLosslessRGBA: bad arguments")
		return
	}
	if err = validatePackedPixels(pix, width, height, stride, 4); err != nil {
		return
	}

	var cptr_size C.size_t
	var cptr = C.webpEncodeLosslessRGBA(
		C.int(exact), (*C.uint8_t)(unsafe.Pointer(&pix[0])), C.int(width), C.int(height),
		C.int(stride), C.int(method),
		C.int(targetSize), C.int(alphaQuality), C.int(autoFilter), C.int(threadLevel),
		&cptr_size,
	)
	if cptr == nil || cptr_size == 0 {
		err = errors.New("webpEncodeLosslessRGBA: failed")
		return
	}
	defer C.WebPFree(unsafe.Pointer(cptr))
	output, err = copyWebPBytes(unsafe.Pointer(cptr), cptr_size)
	return
}

func webpGetEXIF(data []byte) (metadata []byte, err error) {
	if len(data) == 0 {
		err = errors.New("webpGetEXIF: bad arguments")
		return
	}

	var cptr_size C.size_t
	var cptr = C.webpGetEXIF(
		(*C.uint8_t)(unsafe.Pointer(&data[0])), C.size_t(len(data)),
		&cptr_size,
	)
	if cptr == nil || cptr_size == 0 {
		err = errors.New("webpGetEXIF: failed")
		return
	}
	defer C.free(unsafe.Pointer(cptr))
	metadata, err = copyWebPBytes(unsafe.Pointer(cptr), cptr_size)
	return
}
func webpGetICCP(data []byte) (metadata []byte, err error) {
	if len(data) == 0 {
		err = errors.New("webpGetICCP: bad arguments")
		return
	}

	var cptr_size C.size_t
	var cptr = C.webpGetICCP(
		(*C.uint8_t)(unsafe.Pointer(&data[0])), C.size_t(len(data)),
		&cptr_size,
	)
	if cptr == nil || cptr_size == 0 {
		err = errors.New("webpGetICCP: failed")
		return
	}
	defer C.free(unsafe.Pointer(cptr))
	metadata, err = copyWebPBytes(unsafe.Pointer(cptr), cptr_size)
	return
}
func webpGetXMP(data []byte) (metadata []byte, err error) {
	if len(data) == 0 {
		err = errors.New("webpGetXMP: bad arguments")
		return
	}

	var cptr_size C.size_t
	var cptr = C.webpGetXMP(
		(*C.uint8_t)(unsafe.Pointer(&data[0])), C.size_t(len(data)),
		&cptr_size,
	)
	if cptr == nil || cptr_size == 0 {
		err = errors.New("webpGetXMP: failed")
		return
	}
	defer C.free(unsafe.Pointer(cptr))
	metadata, err = copyWebPBytes(unsafe.Pointer(cptr), cptr_size)
	return
}
func webpGetMetadata(data []byte, format string) (metadata []byte, err error) {
	if len(data) == 0 {
		err = errors.New("webpGetMetadata: bad arguments")
		return
	}

	switch format {
	case "EXIF":
		return webpCopyMetadata(data, 1)
	case "ICCP":
		return webpCopyMetadata(data, 2)
	case "XMP":
		return webpCopyMetadata(data, 3)
	default:
		err = errors.New("webpGetMetadata: unknown format")
		return
	}
}

func webpCopyMetadata(data []byte, metadataType int) ([]byte, error) {
	var size C.size_t
	if C.webpGetMetadataSize((*C.uint8_t)(unsafe.Pointer(&data[0])), C.size_t(len(data)), C.int(metadataType), &size) == 0 {
		return nil, errors.New("webp: metadata not found")
	}
	n, ok := checkedSizeToInt(uint64(size))
	if !ok || n == 0 {
		return nil, errors.New("webp: invalid metadata size")
	}
	metadata := make([]byte, n)
	if C.webpCopyMetadata((*C.uint8_t)(unsafe.Pointer(&data[0])), C.size_t(len(data)), C.int(metadataType), (*C.uint8_t)(unsafe.Pointer(&metadata[0])), size) == 0 {
		return nil, errors.New("webp: metadata copy failed")
	}
	return metadata, nil
}

func webpSetEXIF(data, metadata []byte) (newData []byte, err error) {
	if len(data) == 0 || len(metadata) == 0 {
		err = errors.New("webpSetEXIF: bad arguments")
		return
	}

	var cptr_size C.size_t
	var cptr = C.webpSetEXIF(
		(*C.uint8_t)(unsafe.Pointer(&data[0])), C.size_t(len(data)),
		(*C.char)(unsafe.Pointer(&metadata[0])), C.size_t(len(metadata)),
		&cptr_size,
	)
	if cptr == nil || cptr_size == 0 {
		err = errors.New("webpSetEXIF: failed")
		return
	}
	defer C.WebPFree(unsafe.Pointer(cptr))
	newData, err = copyWebPBytes(unsafe.Pointer(cptr), cptr_size)
	return
}
func webpSetICCP(data, metadata []byte) (newData []byte, err error) {
	if len(data) == 0 || len(metadata) == 0 {
		err = errors.New("webpSetICCP: bad arguments")
		return
	}

	var cptr_size C.size_t
	var cptr = C.webpSetICCP(
		(*C.uint8_t)(unsafe.Pointer(&data[0])), C.size_t(len(data)),
		(*C.char)(unsafe.Pointer(&metadata[0])), C.size_t(len(metadata)),
		&cptr_size,
	)
	if cptr == nil || cptr_size == 0 {
		err = errors.New("webpSetICCP: failed")
		return
	}
	defer C.WebPFree(unsafe.Pointer(cptr))
	newData, err = copyWebPBytes(unsafe.Pointer(cptr), cptr_size)
	return
}
func webpSetXMP(data, metadata []byte) (newData []byte, err error) {
	if len(data) == 0 || len(metadata) == 0 {
		err = errors.New("webpSetXMP: bad arguments")
		return
	}

	var cptr_size C.size_t
	var cptr = C.webpSetXMP(
		(*C.uint8_t)(unsafe.Pointer(&data[0])), C.size_t(len(data)),
		(*C.char)(unsafe.Pointer(&metadata[0])), C.size_t(len(metadata)),
		&cptr_size,
	)
	if cptr == nil || cptr_size == 0 {
		err = errors.New("webpSetXMP: failed")
		return
	}
	defer C.WebPFree(unsafe.Pointer(cptr))
	newData, err = copyWebPBytes(unsafe.Pointer(cptr), cptr_size)
	return
}
func webpSetMetadata(data, metadata []byte, format string) (newData []byte, err error) {
	if len(data) == 0 || len(metadata) == 0 {
		err = errors.New("webpSetMetadata: bad arguments")
		return
	}

	switch format {
	case "EXIF":
		return webpSetEXIF(data, metadata)
	case "ICCP":
		return webpSetICCP(data, metadata)
	case "XMP":
		return webpSetXMP(data, metadata)
	default:
		err = errors.New("webpSetMetadata: unknown format")
		return
	}
}

func webpDelEXIF(data []byte) (newData []byte, err error) {
	if len(data) == 0 {
		err = errors.New("webpDelEXIF: bad arguments")
		return
	}

	var cptr_size C.size_t
	var cptr = C.webpDelEXIF(
		(*C.uint8_t)(unsafe.Pointer(&data[0])), C.size_t(len(data)),
		&cptr_size,
	)
	if cptr == nil || cptr_size == 0 {
		err = errors.New("webpDelEXIF: failed")
		return
	}
	defer C.WebPFree(unsafe.Pointer(cptr))
	newData, err = copyWebPBytes(unsafe.Pointer(cptr), cptr_size)
	return
}
func webpDelICCP(data []byte) (newData []byte, err error) {
	if len(data) == 0 {
		err = errors.New("webpDelICCP: bad arguments")
		return
	}

	var cptr_size C.size_t
	var cptr = C.webpDelICCP(
		(*C.uint8_t)(unsafe.Pointer(&data[0])), C.size_t(len(data)),
		&cptr_size,
	)
	if cptr == nil || cptr_size == 0 {
		err = errors.New("webpDelICCP: failed")
		return
	}
	defer C.WebPFree(unsafe.Pointer(cptr))
	newData, err = copyWebPBytes(unsafe.Pointer(cptr), cptr_size)
	return
}
func webpDelXMP(data []byte) (newData []byte, err error) {
	if len(data) == 0 {
		err = errors.New("webpDelXMP: bad arguments")
		return
	}

	var cptr_size C.size_t
	var cptr = C.webpDelXMP(
		(*C.uint8_t)(unsafe.Pointer(&data[0])), C.size_t(len(data)),
		&cptr_size,
	)
	if cptr == nil || cptr_size == 0 {
		err = errors.New("webpDelXMP: failed")
		return
	}
	defer C.WebPFree(unsafe.Pointer(cptr))
	newData, err = copyWebPBytes(unsafe.Pointer(cptr), cptr_size)
	return
}
