// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rawp

import (
	"image/color"
	"reflect"
)

type Pixel struct {
	Channels int
	DataType reflect.Kind
	Pix      DataView
}

func (c Pixel) RGBA() (r, g, b, a uint32) {
	if len(c.Pix) == 0 {
		return
	}
	switch c.Channels {
	case 1:
		switch reflect.Kind(c.DataType) {
		case reflect.Uint8:
			return color.Gray{
				Y: c.Pix[0],
			}.RGBA()
		case reflect.Uint16:
			return color.Gray16{
				Y: c.Pix.Uint16Slice()[0],
			}.RGBA()
		default:
			return color.Gray16{
				Y: uint16(c.Pix.Value(0, reflect.Kind(c.DataType))),
			}.RGBA()
		}
	case 2:
		switch reflect.Kind(c.DataType) {
		case reflect.Uint8:
			return color.RGBA{
				R: c.Pix[0],
				G: c.Pix[1],
				B: 0xFF,
				A: 0xFF,
			}.RGBA()
		case reflect.Uint16:
			return color.RGBA64{
				R: c.Pix.Uint16Slice()[0],
				G: c.Pix.Uint16Slice()[1],
				B: 0xFFFF,
				A: 0xFFFF,
			}.RGBA()
		default:
			return color.RGBA64{
				R: uint16(c.Pix.Value(0, reflect.Kind(c.DataType))),
				G: uint16(c.Pix.Value(1, reflect.Kind(c.DataType))),
				B: 0xFFFF,
				A: 0xFFFF,
			}.RGBA()
		}
	case 3:
		switch reflect.Kind(c.DataType) {
		case reflect.Uint8:
			return color.RGBA{
				R: c.Pix[0],
				G: c.Pix[1],
				B: c.Pix[2],
				A: 0xFF,
			}.RGBA()
		case reflect.Uint16:
			return color.RGBA64{
				R: c.Pix.Uint16Slice()[0],
				G: c.Pix.Uint16Slice()[1],
				B: c.Pix.Uint16Slice()[2],
				A: 0xFFFF,
			}.RGBA()
		default:
			return color.RGBA64{
				R: uint16(c.Pix.Value(0, reflect.Kind(c.DataType))),
				G: uint16(c.Pix.Value(1, reflect.Kind(c.DataType))),
				B: uint16(c.Pix.Value(2, reflect.Kind(c.DataType))),
				A: 0xFFFF,
			}.RGBA()
		}
	case 4:
		switch reflect.Kind(c.DataType) {
		case reflect.Uint8:
			return color.RGBA{
				R: c.Pix[0],
				G: c.Pix[1],
				B: c.Pix[2],
				A: c.Pix[3],
			}.RGBA()
		case reflect.Uint16:
			return color.RGBA64{
				R: c.Pix.Uint16Slice()[0],
				G: c.Pix.Uint16Slice()[1],
				B: c.Pix.Uint16Slice()[2],
				A: c.Pix.Uint16Slice()[3],
			}.RGBA()
		default:
			return color.RGBA64{
				R: uint16(c.Pix.Value(0, reflect.Kind(c.DataType))),
				G: uint16(c.Pix.Value(1, reflect.Kind(c.DataType))),
				B: uint16(c.Pix.Value(2, reflect.Kind(c.DataType))),
				A: uint16(c.Pix.Value(3, reflect.Kind(c.DataType))),
			}.RGBA()
		}
	}
	return
}

func ColorModel(channels int, dataType reflect.Kind) color.Model {
	return color.ModelFunc(func(c color.Color) color.Color {
		return colorModelConvert(channels, dataType, c)
	})
}

func colorModelConvert(channels int, dataType reflect.Kind, c color.Color) color.Color {
	c2 := Pixel{
		Channels: channels,
		DataType: dataType,
		Pix:      make(DataView, channels*SizeofKind(dataType)),
	}

	if c1, ok := c.(Pixel); ok {
		if c1.Channels == c2.Channels && c1.DataType == c2.DataType {
			copy(c2.Pix, c1.Pix)
			return c2
		}
		if c1.DataType == c2.DataType {
			copy(c2.Pix, c1.Pix)
			return c2
		}
		for i := 0; i < c1.Channels && i < c2.Channels; i++ {
			c2.Pix.SetValue(i, reflect.Kind(c2.DataType), c1.Pix.Value(i, reflect.Kind(c1.DataType)))
		}
		return c2
	}

	switch {
	case channels == 1 && reflect.Kind(dataType) == reflect.Uint8:
		v := color.GrayModel.Convert(c).(color.Gray)
		c2.Pix[0] = v.Y
		return c2
	case channels == 1 && reflect.Kind(dataType) == reflect.Uint16:
		v := color.Gray16Model.Convert(c).(color.Gray16)
		c2.Pix[0] = uint8(v.Y >> 8)
		c2.Pix[1] = uint8(v.Y)
		return c2
	case channels == 3 && reflect.Kind(dataType) == reflect.Uint8:
		r, g, b, _ := c.RGBA()
		c2.Pix[0] = uint8(r >> 8)
		c2.Pix[1] = uint8(g >> 8)
		c2.Pix[2] = uint8(b >> 8)
		return c2
	case channels == 3 && reflect.Kind(dataType) == reflect.Uint16:
		r, g, b, _ := c.RGBA()
		c2.Pix[0] = uint8(r >> 8)
		c2.Pix[1] = uint8(r)
		c2.Pix[2] = uint8(g >> 8)
		c2.Pix[3] = uint8(g)
		c2.Pix[4] = uint8(b >> 8)
		c2.Pix[5] = uint8(b)
		return c2
	case channels == 4 && reflect.Kind(dataType) == reflect.Uint8:
		r, g, b, a := c.RGBA()
		c2.Pix[0] = uint8(r >> 8)
		c2.Pix[1] = uint8(g >> 8)
		c2.Pix[2] = uint8(b >> 8)
		c2.Pix[3] = uint8(a >> 8)
		return c2
	case channels == 4 && reflect.Kind(dataType) == reflect.Uint16:
		r, g, b, a := c.RGBA()
		c2.Pix[0] = uint8(r >> 8)
		c2.Pix[1] = uint8(r)
		c2.Pix[2] = uint8(g >> 8)
		c2.Pix[3] = uint8(g)
		c2.Pix[4] = uint8(b >> 8)
		c2.Pix[5] = uint8(b)
		c2.Pix[6] = uint8(a >> 8)
		c2.Pix[7] = uint8(a)
		return c2
	}

	r, g, b, a := c.RGBA()
	rgba := []uint32{r, g, b, a}
	for i := 0; i < c2.Channels && i < len(rgba); i++ {
		c2.Pix.SetValue(i, reflect.Kind(c2.DataType), float64(rgba[i]))
	}
	return c2
}

func SizeofKind(d reflect.Kind) int {
	switch reflect.Kind(d) {
	case reflect.Int8:
		return 1
	case reflect.Int16:
		return 2
	case reflect.Int32:
		return 4
	case reflect.Int64:
		return 8
	case reflect.Uint8:
		return 1
	case reflect.Uint16:
		return 2
	case reflect.Uint32:
		return 4
	case reflect.Uint64:
		return 8
	case reflect.Float32:
		return 4
	case reflect.Float64:
		return 8
	case reflect.Complex64:
		return 8
	case reflect.Complex128:
		return 16
	}
	return 0
}
