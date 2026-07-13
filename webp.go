// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Copyright 2026 github.com/coalaura. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webp

import (
	"image"
	"strings"

	"embed"
)

//go:embed internal
var _ embed.FS

const (
	maxWebpHeaderSize = 32
)

// GetInfo returns image dimensions and WebP bitstream flags.
func GetInfo(data []byte) (width, height int, hasAlpha bool, hasAnimation bool, format int, err error) {
	return webpGetInfo(data)
}

// DecodeGray decodes a WebP payload into a Gray image.
// If opt is nil, defaults are used.
func DecodeGray(data []byte, opt *DecodeOptions) (m *image.Gray, err error) {
	if opt != nil && (opt.Crop != (image.Rectangle{}) || opt.Width != 0 || opt.Height != 0 || opt.UseThreads) {
		width, height, _, _, _, infoErr := webpGetInfo(data)
		if infoErr != nil {
			return nil, infoErr
		}
		transform, transformErr := decodeTransformFor(width, height, opt, 0, 0)
		if transformErr != nil {
			return nil, transformErr
		}
		pix, decodeErr := webpDecodeGrayToSize(data, transform.width, transform.height, opt)
		if decodeErr != nil {
			return nil, decodeErr
		}
		return &image.Gray{Pix: pix, Stride: transform.width, Rect: image.Rect(0, 0, transform.width, transform.height)}, nil
	}
	pix, w, h, err := webpDecodeGray(data)
	if err != nil {
		return
	}
	m = &image.Gray{
		Pix:    pix,
		Stride: 1 * w,
		Rect:   image.Rect(0, 0, w, h),
	}
	return
}

// DecodeRGB decodes a WebP payload into an RGBImage.
// If opt is nil, defaults are used.
func DecodeRGB(data []byte, opt *DecodeOptions) (m *RGBImage, err error) {
	pix, w, h, err := webpDecodeRGB(data, opt)
	if err != nil {
		return
	}
	m = &RGBImage{
		XPix:    pix,
		XStride: 3 * w,
		XRect:   image.Rect(0, 0, w, h),
	}
	return
}

// DecodeRGBA decodes a WebP payload into an RGBA image.
// If opt is nil, defaults are used.
func DecodeRGBA(data []byte, opt *DecodeOptions) (m *image.RGBA, err error) {
	pix, w, h, err := webpDecodeRGBA(data, opt)
	if err != nil {
		return
	}
	m = &image.RGBA{
		Pix:    pix,
		Stride: 4 * w,
		Rect:   image.Rect(0, 0, w, h),
	}
	return
}

// DecodeGrayToSize decodes a Gray image scaled to the given dimensions. For
// large images, the DecodeXXXToSize methods are significantly faster and
// require less memory compared to decoding a full-size image and then resizing it.
// If opt is nil, defaults are used.
func DecodeGrayToSize(data []byte, width, height int, opt *DecodeOptions) (m *image.Gray, err error) {
	pix, err := webpDecodeGrayToSize(data, width, height, opt)
	if err != nil {
		return
	}
	m = &image.Gray{
		Pix:    pix,
		Stride: width,
		Rect:   image.Rect(0, 0, width, height),
	}
	return
}

// DecodeRGBToSize decodes an RGB image scaled to the given dimensions.
// If opt is nil, defaults are used.
func DecodeRGBToSize(data []byte, width, height int, opt *DecodeOptions) (m *RGBImage, err error) {
	pix, err := webpDecodeRGBToSize(data, width, height, opt)
	if err != nil {
		return
	}
	m = &RGBImage{
		XPix:    pix,
		XStride: 3 * width,
		XRect:   image.Rect(0, 0, width, height),
	}
	return
}

// DecodeRGBAToSize decodes an RGBA image scaled to the given dimensions.
// If opt is nil, defaults are used.
func DecodeRGBAToSize(data []byte, width, height int, opt *DecodeOptions) (m *image.RGBA, err error) {
	pix, err := webpDecodeRGBAToSize(data, width, height, opt)
	if err != nil {
		return
	}
	m = &image.RGBA{
		Pix:    pix,
		Stride: 4 * width,
		Rect:   image.Rect(0, 0, width, height),
	}
	return
}

// EncodeGray encodes an image as lossy grayscale WebP.
func EncodeGray(m image.Image, quality float32, method int) (data []byte, err error) {
	return encodeGray(m, quality, method, 0, DefaultAlphaQuality, false, false)
}

// EncodeRGB encodes an image as lossy RGB WebP.
func EncodeRGB(m image.Image, quality float32, method int) (data []byte, err error) {
	return encodeRGB(m, quality, method, 0, DefaultAlphaQuality, false, false)
}

// EncodeRGBA encodes an image as lossy RGBA WebP.
func EncodeRGBA(m image.Image, quality float32, method int) (data []byte, err error) {
	return encodeRGBA(m, quality, method, 0, DefaultAlphaQuality, false, false)
}

// EncodeLosslessGray encodes an image as lossless grayscale WebP.
func EncodeLosslessGray(m image.Image, method int) (data []byte, err error) {
	return encodeLosslessGray(m, method, 0, DefaultAlphaQuality, false, false)
}

// EncodeLosslessRGB encodes an image as lossless RGB WebP.
func EncodeLosslessRGB(m image.Image, method int) (data []byte, err error) {
	return encodeLosslessRGB(m, method, 0, DefaultAlphaQuality, false, false)
}

// EncodeLosslessRGBA encodes an image as lossless RGBA WebP.
func EncodeLosslessRGBA(m image.Image, method int) (data []byte, err error) {
	return encodeLosslessRGBA(m, method, 0, DefaultAlphaQuality, false, false)
}

// EncodeExactLosslessRGBA encodes lossless RGBA WebP and preserves RGB values
// in transparent areas.
func EncodeExactLosslessRGBA(m image.Image, method int) (data []byte, err error) {
	return encodeExactLosslessRGBA(m, method, 0, DefaultAlphaQuality, false, false)
}

func encodeGray(m image.Image, quality float32, method int, targetSize int, alphaQuality int, autoFilter bool, useThreads bool) (data []byte, err error) {
	p := toGrayImage(m)
	data, err = webpEncodeGray(p.Pix, p.Rect.Dx(), p.Rect.Dy(), p.Stride, quality, method, targetSize, alphaQuality, boolToInt(autoFilter), boolToInt(useThreads))
	if err != nil {
		return
	}
	return
}

func encodeRGB(m image.Image, quality float32, method int, targetSize int, alphaQuality int, autoFilter bool, useThreads bool) (data []byte, err error) {
	p := NewRGBImageFrom(m)
	data, err = webpEncodeRGB(p.XPix, p.XRect.Dx(), p.XRect.Dy(), p.XStride, quality, method, targetSize, alphaQuality, boolToInt(autoFilter), boolToInt(useThreads))
	return
}

func encodeRGBA(m image.Image, quality float32, method int, targetSize int, alphaQuality int, autoFilter bool, useThreads bool) (data []byte, err error) {
	p := toRGBAImage(m)
	data, err = webpEncodeRGBA(p.Pix, p.Rect.Dx(), p.Rect.Dy(), p.Stride, quality, method, targetSize, alphaQuality, boolToInt(autoFilter), boolToInt(useThreads))
	return
}

func encodeLosslessGray(m image.Image, method int, targetSize int, alphaQuality int, autoFilter bool, useThreads bool) (data []byte, err error) {
	p := toGrayImage(m)
	data, err = webpEncodeLosslessGray(p.Pix, p.Rect.Dx(), p.Rect.Dy(), p.Stride, method, targetSize, alphaQuality, boolToInt(autoFilter), boolToInt(useThreads))
	return
}

func encodeLosslessRGB(m image.Image, method int, targetSize int, alphaQuality int, autoFilter bool, useThreads bool) (data []byte, err error) {
	p := NewRGBImageFrom(m)
	data, err = webpEncodeLosslessRGB(p.XPix, p.XRect.Dx(), p.XRect.Dy(), p.XStride, method, targetSize, alphaQuality, boolToInt(autoFilter), boolToInt(useThreads))
	return
}

func encodeLosslessRGBA(m image.Image, method int, targetSize int, alphaQuality int, autoFilter bool, useThreads bool) (data []byte, err error) {
	p := toRGBAImage(m)
	data, err = webpEncodeLosslessRGBA(0, p.Pix, p.Rect.Dx(), p.Rect.Dy(), p.Stride, method, targetSize, alphaQuality, boolToInt(autoFilter), boolToInt(useThreads))
	return
}

func encodeExactLosslessRGBA(m image.Image, method int, targetSize int, alphaQuality int, autoFilter bool, useThreads bool) (data []byte, err error) {
	p := toRGBAImage(m)
	data, err = webpEncodeLosslessRGBA(1, p.Pix, p.Rect.Dx(), p.Rect.Dy(), p.Stride, method, targetSize, alphaQuality, boolToInt(autoFilter), boolToInt(useThreads))
	return
}

func useDecodeThreads(opt *DecodeOptions) bool {
	if opt == nil {
		return true
	}
	return opt.UseThreads
}

func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}

// GetMetadata returns EXIF/ICCP/XMP metadata from a WebP payload.
func GetMetadata(data []byte, format string) (metadata []byte, err error) {
	return webpGetMetadata(data, strings.ToUpper(format))
}

// SetMetadata writes EXIF/ICCP/XMP metadata into a WebP payload.
func SetMetadata(data, metadata []byte, format string) (newData []byte, err error) {
	return webpSetMetadata(data, metadata, format)
}
