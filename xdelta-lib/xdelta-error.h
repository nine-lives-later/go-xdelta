#ifndef _XDELTA_ERROR_39847563543
#define _XDELTA_ERROR_39847563543

typedef size_t XdeltaError;
typedef int XdeltaState;

enum {
	/* go-xdelta errors: */
	XdeltaError_OK                  = 0,

	XdeltaError_ArgumentNull = 101,
	XdeltaError_ArgumentOutOfRange = 102,

	/* native Xdelta errors: */
	XdeltaError_Input               = XD3_INPUT,
	XdeltaError_Output              = XD3_OUTPUT,
	XdeltaError_GetSourceBlock      = XD3_GETSRCBLK,
	XdeltaError_GotHeader           = XD3_GOTHEADER,
	XdeltaError_WindowStart         = XD3_WINSTART,
	XdeltaError_WindowFinish        = XD3_WINFINISH,
	XdeltaError_TooFarBack          = XD3_TOOFARBACK,
	XdeltaError_Internal            = XD3_INTERNAL,
	XdeltaError_InvalidConfig       = XD3_INVALID,
	XdeltaError_InvalidInput        = XD3_INVALID_INPUT,
	XdeltaError_NoSecondCompression = XD3_NOSECOND,
	XdeltaError_NotImplemented      = XD3_UNIMPLEMENTED,
};

inline bool isXdeltaStateError(XdeltaState s) {
	switch (s) {
	case XD3_INPUT:
	case XD3_OUTPUT:
	case XD3_GETSRCBLK:
	case XD3_GOTHEADER:
	case XD3_WINSTART:
	case XD3_WINFINISH:
		return false;
	}

	return true;
}

#endif
