// Copyright 2015 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rawp

import (
	"fmt"
	"reflect"
)

type DataType reflect.Kind

func (d DataType) ByteSize() int {
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

func (d DataType) Depth() int {
	return d.ByteSize() * 8
}

func (d DataType) String() string {
	return fmt.Sprintf("DataType(%s)", reflect.Kind(d).String())
}
