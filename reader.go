// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

package rawp

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"io/ioutil"
)

// Options are the encoding and decoding parameters.
type Options struct {
	RawPColorModel color.Model
	UseSnappy      bool
}

// DecodeConfig returns the color model and dimensions of a RawP image without
// decoding the entire image.
func DecodeConfig(r io.Reader) (config image.Config, err error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}
	hdr, err := rawpDecodeHeader(data)
	if err != nil {
		return
	}

	model, err := rawpColorModel(hdr)
	if err != nil {
		return
	}

	config = image.Config{
		ColorModel: model,
		Width:      int(hdr.Width),
		Height:     int(hdr.Height),
	}
	return
}

// Decode reads a RawP image from r and returns it as an image.Image.
// The type of Image returned depends on the contents of the RawP.
func Decode(r io.Reader) (m image.Image, err error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}
	hdr, err := rawpDecodeHeader(data)
	if err != nil {
		return
	}

	// new decoder
	decoder, err := rawpPixDecoder(hdr)
	if err != nil {
		return
	}

	// decode snappy
	pix := hdr.Data
	if hdr.UseSnappy != 0 {
		if pix, err = snappy.Decode(nil, hdr.Data); err != nil {
			err = fmt.Errorf("image/rawp: Decode, snappy err: %v", err)
			return
		}
	}

	// decode raw pix
	m, err = decoder.Decode(pix, nil)
	if err != nil {
		return
	}

	return
}

func init() {
	image.RegisterFormat("rawp", "RAWP\x1B\xF2\x38\x0A", Decode, DecodeConfig)
}