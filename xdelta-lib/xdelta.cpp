// +build !windows

// disable some warnings
#ifdef _MSC_VER
	#pragma warning( disable : 4100 )   // unused parameters
	#pragma warning( disable : 4244 )   // possible loss of data
    #pragma warning( disable : 4267 )   // possible loss of data
#endif

#ifdef __APPLE__
    #pragma clang diagnostic ignored "-Wunused-parameter"
	#pragma clang diagnostic ignored "-Wreserved-user-defined-literal"
#endif

#ifdef __GNUC__
    #pragma GCC diagnostic ignored "-Wunused-parameter"
#endif

#include "xdelta.h"

extern "C" {

#include "xdelta3/xdelta3.c"

}
