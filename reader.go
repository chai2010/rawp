// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rawp

import (
	"image"
	"io"
	"io/ioutil"
	"os"
	"reflect"
)

func LoadConfig(name string) (config image.Config, err error) {
	f, err := os.Open(name)
	if err != nil {
		return image.Config{}, err
	}
	defer f.Close()
	return DecodeConfig(f)
}

func Load(name string) (m image.Image, err error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Decode(f)
}

func LoadImage(name string) (m *MemPImage, err error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return DecodeImage(f)
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

	p := &MemPImage{
		XMemPMagic: MemPMagic,
		XRect:      image.Rect(0, 0, int(hdr.Width), int(hdr.Height)),
		XStride:    int(hdr.Width) * int(hdr.Channels) * SizeofKind(rawpDataType(hdr.Depth, hdr.DataType)),
		XChannels:  int(hdr.Channels),
		XDataType:  rawpDataType(hdr.Depth, hdr.DataType),
		XPix:       hdr.Data,
	}

	if p.XChannels == 1 && p.XDataType == reflect.Uint8 {
		return &image.Gray{
			Pix:    p.XPix,
			Stride: p.XStride,
			Rect:   p.XRect,
		}, nil
	}
	if p.XChannels == 4 && p.XDataType == reflect.Uint8 {
		return &image.RGBA{
			Pix:    p.XPix,
			Stride: p.XStride,
			Rect:   p.XRect,
		}, nil
	}
	if p.XChannels == 1 && p.XDataType == reflect.Uint16 {
		if isLittleEndian {
			p.XPix.SwapEndian(p.XDataType)
		}
		return &image.Gray16{
			Pix:    p.XPix,
			Stride: p.XStride,
			Rect:   p.XRect,
		}, nil
	}
	if p.XChannels == 4 && p.XDataType == reflect.Uint16 {
		if isLittleEndian {
			p.XPix.SwapEndian(p.XDataType)
		}
		return &image.RGBA64{
			Pix:    p.XPix,
			Stride: p.XStride,
			Rect:   p.XRect,
		}, nil
	}

	m = p.StdImage()
	return
}

// DecodeImage reads a RawP image from r and returns it as an Image.
// The type of Image returned depends on the contents of the RawP.
func DecodeImage(r io.Reader) (m *MemPImage, err error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}
	hdr, err := rawpDecodeHeader(data)
	if err != nil {
		return
	}

	m = &MemPImage{
		XMemPMagic: MemPMagic,
		XRect:      image.Rect(0, 0, int(hdr.Width), int(hdr.Height)),
		XStride:    int(hdr.Width) * int(hdr.Channels) * SizeofKind(rawpDataType(hdr.Depth, hdr.DataType)),
		XChannels:  int(hdr.Channels),
		XDataType:  rawpDataType(hdr.Depth, hdr.DataType),
		XPix:       hdr.Data,
	}
	return
}

func init() {
	image.RegisterFormat("rawp", "RAWP\x1B\xF2\x38\x0A", Decode, DecodeConfig)
}
