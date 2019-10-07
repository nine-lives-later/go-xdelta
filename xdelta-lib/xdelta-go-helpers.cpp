#include "xdelta.h"
#include "xdelta-error.h"
#include "xdelta-encoder.h"
#include "xdelta-decls.h"

extern "C" {

DECLSPEC XdeltaError DECL goXdeltaGetStringLength(char* ptr, int* len)
{
    if (ptr == nullptr)
        return XdeltaError_ArgumentNull;
    if (len == nullptr)
        return XdeltaError_ArgumentNull;

    *len = strlen(ptr);

    return XdeltaError_OK;
}

DECLSPEC XdeltaError DECL goXdeltaCopyString(char* ptr, const char* src, int len)
{
    if (ptr == nullptr)
        return XdeltaError_ArgumentNull;
    if (src == nullptr)
        return XdeltaError_ArgumentNull;
    if (len < 0)
        return XdeltaError_ArgumentOutOfRange;
    if (len == 0)
        return XdeltaError_OK;

    memcpy(ptr, src, len);

    return XdeltaError_OK;
}

DECLSPEC XdeltaError DECL goXdeltaFreeString(char* ptr)
{
    if (ptr == nullptr)
        return XdeltaError_ArgumentNull;

    free(ptr);

    return XdeltaError_OK;
}

DECLSPEC XdeltaError DECL goXdeltaTestReturnErrorNotImplemented()
{
    return XD3_UNIMPLEMENTED;
}

}
