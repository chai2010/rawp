// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rawp

import (
	"image"
	"image/color"
	"reflect"
)

var (
	_ image.Image = (*Image)(nil)
)

type Image struct {
	Rect     image.Rectangle
	Channels int
	DataType reflect.Kind
	Pix      []byte

	// Stride is the Pix stride (in bytes, must align with PixelSize)
	// between vertically adjacent pixels.
	Stride int
}

func NewImage(r image.Rectangle, channels int, dataType reflect.Kind) *Image {
	m := &Image{
		Rect:     r,
		Stride:   r.Dx() * channels * SizeofKind(dataType),
		Channels: channels,
		DataType: dataType,
	}
	m.Pix = make([]byte, r.Dy()*m.Stride)
	return m
}

func NewImageFrom(m image.Image) *Image {
	if p, _ := m.(*Image); p != nil {
		return p
	}

	switch m := m.(type) {
	case *image.Gray:
		b := m.Bounds()
		p := NewImage(b, 1, reflect.Uint8)

		for y := b.Min.Y; y < b.Max.Y; y++ {
			off0 := m.PixOffset(0, y)
			off1 := p.PixOffset(0, y)
			copy(p.Pix[off1:][:p.Stride], m.Pix[off0:][:m.Stride])
			off0 += m.Stride
			off1 += p.Stride
		}
		return p

	case *image.Gray16:
		b := m.Bounds()
		p := NewImage(b, 1, reflect.Uint16)

		for y := b.Min.Y; y < b.Max.Y; y++ {
			off0 := m.PixOffset(0, y)
			off1 := p.PixOffset(0, y)
			copy(p.Pix[off1:][:p.Stride], m.Pix[off0:][:m.Stride])
			off0 += m.Stride
			off1 += p.Stride
		}

		if !isBigEndian {
			bigToNativeEndian(p.Pix, SizeofKind(p.DataType))
		}
		return p

	case *image.RGBA:
		b := m.Bounds()
		p := NewImage(b, 4, reflect.Uint8)

		for y := b.Min.Y; y < b.Max.Y; y++ {
			off0 := m.PixOffset(0, y)
			off1 := p.PixOffset(0, y)
			copy(p.Pix[off1:][:p.Stride], m.Pix[off0:][:m.Stride])
			off0 += m.Stride
			off1 += p.Stride
		}
		return p

	case *image.RGBA64:
		b := m.Bounds()
		p := NewImage(b, 4, reflect.Uint16)

		for y := b.Min.Y; y < b.Max.Y; y++ {
			off0 := m.PixOffset(0, y)
			off1 := p.PixOffset(0, y)
			copy(p.Pix[off1:][:p.Stride], m.Pix[off0:][:m.Stride])
			off0 += m.Stride
			off1 += p.Stride
		}
		if !isBigEndian {
			bigToNativeEndian(p.Pix, SizeofKind(p.DataType))
		}
		return p

	case *image.YCbCr:
		b := m.Bounds()
		p := NewImage(b, 4, reflect.Uint8)
		for y := b.Min.Y; y < b.Max.Y; y++ {
			for x := b.Min.X; x < b.Max.X; x++ {
				R, G, B, A := m.At(x, y).RGBA()

				i := p.PixOffset(x, y)
				p.Pix[i+0] = uint8(R >> 8)
				p.Pix[i+1] = uint8(G >> 8)
				p.Pix[i+2] = uint8(B >> 8)
				p.Pix[i+3] = uint8(A >> 8)
			}
		}
		return p

	default:
		b := m.Bounds()
		p := NewImage(b, 4, reflect.Uint16)
		for y := b.Min.Y; y < b.Max.Y; y++ {
			for x := b.Min.X; x < b.Max.X; x++ {
				R, G, B, A := m.At(x, y).RGBA()

				i := p.PixOffset(x, y)
				p.Pix[i+0] = uint8(R >> 8)
				p.Pix[i+1] = uint8(R)
				p.Pix[i+2] = uint8(G >> 8)
				p.Pix[i+3] = uint8(G)
				p.Pix[i+4] = uint8(B >> 8)
				p.Pix[i+5] = uint8(B)
				p.Pix[i+6] = uint8(A >> 8)
				p.Pix[i+7] = uint8(A)
			}
		}
		return p
	}
}

func (p *Image) Bounds() image.Rectangle {
	return p.Rect
}

func (p *Image) ColorModel() color.Model {
	return ColorModel(p.Channels, p.DataType)
}

func (p *Image) At(x, y int) color.Color {
	if !(image.Point{x, y}.In(p.Rect)) {
		return Pixel{
			Channels: p.Channels,
			DataType: p.DataType,
		}
	}
	i, n := p.PixOffset(x, y), p.PixSize()
	return Pixel{
		Channels: p.Channels,
		DataType: p.DataType,
		Pix:      p.Pix[i:][:n],
	}
}

func (p *Image) PixelAt(x, y int) []byte {
	if !(image.Point{x, y}.In(p.Rect)) {
		return nil
	}
	i, n := p.PixOffset(x, y), p.PixSize()
	return p.Pix[i:][:n]
}

func (p *Image) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i, n := p.PixOffset(x, y), p.PixSize()
	v := p.ColorModel().Convert(c).(Pixel)
	copy(p.Pix[i:][:n], v.Pix)
}

func (p *Image) SetPixel(x, y int, c []byte) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i, n := p.PixOffset(x, y), p.PixSize()
	copy(p.Pix[i:][:n], c)
}

func (p *Image) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*p.PixSize()
}

func (p *Image) PixSize() int {
	return p.Channels * SizeofKind(p.DataType)
}

func (p *Image) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to be inside
	// either r1 or r2 if the intersection is empty. Without explicitly checking for
	// this, the Pix[i:] expression below can panic.
	if r.Empty() {
		return &Image{}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &Image{
		Pix:      p.Pix[i:],
		Stride:   p.Stride,
		Rect:     r,
		Channels: p.Channels,
		DataType: p.DataType,
	}
}

func (p *Image) StdImage() image.Image {
	switch {
	case p.Channels == 1 && reflect.Kind(p.DataType) == reflect.Uint8:
		return &image.Gray{
			Pix:    p.Pix,
			Stride: p.Stride,
			Rect:   p.Rect,
		}
	case p.Channels == 1 && reflect.Kind(p.DataType) == reflect.Uint16:
		m := &image.Gray16{
			Pix:    p.Pix,
			Stride: p.Stride,
			Rect:   p.Rect,
		}
		if !isBigEndian {
			m.Pix = append([]byte(nil), m.Pix...)
			nativeToBigEndian(m.Pix, SizeofKind(p.DataType))
		}
		return m
	case p.Channels == 4 && reflect.Kind(p.DataType) == reflect.Uint8:
		return &image.RGBA{
			Pix:    p.Pix,
			Stride: p.Stride,
			Rect:   p.Rect,
		}
	case p.Channels == 4 && reflect.Kind(p.DataType) == reflect.Uint16:
		m := &image.RGBA64{
			Pix:    p.Pix,
			Stride: p.Stride,
			Rect:   p.Rect,
		}
		if !isBigEndian {
			m.Pix = append([]byte(nil), m.Pix...)
			nativeToBigEndian(m.Pix, SizeofKind(p.DataType))
		}
		return m
	}

	return p
}
