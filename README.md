rawp
=====

PkgDoc: [http://godoc.org/github.com/chai2010/rawp](http://godoc.org/github.com/chai2010/rawp)

Install
=======

1. `go get github.com/chai2010/rawp`
2. `go run hello.go`

Note: Only support `Decode` and `DecodeConfig`, `go test` will failed on some other api test.


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

	// Decode rawp
	m, err := jpeg.Decode(bytes.NewReader(data))
	if err != nil {
		log.Println(err)
	}

	// Encode snappy rawp
	if err := rawp.Encode(&buf, m, &rawp.Options{UseSnappy: true}); err != nil {
		log.Println(err)
	}

	// Decode rawp
	m, err = rawp.Decode(&buf)
	if err != nil {
		log.Println(err)
	}

	fmt.Println("Done")
}
```

BUGS
====

Report bugs to <chaishushan@gmail.com>.

Thanks!
