package lib

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
	err := callToError(goXdeltaNewEncoder.Call(uintptr(&outHandle)))
	if err != nil {
		return nil, err
	}

	return outHandle, nil
}

func FreeEncoder(handle unsafe.Pointer) error {
	return callToError(goXdeltaFreeEncoder.Call(uintptr(&handle)))
}

func EncoderGetStreamError(handle unsafe.Pointer) error {
	var s uintptr
	err := callToError(goXdeltaEncoderGetStreamError.Call(enc.handle, uintptr(unsafe.Pointer(&s))))
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

	return lib.callToError(goXdeltaEncoderInit.Call(handle, uintptr(options.BlockSizeKB), fromString(options.FileID), hasSourceInt))
}

func EncoderSetHeader(handle unsafe.Pointer, data unsafe.Pointer, dataLen int) error {
	return lib.callToError(goXdeltaEncoderSetHeader.Call(enc.handle, uintptr(data), uintptr(dataLen)))
}

func EncoderProcess(handle unsafe.Pointer) (XdeltaState, error) {
	var state lib.XdeltaState
	err := lib.callToError(goXdeltaEncoderProcess.Call(handle, uintptr(unsafe.Pointer(&state))))
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

	return lib.callToError(lib.goXdeltaEncoderProvideInputData.Call(handle, uintptr(unsafe.Pointer(&data[0])), uintptr(dataLen), isFinal))
}

func EncoderGetOutputRequest(handle unsafe.Pointer) (int, error) {
	var length int
	err := lib.callToError(goXdeltaEncoderGetOutputRequest.Call(handle, uintptr(unsafe.Pointer(&length))))
	if err != nil {
		return -1, err
	}

	return length, nil
}

func EncoderCopyOutputData(handle unsafe.Pointer, destBuffer unsafe.Pointer) error {
	return lib.callToError(lib.goXdeltaEncoderCopyOutputData.Call(handle, uintptr(destBuffer)))
}

func EncoderGetSourceRequest(handle unsafe.Pointer) (int, int, error) {
	var blockno, blocksize int
	err := lib.callToError(goXdeltaEncoderGetSourceRequest.Call(enc.handle, uintptr(unsafe.Pointer(&blockno)), uintptr(unsafe.Pointer(&blocksize))))
	if err != nil {
		return -1, 0, err
	}

	return blockno, blocksize, nil
}

func EncoderProvideSourceData(handle unsafe.Pointer, data unsafe.Pointer, dataLen int) error {
	return lib.callToError(goXdeltaEncoderProvideSourceData.Call(enc.handle, uintptr(data), uintptr(dataLen)))
}

func NewDecoder() (unsafe.Pointer, error) {
	var outHandle unsafe.Pointer
	err := callToError(goXdeltaNewDecoder.Call(uintptr(&outHandle)))
	if err != nil {
		return nil, err
	}

	return outHandle, nil
}

func FreeDecoder(handle unsafe.Pointer) error {
	return callToError(goXdeltaFreeDecoder.Call(uintptr(&handle)))
}

func DecoderGetStreamError(handle unsafe.Pointer) error {
	var s uintptr
	err := callToError(lib.goXdeltaDecoderGetStreamError.Call(enc.handle, uintptr(unsafe.Pointer(&s))))
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

	return lib.callToError(goXdeltaDecoderInit.Call(handle, uintptr(options.BlockSizeKB), fromString(options.FileID), hasSourceInt))
}

func DecoderProcess(handle unsafe.Pointer) (XdeltaState, error) {
	var state lib.XdeltaState
	err := lib.callToError(goXdeltaDecoderProcess.Call(handle, uintptr(unsafe.Pointer(&state))))
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

	return lib.callToError(lib.goXdeltaDecoderProvideInputData.Call(handle, uintptr(unsafe.Pointer(&data[0])), uintptr(dataLen), isFinal))
}

func DecoderGetOutputRequest(handle unsafe.Pointer) (int, error) {
	var length int
	err := lib.callToError(goXdeltaDecoderGetOutputRequest.Call(handle, uintptr(unsafe.Pointer(&length))))
	if err != nil {
		return -1, err
	}

	return length, nil
}

func DecoderCopyOutputData(handle unsafe.Pointer, destBuffer unsafe.Pointer) error {
	return lib.callToError(lib.goXdeltaDecoderCopyOutputData.Call(handle, uintptr(destBuffer)))
}

func DecoderGetHeaderRequest(handle unsafe.Pointer) (int, error) {
	var length int
	err := lib.callToError(goXdeltaDecoderGetHeaderRequest.Call(handle, uintptr(unsafe.Pointer(&length))))
	if err != nil {
		return -1, err
	}

	return length, nil
}

func DecoderCopyHeaderData(handle unsafe.Pointer, destBuffer unsafe.Pointer) error {
	return lib.callToError(goXdeltaDecoderCopyHeaderData.Call(handle, uintptr(destBuffer)))
}

func DecoderGetSourceRequest(handle unsafe.Pointer) (int, int, error) {
	var blockno, blocksize int
	err := lib.callToError(goXdeltaDecoderGetSourceRequest.Call(enc.handle, uintptr(unsafe.Pointer(&blockno)), uintptr(unsafe.Pointer(&blocksize))))
	if err != nil {
		return -1, 0, err
	}

	return blockno, blocksize, nil
}

func DecoderProvideSourceData(handle unsafe.Pointer, data unsafe.Pointer, dataLen int) error {
	return callToError(goXdeltaDecoderProvideSourceData.Call(enc.handle, uintptr(data), uintptr(dataLen)))
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
