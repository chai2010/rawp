// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rawp

import (
	"reflect"
	"unsafe"
)

type DataView []byte

// NewDataView convert a normal slice to byte slice.
//
// Convert []X to []byte:
//
//	x := make([]X, xLen)
//	y := NewDataView(x)
//
// or
//
//	x := make([]X, xLen)
//	y := ((*[1 << 30]byte)(unsafe.Pointer(&x[0])))[:yLen:yLen]
//
func NewDataView(slice interface{}) (data DataView) {
	sv := reflect.ValueOf(slice)
	h := (*reflect.SliceHeader)((unsafe.Pointer(&data)))
	h.Cap = sv.Cap() * int(sv.Type().Elem().Size())
	h.Len = sv.Len() * int(sv.Type().Elem().Size())
	h.Data = sv.Pointer()
	return
}

// Slice convert a normal slice to new type slice.
//
// Convert []byte to []Y:
//	x := make([]byte, xLen)
//	y := DataView(x).Slice(reflect.TypeOf([]Y(nil))).([]Y)
//
// or
//
//	x := make([]X, xLen)
//	y := ((*[1 << 30]Y)(unsafe.Pointer(&x[0])))[:yLen]
//
func (d DataView) Slice(newSliceType reflect.Type) interface{} {
	sv := reflect.ValueOf(d)
	newSlice := reflect.New(newSliceType)
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(newSlice.Pointer()))
	hdr.Cap = sv.Cap() * int(sv.Type().Elem().Size()) / int(newSliceType.Elem().Size())
	hdr.Len = sv.Len() * int(sv.Type().Elem().Size()) / int(newSliceType.Elem().Size())
	hdr.Data = uintptr(sv.Pointer())
	return newSlice.Elem().Interface()
}

func (d DataView) ByteSlice() (v []byte) {
	return d
}

func (d DataView) Int8Slice() (v []int8) {
	h0 := (*reflect.SliceHeader)(unsafe.Pointer(&d))
	h1 := (*reflect.SliceHeader)(unsafe.Pointer(&v))

	h1.Cap = h0.Cap
	h1.Len = h0.Len
	h1.Data = h0.Data
	return
}

func (d DataView) Int16Slice() (v []int16) {
	h0 := (*reflect.SliceHeader)(unsafe.Pointer(&d))
	h1 := (*reflect.SliceHeader)(unsafe.Pointer(&v))

	h1.Cap = h0.Cap / 2
	h1.Len = h0.Len / 2
	h1.Data = h0.Data
	return
}

func (d DataView) Int32Slice() (v []int32) {
	h0 := (*reflect.SliceHeader)(unsafe.Pointer(&d))
	h1 := (*reflect.SliceHeader)(unsafe.Pointer(&v))

	h1.Cap = h0.Cap / 4
	h1.Len = h0.Len / 4
	h1.Data = h0.Data
	return
}

func (d DataView) Int64Slice() (v []int64) {
	h0 := (*reflect.SliceHeader)(unsafe.Pointer(&d))
	h1 := (*reflect.SliceHeader)(unsafe.Pointer(&v))

	h1.Cap = h0.Cap / 8
	h1.Len = h0.Len / 8
	h1.Data = h0.Data
	return
}

func (d DataView) Uint8Slice() []uint8 {
	return d
}

func (d DataView) Uint16Slice() (v []uint16) {
	h0 := (*reflect.SliceHeader)(unsafe.Pointer(&d))
	h1 := (*reflect.SliceHeader)(unsafe.Pointer(&v))

	h1.Cap = h0.Cap / 2
	h1.Len = h0.Len / 2
	h1.Data = h0.Data
	return
}

func (d DataView) Uint32Slice() (v []uint32) {
	h0 := (*reflect.SliceHeader)(unsafe.Pointer(&d))
	h1 := (*reflect.SliceHeader)(unsafe.Pointer(&v))

	h1.Cap = h0.Cap / 4
	h1.Len = h0.Len / 4
	h1.Data = h0.Data
	return
}

func (d DataView) Uint64Slice() (v []uint64) {
	h0 := (*reflect.SliceHeader)(unsafe.Pointer(&d))
	h1 := (*reflect.SliceHeader)(unsafe.Pointer(&v))

	h1.Cap = h0.Cap / 8
	h1.Len = h0.Len / 8
	h1.Data = h0.Data
	return
}

func (d DataView) Float32Slice() (v []float32) {
	h0 := (*reflect.SliceHeader)(unsafe.Pointer(&d))
	h1 := (*reflect.SliceHeader)(unsafe.Pointer(&v))

	h1.Cap = h0.Cap / 4
	h1.Len = h0.Len / 4
	h1.Data = h0.Data
	return
}

func (d DataView) Float64Slice() (v []float64) {
	h0 := (*reflect.SliceHeader)(unsafe.Pointer(&d))
	h1 := (*reflect.SliceHeader)(unsafe.Pointer(&v))

	h1.Cap = h0.Cap / 8
	h1.Len = h0.Len / 8
	h1.Data = h0.Data
	return
}

func (d DataView) Complex64Slice() (v []complex64) {
	h0 := (*reflect.SliceHeader)(unsafe.Pointer(&d))
	h1 := (*reflect.SliceHeader)(unsafe.Pointer(&v))

	h1.Cap = h0.Cap / 16
	h1.Len = h0.Len / 16
	h1.Data = h0.Data
	return
}

func (d DataView) Complex128Slice() (v []complex128) {
	h0 := (*reflect.SliceHeader)(unsafe.Pointer(&d))
	h1 := (*reflect.SliceHeader)(unsafe.Pointer(&v))

	h1.Cap = h0.Cap / 32
	h1.Len = h0.Len / 32
	h1.Data = h0.Data
	return
}

func (d DataView) Value(i int, dataType reflect.Kind) float64 {
	switch dataType {
	case reflect.Uint8:
		return float64(d[i])
	case reflect.Uint16:
		return float64(d.Uint16Slice()[i])
	case reflect.Uint32:
		return float64(d.Uint32Slice()[i])
	case reflect.Uint64:
		return float64(d.Uint64Slice()[i])
	case reflect.Float32:
		return float64(d.Float32Slice()[i])
	case reflect.Float64:
		return float64(d.Float64Slice()[i])
	}
	return 0
}

func (d DataView) SetValue(i int, dataType reflect.Kind, v float64) {
	switch dataType {
	case reflect.Uint8:
		d[i] = byte(v)
	case reflect.Uint16:
		d.Uint16Slice()[i] = uint16(v)
	case reflect.Uint32:
		d.Uint32Slice()[i] = uint32(v)
	case reflect.Uint64:
		d.Uint64Slice()[i] = uint64(v)
	case reflect.Float32:
		d.Float32Slice()[i] = float32(v)
	case reflect.Float64:
		d.Float64Slice()[i] = float64(v)
	}
}
