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
```

BUGS
====

Report bugs to <chaishushan@gmail.com>.

Thanks!
