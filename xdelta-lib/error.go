package lib

import (
	"fmt"
)

type XdeltaState int32

const (
	XdeltaState_INPUT         XdeltaState = -17703 /* need input */
	XdeltaState_OUTPUT        XdeltaState = -17704 /* have output */
	XdeltaState_GETSRCBLK     XdeltaState = -17705 /* need a block of source input */
	XdeltaState_GOTHEADER     XdeltaState = -17706 /* (decode-only) after the initial VCDIFF & first window header */
	XdeltaState_WINSTART      XdeltaState = -17707 /* notification: returned before a window is processed */
	XdeltaState_WINFINISH     XdeltaState = -17708 /* notification: returned after encode/decode & output for a window */
	XdeltaState_TOOFARBACK    XdeltaState = -17709 /* (encoder only) may be returned by getblk() if the block is too old */
	XdeltaState_INTERNAL      XdeltaState = -17710 /* internal error */
	XdeltaState_INVALID       XdeltaState = -17711 /* invalid config */
	XdeltaState_INVALID_INPUT XdeltaState = -17712 /* invalid input/decoder error */
	XdeltaState_NOSECOND      XdeltaState = -17713 /* when secondary compression finds no improvement. */
	XdeltaState_UNIMPLEMENTED XdeltaState = -17714 /* currently VCD_TARGET, VCD_CODETABLE */

	XdeltaState_SeeGoError XdeltaState = -17800 /* an error happend, see the returned Go error */
)


func (e XdeltaState) String() string {
	switch e {
	case XdeltaState_INPUT:
		return "INPUT"
	case XdeltaState_OUTPUT:
		return "OUTPUT"
	case XdeltaState_GETSRCBLK:
		return "GETSRCBLK"
	case XdeltaState_GOTHEADER:
		return "GOTHEADER"
	case XdeltaState_WINSTART:
		return "WINSTART"
	case XdeltaState_WINFINISH:
		return "WINFINISH"
	}

	return fmt.Sprintf("Unknown state: %v", int(e))
}

type XdeltaError int32

const (
	XdeltaError_OK                  XdeltaError = XdeltaError(0)
	XdeltaError_ArgumentNull        XdeltaError = XdeltaError(101)
	XdeltaError_ArgumentOutOfRange  XdeltaError = XdeltaError(102)
	XdeltaError_Input               XdeltaError = XdeltaError(XdeltaState_INPUT)
	XdeltaError_Output              XdeltaError = XdeltaError(XdeltaState_OUTPUT)
	XdeltaError_GetSourceBlock      XdeltaError = XdeltaError(XdeltaState_GETSRCBLK)
	XdeltaError_GotHeader           XdeltaError = XdeltaError(XdeltaState_GOTHEADER)
	XdeltaError_WindowStart         XdeltaError = XdeltaError(XdeltaState_WINSTART)
	XdeltaError_WindowFinish        XdeltaError = XdeltaError(XdeltaState_WINFINISH)
	XdeltaError_TooFarBack          XdeltaError = XdeltaError(XdeltaState_TOOFARBACK)
	XdeltaError_Internal            XdeltaError = XdeltaError(XdeltaState_INTERNAL)
	XdeltaError_InvalidConfig       XdeltaError = XdeltaError(XdeltaState_INVALID)
	XdeltaError_InvalidInput        XdeltaError = XdeltaError(XdeltaState_INVALID_INPUT)
	XdeltaError_NoSecondCompression XdeltaError = XdeltaError(XdeltaState_NOSECOND)
	XdeltaError_NotImplemented      XdeltaError = XdeltaError(XdeltaState_UNIMPLEMENTED)
	XdeltaError_SeeGoError    XdeltaError = XdeltaError(XdeltaState_SeeGoError)
)

func (e XdeltaError) Error() string {
	switch e {
	case XdeltaError_OK:
		return "OK"
	case XdeltaError_ArgumentNull:
		return "argument is null/nil"
	case XdeltaError_ArgumentOutOfRange:
		return "argument is out of range"
	case XdeltaError_Input:
		return "need input"
	case XdeltaError_Output:
		return "have output"
	case XdeltaError_GetSourceBlock:
		return "need a block of source input"
	case XdeltaError_GotHeader:
		return "after the initial VCDIFF & first window header"
	case XdeltaError_WindowStart:
		return "before a window is processed"
	case XdeltaError_WindowFinish:
		return "after encode/decode & output for a window"
	case XdeltaError_TooFarBack:
		return "block is too old"
	case XdeltaError_Internal:
		return "internal error"
	case XdeltaError_InvalidConfig:
		return "invalid config"
	case XdeltaError_InvalidInput:
		return "invalid input/decoder error"
	case XdeltaError_NoSecondCompression:
		return "secondary compression finds no improvement"
	case XdeltaError_NotImplemented:
		return "not implemented (VCD_TARGET, VCD_CODETABLE)"
	case XdeltaError_SeeGoError:
		return "see Go returned error"
	}

	return fmt.Sprintf("Unknown error code: %v", int(e))
}
