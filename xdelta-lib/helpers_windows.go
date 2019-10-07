package lib

import (
	"fmt"
	"unsafe"
)

func fromString(s string) uintptr {
	b := append([]byte(s), 0)

	return uintptr(unsafe.Pointer(&b[0]))
}

func toString(s uintptr, freeString bool) string {
	if s == 0 {
		return ""
	}

	var l int
	err := callToError(GetStringLength.Call(s, uintptr(unsafe.Pointer(&l))))
	if err != nil {
		return fmt.Sprintf("STRING_ERROR: %v", err)
	}

	b := make([]byte, l)
	err = callToError(CopyString.Call(uintptr(unsafe.Pointer(&b[0])), s, uintptr(l)))
	if err != nil {
		return fmt.Sprintf("STRING_ERROR: %v", err)
	}

	if freeString {
		err = callToError(FreeString.Call(s))
		if err != nil {
			return fmt.Sprintf("STRING_ERROR: %v", err)
		}
	}

	return string(b)
}
