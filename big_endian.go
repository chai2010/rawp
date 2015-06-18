// Copyright 2015 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build 386 amd64 arm

package rawp

const isBigEndian = false

func nativeToBigEndian(data []byte, elemSize int) {
	if isBigEndian {
		return
	}
	switch elemSize {
	case 2:
		for i := 0; i+2-1 < len(data); i = i + 2 {
			data[i+0], data[i+1] = data[i+1], data[i+0]
		}
	case 4:
		for i := 0; i+4-1 < len(data); i = i + 4 {
			data[i+0], data[i+3] = data[i+3], data[i+0]
			data[i+1], data[i+2] = data[i+2], data[i+1]
		}
	case 8:
		for i := 0; i+8-1 < len(data); i = i + 8 {
			data[i+0], data[i+7] = data[i+7], data[i+0]
			data[i+1], data[i+6] = data[i+6], data[i+1]
			data[i+2], data[i+5] = data[i+5], data[i+2]
			data[i+3], data[i+4] = data[i+4], data[i+3]
		}
	}
	return
}

func bigToNativeEndian(data []byte, elemSize int) {
	if isBigEndian {
		return
	}
	switch elemSize {
	case 2:
		for i := 0; i+2-1 < len(data); i = i + 2 {
			data[i+0], data[i+1] = data[i+1], data[i+0]
		}
	case 4:
		for i := 0; i+4-1 < len(data); i = i + 4 {
			data[i+0], data[i+3] = data[i+3], data[i+0]
			data[i+1], data[i+2] = data[i+2], data[i+1]
		}
	case 8:
		for i := 0; i+8-1 < len(data); i = i + 8 {
			data[i+0], data[i+7] = data[i+7], data[i+0]
			data[i+1], data[i+6] = data[i+6], data[i+1]
			data[i+2], data[i+5] = data[i+5], data[i+2]
			data[i+3], data[i+4] = data[i+4], data[i+3]
		}
	}
	return
}
