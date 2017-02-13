// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rawp

import (
	"fmt"
	"hash/crc32"
	"image/color"
	"math"
	"reflect"
	"unsafe"

	"github.com/golang/snappy"
)

const (
	rawpHeaderSize = 24
	rawpSig        = "RAWP"
	rawpMagic      = 0x1BF2380A // CRC32("RAWP")
)

// data type
const (
	rawpDataType_UInt  = 1
	rawpDataType_Int   = 2
	rawpDataType_Float = 3
)

// RawP Image Spec (Little Endian), 24Bytes.
type rawpHeader struct {
	Sig          [4]byte // 4Bytes, RAWP
	Magic        uint32  // 4Bytes, 0x1BF2380A, CRC32("RAWP")
	Width        uint16  // 2Bytes, image Width
	Height       uint16  // 2Bytes, image Height
	Channels     byte    // 1Bytes, 1=Gray, 3=RGB, 4=RGBA
	Depth        byte    // 1Bytes, 8/16/32/64 bits
	DataType     byte    // 1Bytes, 1=Uint, 2=Int, 3=Float
	UseSnappy    byte    // 1Bytes, 0=disabled, 1=enabled (Header.Data)
	DataSize     uint32  // 4Bytes, image data size (Header.Data)
	DataCheckSum uint32  // 4Bytes, CRC32(RawPHeader.Data[RawPHeader.DataSize])
	Data         []byte  // ?Bytes, image data (RawPHeader.DataSize)
}

func (p *rawpHeader) String() string {
	return fmt.Sprintf(`
rawp.rawpHeader{
	Sig:          %q
	Magic:        0x%x
	Width:        %d
	Height:       %d
	Channels:     %d
	Depth:        %d
	DataType:     %d
	UseSnappy:    %d
	DataSize:     %d
	DataCheckSum: 0x%x
	Data:         ?
}
`[1:],
		p.Sig,
		p.Magic,
		p.Width,
		p.Height,
		p.Channels,
		p.Depth,
		p.DataType,
		p.UseSnappy,
		p.DataSize,
		p.DataCheckSum,
	)
}

func rawpDataType(depth, dataType byte) reflect.Kind {
	switch depth {
	case 8:
		return reflect.Uint8
	case 16:
		return reflect.Uint16
	case 32:
		switch dataType {
		case rawpDataType_UInt:
			return reflect.Uint32
		case rawpDataType_Float:
			return reflect.Float32
		}
	case 64:
		switch dataType {
		case rawpDataType_UInt:
			return reflect.Uint64
		case rawpDataType_Float:
			return reflect.Float64
		}
	}
	return (reflect.Invalid)
}

func rawpIsValidChannels(channels byte) bool {
	return channels > 0
}

func rawpIsValidDepth(depth byte) bool {
	return depth > 0 && (depth%8) == 0
}

func rawpIsValidDataType(t byte) bool {
	return t == rawpDataType_UInt || t == rawpDataType_Int || t == rawpDataType_Float
}

func rawpIsValidHeader(hdr *rawpHeader) error {
	if string(hdr.Sig[:]) != rawpSig {
		return fmt.Errorf("rawp: bad Sig, %v", hdr.Sig)
	}
	if hdr.Magic != rawpMagic {
		return fmt.Errorf("rawp: bad Magic, %x", hdr.Magic)
	}

	if hdr.Width <= 0 || hdr.Height <= 0 {
		return fmt.Errorf("rawp: bad size, width = %v, height = %v", hdr.Width, hdr.Height)
	}
	if !rawpIsValidChannels(hdr.Channels) {
		return fmt.Errorf("rawp: bad Channels, %v", hdr.Channels)
	}
	if !rawpIsValidDepth(hdr.Depth) {
		return fmt.Errorf("rawp: bad Depth, %v", hdr.Depth)
	}
	if !rawpIsValidDataType(hdr.DataType) {
		return fmt.Errorf("rawp: bad DataType, %v", hdr.DataType)
	}

	if hdr.UseSnappy != 0 && hdr.UseSnappy != 1 {
		return fmt.Errorf("rawp: bad UseSnappy, %v", hdr.UseSnappy)
	}
	if hdr.DataSize <= 0 {
		return fmt.Errorf("rawp: bad DataSize, %v", hdr.DataSize)
	}

	// check type more ...
	if hdr.Depth == 8 || hdr.Depth == 16 {
		if hdr.DataType == rawpDataType_Float {
			return fmt.Errorf("rawp: bad Depth, %v", hdr.Depth)
		}
	}

	// check data size more ...
	if x := int(hdr.Width) * int(hdr.Height) * int(hdr.Channels) * int(hdr.Depth) / 8; x < int(hdr.DataSize) {
		return fmt.Errorf("rawp: bad DataSize, %v", hdr.DataSize)
	}

	return nil
}

func rawpColorModel(hdr *rawpHeader) (color.Model, error) {
	if v := hdr.Channels; v != 1 && v != 3 && v != 4 {
		return nil, fmt.Errorf("rawp: unsupport color model, hdr = %v", hdr)
	}
	dataType := rawpDataType(hdr.Depth, hdr.DataType)
	if reflect.Kind(dataType) == reflect.Invalid {
		return nil, fmt.Errorf("rawp: unsupport color model, hdr = %v", hdr)
	}
	return ColorModel(int(hdr.Channels), dataType), nil
}

func rawpMakeHeader(width, height, channels int, dataType reflect.Kind, useSnappy bool) (hdr *rawpHeader, err error) {
	if width <= 0 || width > math.MaxUint16 {
		err = fmt.Errorf("rawp: image size overflow: width = %v, height = %v", width, height)
		return
	}
	if height <= 0 || height > math.MaxUint16 {
		err = fmt.Errorf("rawp: image size overflow: width = %v, height = %v", width, height)
		return
	}
	if v := channels; v != 1 && v != 3 && v != 4 {
		err = fmt.Errorf("rawp: invalid channels: %v", channels)
		return
	}

	hdr = &rawpHeader{
		Sig:      [4]byte{'R', 'A', 'W', 'P'},
		Magic:    rawpMagic,
		Width:    uint16(width),
		Height:   uint16(height),
		Channels: byte(channels),
	}
	if useSnappy {
		hdr.UseSnappy = 1
	}

	switch dataType {
	case reflect.Uint8:
		hdr.Depth = 1 * 8
		hdr.DataType = rawpDataType_UInt
		return
	case reflect.Uint16:
		hdr.Depth = 2 * 8
		hdr.DataType = rawpDataType_UInt
		return
	case reflect.Uint32:
		hdr.Depth = 4 * 8
		hdr.DataType = rawpDataType_UInt
		return
	case reflect.Uint64:
		hdr.Depth = 8 * 8
		hdr.DataType = rawpDataType_UInt
		return
	case reflect.Float32:
		hdr.Depth = 4 * 8
		hdr.DataType = rawpDataType_Float
		return
	case reflect.Float64:
		hdr.Depth = 8 * 8
		hdr.DataType = rawpDataType_Float
		return
	}

	return nil, fmt.Errorf("rawp: unsupport DataType, %V", dataType)
}

func rawpDecodeHeader(data []byte) (hdr *rawpHeader, err error) {
	if len(data) < rawpHeaderSize {
		err = fmt.Errorf("rawp: bad header.")
		return
	}

	// reader header
	hdr = new(rawpHeader)
	copy(((*[1 << 30]byte)(unsafe.Pointer(hdr)))[:rawpHeaderSize], data)
	hdr.Data = data[rawpHeaderSize:]

	// Check CRC32
	if v := crc32.ChecksumIEEE(hdr.Data); v != hdr.DataCheckSum {
		return nil, fmt.Errorf("rawp: bad DataCheckSum, expect = %x, got = %x", hdr.DataCheckSum, v)
	}

	// uncompress
	if hdr.UseSnappy != 0 {
		pix, err := snappy.Decode(nil, hdr.Data)
		if err != nil {
			return nil, fmt.Errorf("rawp: snappyDecode, err = %v", err)
		}
		if len(pix) < int(hdr.DataSize) {
			return nil, fmt.Errorf("rawp: snappyDecode, bad DataSize: %v != %v", len(pix), hdr.DataSize)
		}
		hdr.Data = pix
	}

	// check header
	if err = rawpIsValidHeader(hdr); err != nil {
		return
	}

	return
}
