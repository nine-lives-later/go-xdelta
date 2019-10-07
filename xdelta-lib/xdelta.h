#ifndef XDELTA_H__345345
#define XDELTA_H__345345

#undef _POSIX_SOURCE
#undef _ISOC99_SOURCE
#undef _C99_SOURCE

#define XD3_USE_LARGEFILE64 1
#define SECONDARY_DJW 0
#define SECONDARY_FGK 0
#define SECONDARY_LZMA 0

#ifdef _DEBUG
    #define XD3_DEBUG 1
#endif

#if defined(_WIN64) || defined(__LP64__)
    #define SIZEOF_SIZE_T 8
    #define SIZEOF_UNSIGNED_LONG_LONG 8
#else
    #define SIZEOF_SIZE_T 4
    #define SIZEOF_UNSIGNED_LONG_LONG 8
#endif

extern "C" {

#include "xdelta3/xdelta3.h"

}

#endif
