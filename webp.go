// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
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

func GetInfo(data []byte) (width, height int, hasAlpha bool, hasAnimation bool, format int, err error) {
	return webpGetInfo(data)
}

func DecodeGray(data []byte) (m *image.Gray, err error) {
	return DecodeGrayWithOptions(data, nil)
}

func DecodeGrayWithOptions(data []byte, opt *DecodeOptions) (m *image.Gray, err error) {
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

func DecodeRGB(data []byte) (m *RGBImage, err error) {
	return DecodeRGBWithOptions(data, nil)
}

func DecodeRGBWithOptions(data []byte, opt *DecodeOptions) (m *RGBImage, err error) {
	useThreads := false
	if opt != nil {
		useThreads = opt.UseThreads
	}

	pix, w, h, err := webpDecodeRGB(data, useThreads)
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

func DecodeRGBA(data []byte) (m *image.RGBA, err error) {
	return DecodeRGBAWithOptions(data, nil)
}

func DecodeRGBAWithOptions(data []byte, opt *DecodeOptions) (m *image.RGBA, err error) {
	useThreads := false
	if opt != nil {
		useThreads = opt.UseThreads
	}

	pix, w, h, err := webpDecodeRGBA(data, useThreads)
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
func DecodeGrayToSize(data []byte, width, height int) (m *image.Gray, err error) {
	return DecodeGrayToSizeWithOptions(data, width, height, nil)
}

func DecodeGrayToSizeWithOptions(data []byte, width, height int, opt *DecodeOptions) (m *image.Gray, err error) {
	pix, err := webpDecodeGrayToSize(data, width, height, opt.UseThreads)
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
func DecodeRGBToSize(data []byte, width, height int) (m *RGBImage, err error) {
	return DecodeRGBToSizeWithOptions(data, width, height, nil)
}

func DecodeRGBToSizeWithOptions(data []byte, width, height int, opt *DecodeOptions) (m *RGBImage, err error) {
	pix, err := webpDecodeRGBToSize(data, width, height, opt.UseThreads)
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

// DecodeRGBAToSize decodes a Gray image scaled to the given dimensions.
func DecodeRGBAToSize(data []byte, width, height int) (m *image.RGBA, err error) {
	return DecodeRGBAToSizeWithOptions(data, width, height, nil)
}

func DecodeRGBAToSizeWithOptions(data []byte, width, height int, opt *DecodeOptions) (m *image.RGBA, err error) {
	pix, err := webpDecodeRGBAToSize(data, width, height, opt.UseThreads)
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

func EncodeGray(m image.Image, quality float32, method int) (data []byte, err error) {
	return encodeGray(m, quality, method, 0, DefaultAlphaQuality, false, false)
}

func EncodeRGB(m image.Image, quality float32, method int) (data []byte, err error) {
	return encodeRGB(m, quality, method, 0, DefaultAlphaQuality, false, false)
}

func EncodeRGBA(m image.Image, quality float32, method int) (data []byte, err error) {
	return encodeRGBA(m, quality, method, 0, DefaultAlphaQuality, false, false)
}

func EncodeLosslessGray(m image.Image, method int) (data []byte, err error) {
	return encodeLosslessGray(m, method, 0, DefaultAlphaQuality, false, false)
}

func EncodeLosslessRGB(m image.Image, method int) (data []byte, err error) {
	return encodeLosslessRGB(m, method, 0, DefaultAlphaQuality, false, false)
}

func EncodeLosslessRGBA(m image.Image, method int) (data []byte, err error) {
	return encodeLosslessRGBA(m, method, 0, DefaultAlphaQuality, false, false)
}

// EncodeExactLosslessRGBA Encode lossless RGB mode with exact.
// exact: preserve RGB values in transparent area.
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

func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}

// GetMetadata return EXIF/ICCP/XMP format metadata.
func GetMetadata(data []byte, format string) (metadata []byte, err error) {
	return webpGetMetadata(data, strings.ToUpper(format))
}

// SetMetadata set EXIF/ICCP/XMP format metadata.
func SetMetadata(data, metadata []byte, format string) (newData []byte, err error) {
	return webpSetMetadata(data, metadata, format)
}
