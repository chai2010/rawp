// Copyright 2011 The Snappy-Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rawp

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var download = flag.Bool("download", false, "If true, download any missing files before running benchmarks")

func tSnappyRoundtrip(b, ebuf, dbuf []byte) error {
	e, err := snappyEncode(ebuf, b)
	if err != nil {
		return fmt.Errorf("encoding error: %v", err)
	}
	d, err := snappyDecode(dbuf, e)
	if err != nil {
		return fmt.Errorf("decoding error: %v", err)
	}
	if !bytes.Equal(b, d) {
		return fmt.Errorf("tSnappyRoundtrip mismatch:\n\twant %v\n\tgot  %v", b, d)
	}
	return nil
}

func TestSnappyEmpty(t *testing.T) {
	if err := tSnappyRoundtrip(nil, nil, nil); err != nil {
		t.Fatal(err)
	}
}

func TestSnappySmallCopy(t *testing.T) {
	for _, ebuf := range [][]byte{nil, make([]byte, 20), make([]byte, 64)} {
		for _, dbuf := range [][]byte{nil, make([]byte, 20), make([]byte, 64)} {
			for i := 0; i < 32; i++ {
				s := "aaaa" + strings.Repeat("b", i) + "aaaabbbb"
				if err := tSnappyRoundtrip([]byte(s), ebuf, dbuf); err != nil {
					t.Errorf("len(ebuf)=%d, len(dbuf)=%d, i=%d: %v", len(ebuf), len(dbuf), i, err)
				}
			}
		}
	}
}

func TestSnappySmallRand(t *testing.T) {
	rand.Seed(27354294)
	for n := 1; n < 20000; n += 23 {
		b := make([]byte, n)
		for i, _ := range b {
			b[i] = uint8(rand.Uint32())
		}
		if err := tSnappyRoundtrip(b, nil, nil); err != nil {
			t.Fatal(err)
		}
	}
}

func TestSnappySmallRegular(t *testing.T) {
	for n := 1; n < 20000; n += 23 {
		b := make([]byte, n)
		for i, _ := range b {
			b[i] = uint8(i%10 + 'a')
		}
		if err := tSnappyRoundtrip(b, nil, nil); err != nil {
			t.Fatal(err)
		}
	}
}

func tSnappyBenchDecode(b *testing.B, src []byte) {
	encoded, err := snappyEncode(nil, src)
	if err != nil {
		b.Fatal(err)
	}
	// Bandwidth is in amount of uncompressed data.
	b.SetBytes(int64(len(src)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		snappyDecode(src, encoded)
	}
}

func tSnappyBenchEncode(b *testing.B, src []byte) {
	// Bandwidth is in amount of uncompressed data.
	b.SetBytes(int64(len(src)))
	dst := make([]byte, snappyMaxEncodedLen(len(src)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		snappyEncode(dst, src)
	}
}

func tSnappyReadFile(b *testing.B, filename string) []byte {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		b.Fatalf("failed reading %s: %s", filename, err)
	}
	if len(src) == 0 {
		b.Fatalf("%s has zero length", filename)
	}
	return src
}

// tSnappyExpand returns a slice of length n containing repeated copies of src.
func tSnappyExpand(src []byte, n int) []byte {
	dst := make([]byte, n)
	for x := dst; len(x) > 0; {
		i := copy(x, src)
		x = x[i:]
	}
	return dst
}

func tSnappyBenchWords(b *testing.B, n int, decode bool) {
	// Note: the file is OS-language dependent so the resulting values are not
	// directly comparable for non-US-English OS installations.
	data := tSnappyExpand(tSnappyReadFile(b, "/usr/share/dict/words"), n)
	if decode {
		tSnappyBenchDecode(b, data)
	} else {
		tSnappyBenchEncode(b, data)
	}
}

func _BenchmarkWordsSnappyDecode1e3(b *testing.B) { tSnappyBenchWords(b, 1e3, true) }
func _BenchmarkWordsSnappyDecode1e4(b *testing.B) { tSnappyBenchWords(b, 1e4, true) }
func _BenchmarkWordsSnappyDecode1e5(b *testing.B) { tSnappyBenchWords(b, 1e5, true) }
func _BenchmarkWordsSnappyDecode1e6(b *testing.B) { tSnappyBenchWords(b, 1e6, true) }
func _BenchmarkWordsSnappyEncode1e3(b *testing.B) { tSnappyBenchWords(b, 1e3, false) }
func _BenchmarkWordsSnappyEncode1e4(b *testing.B) { tSnappyBenchWords(b, 1e4, false) }
func _BenchmarkWordsSnappyEncode1e5(b *testing.B) { tSnappyBenchWords(b, 1e5, false) }
func _BenchmarkWordsSnappyEncode1e6(b *testing.B) { tSnappyBenchWords(b, 1e6, false) }

// tSnappyTestFiles' values are copied directly from
// https://code.google.com/p/snappy/source/browse/trunk/snappy_unittest.cc.
// The label field is unused in snappy-go.
var tSnappyTestFiles = []struct {
	label    string
	filename string
}{
	{"html", "html"},
	{"urls", "urls.10K"},
	{"jpg", "house.jpg"},
	{"pdf", "mapreduce-osdi-1.pdf"},
	{"html4", "html_x_4"},
	{"cp", "cp.html"},
	{"c", "fields.c"},
	{"lsp", "grammar.lsp"},
	{"xls", "kennedy.xls"},
	{"txt1", "alice29.txt"},
	{"txt2", "asyoulik.txt"},
	{"txt3", "lcet10.txt"},
	{"txt4", "plrabn12.txt"},
	{"bin", "ptt5"},
	{"sum", "sum"},
	{"man", "xargs.1"},
	{"pb", "geo.protodata"},
	{"gaviota", "kppkn.gtb"},
}

// The test data files are present at this canonical URL.
const tSnappyBaseURL = "https://snappy.googlecode.com/svn/trunk/testdata/"

func tSnappyDownloadTestdata(basename string) (errRet error) {
	filename := filepath.Join("testdata", basename)
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create %s: %s", filename, err)
	}
	defer f.Close()
	defer func() {
		if errRet != nil {
			os.Remove(filename)
		}
	}()
	resp, err := http.Get(tSnappyBaseURL + basename)
	if err != nil {
		return fmt.Errorf("failed to download %s: %s", tSnappyBaseURL+basename, err)
	}
	defer resp.Body.Close()
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write %s: %s", filename, err)
	}
	return nil
}

func tSnappyBenchFile(b *testing.B, n int, decode bool) {
	filename := filepath.Join("testdata", tSnappyTestFiles[n].filename)
	if stat, err := os.Stat(filename); err != nil || stat.Size() == 0 {
		if !*download {
			b.Fatal("test data not found; skipping benchmark without the -download flag")
		}
		// Download the official snappy C++ implementation reference test data
		// files for benchmarking.
		if err := os.Mkdir("testdata", 0777); err != nil && !os.IsExist(err) {
			b.Fatalf("failed to create testdata: %s", err)
		}
		for _, tf := range tSnappyTestFiles {
			if err := tSnappyDownloadTestdata(tf.filename); err != nil {
				b.Fatalf("failed to download testdata: %s", err)
			}
		}
	}
	data := tSnappyReadFile(b, filename)
	if decode {
		tSnappyBenchDecode(b, data)
	} else {
		tSnappyBenchEncode(b, data)
	}
}

// Naming convention is kept similar to what snappy's C++ implementation uses.
func _Benchmark_Snappy_UFlat0(b *testing.B)  { tSnappyBenchFile(b, 0, true) }
func _Benchmark_Snappy_UFlat1(b *testing.B)  { tSnappyBenchFile(b, 1, true) }
func _Benchmark_Snappy_UFlat2(b *testing.B)  { tSnappyBenchFile(b, 2, true) }
func _Benchmark_Snappy_UFlat3(b *testing.B)  { tSnappyBenchFile(b, 3, true) }
func _Benchmark_Snappy_UFlat4(b *testing.B)  { tSnappyBenchFile(b, 4, true) }
func _Benchmark_Snappy_UFlat5(b *testing.B)  { tSnappyBenchFile(b, 5, true) }
func _Benchmark_Snappy_UFlat6(b *testing.B)  { tSnappyBenchFile(b, 6, true) }
func _Benchmark_Snappy_UFlat7(b *testing.B)  { tSnappyBenchFile(b, 7, true) }
func _Benchmark_Snappy_UFlat8(b *testing.B)  { tSnappyBenchFile(b, 8, true) }
func _Benchmark_Snappy_UFlat9(b *testing.B)  { tSnappyBenchFile(b, 9, true) }
func _Benchmark_Snappy_UFlat10(b *testing.B) { tSnappyBenchFile(b, 10, true) }
func _Benchmark_Snappy_UFlat11(b *testing.B) { tSnappyBenchFile(b, 11, true) }
func _Benchmark_Snappy_UFlat12(b *testing.B) { tSnappyBenchFile(b, 12, true) }
func _Benchmark_Snappy_UFlat13(b *testing.B) { tSnappyBenchFile(b, 13, true) }
func _Benchmark_Snappy_UFlat14(b *testing.B) { tSnappyBenchFile(b, 14, true) }
func _Benchmark_Snappy_UFlat15(b *testing.B) { tSnappyBenchFile(b, 15, true) }
func _Benchmark_Snappy_UFlat16(b *testing.B) { tSnappyBenchFile(b, 16, true) }
func _Benchmark_Snappy_UFlat17(b *testing.B) { tSnappyBenchFile(b, 17, true) }
func _Benchmark_Snappy_ZFlat0(b *testing.B)  { tSnappyBenchFile(b, 0, false) }
func _Benchmark_Snappy_ZFlat1(b *testing.B)  { tSnappyBenchFile(b, 1, false) }
func _Benchmark_Snappy_ZFlat2(b *testing.B)  { tSnappyBenchFile(b, 2, false) }
func _Benchmark_Snappy_ZFlat3(b *testing.B)  { tSnappyBenchFile(b, 3, false) }
func _Benchmark_Snappy_ZFlat4(b *testing.B)  { tSnappyBenchFile(b, 4, false) }
func _Benchmark_Snappy_ZFlat5(b *testing.B)  { tSnappyBenchFile(b, 5, false) }
func _Benchmark_Snappy_ZFlat6(b *testing.B)  { tSnappyBenchFile(b, 6, false) }
func _Benchmark_Snappy_ZFlat7(b *testing.B)  { tSnappyBenchFile(b, 7, false) }
func _Benchmark_Snappy_ZFlat8(b *testing.B)  { tSnappyBenchFile(b, 8, false) }
func _Benchmark_Snappy_ZFlat9(b *testing.B)  { tSnappyBenchFile(b, 9, false) }
func _Benchmark_Snappy_ZFlat10(b *testing.B) { tSnappyBenchFile(b, 10, false) }
func _Benchmark_Snappy_ZFlat11(b *testing.B) { tSnappyBenchFile(b, 11, false) }
func _Benchmark_Snappy_ZFlat12(b *testing.B) { tSnappyBenchFile(b, 12, false) }
func _Benchmark_Snappy_ZFlat13(b *testing.B) { tSnappyBenchFile(b, 13, false) }
func _Benchmark_Snappy_ZFlat14(b *testing.B) { tSnappyBenchFile(b, 14, false) }
func _Benchmark_Snappy_ZFlat15(b *testing.B) { tSnappyBenchFile(b, 15, false) }
func _Benchmark_Snappy_ZFlat16(b *testing.B) { tSnappyBenchFile(b, 16, false) }
func _Benchmark_Snappy_ZFlat17(b *testing.B) { tSnappyBenchFile(b, 17, false) }
