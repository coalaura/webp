// Copyright 2016 <chaishushan{AT}gmail.com>. All rights reserved.
// Copyright 2026 github.com/coalaura. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webp

//#include <webp/decode.h>
import "C"
import "unsafe"

// WEBP_DECODER_ABI_VERSION is the libwebp decoder ABI version.
const (
	WEBP_DECODER_ABI_VERSION = 0x0210 // MAJOR(8b) + MINOR(8b)
)

// _C_WEBP_DECODER_ABI_VERSION exposes the C ABI version for tests.
const (
	_C_WEBP_DECODER_ABI_VERSION = C.WEBP_DECODER_ABI_VERSION // for test
)

// WebPGetDecoderVersion returns the decoder version packed in hex using
// 8 bits each for major/minor/revision.
func WebPGetDecoderVersion() uint {
	return uint(C.WebPGetDecoderVersion())
}

// WebPGetInfo returns width and height from a WebP header and reports whether
// the header is valid.
func WebPGetInfo(data []byte) (width, height int, ok bool) {
	if len(data) == 0 {
		return 0, 0, false
	}
	var cw, ch C.int
	if C.WebPGetInfo((*C.uint8_t)(unsafe.Pointer(&data[0])), C.size_t(len(data)), &cw, &ch) == 0 {
		return 0, 0, false
	}
	return int(cw), int(ch), true
}
