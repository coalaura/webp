// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webp

import (
	"bytes"
	"io"
	"sync"
)

const (
	readAllPoolInitSize = 32 * 1024
	readAllPoolMaxSize  = 1 << 20
)

var readAllPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 0, readAllPoolInitSize)
	},
}

func readAllPooled(r io.Reader) ([]byte, func(), error) {
	buf := readAllPool.Get().([]byte)
	b := bytes.NewBuffer(buf[:0])
	_, err := b.ReadFrom(r)
	data := b.Bytes()
	if err != nil {
		releaseReadAllBuffer(data)
		return nil, nil, err
	}
	release := func() {
		releaseReadAllBuffer(data)
	}
	return data, release, nil
}

func releaseReadAllBuffer(buf []byte) {
	if cap(buf) > readAllPoolMaxSize {
		return
	}
	readAllPool.Put(buf[:0])
}
