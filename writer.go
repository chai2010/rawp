// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

package rawp

import (
	"hash/crc32"
	"image"
	"image/color"
	"io"
	"unsafe"
)

// Encode writes the image m to w in RawP format.
func Encode(w io.Writer, m image.Image, opt *Options) (err error) {
	if opt != nil && opt.RawPColorModel != nil {
		m = convert.ColorModel(m, opt.RawPColorModel)
	}
	m = adjustImage(m)

	var useSnappy bool
	if opt != nil {
		useSnappy = opt.UseSnappy
	}

	hdr, err := rawpMakeHeader(m.Bounds().Dx(), m.Bounds().Dy(), m.ColorModel(), useSnappy)
	if err != nil {
		return
	}

	// encode raw pix
	encoder, err := rawpPixEncoder(hdr)
	if err != nil {
		return
	}
	pix, err := encoder.Encode(m, nil)
	if err != nil {
		return
	}
	if useSnappy {
		pix, err = snappy.Encode(nil, pix)
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
