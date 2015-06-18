// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rawp

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestEncodeAndDecode(t *testing.T) {
	var buf bytes.Buffer
	var data []byte
	var err error

	// Load file data
	if data, err = ioutil.ReadFile("./testdata/lena.jpg"); err != nil {
		t.Fatal(err)
	}

	// Decode jpeg
	m0, err := jpeg.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}

	// Encode snappy rawp
	if err := Encode(&buf, m0, &Options{UseSnappy: true}); err != nil {
		t.Fatal(err)
	}

	// Decode rawp
	m1, err := Decode(&buf)
	if err != nil {
		t.Fatal(err)
	}

	// compare image size
	if b0, b1 := m0.Bounds(), m1.Bounds(); b0 != b1 {
		t.Fatalf("bounds: %v, %v", b0, b1)
	}

	// compare pixel
	b := m0.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c0 := color.RGBAModel.Convert(m0.At(x, y))
			c1 := color.RGBAModel.Convert(m1.At(x, y))
			if c0 != c1 {
				t.Fatalf("(%d,%d): %v, %v", x, y, c0, c1)
			}
		}
	}
}

func BenchmarkDecode(b *testing.B) {
	data, err := ioutil.ReadFile("./testdata/lena.jpg")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Decode(bytes.NewReader(data))
	}
}

func BenchmarkEncode(b *testing.B) {
	m := tLoadImage("./testdata/lena.jpg")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Encode(ioutil.Discard, m, nil)
	}
}

func BenchmarkEncode_snappy(b *testing.B) {
	m := tLoadImage("./testdata/lena.jpg")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Encode(ioutil.Discard, m, &Options{UseSnappy: true})
	}
}

func tLoadImage(filename string) image.Image {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("loadImage: os.Open(%q), err= %v", filename, err)
	}
	defer f.Close()

	m, _, err := image.Decode(f)
	if err != nil {
		log.Fatalf("loadImage: image.Decode, err= %v", err)
	}
	return m
}
