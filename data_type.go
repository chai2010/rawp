// Copyright 2015 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rawp

import (
	"fmt"
)

type DataType byte

const (
	Invalid DataType = iota
	Uint8
	Uint16
	Uint32
	Uint64
	Float32
	Float64
)

func (d DataType) Valid() bool {
	return d <= Float64
}

func (d DataType) Depth() int {
	switch d {
	case Uint8:
		return 1 * 8
	case Uint16:
		return 2 * 8
	case Uint32:
		return 4 * 8
	case Uint64:
		return 8 * 8
	case Float32:
		return 4 * 8
	case Float64:
		return 8 * 8
	}
	return 0
}

func (d DataType) ByteSize() int {
	switch d {
	case Uint8:
		return 1
	case Uint16:
		return 2
	case Uint32:
		return 4
	case Uint64:
		return 8
	case Float32:
		return 4
	case Float64:
		return 8
	}
	return 0
}

func (d DataType) String() string {
	switch d {
	case Uint8:
		return "Uint8"
	case Uint16:
		return "Uint16"
	case Uint32:
		return "Uint32"
	case Uint64:
		return "Uint64"
	case Float32:
		return "Float32"
	case Float64:
		return "Float64"
	}
	return fmt.Sprintf("DataType(%d)", int(d))
}
