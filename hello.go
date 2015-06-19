// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

package main

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"log"

	"github.com/chai2010/rawp"
)

func main() {
	var buf bytes.Buffer
	var data []byte
	var err error

	// Load file data
	if data, err = ioutil.ReadFile("./testdata/lena.jpg"); err != nil {
		log.Println(err)
	}

	// Decode jpeg
	m0, err := jpeg.Decode(bytes.NewReader(data))
	if err != nil {
		log.Println(err)
	}

	// Encode rawp with snappy
	if err = rawp.Encode(&buf, m0, &rawp.Options{UseSnappy: true}); err != nil {
		log.Println(err)
	}

	// Decode rawp
	m1, err := rawp.Decode(&buf)
	if err != nil {
		log.Println(err)
	}

	// save as jpeg
	if err = jpeg.Encode(&buf, m1, nil); err != nil {
		log.Println(err)
	}
	if err = ioutil.WriteFile("output.jpg", buf.Bytes(), 0666); err != nil {
		log.Println(err)
	}

	fmt.Println("Done")
}
