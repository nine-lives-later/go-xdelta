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

	GetStringLength = xdeltaDLL.NewProc("goXdeltaGetStringLength")
	CopyString      = xdeltaDLL.NewProc("goXdeltaCopyString")
	FreeString      = xdeltaDLL.NewProc("goXdeltaFreeString")

	testReturnErrorNotImplemented = xdeltaDLL.NewProc("goXdeltaTestReturnErrorNotImplemented")
)
