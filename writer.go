// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Copyright 2026 github.com/coalaura. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build cgo
// +build cgo

package webp

import (
	"errors"
	"image"
	"image/color"
	"io"
	"os"
	"reflect"
)

func writeAll(w io.Writer, data []byte) error {
	for len(data) > 0 {
		n, err := w.Write(data)
		if n < 0 || n > len(data) {
			return errors.New("webp: invalid writer count")
		}
		if n > 0 {
			data = data[n:]
		}
		if err != nil {
			return err
		}
		if n == 0 {
			return io.ErrShortWrite
		}
	}
	return nil
}

// DefaultQuality is the default lossy quality used by Encode.
const DefaultQuality = 90

// DefaultMethod is the default encoding method used by Encode.
const DefaultMethod = 4

// DefaultAlphaQuality is the default alpha plane quality used by Encode.
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
	UseThreads   bool    // Enable libwebp multi-threading (default true)
	// LosslessLevel selects libwebp's lossless preset (0 fastest through 9
	// smallest). Nil preserves the legacy quality-100/method configuration.
	LosslessLevel *int
}

type colorModeler interface {
	ColorModel() color.Model
}

// Save encodes m as WebP and writes it to the named file.
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
	if opt != nil && opt.Lossless && opt.LosslessLevel != nil {
		return EncodeTo(w, m, opt)
	}
	return encode(w, m, opt)
}

func encodeSettings(opt *Options) (method, targetSize, alphaQuality int, autoFilter, useThreads bool, quality float32, lossless, exact bool, losslessLevel int, err error) {
	method = DefaultMethod
	alphaQuality = DefaultAlphaQuality
	useThreads = true
	quality = DefaultQuality
	losslessLevel = -1
	if opt == nil {
		return
	}
	method = opt.Method
	targetSize = opt.TargetSize
	alphaQuality = opt.AlphaQuality
	if alphaQuality == 0 {
		alphaQuality = DefaultAlphaQuality
	}
	autoFilter = opt.AutoFilter
	useThreads = opt.UseThreads
	lossless = opt.Lossless
	exact = opt.Exact
	if targetSize == 0 {
		quality = opt.Quality
	}
	if opt.LosslessLevel != nil {
		losslessLevel = *opt.LosslessLevel
		if losslessLevel < 0 || losslessLevel > 9 {
			err = errors.New("webp: invalid lossless level")
			return
		}
	}
	if !validCInts(method, targetSize, alphaQuality) {
		err = errors.New("webp: invalid encoder options")
	}
	return
}

func encode(w io.Writer, m image.Image, opt *Options) (err error) {
	var output []byte
	method := DefaultMethod
	targetSize := 0
	alphaQuality := DefaultAlphaQuality
	autoFilter := false
	useThreads := true
	if opt != nil {
		method = opt.Method
		targetSize = opt.TargetSize
		alphaQuality = opt.AlphaQuality
		if alphaQuality == 0 {
			alphaQuality = DefaultAlphaQuality
		}
		autoFilter = opt.AutoFilter
		useThreads = opt.UseThreads
	}
	if opt != nil && opt.Lossless {
		switch m := adjustImage(m).(type) {
		case *image.Gray:
			if output, err = encodeLosslessGray(m, method, targetSize, alphaQuality, autoFilter, useThreads); err != nil {
				return
			}
		case *RGBImage:
			if output, err = encodeLosslessRGB(m, method, targetSize, alphaQuality, autoFilter, useThreads); err != nil {
				return
			}
		case *image.RGBA:
			if opt.Exact {
				output, err = encodeExactLosslessRGBA(m, method, targetSize, alphaQuality, autoFilter, useThreads)
			} else {
				output, err = encodeLosslessRGBA(m, method, targetSize, alphaQuality, autoFilter, useThreads)
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
			if output, err = encodeGray(m, quality, method, targetSize, alphaQuality, autoFilter, useThreads); err != nil {
				return
			}
		case *RGBImage:
			if output, err = encodeRGB(m, quality, method, targetSize, alphaQuality, autoFilter, useThreads); err != nil {
				return
			}
		case *image.RGBA:
			if output, err = encodeRGBA(m, quality, method, targetSize, alphaQuality, autoFilter, useThreads); err != nil {
				return
			}
		default:
			panic("image/webp: Encode, unreachable!")
		}
	}
	return writeAll(w, output)
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
	switch m := m.(type) {
	case *image.NRGBA:
		for y := b.Min.Y; y < b.Max.Y; y++ {
			src, dst := m.PixOffset(b.Min.X, y), rgba.PixOffset(b.Min.X, y)
			for x := 0; x < b.Dx(); x++ {
				r, g, bl, a := m.Pix[src], m.Pix[src+1], m.Pix[src+2], m.Pix[src+3]
				if a == 0xff {
					rgba.Pix[dst], rgba.Pix[dst+1], rgba.Pix[dst+2], rgba.Pix[dst+3] = r, g, bl, a
				} else if a != 0 {
					rgba.Pix[dst], rgba.Pix[dst+1], rgba.Pix[dst+2], rgba.Pix[dst+3] = mulAlpha8(r, a), mulAlpha8(g, a), mulAlpha8(bl, a), a
				} else {
					rgba.Pix[dst], rgba.Pix[dst+1], rgba.Pix[dst+2], rgba.Pix[dst+3] = 0, 0, 0, 0
				}
				src, dst = src+4, dst+4
			}
		}
		return rgba
	case *image.YCbCr:
		for y := b.Min.Y; y < b.Max.Y; y++ {
			for x := b.Min.X; x < b.Max.X; x++ {
				si, di := m.YOffset(x, y), rgba.PixOffset(x, y)
				r, g, bl := color.YCbCrToRGB(m.Y[si], m.Cb[m.COffset(x, y)], m.Cr[m.COffset(x, y)])
				rgba.Pix[di], rgba.Pix[di+1], rgba.Pix[di+2], rgba.Pix[di+3] = r, g, bl, 0xff
			}
		}
		return rgba
	case *image.Gray16:
		for y := b.Min.Y; y < b.Max.Y; y++ {
			src, dst := m.PixOffset(b.Min.X, y), rgba.PixOffset(b.Min.X, y)
			for x := 0; x < b.Dx(); x++ {
				v := m.Pix[src]
				rgba.Pix[dst], rgba.Pix[dst+1], rgba.Pix[dst+2], rgba.Pix[dst+3] = v, v, v, 0xff
				src, dst = src+2, dst+4
			}
		}
		return rgba
	case *image.RGBA64:
		for y := b.Min.Y; y < b.Max.Y; y++ {
			src, dst := m.PixOffset(b.Min.X, y), rgba.PixOffset(b.Min.X, y)
			for x := 0; x < b.Dx(); x++ {
				rgba.Pix[dst], rgba.Pix[dst+1], rgba.Pix[dst+2], rgba.Pix[dst+3] = m.Pix[src], m.Pix[src+2], m.Pix[src+4], m.Pix[src+6]
				src, dst = src+8, dst+4
			}
		}
		return rgba
	case *image.NRGBA64:
		for y := b.Min.Y; y < b.Max.Y; y++ {
			src, dst := m.PixOffset(b.Min.X, y), rgba.PixOffset(b.Min.X, y)
			for x := 0; x < b.Dx(); x++ {
				r, g, bl, a := uint16(m.Pix[src])<<8|uint16(m.Pix[src+1]), uint16(m.Pix[src+2])<<8|uint16(m.Pix[src+3]), uint16(m.Pix[src+4])<<8|uint16(m.Pix[src+5]), uint16(m.Pix[src+6])<<8|uint16(m.Pix[src+7])
				rgba.Pix[dst], rgba.Pix[dst+1], rgba.Pix[dst+2], rgba.Pix[dst+3] = mulAlpha16To8(r, a), mulAlpha16To8(g, a), mulAlpha16To8(bl, a), uint8(a>>8)
				src, dst = src+8, dst+4
			}
		}
		return rgba
	case *image.Paletted:
		palette := make([]color.RGBA, len(m.Palette))
		for i, c := range m.Palette {
			palette[i] = color.RGBAModel.Convert(c).(color.RGBA)
		}
		for y := b.Min.Y; y < b.Max.Y; y++ {
			src, dst := m.PixOffset(b.Min.X, y), rgba.PixOffset(b.Min.X, y)
			for x := 0; x < b.Dx(); x++ {
				c := palette[m.Pix[src]]
				rgba.Pix[dst], rgba.Pix[dst+1], rgba.Pix[dst+2], rgba.Pix[dst+3] = c.R, c.G, c.B, c.A
				src, dst = src+1, dst+4
			}
		}
		return rgba
	}
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

func mulAlpha8(v, a uint8) uint8 {
	// Match image.NRGBA.RGBA followed by image.RGBA.Set's 8-bit conversion.
	return uint8((uint32(v) * 0x101 * uint32(a) / 0xff) >> 8)
}

func mulAlpha16To8(v, a uint16) uint8 {
	return uint8((uint32(v) * uint32(a) / 0xffff) >> 8)
}
