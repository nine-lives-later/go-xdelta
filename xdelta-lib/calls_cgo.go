// +build cgo

package lib

// #cgo CFLAGS: -I src -D_POSIX_C_SOURCE=200809L -D_XOPEN_SOURCE=700 -DWINVER=0x0601 -D_WIN32_WINNT=0x0601
// #cgo CXXFLAGS: -I src -Wno-literal-suffix -D_POSIX_C_SOURCE=200809L -D_XOPEN_SOURCE=700 -DWINVER=0x0601 -D_WIN32_WINNT=0x0601
//
// #include <stdlib.h>
//
// typedef size_t XdeltaError;
// typedef void XdeltaEncoder;
// typedef void XdeltaDecoder;
//
// #define DECLSPEC extern
// #define DECL
//
// DECLSPEC XdeltaError DECL goXdeltaNewEncoder(XdeltaEncoder** ptr);
// DECLSPEC XdeltaError DECL goXdeltaFreeEncoder(XdeltaEncoder** ptr);
// DECLSPEC XdeltaError DECL goXdeltaEncoderGetStreamError(XdeltaEncoder* ptr, char** str);
// DECLSPEC XdeltaError DECL goXdeltaEncoderInit(XdeltaEncoder* ptr, int blockSizeKB, const char* fileId, int hasSource);
// DECLSPEC XdeltaError DECL goXdeltaEncoderSetHeader(XdeltaEncoder* ptr, const char* data, int length);
// DECLSPEC XdeltaError DECL goXdeltaEncoderProcess(XdeltaEncoder* ptr, int* state);
// DECLSPEC XdeltaError DECL goXdeltaEncoderProvideInputData(XdeltaEncoder* ptr, const char* data, int length, int finalInput);
// DECLSPEC XdeltaError DECL goXdeltaEncoderGetSourceRequest(XdeltaEncoder* ptr, int* block, int* blockSize);
// DECLSPEC XdeltaError DECL goXdeltaEncoderProvideSourceData(XdeltaEncoder* ptr, const char* data, int length);
// DECLSPEC XdeltaError DECL goXdeltaEncoderGetOutputRequest(XdeltaEncoder* ptr, int* size);
// DECLSPEC XdeltaError DECL goXdeltaEncoderCopyOutputData(XdeltaEncoder* ptr, char* data);
//
// DECLSPEC XdeltaError DECL goXdeltaNewDecoder(XdeltaDecoder** ptr);
// DECLSPEC XdeltaError DECL goXdeltaFreeDecoder(XdeltaDecoder** ptr);
// DECLSPEC XdeltaError DECL goXdeltaDecoderGetStreamError(XdeltaDecoder* ptr, char** str);
// DECLSPEC XdeltaError DECL goXdeltaDecoderInit(XdeltaDecoder* ptr, int blockSizeKB, const char* fileId, int hasSource);
// DECLSPEC XdeltaError DECL goXdeltaDecoderGetHeaderRequest(XdeltaDecoder* ptr, int* size);
// DECLSPEC XdeltaError DECL goXdeltaDecoderCopyHeaderData(XdeltaDecoder* ptr, char* data);
// DECLSPEC XdeltaError DECL goXdeltaDecoderProcess(XdeltaDecoder* ptr, int* state);
// DECLSPEC XdeltaError DECL goXdeltaDecoderProvideInputData(XdeltaDecoder* ptr, const char* data, int length, int finalInput);
// DECLSPEC XdeltaError DECL goXdeltaDecoderGetSourceRequest(XdeltaDecoder* ptr, int* block, int* blockSize);
// DECLSPEC XdeltaError DECL goXdeltaDecoderProvideSourceData(XdeltaDecoder* ptr, const char* data, int length);
// DECLSPEC XdeltaError DECL goXdeltaDecoderGetOutputRequest(XdeltaDecoder* ptr, int* size);
// DECLSPEC XdeltaError DECL goXdeltaDecoderCopyOutputData(XdeltaDecoder* ptr, char* data);
//
// DECLSPEC XdeltaError DECL goXdeltaGetStringLength(char* ptr, int* len);
// DECLSPEC XdeltaError DECL goXdeltaCopyString(char* ptr, const char* src, int len);
// DECLSPEC XdeltaError DECL goXdeltaFreeString(char* ptr);
// DECLSPEC XdeltaError DECL goXdeltaTestReturnErrorNotImplemented();
//
// void CFree(char* c) { free(c); }
//
import "C"
import (
	"fmt"
	"unsafe"
)

func NewEncoder() (unsafe.Pointer, error) {
	var outHandle unsafe.Pointer
	err := toError(C.goXdeltaNewEncoder(&outHandle))
	if err != nil {
		return nil, err
	}

	return outHandle, nil
}

func FreeEncoder(handle unsafe.Pointer) error {
	return toError(C.goXdeltaFreeEncoder(&handle))
}

func EncoderGetStreamError(handle unsafe.Pointer) error {
	var str *C.char
	err := toError(C.goXdeltaEncoderGetStreamError(handle, &str))
	if err != nil {
		return err
	}
	if unsafe.Pointer(str) == C.NULL {
		return nil
	}

	defer C.CFree(str)

	return fmt.Errorf("%v", C.GoString(str))
}

func EncoderInit(handle unsafe.Pointer, blockSizeKB int, fileId string, hasSource bool) error {
	fileIdStr := C.CString(fileId)
	defer C.CFree(fileIdStr)

	var hasSourceInt C.int = 0
	if hasSource {
		hasSourceInt = 1
	}

	return toError(C.goXdeltaEncoderInit(handle, C.int(blockSizeKB), fileIdStr, hasSourceInt))
}

func EncoderSetHeader(handle unsafe.Pointer, data unsafe.Pointer, dataLen int) error {
	return toError(C.goXdeltaEncoderSetHeader(handle, (*C.char)(data), C.int(dataLen)))
}

func EncoderProcess(handle unsafe.Pointer) (XdeltaState, error) {
	var s C.int
	err := toError(C.goXdeltaEncoderProcess(handle, &s))
	if err != nil {
		return XdeltaState_SeeGoError, err
	}

	return XdeltaState(s), nil
}

func EncoderProvideInputData(handle unsafe.Pointer, data unsafe.Pointer, dataLen int, isFinal bool) error {
	var isFinalInt C.int = 0
	if isFinal {
		isFinalInt = 1
	}

	return toError(C.goXdeltaEncoderProvideInputData(handle, (*C.char)(data), C.int(dataLen), isFinalInt))
}

func EncoderGetOutputRequest(handle unsafe.Pointer) (int, error) {
	var length C.int
	err := toError(C.goXdeltaEncoderGetOutputRequest(handle, &length))
	if err != nil {
		return 0, err
	}

	return int(length), nil
}

func EncoderCopyOutputData(handle unsafe.Pointer, destBuffer unsafe.Pointer) error {
	return toError(C.goXdeltaEncoderCopyOutputData(handle, (*C.char)(destBuffer)))
}

func EncoderGetSourceRequest(handle unsafe.Pointer) (int, int, error) {
	var blockno, blocksize C.int
	err := toError(C.goXdeltaEncoderGetSourceRequest(handle, &blockno, &blocksize))
	if err != nil {
		return -1, 0, err
	}

	return int(blockno), int(blocksize), nil
}

func EncoderProvideSourceData(handle unsafe.Pointer, data unsafe.Pointer, dataLen int) error {
	err := toError(C.goXdeltaEncoderProvideSourceData(handle, (*C.char)(data), C.int(dataLen)))
	if err != nil {
		return err
	}

	return nil
}

func NewDecoder() (unsafe.Pointer, error) {
	var outHandle unsafe.Pointer
	err := toError(C.goXdeltaNewDecoder(&outHandle))
	if err != nil {
		return nil, err
	}

	return outHandle, nil
}

func FreeDecoder(handle unsafe.Pointer) error {
	return toError(C.goXdeltaFreeDecoder(&handle))
}

func DecoderGetStreamError(handle unsafe.Pointer) error {
	var str *C.char
	err := toError(C.goXdeltaDecoderGetStreamError(handle, &str))
	if err != nil {
		return err
	}
	if unsafe.Pointer(str) == C.NULL {
		return nil
	}

	defer C.CFree(str)

	return fmt.Errorf("%v", C.GoString(str))
}

func DecoderInit(handle unsafe.Pointer, blockSizeKB int, fileId string, hasSource bool) error {
	fileIdStr := C.CString(fileId)
	defer C.CFree(fileIdStr)

	var hasSourceInt C.int = 0
	if hasSource {
		hasSourceInt = 1
	}

	return toError(C.goXdeltaDecoderInit(handle, C.int(blockSizeKB), fileIdStr, hasSourceInt))
}

func DecoderProcess(handle unsafe.Pointer) (XdeltaState, error) {
	var s C.int
	err := toError(C.goXdeltaDecoderProcess(handle, &s))
	if err != nil {
		return XdeltaState_SeeGoError, err
	}

	return XdeltaState(s), nil
}

func DecoderProvideInputData(handle unsafe.Pointer, data unsafe.Pointer, dataLen int, isFinal bool) error {
	var isFinalInt C.int = 0
	if isFinal {
		isFinalInt = 1
	}

	err := toError(C.goXdeltaDecoderProvideInputData(handle, (*C.char)(data), C.int(dataLen), isFinalInt))
	if err != nil {
		return err
	}

	return nil
}

func DecoderGetOutputRequest(handle unsafe.Pointer) (int, error) {
	var length C.int
	err := toError(C.goXdeltaDecoderGetOutputRequest(handle, &length))
	if err != nil {
		return 0, err
	}

	return int(length), nil
}

func DecoderCopyOutputData(handle unsafe.Pointer, destBuffer unsafe.Pointer) error {
	return toError(C.goXdeltaDecoderCopyOutputData(handle, (*C.char)(destBuffer)))
}

func DecoderGetHeaderRequest(handle unsafe.Pointer) (int, error) {
	var length C.int
	err := toError(C.goXdeltaDecoderGetHeaderRequest(handle, &length))
	if err != nil {
		return 0, err
	}

	return int(length), nil
}

func DecoderCopyHeaderData(handle unsafe.Pointer, destBuffer unsafe.Pointer) error {
	return toError(C.goXdeltaDecoderCopyHeaderData(handle, (*C.char)(destBuffer)))
}

func DecoderGetSourceRequest(handle unsafe.Pointer) (int, int, error) {
	var blockno, blocksize C.int
	err := toError(C.goXdeltaDecoderGetSourceRequest(handle, &blockno, &blocksize))
	if err != nil {
		return -1, 0, err
	}

	return int(blockno), int(blocksize), nil
}

func DecoderProvideSourceData(handle unsafe.Pointer, data unsafe.Pointer, dataLen int) error {
	return toError(C.goXdeltaDecoderProvideSourceData(handle, (*C.char)(data), C.int(dataLen)))
}

func toError(ret C.XdeltaError) error {
	r := XdeltaError(ret)

	if r == XdeltaError_OK {
		return nil
	}

	return r
}
