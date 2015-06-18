// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rawp

import (
	"hash/crc32"
	"image"
	"io"
	"unsafe"
)

// Options are the encoding parameters.
type Options struct {
	UseSnappy bool
}

// Encode writes the image m to w in RawP format.
func Encode(w io.Writer, m image.Image, opt *Options) (err error) {
	p := NewImageFrom(m)

	var useSnappy bool
	if opt != nil {
		useSnappy = opt.UseSnappy
	}

	hdr, err := rawpMakeHeader(p.Bounds().Dx(), p.Bounds().Dy(), p.Channels, p.DataType, useSnappy)
	if err != nil {
		return
	}

	stride := p.Rect.Dx() * p.Channels * p.DataType.ByteSize()
	pix := make([]byte, stride*p.Rect.Dy())

	off := 0
	for y := p.Rect.Min.Y; y < p.Rect.Max.Y; y++ {
		copy(pix[off:][:stride], p.Pix[p.PixOffset(0, y):])
		off += stride
	}

	if useSnappy {
		pix, err = snappyEncode(nil, pix)
		if err != nil {
			return
		}
	}

	hdr.DataSize = uint32(len(pix))
	hdr.DataCheckSum = crc32.ChecksumIEEE(pix)
	hdr.Data = pix

	if _, err = w.Write(((*[1 << 30]byte)(unsafe.Pointer(hdr)))[:rawpHeaderSize]); err != nil {
		return
	}
	if _, err = w.Write(hdr.Data); err != nil {
		return
	}
	return
}
