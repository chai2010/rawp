// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rawp

import (
	"image"
	"io"
	"io/ioutil"
	"reflect"
)

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

	dataType := rawpDataType(hdr.Depth, hdr.DataType)
	p := NewImage(image.Rect(0, 0, int(hdr.Width), int(hdr.Height)), int(hdr.Channels), reflect.Kind(dataType))
	copy(p.Pix, hdr.Data)

	if p.Channels == 1 && p.DataType == reflect.Uint8 {
		return &image.Gray{
			Pix:    p.Pix,
			Stride: p.Stride,
			Rect:   p.Rect,
		}, nil
	}
	if p.Channels == 4 && p.DataType == reflect.Uint8 {
		return &image.RGBA{
			Pix:    p.Pix,
			Stride: p.Stride,
			Rect:   p.Rect,
		}, nil
	}
	if p.Channels == 1 && p.DataType == reflect.Uint16 {
		if isLittleEndian {
			p.Pix.SwapEndian(p.DataType)
		}
		return &image.Gray16{
			Pix:    p.Pix,
			Stride: p.Stride,
			Rect:   p.Rect,
		}, nil
	}
	if p.Channels == 4 && p.DataType == reflect.Uint16 {
		if isLittleEndian {
			p.Pix.SwapEndian(p.DataType)
		}
		return &image.RGBA64{
			Pix:    p.Pix,
			Stride: p.Stride,
			Rect:   p.Rect,
		}, nil
	}

	m = p.StdImage()
	return
}

// DecodeImage reads a RawP image from r and returns it as an Image.
// The type of Image returned depends on the contents of the RawP.
func DecodeImage(r io.Reader) (m *Image, err error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}
	hdr, err := rawpDecodeHeader(data)
	if err != nil {
		return
	}

	dataType := rawpDataType(hdr.Depth, hdr.DataType)
	m = NewImage(image.Rect(0, 0, int(hdr.Width), int(hdr.Height)), int(hdr.Channels), reflect.Kind(dataType))
	copy(m.Pix, hdr.Data)

	return
}

func init() {
	image.RegisterFormat("rawp", "RAWP\x1B\xF2\x38\x0A", Decode, DecodeConfig)
}
