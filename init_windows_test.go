package xdelta

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const (
	xdeltaLibDLLSource = "./xdelta-lib/go-xdelta-lib.dll"
)

func TestMain(m *testing.M) {
	dll, err := ioutil.ReadFile(xdeltaLibDLLSource)
	if err != nil {
		panic(fmt.Errorf("Failed to read xdelta-lib DLL"))
	}
	err = ioutil.WriteFile(filepath.Base(xdeltaLibDLLSource), dll, 0644)
	if err != nil {
		panic(fmt.Errorf("Failed to write xdelta-lib DLL"))
	}

	os.Exit(m.Run())
}
