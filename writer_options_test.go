// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build cgo
// +build cgo

package webp

import (
	"bytes"
	"testing"
)

func TestEncodeTargetSizeOption(t *testing.T) {
	img, err := loadImage("video-001.png")
	if err != nil {
		t.Fatal(err)
	}

	buf := new(bytes.Buffer)
	err = Encode(buf, img, &Options{
		Quality:    100,
		TargetSize: 4096,
		Method:     DefaultMethod,
	})
	if err != nil {
		t.Fatal(err)
	}
	if buf.Len() == 0 {
		t.Fatal("expected encoded data")
	}
	if _, err := Decode(bytes.NewReader(buf.Bytes())); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
}

func TestEncodeTargetSizeInvalid(t *testing.T) {
	img, err := loadImage("video-001.png")
	if err != nil {
		t.Fatal(err)
	}

	buf := new(bytes.Buffer)
	err = Encode(buf, img, &Options{
		Quality:    75,
		TargetSize: -1,
		Method:     DefaultMethod,
	})
	if err == nil {
		t.Fatal("expected error for negative target size")
	}
}

func TestEncodeAlphaQualityInvalid(t *testing.T) {
	img, err := loadImage("video-001.png")
	if err != nil {
		t.Fatal(err)
	}

	buf := new(bytes.Buffer)
	err = Encode(buf, img, &Options{
		Quality:      75,
		AlphaQuality: 101,
		Method:       DefaultMethod,
	})
	if err == nil {
		t.Fatal("expected error for alpha quality out of range")
	}
}

func TestEncodeAutoFilterOption(t *testing.T) {
	img, err := loadImage("video-001.png")
	if err != nil {
		t.Fatal(err)
	}

	buf := new(bytes.Buffer)
	err = Encode(buf, img, &Options{
		Quality:    80,
		AutoFilter: true,
		Method:     DefaultMethod,
	})
	if err != nil {
		t.Fatal(err)
	}
	if buf.Len() == 0 {
		t.Fatal("expected encoded data")
	}
	if _, err := Decode(bytes.NewReader(buf.Bytes())); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
}

func TestEncodeTargetSizeOverridesQuality(t *testing.T) {
	img, err := loadImage("video-001.png")
	if err != nil {
		t.Fatal(err)
	}

	buf := new(bytes.Buffer)
	err = Encode(buf, img, &Options{
		Quality:    -1,
		TargetSize: 4096,
		Method:     DefaultMethod,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEncodeAlphaQualityDefault(t *testing.T) {
	img, err := loadImage("video-001.png")
	if err != nil {
		t.Fatal(err)
	}

	buf := new(bytes.Buffer)
	err = Encode(buf, img, &Options{
		Quality:      75,
		AlphaQuality: 0,
		Method:       DefaultMethod,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEncodeMethodOption(t *testing.T) {
	img, err := loadImage("video-001.png")
	if err != nil {
		t.Fatal(err)
	}

	buf := new(bytes.Buffer)
	err = Encode(buf, img, &Options{
		Quality: 80,
		Method:  0,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEncodeQualityZero(t *testing.T) {
	img, err := loadImage("video-001.png")
	if err != nil {
		t.Fatal(err)
	}

	buf := new(bytes.Buffer)
	err = Encode(buf, img, &Options{
		Quality: 0,
		Method:  DefaultMethod,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
