rawp
=====

[![Build Status](https://travis-ci.org/chai2010/rawp.svg)](https://travis-ci.org/chai2010/rawp)
[![GoDoc](https://godoc.org/github.com/chai2010/rawp?status.svg)](https://godoc.org/github.com/chai2010/rawp)

Install
=======

1. `go get github.com/chai2010/rawp`
2. `go run hello.go`


Example
=======

This is a simple example:

```Go
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
```

BUGS
====

Report bugs to <chaishushan@gmail.com>.

Thanks!
