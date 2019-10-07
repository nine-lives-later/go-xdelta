// +build !windows

#include "xdelta.h"
#include "xdelta-error.h"
#include "xdelta-encoder.h"
#include "xdelta-decls.h"

extern "C" {

DECLSPEC XdeltaError DECL goXdeltaNewEncoder(XdeltaEncoder** ptr)
{
    if (ptr == nullptr)
        return XdeltaError_ArgumentNull;

    *ptr = new XdeltaEncoder();

    return XdeltaError_OK;
}

DECLSPEC XdeltaError DECL goXdeltaFreeEncoder(XdeltaEncoder** ptr)
{
    if (ptr == nullptr)
        return XdeltaError_ArgumentNull;

    if (*ptr == nullptr)
        return XdeltaError_OK;

    delete *ptr;
    *ptr = nullptr;

    return XdeltaError_OK;
}

DECLSPEC XdeltaError DECL goXdeltaEncoderGetStreamError(XdeltaEncoder* ptr, char** str)
{
    if (ptr == nullptr)
        return XdeltaError_ArgumentNull;

    *str = ptr->getStreamError();

    return XdeltaError_OK;
}

DECLSPEC XdeltaError DECL goXdeltaEncoderInit(XdeltaEncoder* ptr, int blockSizeKB, const char* fileId, int hasSource)
{
    if (ptr == nullptr)
        return XdeltaError_ArgumentNull;

    return ptr->init(blockSizeKB, fileId, hasSource);
}

DECLSPEC XdeltaError DECL goXdeltaEncoderSetHeader(XdeltaEncoder* ptr, const char* data, int length)
{
    if (ptr == nullptr)
        return XdeltaError_ArgumentNull;

    return ptr->setHeader(data, length);
}

DECLSPEC XdeltaError DECL goXdeltaEncoderProcess(XdeltaEncoder* ptr, int* state)
{
    if (ptr == nullptr)
        return XdeltaError_ArgumentNull;

    return ptr->process(state);
}

DECLSPEC XdeltaError DECL goXdeltaEncoderProvideInputData(XdeltaEncoder* ptr, const char* data, int length, int finalInput)
{
    if (ptr == nullptr)
        return XdeltaError_ArgumentNull;

    return ptr->provideInputData(data, length, finalInput);
}

DECLSPEC XdeltaError DECL goXdeltaEncoderGetSourceRequest(XdeltaEncoder* ptr, int* block, int* blockSize)
{
    if (ptr == nullptr)
        return XdeltaError_ArgumentNull;

    return ptr->getSourceRequest(block, blockSize);
}

DECLSPEC XdeltaError DECL goXdeltaEncoderProvideSourceData(XdeltaEncoder* ptr, const char* data, int length)
{
    if (ptr == nullptr)
        return XdeltaError_ArgumentNull;

    return ptr->provideSourceData(data, length);
}

DECLSPEC XdeltaError DECL goXdeltaEncoderGetOutputRequest(XdeltaEncoder* ptr, int* size)
{
    if (ptr == nullptr)
        return XdeltaError_ArgumentNull;

    return ptr->getOutputRequest(size);
}

DECLSPEC XdeltaError DECL goXdeltaEncoderCopyOutputData(XdeltaEncoder* ptr, char* data)
{
    if (ptr == nullptr)
        return XdeltaError_ArgumentNull;

    return ptr->copyOutputData(data);
}

}
