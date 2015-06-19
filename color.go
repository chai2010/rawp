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
	DataType DataType
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
				Y: c.Pix.Byte(0),
			}.RGBA()
		case reflect.Uint16:
			return color.Gray16{
				Y: c.Pix.Uint16(0),
			}.RGBA()
		default:
			return color.Gray16{
				Y: uint16(c.Pix.FloatValue(0, reflect.Kind(c.DataType))),
			}.RGBA()
		}
	case 2:
		switch reflect.Kind(c.DataType) {
		case reflect.Uint8:
			return color.RGBA{
				R: c.Pix.Byte(0),
				G: c.Pix.Byte(1),
				B: 0xFF,
				A: 0xFF,
			}.RGBA()
		case reflect.Uint16:
			return color.RGBA64{
				R: c.Pix.Uint16(0),
				G: c.Pix.Uint16(1),
				B: 0xFFFF,
				A: 0xFFFF,
			}.RGBA()
		default:
			return color.RGBA64{
				R: uint16(c.Pix.FloatValue(0, reflect.Kind(c.DataType))),
				G: uint16(c.Pix.FloatValue(1, reflect.Kind(c.DataType))),
				B: 0xFFFF,
				A: 0xFFFF,
			}.RGBA()
		}
	case 3:
		switch reflect.Kind(c.DataType) {
		case reflect.Uint8:
			return color.RGBA{
				R: c.Pix.Byte(0),
				G: c.Pix.Byte(1),
				B: c.Pix.Byte(2),
				A: 0xFF,
			}.RGBA()
		case reflect.Uint16:
			return color.RGBA64{
				R: c.Pix.Uint16(0),
				G: c.Pix.Uint16(1),
				B: c.Pix.Uint16(2),
				A: 0xFFFF,
			}.RGBA()
		default:
			return color.RGBA64{
				R: uint16(c.Pix.FloatValue(0, reflect.Kind(c.DataType))),
				G: uint16(c.Pix.FloatValue(1, reflect.Kind(c.DataType))),
				B: uint16(c.Pix.FloatValue(2, reflect.Kind(c.DataType))),
				A: 0xFFFF,
			}.RGBA()
		}
	case 4:
		switch reflect.Kind(c.DataType) {
		case reflect.Uint8:
			return color.RGBA{
				R: c.Pix.Byte(0),
				G: c.Pix.Byte(1),
				B: c.Pix.Byte(2),
				A: c.Pix.Byte(3),
			}.RGBA()
		case reflect.Uint16:
			return color.RGBA64{
				R: c.Pix.Uint16(0),
				G: c.Pix.Uint16(1),
				B: c.Pix.Uint16(2),
				A: c.Pix.Uint16(3),
			}.RGBA()
		default:
			return color.RGBA64{
				R: uint16(c.Pix.FloatValue(0, reflect.Kind(c.DataType))),
				G: uint16(c.Pix.FloatValue(1, reflect.Kind(c.DataType))),
				B: uint16(c.Pix.FloatValue(2, reflect.Kind(c.DataType))),
				A: uint16(c.Pix.FloatValue(3, reflect.Kind(c.DataType))),
			}.RGBA()
		}
	}
	return
}

func ColorModel(channels int, dataType DataType) color.Model {
	return color.ModelFunc(func(c color.Color) color.Color {
		return colorModelConvert(channels, dataType, c)
	})
}

func colorModelConvert(channels int, dataType DataType, c color.Color) color.Color {
	c2 := Pixel{
		Channels: channels,
		DataType: dataType,
		Pix:      make(DataView, channels*dataType.ByteSize()),
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
			c2.Pix.SetFloatValue(i, reflect.Kind(c2.DataType), c1.Pix.FloatValue(i, reflect.Kind(c1.DataType)))
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
		c2.Pix.SetFloatValue(i, reflect.Kind(c2.DataType), float64(rgba[i]))
	}
	return c2
}
