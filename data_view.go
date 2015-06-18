// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rawp

import (
	"fmt"
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
	if sv.Kind() != reflect.Slice {
		panic(fmt.Sprintf("rawp: NewDataView called with non-slice value of type %T", slice))
	}
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
	if sv.Kind() != reflect.Slice {
		panic(fmt.Sprintf("rawp: DataView.Slice called with non-slice value of type %T", d))
	}
	if newSliceType.Kind() != reflect.Slice {
		panic(fmt.Sprintf("rawp: DataView.Slice called with non-slice type of type %T", newSliceType))
	}
	newSlice := reflect.New(newSliceType)
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(newSlice.Pointer()))
	hdr.Cap = sv.Cap() * int(sv.Type().Elem().Size()) / int(newSliceType.Elem().Size())
	hdr.Len = sv.Len() * int(sv.Type().Elem().Size()) / int(newSliceType.Elem().Size())
	hdr.Data = uintptr(sv.Pointer())
	return newSlice.Elem().Interface()
}

func (d DataView) Byte(i int) byte {
	return d[i]
}

func (d DataView) Uint16(i int) uint16 {
	return d.UInt16Slice()[i]
}

func (d DataView) Uint32(i int) uint32 {
	return d.UInt32Slice()[i]
}

func (d DataView) Uint64(i int) uint64 {
	return d.UInt64Slice()[i]
}

func (d DataView) Float32(i int) float32 {
	return d.Float32Slice()[i]
}

func (d DataView) Float64(i int) float64 {
	return d.Float64Slice()[i]
}

func (d DataView) FloatValue(i int, dataType DataType) float64 {
	switch dataType {
	case Uint8:
		return float64(d.Byte(i))
	case Uint16:
		return float64(d.Uint16(i))
	case Uint32:
		return float64(d.Uint32(i))
	case Uint64:
		return float64(d.Uint64(i))
	case Float32:
		return float64(d.Float32(i))
	case Float64:
		return float64(d.Float64(i))
	}
	return 0
}

func (d DataView) SetByte(i int, v ...byte) {
	copy(d[i:], v)
}

func (d DataView) SetUInt16(i int, v ...uint16) {
	copy(d.UInt16Slice()[i:], v)
}

func (d DataView) SetUInt32(i int, v ...uint32) {
	copy(d.UInt32Slice()[i:], v)
}

func (d DataView) SetUInt64(i int, v ...uint64) {
	copy(d.UInt64Slice()[i:], v)
}

func (d DataView) SetFloat32(i int, v ...float32) {
	copy(d.Float32Slice()[i:], v)
}

func (d DataView) SetFloat64(i int, v ...float64) {
	copy(d.Float64Slice()[i:], v)
}

func (d DataView) SetFloatValue(i int, dataType DataType, v float64) {
	switch dataType {
	case Uint8:
		d.SetByte(i, byte(v))
	case Uint16:
		d.SetUInt16(i, uint16(v))
	case Uint32:
		d.SetUInt32(i, uint32(v))
	case Uint64:
		d.SetUInt64(i, uint64(v))
	case Float32:
		d.SetFloat32(i, float32(v))
	case Float64:
		d.SetFloat64(i, float64(v))
	}
}

func (d DataView) ByteSlice() []byte {
	return d
}

func (d DataView) UInt16Slice() []uint16 {
	return ((*[1 << 30]uint16)(unsafe.Pointer(&d[0])))[0 : len(d)/2 : len(d)/2]
}

func (d DataView) UInt32Slice() []uint32 {
	return ((*[1 << 30]uint32)(unsafe.Pointer(&d[0])))[0 : len(d)/4 : len(d)/4]
}

func (d DataView) UInt64Slice() []uint64 {
	return ((*[1 << 30]uint64)(unsafe.Pointer(&d[0])))[0 : len(d)/8 : len(d)/8]
}

func (d DataView) Float32Slice() []float32 {
	return ((*[1 << 30]float32)(unsafe.Pointer(&d[0])))[0 : len(d)/4 : len(d)/4]
}

func (d DataView) Float64Slice() []float64 {
	return ((*[1 << 30]float64)(unsafe.Pointer(&d[0])))[0 : len(d)/8 : len(d)/8]
}
