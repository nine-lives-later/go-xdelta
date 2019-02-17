// +build ignore

#include "xdelta.h"
#include "xdelta-error.h"
#include "xdelta-decoder.h"

#ifdef _WIN32
    #define DECLSPEC extern "C" __declspec(dllexport)
    #define DECL __cdecl
#else
    #define DECLSPEC extern "C" 
    #define DECL __cdecl
#endif

DECLSPEC XdeltaError DECL goXdeltaNewDecoder(XdeltaDecoder** ptr)
{
    if (ptr == nullptr)
        return XdeltaError_ArgumentNull;

    *ptr = new XdeltaDecoder();

    return XdeltaError_OK;
}

DECLSPEC XdeltaError DECL goXdeltaFreeDecoder(XdeltaDecoder** ptr)
{
    if (ptr == nullptr)
        return XdeltaError_ArgumentNull;

    if (*ptr == nullptr)
        return XdeltaError_OK;

    delete *ptr;
    *ptr = nullptr;

    return XdeltaError_OK;
}

DECLSPEC XdeltaError DECL goXdeltaDecoderGetStreamError(XdeltaDecoder* ptr, char** str)
{
    if (ptr == nullptr)
        return XdeltaError_ArgumentNull;

    *str = ptr->getStreamError();

    return XdeltaError_OK;
}

DECLSPEC XdeltaError DECL goXdeltaDecoderInit(XdeltaDecoder* ptr, int blockSizeKB, const char* fileId, int hasSource)
{
    if (ptr == nullptr)
        return XdeltaError_ArgumentNull;

    return ptr->init(blockSizeKB, fileId, hasSource);
}

DECLSPEC XdeltaError DECL goXdeltaDecoderGetHeaderRequest(XdeltaDecoder* ptr, int* size)
{
    if (ptr == nullptr)
        return XdeltaError_ArgumentNull;

    return ptr->getHeaderRequest(size);
}

DECLSPEC XdeltaError DECL goXdeltaDecoderCopyHeaderData(XdeltaDecoder* ptr, char* data)
{
    if (ptr == nullptr)
        return XdeltaError_ArgumentNull;

    return ptr->copyHeaderData(data);
}

DECLSPEC XdeltaError DECL goXdeltaDecoderProcess(XdeltaDecoder* ptr, int* state)
{
    if (ptr == nullptr)
        return XdeltaError_ArgumentNull;

    return ptr->process(state);
}

DECLSPEC XdeltaError DECL goXdeltaDecoderProvideInputData(XdeltaDecoder* ptr, const char* data, int length, int finalInput)
{
    if (ptr == nullptr)
        return XdeltaError_ArgumentNull;

    return ptr->provideInputData(data, length, finalInput);
}

DECLSPEC XdeltaError DECL goXdeltaDecoderGetSourceRequest(XdeltaDecoder* ptr, int* block, int* blockSize)
{
    if (ptr == nullptr)
        return XdeltaError_ArgumentNull;

    return ptr->getSourceRequest(block, blockSize);
}

DECLSPEC XdeltaError DECL goXdeltaDecoderProvideSourceData(XdeltaDecoder* ptr, const char* data, int length)
{
    if (ptr == nullptr)
        return XdeltaError_ArgumentNull;

    return ptr->provideSourceData(data, length);
}

DECLSPEC XdeltaError DECL goXdeltaDecoderGetOutputRequest(XdeltaDecoder* ptr, int* size)
{
    if (ptr == nullptr)
        return XdeltaError_ArgumentNull;

    return ptr->getOutputRequest(size);
}

DECLSPEC XdeltaError DECL goXdeltaDecoderCopyOutputData(XdeltaDecoder* ptr, char* data)
{
    if (ptr == nullptr)
        return XdeltaError_ArgumentNull;

    return ptr->copyOutputData(data);
}
