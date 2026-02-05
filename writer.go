// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build cgo
// +build cgo

package webp

import (
	"image"
	"image/color"
	"io"
	"os"
	"reflect"
)

const DefaultQuality = 90
const DefaultMethod = 4
const DefaultAlphaQuality = 100

// Options are the encoding parameters.
type Options struct {
	Lossless     bool
	Quality      float32 // 0 ~ 100
	TargetSize   int     // Target byte size (overrides Quality)
	AlphaQuality int     // Alpha plane quality 0-100 (default 100)
	AutoFilter   bool    // Auto-adjust filter for each image
	Exact        bool    // Preserve RGB values in transparent area.
	Method       int     // Quality/speed trade-off (0=fast, 6=slower-better)
}

type colorModeler interface {
	ColorModel() color.Model
}

func Save(name string, m image.Image, opt *Options) (err error) {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()

	return encode(f, m, opt)
}

// Encode writes the image m to w in WEBP format.
func Encode(w io.Writer, m image.Image, opt *Options) (err error) {
	return encode(w, m, opt)
}

func encode(w io.Writer, m image.Image, opt *Options) (err error) {
	var output []byte
	method := DefaultMethod
	targetSize := 0
	alphaQuality := DefaultAlphaQuality
	autoFilter := false
	if opt != nil {
		method = opt.Method
		targetSize = opt.TargetSize
		alphaQuality = opt.AlphaQuality
		if alphaQuality == 0 {
			alphaQuality = DefaultAlphaQuality
		}
		autoFilter = opt.AutoFilter
	}
	if opt != nil && opt.Lossless {
		switch m := adjustImage(m).(type) {
		case *image.Gray:
			if output, err = encodeLosslessGray(m, method, targetSize, alphaQuality, autoFilter); err != nil {
				return
			}
		case *RGBImage:
			if output, err = encodeLosslessRGB(m, method, targetSize, alphaQuality, autoFilter); err != nil {
				return
			}
		case *image.RGBA:
			if opt.Exact {
				output, err = encodeExactLosslessRGBA(m, method, targetSize, alphaQuality, autoFilter)
			} else {
				output, err = encodeLosslessRGBA(m, method, targetSize, alphaQuality, autoFilter)
			}
			if err != nil {
				return
			}
		default:
			panic("image/webp: Encode, unreachable!")
		}
	} else {
		quality := float32(DefaultQuality)
		if opt != nil {
			if targetSize == 0 {
				quality = opt.Quality
			}
		}

		switch m := adjustImage(m).(type) {
		case *image.Gray:
			if output, err = encodeGray(m, quality, method, targetSize, alphaQuality, autoFilter); err != nil {
				return
			}
		case *RGBImage:
			if output, err = encodeRGB(m, quality, method, targetSize, alphaQuality, autoFilter); err != nil {
				return
			}
		case *image.RGBA:
			if output, err = encodeRGBA(m, quality, method, targetSize, alphaQuality, autoFilter); err != nil {
				return
			}
		default:
			panic("image/webp: Encode, unreachable!")
		}
	}
	_, err = w.Write(output)
	return
}

func adjustImage(m image.Image) image.Image {
	if p, ok := AsMemPImage(m); ok {
		switch {
		case p.XChannels == 1 && p.XDataType == reflect.Uint8:
			m = &image.Gray{
				Pix:    p.XPix,
				Stride: p.XStride,
				Rect:   p.XRect,
			}
		case p.XChannels == 1 && p.XDataType == reflect.Uint16:
			m = toGrayImage(m) // MemP is little endian
		case p.XChannels == 3 && p.XDataType == reflect.Uint8:
			m = &RGBImage{
				XPix:    p.XPix,
				XStride: p.XStride,
				XRect:   p.XRect,
			}
		case p.XChannels == 3 && p.XDataType == reflect.Uint16:
			m = NewRGBImageFrom(m) // MemP is little endian
		case p.XChannels == 4 && p.XDataType == reflect.Uint8:
			m = &image.RGBA{
				Pix:    p.XPix,
				Stride: p.XStride,
				Rect:   p.XRect,
			}
		case p.XChannels == 4 && p.XDataType == reflect.Uint16:
			m = toRGBAImage(m) // MemP is little endian
		}
	}
	switch m := m.(type) {
	case *image.Gray:
		return m
	case *RGBImage:
		return m
	case *RGB48Image:
		return NewRGBImageFrom(m)
	case *image.RGBA:
		return m
	case *image.YCbCr:
		return NewRGBImageFrom(m)

	case *image.Gray16:
		return toGrayImage(m)
	case *image.RGBA64:
		return toRGBAImage(m)
	case *image.NRGBA:
		return toRGBAImage(m)
	case *image.NRGBA64:
		return toRGBAImage(m)

	default:
		return toRGBAImage(m)
	}
}

func toGrayImage(m image.Image) *image.Gray {
	if m, ok := m.(*image.Gray); ok {
		return m
	}
	b := m.Bounds()
	gray := image.NewGray(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := color.GrayModel.Convert(m.At(x, y)).(color.Gray)
			gray.SetGray(x, y, c)
		}
	}
	return gray
}

func toRGBAImage(m image.Image) *image.RGBA {
	if m, ok := m.(*image.RGBA); ok {
		return m
	}
	b := m.Bounds()
	rgba := image.NewRGBA(b)
	dstColorRGBA64 := &color.RGBA64{}
	dstColor := color.Color(dstColorRGBA64)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			pr, pg, pb, pa := m.At(x, y).RGBA()
			dstColorRGBA64.R = uint16(pr)
			dstColorRGBA64.G = uint16(pg)
			dstColorRGBA64.B = uint16(pb)
			dstColorRGBA64.A = uint16(pa)
			rgba.Set(x, y, dstColor)
		}
	}
	return rgba
}
