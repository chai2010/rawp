// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package rawp implements a decoder and encoder for RawP images.
//
// RawP Image Structs (Little Endian):
//	type RawPImage struct {
//		Sig          [4]byte // 4Bytes, RAWP
//		Magic        uint32  // 4Bytes, 0x1BF2380A
//		Width        uint16  // 2Bytes, image Width
//		Height       uint16  // 2Bytes, image Height
//		Channels     byte    // 1Bytes, 1=Gray, 3=RGB, 4=RGBA
//		Depth        byte    // 1Bytes, 8/16/32/64 bits
//		DataType     byte    // 1Bytes, 1=Uint, 2=Int, 3=Float
//		UseSnappy    byte    // 1Bytes, 0=disabled, 1=enabled (RawPImage.Data)
//		DataSize     uint32  // 4Bytes, image data size (RawPImage.Data)
//		DataCheckSum uint32  // 4Bytes, CRC32(RawPImage.Data[RawPImage.DataSize])
//		Data         []byte  // ?Bytes, image data (RawPImage.DataSize)
//	}
//
// Please report bugs to chaishushan{AT}gmail.com.
//
// Thanks!
package rawp
