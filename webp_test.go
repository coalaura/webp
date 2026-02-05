// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webp

import (
	"os"
	"testing"
)

type tGetInfoTester struct {
	Filename     string
	HdrSize      int
	Width        int
	Height       int
	HasAlpha     bool
	HasAnimation bool
	Format       int
}

func TestGetInfo(t *testing.T) {
	for i, v := range tGetInfoTesterList {
		data, err := os.ReadFile(testdataDir + v.Filename)
		if err != nil {
			t.Fatalf("%d: %v", i, err)
		}
		width, height, hasAlpha, hasAnimation, format, err := GetInfo(data)
		if err != nil {
			t.Fatalf("%d: %v", i, err)
		}
		if width != v.Width {
			t.Fatalf("%d: expect = %v, got = %v", i, v.Width, width)
		}
		if height != v.Height {
			t.Fatalf("%d: expect = %v, got = %v", i, v.Height, height)
		}
		if hasAlpha != v.HasAlpha {
			t.Fatalf("%d: expect = %v, got = %v", i, v.HasAlpha, hasAlpha)
		}
		if hasAnimation != v.HasAnimation {
			t.Fatalf("%d: expect = %v, got = %v", i, v.HasAnimation, hasAnimation)
		}
		if format != v.Format {
			t.Fatalf("%d: expect = %v, got = %v", i, v.Format, format)
		}
	}
}

var tGetInfoTesterList = []tGetInfoTester{
	{
		Filename:     "1_webp_ll.webp",
		Width:        400,
		Height:       301,
		HasAlpha:     true,
		HasAnimation: false,
		Format:       2,
	},
	{
		Filename:     "animated_1.webp",
		Width:        400,
		Height:       400,
		HasAlpha:     true,
		HasAnimation: true,
		Format:       0,
	},
}
