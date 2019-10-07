// build +windows

package lib

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	xdeltaDLL = syscall.NewLazyDLL(`.\go-xdelta-lib.dll`)

	goXdeltaNewEncoder               = xdeltaDLL.NewProc("goXdeltaNewEncoder")
	goXdeltaFreeEncoder              = xdeltaDLL.NewProc("goXdeltaFreeEncoder")
	goXdeltaEncoderGetStreamError    = xdeltaDLL.NewProc("goXdeltaEncoderGetStreamError")
	goXdeltaEncoderInit              = xdeltaDLL.NewProc("goXdeltaEncoderInit")
	goXdeltaEncoderSetHeader         = xdeltaDLL.NewProc("goXdeltaEncoderSetHeader")
	goXdeltaEncoderProcess           = xdeltaDLL.NewProc("goXdeltaEncoderProcess")
	goXdeltaEncoderProvideInputData  = xdeltaDLL.NewProc("goXdeltaEncoderProvideInputData")
	goXdeltaEncoderGetSourceRequest  = xdeltaDLL.NewProc("goXdeltaEncoderGetSourceRequest")
	goXdeltaEncoderProvideSourceData = xdeltaDLL.NewProc("goXdeltaEncoderProvideSourceData")
	goXdeltaEncoderGetOutputRequest  = xdeltaDLL.NewProc("goXdeltaEncoderGetOutputRequest")
	goXdeltaEncoderCopyOutputData    = xdeltaDLL.NewProc("goXdeltaEncoderCopyOutputData")

	goXdeltaNewDecoder               = xdeltaDLL.NewProc("goXdeltaNewDecoder")
	goXdeltaFreeDecoder              = xdeltaDLL.NewProc("goXdeltaFreeDecoder")
	goXdeltaDecoderGetStreamError    = xdeltaDLL.NewProc("goXdeltaDecoderGetStreamError")
	goXdeltaDecoderInit              = xdeltaDLL.NewProc("goXdeltaDecoderInit")
	goXdeltaDecoderGetHeaderRequest  = xdeltaDLL.NewProc("goXdeltaDecoderGetHeaderRequest")
	goXdeltaDecoderCopyHeaderData    = xdeltaDLL.NewProc("goXdeltaDecoderCopyHeaderData")
	goXdeltaDecoderProcess           = xdeltaDLL.NewProc("goXdeltaDecoderProcess")
	goXdeltaDecoderProvideInputData  = xdeltaDLL.NewProc("goXdeltaDecoderProvideInputData")
	goXdeltaDecoderGetSourceRequest  = xdeltaDLL.NewProc("goXdeltaDecoderGetSourceRequest")
	goXdeltaDecoderProvideSourceData = xdeltaDLL.NewProc("goXdeltaDecoderProvideSourceData")
	goXdeltaDecoderGetOutputRequest  = xdeltaDLL.NewProc("goXdeltaDecoderGetOutputRequest")
	goXdeltaDecoderCopyOutputData    = xdeltaDLL.NewProc("goXdeltaDecoderCopyOutputData")

	goXdeltaGetStringLength = xdeltaDLL.NewProc("goXdeltaGetStringLength")
	goXdeltaCopyString      = xdeltaDLL.NewProc("goXdeltaCopyString")
	goXdeltaFreeString      = xdeltaDLL.NewProc("goXdeltaFreeString")

	goXdeltaTestReturnErrorNotImplemented = xdeltaDLL.NewProc("goXdeltaTestReturnErrorNotImplemented")
)

func NewEncoder() (unsafe.Pointer, error) {
	var outHandle unsafe.Pointer
	err := callToError(goXdeltaNewEncoder.Call(uintptr(unsafe.Pointer(&outHandle))))
	if err != nil {
		return nil, err
	}

	return outHandle, nil
}

func FreeEncoder(handle unsafe.Pointer) error {
	return callToError(goXdeltaFreeEncoder.Call(uintptr(unsafe.Pointer(&handle))))
}

func EncoderGetStreamError(handle unsafe.Pointer) error {
	var s uintptr
	err := callToError(goXdeltaEncoderGetStreamError.Call(uintptr(handle), uintptr(unsafe.Pointer(&s))))
	if err != nil {
		return err
	}
	if s == 0 {
		return nil
	}

	return fmt.Errorf("%v", toString(s, true))
}

func EncoderInit(handle unsafe.Pointer, blockSizeKB int, fileId string, hasSource bool) error {
	var hasSourceInt uintptr = 0
	if hasSource {
		hasSourceInt = 1
	}

	return callToError(goXdeltaEncoderInit.Call(uintptr(handle), uintptr(blockSizeKB), fromString(fileId), hasSourceInt))
}

func EncoderSetHeader(handle unsafe.Pointer, data unsafe.Pointer, dataLen int) error {
	return callToError(goXdeltaEncoderSetHeader.Call(uintptr(handle), uintptr(data), uintptr(dataLen)))
}

func EncoderProcess(handle unsafe.Pointer) (XdeltaState, error) {
	var state XdeltaState
	err := callToError(goXdeltaEncoderProcess.Call(uintptr(handle), uintptr(unsafe.Pointer(&state))))
	if err != nil {
		return XdeltaState_SeeGoError, err
	}

	return state, nil
}

func EncoderProvideInputData(handle unsafe.Pointer, data unsafe.Pointer, dataLen int, isFinal bool) error {
	var isFinalInt uintptr = 0
	if isFinal {
		isFinalInt = 1
	}

	return callToError(goXdeltaEncoderProvideInputData.Call(uintptr(handle), uintptr(data), uintptr(dataLen), isFinalInt))
}

func EncoderGetOutputRequest(handle unsafe.Pointer) (int, error) {
	var length int
	err := callToError(goXdeltaEncoderGetOutputRequest.Call(uintptr(handle), uintptr(unsafe.Pointer(&length))))
	if err != nil {
		return -1, err
	}

	return length, nil
}

func EncoderCopyOutputData(handle unsafe.Pointer, destBuffer unsafe.Pointer) error {
	return callToError(goXdeltaEncoderCopyOutputData.Call(uintptr(handle), uintptr(destBuffer)))
}

func EncoderGetSourceRequest(handle unsafe.Pointer) (int, int, error) {
	var blockno, blocksize int
	err := callToError(goXdeltaEncoderGetSourceRequest.Call(uintptr(handle), uintptr(unsafe.Pointer(&blockno)), uintptr(unsafe.Pointer(&blocksize))))
	if err != nil {
		return -1, 0, err
	}

	return blockno, blocksize, nil
}

func EncoderProvideSourceData(handle unsafe.Pointer, data unsafe.Pointer, dataLen int) error {
	return callToError(goXdeltaEncoderProvideSourceData.Call(uintptr(handle), uintptr(data), uintptr(dataLen)))
}

func NewDecoder() (unsafe.Pointer, error) {
	var outHandle unsafe.Pointer
	err := callToError(goXdeltaNewDecoder.Call(uintptr(unsafe.Pointer(&outHandle))))
	if err != nil {
		return nil, err
	}

	return outHandle, nil
}

func FreeDecoder(handle unsafe.Pointer) error {
	return callToError(goXdeltaFreeDecoder.Call(uintptr(unsafe.Pointer(&handle))))
}

func DecoderGetStreamError(handle unsafe.Pointer) error {
	var s uintptr
	err := callToError(goXdeltaDecoderGetStreamError.Call(uintptr(handle), uintptr(unsafe.Pointer(&s))))
	if err != nil {
		return err
	}
	if s == 0 {
		return nil
	}

	return fmt.Errorf("%v", toString(s, true))
}

func DecoderInit(handle unsafe.Pointer, blockSizeKB int, fileId string, hasSource bool) error {
	var hasSourceInt uintptr = 0
	if hasSource {
		hasSourceInt = 1
	}

	return callToError(goXdeltaDecoderInit.Call(uintptr(handle), uintptr(blockSizeKB), fromString(fileId), hasSourceInt))
}

func DecoderProcess(handle unsafe.Pointer) (XdeltaState, error) {
	var state XdeltaState
	err := callToError(goXdeltaDecoderProcess.Call(uintptr(handle), uintptr(unsafe.Pointer(&state))))
	if err != nil {
		return XdeltaState_SeeGoError, err
	}

	return state, nil
}

func DecoderProvideInputData(handle unsafe.Pointer, data unsafe.Pointer, dataLen int, isFinal bool) error {
	var isFinalInt uintptr = 0
	if isFinal {
		isFinalInt = 1
	}

	return callToError(goXdeltaDecoderProvideInputData.Call(uintptr(handle), uintptr(data), uintptr(dataLen), isFinalInt))
}

func DecoderGetOutputRequest(handle unsafe.Pointer) (int, error) {
	var length int
	err := callToError(goXdeltaDecoderGetOutputRequest.Call(uintptr(handle), uintptr(unsafe.Pointer(&length))))
	if err != nil {
		return -1, err
	}

	return length, nil
}

func DecoderCopyOutputData(handle unsafe.Pointer, destBuffer unsafe.Pointer) error {
	return callToError(goXdeltaDecoderCopyOutputData.Call(uintptr(handle), uintptr(destBuffer)))
}

func DecoderGetHeaderRequest(handle unsafe.Pointer) (int, error) {
	var length int
	err := callToError(goXdeltaDecoderGetHeaderRequest.Call(uintptr(handle), uintptr(unsafe.Pointer(&length))))
	if err != nil {
		return -1, err
	}

	return length, nil
}

func DecoderCopyHeaderData(handle unsafe.Pointer, destBuffer unsafe.Pointer) error {
	return callToError(goXdeltaDecoderCopyHeaderData.Call(uintptr(handle), uintptr(destBuffer)))
}

func DecoderGetSourceRequest(handle unsafe.Pointer) (int, int, error) {
	var blockno, blocksize int
	err := callToError(goXdeltaDecoderGetSourceRequest.Call(uintptr(handle), uintptr(unsafe.Pointer(&blockno)), uintptr(unsafe.Pointer(&blocksize))))
	if err != nil {
		return -1, 0, err
	}

	return blockno, blocksize, nil
}

func DecoderProvideSourceData(handle unsafe.Pointer, data unsafe.Pointer, dataLen int) error {
	return callToError(goXdeltaDecoderProvideSourceData.Call(uintptr(handle), uintptr(data), uintptr(dataLen)))
}

func toError(ret uintptr) error {
	r := XdeltaError(ret)

	if r == XdeltaError_OK {
		return nil
	}

	return r
}

func callToError(r1, r2 uintptr, err error) error {
	if err != nil {
		if errNo, ok := err.(syscall.Errno); ok {
			if errNo == 0 /* no error */ {
				return toError(r1)
			}
		}
		return err
	}

	return toError(r1)
}
