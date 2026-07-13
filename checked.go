// Copyright 2026 github.com/coalaura. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webp

import (
	"errors"
	"image"
	"math"
)

type decodeTransform struct {
	cropLeft, cropTop, cropWidth, cropHeight int
	width, height                            int
}

func decodeTransformFor(sourceWidth, sourceHeight int, opt *DecodeOptions, width, height int) (decodeTransform, error) {
	if sourceWidth <= 0 || sourceHeight <= 0 {
		return decodeTransform{}, errInvalidBuffer
	}
	transform := decodeTransform{width: sourceWidth, height: sourceHeight}
	if opt != nil && opt.Crop != (image.Rectangle{}) {
		crop := opt.Crop
		if crop.Empty() || crop.Min.X < 0 || crop.Min.Y < 0 || crop.Max.X > sourceWidth || crop.Max.Y > sourceHeight {
			return decodeTransform{}, errInvalidBuffer
		}
		transform.cropLeft = crop.Min.X
		transform.cropTop = crop.Min.Y
		transform.cropWidth = crop.Dx()
		transform.cropHeight = crop.Dy()
		transform.width = crop.Dx()
		transform.height = crop.Dy()
	}
	if width == 0 && height == 0 && opt != nil {
		width, height = opt.Width, opt.Height
	}
	if (width == 0) != (height == 0) || width < 0 || height < 0 {
		return decodeTransform{}, errInvalidBuffer
	}
	if width > 0 {
		transform.width = width
		transform.height = height
	}
	if _, _, err := decodeBufferSize(transform.width, transform.height, 1); err != nil {
		return decodeTransform{}, err
	}
	for _, value := range []int{transform.cropLeft, transform.cropTop, transform.cropWidth, transform.cropHeight} {
		if _, ok := checkedCInt(value); !ok {
			return decodeTransform{}, errInvalidBuffer
		}
	}
	return transform, nil
}

const maxCInt = math.MaxInt32

var errInvalidBuffer = errors.New("webp: invalid pixel buffer")

func checkedAdd(a, b int) (int, bool) {
	if b > 0 && a > math.MaxInt-b {
		return 0, false
	}
	return a + b, true
}

func checkedMul(a, b int) (int, bool) {
	if a < 0 || b < 0 || (a != 0 && b > math.MaxInt/a) {
		return 0, false
	}
	return a * b, true
}

func checkedPixels(width, height, channels int) (int, int, bool) {
	row, ok := checkedMul(width, channels)
	if !ok {
		return 0, 0, false
	}
	size, ok := checkedMul(row, height)
	return row, size, ok
}

func checkedCInt(v int) (int, bool) {
	return v, v >= 0 && v <= maxCInt
}

func validCInts(values ...int) bool {
	for _, v := range values {
		if _, ok := checkedCInt(v); !ok {
			return false
		}
	}
	return true
}

func checkedSizeToInt(v uint64) (int, bool) {
	if v > uint64(math.MaxInt) {
		return 0, false
	}
	return int(v), true
}

func validatePackedPixels(pix []byte, width, height, stride, channels int) error {
	if width <= 0 || height <= 0 || stride <= 0 || channels <= 0 {
		return errInvalidBuffer
	}
	row, ok := checkedMul(width, channels)
	if !ok || stride < row {
		return errInvalidBuffer
	}
	lastRow, ok := checkedMul(height-1, stride)
	if !ok {
		return errInvalidBuffer
	}
	required, ok := checkedAdd(lastRow, row)
	if !ok || required > len(pix) {
		return errInvalidBuffer
	}
	if _, ok := checkedCInt(width); !ok {
		return errInvalidBuffer
	}
	if _, ok := checkedCInt(height); !ok {
		return errInvalidBuffer
	}
	if _, ok := checkedCInt(stride); !ok {
		return errInvalidBuffer
	}
	return nil
}

func decodeBufferSize(width, height, channels int) (row, size int, err error) {
	if width <= 0 || height <= 0 || channels <= 0 {
		return 0, 0, errInvalidBuffer
	}
	row, size, ok := checkedPixels(width, height, channels)
	if !ok {
		return 0, 0, errInvalidBuffer
	}
	if _, ok := checkedCInt(width); !ok {
		return 0, 0, errInvalidBuffer
	}
	if _, ok := checkedCInt(height); !ok {
		return 0, 0, errInvalidBuffer
	}
	if _, ok := checkedCInt(row); !ok {
		return 0, 0, errInvalidBuffer
	}
	return row, size, nil
}
