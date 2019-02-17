#ifndef _XDELTA_ENCODER_H_34535435
#define _XDELTA_ENCODER_H_34535435

#include "xdelta-error.h"

class XdeltaEncoder {
private:

    xd3_config _config;
    mutable xd3_stream _stream;
    xd3_source _source;
    void* _headerData = nullptr;

public:

    // free the returned string!
    char* getStreamError() const;

    XdeltaError init(int blockSizeKB, const char* fileId, bool hasSource);
    XdeltaError setHeader(const char* ptr, int length);

    XdeltaError process(XdeltaState* state);
    
    XdeltaError provideInputData(const char* ptr, int length, bool finalInput);

    XdeltaError getSourceRequest(int* block, int* blockSize);
    XdeltaError provideSourceData(const char* ptr, int length);

    XdeltaError getOutputRequest(int* size);
    XdeltaError copyOutputData(char* ptr);

public:

    XdeltaEncoder& operator = (const XdeltaEncoder&) = delete;
    XdeltaEncoder(const XdeltaEncoder&) = delete;

    XdeltaEncoder();
    ~XdeltaEncoder();

};

#endif
