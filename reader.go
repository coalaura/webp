// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
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
)

// Config is the same as image.Config but includes alpha and animation metadata.
type Config struct {
	ColorModel   color.Model
	Width        int  // Width in pixels, as read from the bitstream.
	Height       int  // Height in pixels, as read from the bitstream.
	HasAlpha     bool // True if the bitstream contains an alpha channel.
	HasAnimation bool // True if the bitstream is an animation.
	Format       int  // 0 = undefined (/mixed), 1 = lossy, 2 = lossless
}

type DecodeOptions struct {
	UseThreads bool // Enable libwebp multi-threading (default true)
}

func LoadConfigEx(name string) (config Config, err error) {
	f, err := os.Open(name)
	if err != nil {
		return Config{}, err
	}
	defer f.Close()

	var header [maxWebpHeaderSize]byte
	n, err := f.Read(header[:])
	if err != nil && err != io.EOF {
		return
	}
	headerSlice := header[:n]
	width, height, hasAlpha, hasAnimation, format, err := GetInfo(headerSlice)
	if err != nil {
		return
	}

	config.Width = width
	config.Height = height
	config.ColorModel = color.RGBAModel
	config.HasAlpha = hasAlpha
	config.HasAnimation = hasAnimation
	config.Format = format
	return
}

func LoadConfig(name string) (config image.Config, err error) {
	f, err := os.Open(name)
	if err != nil {
		return image.Config{}, err
	}
	defer f.Close()

	var header [maxWebpHeaderSize]byte
	n, err := f.Read(header[:])
	if err != nil && err != io.EOF {
		return
	}
	headerSlice := header[:n]
	width, height, _, _, _, err := GetInfo(headerSlice)
	if err != nil {
		return
	}

	config.Width = width
	config.Height = height
	config.ColorModel = color.RGBAModel
	return
}

func Load(name string) (m image.Image, err error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if fi.Size() > (2 << 30) {
		return nil, errors.New("webp: Load, file size is too large (> 2GB)")
	}

	data := make([]byte, int(fi.Size()))
	if _, err = f.Read(data); err != nil {
		return nil, err
	}
	if m, err = DecodeRGBA(data); err != nil {
		return
	}
	return
}

// DecodeConfigEx returns the color model, dimensions and animation flag of a WEBP image
// without decoding the entire image.
func DecodeConfigEx(r io.Reader) (config Config, err error) {
	var header [maxWebpHeaderSize]byte
	n, err := r.Read(header[:])
	if err != nil && err != io.EOF {
		return
	}
	headerSlice := header[:n]
	width, height, hasAlpha, hasAnimation, format, err := GetInfo(headerSlice)
	if err != nil {
		return
	}
	config.Width = width
	config.Height = height
	config.ColorModel = color.RGBAModel
	config.HasAlpha = hasAlpha
	config.HasAnimation = hasAnimation
	config.Format = format
	return
}

// DecodeConfig returns the color model and dimensions of a WEBP image without
// decoding the entire image.
func DecodeConfig(r io.Reader) (config image.Config, err error) {
	var header [maxWebpHeaderSize]byte
	n, err := r.Read(header[:])
	if err != nil && err != io.EOF {
		return
	}
	headerSlice := header[:n]
	width, height, _, _, _, err := GetInfo(headerSlice)
	if err != nil {
		return
	}
	config.Width = width
	config.Height = height
	config.ColorModel = color.RGBAModel
	return
}

// Decode reads a WEBP image from r and returns it as an image.Image.
func Decode(r io.Reader) (m image.Image, err error) {
	return DecodeWithOptions(r, nil)
}

func DecodeWithOptions(r io.Reader, opt *DecodeOptions) (m image.Image, err error) {
	data, release, err := readAllPooled(r)
	if err != nil {
		return
	}
	if release != nil {
		defer release()
	}
	if m, err = DecodeRGBAWithOptions(data, opt); err != nil {
		return
	}
	return
}

func init() {
	image.RegisterFormat("webp", "RIFF????WEBPVP8", Decode, DecodeConfig)
}
