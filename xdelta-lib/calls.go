package lib

var (
	NewEncoder               = xdeltaDLL.NewProc("goXdeltaNewEncoder")
	FreeEncoder              = xdeltaDLL.NewProc("goXdeltaFreeEncoder")
	EncoderGetStreamError    = xdeltaDLL.NewProc("goXdeltaEncoderGetStreamError")
	EncoderInit              = xdeltaDLL.NewProc("goXdeltaEncoderInit")
	EncoderSetHeader         = xdeltaDLL.NewProc("goXdeltaEncoderSetHeader")
	EncoderProcess           = xdeltaDLL.NewProc("goXdeltaEncoderProcess")
	EncoderProvideInputData  = xdeltaDLL.NewProc("goXdeltaEncoderProvideInputData")
	EncoderGetSourceRequest  = xdeltaDLL.NewProc("goXdeltaEncoderGetSourceRequest")
	EncoderProvideSourceData = xdeltaDLL.NewProc("goXdeltaEncoderProvideSourceData")
	EncoderGetOutputRequest  = xdeltaDLL.NewProc("goXdeltaEncoderGetOutputRequest")
	EncoderCopyOutputData    = xdeltaDLL.NewProc("goXdeltaEncoderCopyOutputData")

	NewDecoder               = xdeltaDLL.NewProc("goXdeltaNewDecoder")
	FreeDecoder              = xdeltaDLL.NewProc("goXdeltaFreeDecoder")
	DecoderGetStreamError    = xdeltaDLL.NewProc("goXdeltaDecoderGetStreamError")
	DecoderInit              = xdeltaDLL.NewProc("goXdeltaDecoderInit")
	DecoderGetHeaderRequest  = xdeltaDLL.NewProc("goXdeltaDecoderGetHeaderRequest")
	DecoderCopyHeaderData    = xdeltaDLL.NewProc("goXdeltaDecoderCopyHeaderData")
	DecoderProcess           = xdeltaDLL.NewProc("goXdeltaDecoderProcess")
	DecoderProvideInputData  = xdeltaDLL.NewProc("goXdeltaDecoderProvideInputData")
	DecoderGetSourceRequest  = xdeltaDLL.NewProc("goXdeltaDecoderGetSourceRequest")
	DecoderProvideSourceData = xdeltaDLL.NewProc("goXdeltaDecoderProvideSourceData")
	DecoderGetOutputRequest  = xdeltaDLL.NewProc("goXdeltaDecoderGetOutputRequest")
	DecoderCopyOutputData    = xdeltaDLL.NewProc("goXdeltaDecoderCopyOutputData")

	GetStringLength = xdeltaDLL.NewProc("goXdeltaGetStringLength")
	CopyString      = xdeltaDLL.NewProc("goXdeltaCopyString")
	FreeString      = xdeltaDLL.NewProc("goXdeltaFreeString")

	testReturnErrorNotImplemented = xdeltaDLL.NewProc("goXdeltaTestReturnErrorNotImplemented")
)
