// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Copyright 2026 github.com/coalaura. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webp

import (
	"image"
	"image/color"
	"reflect"
)

var (
	_ image.Image = (*RGBImage)(nil)
	_ MemP        = (*RGBImage)(nil)
)

// RGBImage stores an RGB image in packed 8-bit format.
type RGBImage struct {
	XPix    []uint8
	XStride int
	XRect   image.Rectangle
}

// MemPMagic returns the MemP magic string.
func (p *RGBImage) MemPMagic() string {
	return MemPMagic
}

// Bounds returns the image bounds.
func (p *RGBImage) Bounds() image.Rectangle {
	return p.XRect
}

// Channels returns the number of channels.
func (p *RGBImage) Channels() int {
	return 3
}

// DataType returns the underlying element kind.
func (p *RGBImage) DataType() reflect.Kind {
	return reflect.Uint8
}

// Pix returns the raw pixel buffer.
func (p *RGBImage) Pix() []byte {
	return p.XPix
}

// Stride returns the byte stride between rows.
func (p *RGBImage) Stride() int {
	return p.XStride
}

// ColorModel returns the image's color model.
func (p *RGBImage) ColorModel() color.Model { return color.RGBAModel }

// At returns the color of the pixel at (x, y).
func (p *RGBImage) At(x, y int) color.Color {
	if !(image.Point{x, y}.In(p.XRect)) {
		return color.RGBA{}
	}
	i := p.PixOffset(x, y)
	return color.RGBA{
		R: p.XPix[i+0],
		G: p.XPix[i+1],
		B: p.XPix[i+2],
		A: 0xff,
	}
}

// RGBAt returns the RGB values at (x, y).
func (p *RGBImage) RGBAt(x, y int) [3]uint8 {
	if !(image.Point{x, y}.In(p.XRect)) {
		return [3]uint8{}
	}
	i := p.PixOffset(x, y)
	return [3]uint8{
		p.XPix[i+0],
		p.XPix[i+1],
		p.XPix[i+2],
	}
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *RGBImage) PixOffset(x, y int) int {
	return (y-p.XRect.Min.Y)*p.XStride + (x-p.XRect.Min.X)*3
}

// Set sets the pixel at (x, y) to c.
func (p *RGBImage) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.XRect)) {
		return
	}
	i := p.PixOffset(x, y)
	c1 := color.RGBAModel.Convert(c).(color.RGBA)
	p.XPix[i+0] = c1.R
	p.XPix[i+1] = c1.G
	p.XPix[i+2] = c1.B
}

// SetRGB sets the pixel at (x, y) to an RGB triplet.
func (p *RGBImage) SetRGB(x, y int, c [3]uint8) {
	if !(image.Point{x, y}.In(p.XRect)) {
		return
	}
	i := p.PixOffset(x, y)
	p.XPix[i+0] = c[0]
	p.XPix[i+1] = c[1]
	p.XPix[i+2] = c[2]
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *RGBImage) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.XRect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to be inside
	// either r1 or r2 if the intersection is empty. Without explicitly checking for
	// this, the Pix[i:] expression below can panic.
	if r.Empty() {
		return &RGBImage{}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &RGBImage{
		XPix:    p.XPix[i:],
		XStride: p.XStride,
		XRect:   r,
	}
}

// Opaque scans the entire image and reports whether it is fully opaque.
func (p *RGBImage) Opaque() bool {
	return true
}

// NewRGBImage returns a new RGBImage with the given bounds.
func NewRGBImage(r image.Rectangle) *RGBImage {
	w, h := r.Dx(), r.Dy()
	pix := make([]uint8, 3*w*h)
	return &RGBImage{
		XPix:    pix,
		XStride: 3 * w,
		XRect:   r,
	}
}

// NewRGBImageFrom converts m into an RGBImage.
func NewRGBImageFrom(m image.Image) *RGBImage {
	if m, ok := m.(*RGBImage); ok {
		return m
	}

	// convert to RGBImage
	b := m.Bounds()
	rgb := NewRGBImage(b)
	switch m := m.(type) {
	case *image.NRGBA:
		for y := b.Min.Y; y < b.Max.Y; y++ {
			src, dst := m.PixOffset(b.Min.X, y), rgb.PixOffset(b.Min.X, y)
			for x := 0; x < b.Dx(); x++ {
				a := m.Pix[src+3]
				rgb.XPix[dst], rgb.XPix[dst+1], rgb.XPix[dst+2] = mulAlpha8(m.Pix[src], a), mulAlpha8(m.Pix[src+1], a), mulAlpha8(m.Pix[src+2], a)
				src, dst = src+4, dst+3
			}
		}
		return rgb
	case *image.YCbCr:
		for y := b.Min.Y; y < b.Max.Y; y++ {
			for x := b.Min.X; x < b.Max.X; x++ {
				si, di := m.YOffset(x, y), rgb.PixOffset(x, y)
				rgb.XPix[di], rgb.XPix[di+1], rgb.XPix[di+2] = color.YCbCrToRGB(m.Y[si], m.Cb[m.COffset(x, y)], m.Cr[m.COffset(x, y)])
			}
		}
		return rgb
	case *image.Gray16:
		for y := b.Min.Y; y < b.Max.Y; y++ {
			src, dst := m.PixOffset(b.Min.X, y), rgb.PixOffset(b.Min.X, y)
			for x := 0; x < b.Dx(); x++ {
				v := m.Pix[src]
				rgb.XPix[dst], rgb.XPix[dst+1], rgb.XPix[dst+2] = v, v, v
				src, dst = src+2, dst+3
			}
		}
		return rgb
	case *image.RGBA64:
		for y := b.Min.Y; y < b.Max.Y; y++ {
			src, dst := m.PixOffset(b.Min.X, y), rgb.PixOffset(b.Min.X, y)
			for x := 0; x < b.Dx(); x++ {
				rgb.XPix[dst], rgb.XPix[dst+1], rgb.XPix[dst+2] = m.Pix[src], m.Pix[src+2], m.Pix[src+4]
				src, dst = src+8, dst+3
			}
		}
		return rgb
	case *image.NRGBA64:
		for y := b.Min.Y; y < b.Max.Y; y++ {
			src, dst := m.PixOffset(b.Min.X, y), rgb.PixOffset(b.Min.X, y)
			for x := 0; x < b.Dx(); x++ {
				r, g, bl, a := uint16(m.Pix[src])<<8|uint16(m.Pix[src+1]), uint16(m.Pix[src+2])<<8|uint16(m.Pix[src+3]), uint16(m.Pix[src+4])<<8|uint16(m.Pix[src+5]), uint16(m.Pix[src+6])<<8|uint16(m.Pix[src+7])
				rgb.XPix[dst], rgb.XPix[dst+1], rgb.XPix[dst+2] = mulAlpha16To8(r, a), mulAlpha16To8(g, a), mulAlpha16To8(bl, a)
				src, dst = src+8, dst+3
			}
		}
		return rgb
	case *image.Paletted:
		palette := make([]color.RGBA, len(m.Palette))
		for i, c := range m.Palette {
			palette[i] = color.RGBAModel.Convert(c).(color.RGBA)
		}
		for y := b.Min.Y; y < b.Max.Y; y++ {
			src, dst := m.PixOffset(b.Min.X, y), rgb.PixOffset(b.Min.X, y)
			for x := 0; x < b.Dx(); x++ {
				c := palette[m.Pix[src]]
				rgb.XPix[dst], rgb.XPix[dst+1], rgb.XPix[dst+2] = c.R, c.G, c.B
				src, dst = src+1, dst+3
			}
		}
		return rgb
	}
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			pr, pg, pb, _ := m.At(x, y).RGBA()
			rgb.SetRGB(x, y, [3]uint8{
				uint8(pr >> 8),
				uint8(pg >> 8),
				uint8(pb >> 8),
			})
		}
	}
	return rgb
}
