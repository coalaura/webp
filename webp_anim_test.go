// Copyright 2026 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webp

import (
	"bytes"
	"os"
	"testing"
)

func TestDecodeAll(t *testing.T) {
	files := []string{
		"./testdata/animated_1.webp",
		"./testdata/animated_2.webp",
	}
	for _, name := range files {
		data, err := os.ReadFile(name)
		tAssertNil(t, err)

		anim, err := DecodeAll(bytes.NewReader(data), nil)
		tAssertNil(t, err)
		tAssert(t, anim != nil)
		tAssert(t, len(anim.Image) > 1, name)
		tAssertEQ(t, len(anim.Image), len(anim.Delay), name)
		tAssert(t, anim.LoopCount >= 0, name)

		first := anim.Image[0]
		b := first.Bounds()
		tAssert(t, b.Dx() > 0 && b.Dy() > 0, name)
		for i, img := range anim.Image {
			ib := img.Bounds()
			tAssertEQ(t, b.Dx(), ib.Dx(), name)
			tAssertEQ(t, b.Dy(), ib.Dy(), name)
			tAssert(t, anim.Delay[i] >= 0, name)
		}
	}
}

func TestEncodeAll(t *testing.T) {
	data, err := os.ReadFile("./testdata/animated_1.webp")
	tAssertNil(t, err)

	anim, err := DecodeAll(bytes.NewReader(data), nil)
	tAssertNil(t, err)

	var buf bytes.Buffer
	err = EncodeAll(&buf, anim, &Options{Lossless: true})
	tAssertNil(t, err)

	encoded, err := DecodeAll(bytes.NewReader(buf.Bytes()), nil)
	tAssertNil(t, err)
	if encoded == nil {
		t.Fatalf("DecodeAll returned nil animation")
	}

	tAssertEQ(t, len(anim.Image), len(encoded.Image))
	tAssertEQ(t, len(anim.Delay), len(encoded.Delay))
	tAssertEQ(t, anim.LoopCount, encoded.LoopCount)

	base := anim.Image[0].Bounds()
	for i, img := range encoded.Image {
		b := img.Bounds()
		tAssertEQ(t, base.Dx(), b.Dx())
		tAssertEQ(t, base.Dy(), b.Dy())
		tAssertEQ(t, anim.Delay[i], encoded.Delay[i])
	}
}

func TestDecodeConfigAnimated(t *testing.T) {
	files := []string{
		"./testdata/animated_1.webp",
		"./testdata/animated_2.webp",
	}
	for _, name := range files {
		f, err := os.Open(name)
		tAssertNil(t, err)
		config, err := DecodeConfigEx(f)
		f.Close()
		tAssertNil(t, err)
		tAssert(t, config.Width > 0 && config.Height > 0, name)
		tAssert(t, config.HasAnimation, name)
	}
}

func TestDecodeConfigStill(t *testing.T) {
	f, err := os.Open("./testdata/1_webp_ll.webp")
	tAssertNil(t, err)
	config, err := DecodeConfigEx(f)
	f.Close()
	tAssertNil(t, err)
	tAssert(t, config.Width > 0 && config.Height > 0)
	tAssert(t, !config.HasAnimation)
}
