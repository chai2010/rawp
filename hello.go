// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

package main

import (
	"bytes"
	"io/ioutil"
	"log"

	"github.com/chai2010/rawp"
)

func main() {
	var buf bytes.Buffer
	var data []byte
	var err error

	// Load file data
	if data, err = ioutil.ReadFile("./testdata/lena.rawp"); err != nil {
		log.Println(err)
	}

	// Decode rawp
	m, err := rawp.Decode(bytes.NewReader(data))
	if err != nil {
		log.Println(err)
	}

	// Encode lossless rawp
	if err = rawp.Encode(&buf, m, &rawp.Options{UseSnappy: true}); err != nil {
		log.Println(err)
	}
	if err = ioutil.WriteFile("output.rawp", buf.Bytes(), 0666); err != nil {
		log.Println(err)
	}
}
