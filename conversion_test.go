package webp

import (
	"bytes"
	"image"
	"image/color"
	"testing"
)

type imageWrapper struct{ image.Image }

func TestNRGBAConversionsPreservePremultipliedColor(t *testing.T) {
	src := image.NewNRGBA(image.Rect(0, 0, 3, 1))
	src.SetNRGBA(0, 0, color.NRGBA{R: 255, G: 101, B: 1, A: 127})
	src.SetNRGBA(1, 0, color.NRGBA{R: 3, G: 199, B: 255, A: 0})
	src.SetNRGBA(2, 0, color.NRGBA{R: 17, G: 91, B: 240, A: 255})

	rgb := NewRGBImageFrom(src)
	rgba := toRGBAImage(src)
	for x := 0; x < 3; x++ {
		r, g, b, a := src.At(x, 0).RGBA()
		rgbOffset := rgb.PixOffset(x, 0)
		rgbaOffset := rgba.PixOffset(x, 0)
		if got, want := rgb.XPix[rgbOffset:rgbOffset+3], []byte{byte(r >> 8), byte(g >> 8), byte(b >> 8)}; !bytes.Equal(got, want) {
			t.Fatalf("RGB pixel %d = %v, want %v", x, got, want)
		}
		if got, want := rgba.Pix[rgbaOffset:rgbaOffset+4], []byte{byte(r >> 8), byte(g >> 8), byte(b >> 8), byte(a >> 8)}; !bytes.Equal(got, want) {
			t.Fatalf("RGBA pixel %d = %v, want %v", x, got, want)
		}
	}
}

func TestMulAlpha8MatchesNRGBAConversion(t *testing.T) {
	for v := 0; v <= 255; v++ {
		for a := 0; a <= 255; a++ {
			got := mulAlpha8(uint8(v), uint8(a))
			want16, _, _, _ := color.NRGBA{R: uint8(v), A: uint8(a)}.RGBA()
			want := uint8(want16 >> 8)
			if got != want {
				t.Fatalf("mulAlpha8(%d, %d) = %d, want %d", v, a, got, want)
			}
		}
	}
}

func TestFastImageConversionsMatchGeneric(t *testing.T) {
	bounds := image.Rect(3, 5, 5, 6)

	nrgba := image.NewNRGBA(bounds)
	nrgba.SetNRGBA(3, 5, color.NRGBA{R: 250, G: 128, B: 1, A: 127})
	nrgba.SetNRGBA(4, 5, color.NRGBA{R: 3, G: 199, B: 255, A: 0})

	ycbcr := image.NewYCbCr(bounds, image.YCbCrSubsampleRatio444)
	ycbcr.Y[ycbcr.YOffset(3, 5)] = 32
	ycbcr.Cb[ycbcr.COffset(3, 5)] = 224
	ycbcr.Cr[ycbcr.COffset(3, 5)] = 96
	ycbcr.Y[ycbcr.YOffset(4, 5)] = 200
	ycbcr.Cb[ycbcr.COffset(4, 5)] = 64
	ycbcr.Cr[ycbcr.COffset(4, 5)] = 192

	gray16 := image.NewGray16(bounds)
	gray16.SetGray16(3, 5, color.Gray16{Y: 0x12ff})
	gray16.SetGray16(4, 5, color.Gray16{Y: 0xef01})

	rgba64 := image.NewRGBA64(bounds)
	rgba64.SetRGBA64(3, 5, color.RGBA64{R: 0x1234, G: 0x5678, B: 0x9abc, A: 0xcdef})
	rgba64.SetRGBA64(4, 5, color.RGBA64{R: 0, G: 0, B: 0, A: 0})

	nrgba64 := image.NewNRGBA64(bounds)
	nrgba64.SetNRGBA64(3, 5, color.NRGBA64{R: 0xfffe, G: 0x1234, B: 0x5678, A: 0x8001})
	nrgba64.SetNRGBA64(4, 5, color.NRGBA64{R: 0x1234, G: 0x5678, B: 0x9abc, A: 0})

	paletted := image.NewPaletted(bounds, color.Palette{
		color.NRGBA{R: 250, G: 128, B: 1, A: 127},
		color.RGBA{R: 3, G: 199, B: 255, A: 255},
	})
	paletted.SetColorIndex(3, 5, 0)
	paletted.SetColorIndex(4, 5, 1)

	for _, src := range []image.Image{nrgba, ycbcr, gray16, rgba64, nrgba64, paletted} {
		t.Run("RGB", func(t *testing.T) {
			got, want := NewRGBImageFrom(src), NewRGBImageFrom(imageWrapper{src})
			if !bytes.Equal(got.XPix, want.XPix) {
				t.Fatalf("RGB conversion = %v, want %v", got.XPix, want.XPix)
			}
		})
		t.Run("RGBA", func(t *testing.T) {
			got, want := toRGBAImage(src), toRGBAImage(imageWrapper{src})
			if !bytes.Equal(got.Pix, want.Pix) {
				t.Fatalf("RGBA conversion = %v, want %v", got.Pix, want.Pix)
			}
		})
	}
}
