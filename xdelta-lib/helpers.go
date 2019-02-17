package lib

import (
	"fmt"
	"unsafe"
)

func FromString(s string) uintptr {
	b := append([]byte(s), 0)

	return uintptr(unsafe.Pointer(&b[0]))
}

func ToString(s uintptr, freeString bool) string {
	if s == 0 {
		return ""
	}

	var l int
	err := CallToError(GetStringLength.Call(s, uintptr(unsafe.Pointer(&l))))
	if err != nil {
		return fmt.Sprintf("STRING_ERROR: %v", err)
	}

	b := make([]byte, l)
	err = CallToError(CopyString.Call(uintptr(unsafe.Pointer(&b[0])), s, uintptr(l)))
	if err != nil {
		return fmt.Sprintf("STRING_ERROR: %v", err)
	}

	if freeString {
		err = CallToError(FreeString.Call(s))
		if err != nil {
			return fmt.Sprintf("STRING_ERROR: %v", err)
		}
	}

	return string(b)
}
