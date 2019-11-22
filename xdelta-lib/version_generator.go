// +build ignore

package main

import (
	"fmt"
	"regexp"
	"io/ioutil"
)

func main() {
	mainHeader, _ := ioutil.ReadFile("src/xdelta3/xdelta3-main.h")

	r := regexp.MustCompile("Xdelta version ([0-9.]+),")
	f := r.FindSubmatch(mainHeader)

	version := (string)(f[1])

	fmt.Printf("Version: %v\n", version)

	ioutil.WriteFile("version.go", ([]byte)(fmt.Sprintf("// Code generated version.go DO NOT EDIT.\n\npackage lib\n\nconst Version = \"%v\"", version)), 0644)
}
