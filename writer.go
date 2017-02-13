// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rawp

import (
	"hash/crc32"
	"image"
	"io"
	"os"
	"unsafe"

	"github.com/golang/snappy"
)

// Options are the encoding parameters.
type Options struct {
	UseSnappy bool
}

func Save(name string, m image.Image, opt *Options) (err error) {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()
	return Encode(f, m, opt)
}

// Encode writes the image m to w in RawP format.
func Encode(w io.Writer, m image.Image, opt *Options) (err error) {
	p, ok := AsMemPImage(m)
	if !ok {
		p = NewMemPImageFrom(m)
	}

	var useSnappy bool
	if opt != nil {
		useSnappy = opt.UseSnappy
	}

	hdr, err := rawpMakeHeader(p.Bounds().Dx(), p.Bounds().Dy(), p.XChannels, p.XDataType, useSnappy)
	if err != nil {
		return
	}

	stride := p.XRect.Dx() * p.XChannels * SizeofKind(p.XDataType)
	pix := make([]byte, stride*p.XRect.Dy())

	off := 0
	for y := p.XRect.Min.Y; y < p.XRect.Max.Y; y++ {
		copy(pix[off:][:stride], p.XPix[p.PixOffset(p.XRect.Min.X, y):])
		off += stride
	}

	if useSnappy {
		pix = snappy.Encode(nil, pix)
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
