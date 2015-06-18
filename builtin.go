// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rawp

import (
	"fmt"
)

func errorf(format string, a ...interface{}) error {
	return fmt.Errorf(format, a...)
}

func print(a ...interface{}) (n int, err error) {
	return fmt.Print(a...)
}

func printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Printf(format, a...)
}

func println(a ...interface{}) (n int, err error) {
	return fmt.Println(a...)
}

func scan(a ...interface{}) (n int, err error) {
	return fmt.Scan(a...)
}

func scanf(format string, a ...interface{}) (n int, err error) {
	return fmt.Scanf(format, a...)
}

func scanln(a ...interface{}) (n int, err error) {
	return fmt.Scanln(a...)
}

func sprint(a ...interface{}) string {
	return fmt.Sprint(a...)
}

func sprintf(format string, a ...interface{}) string {
	return fmt.Sprintf(format, a...)
}

func sprintln(a ...interface{}) string {
	return fmt.Sprintln(a...)
}

func sscan(str string, a ...interface{}) (n int, err error) {
	return fmt.Sscan(str, a...)
}

func sscanf(str string, format string, a ...interface{}) (n int, err error) {
	return fmt.Sscanf(str, format, a...)
}

func sscanln(str string, a ...interface{}) (n int, err error) {
	return fmt.Sscanln(str, a...)
}
