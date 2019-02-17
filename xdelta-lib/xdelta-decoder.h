#ifndef _XDELTA_DECODER_H_34535435
#define _XDELTA_DECODER_H_34535435

#include "xdelta-error.h"

class XdeltaDecoder {
private:

    xd3_config _config;
    mutable xd3_stream _stream;
    xd3_source _source;

public:

    // free the returned string!
    char* getStreamError() const;

    XdeltaError init(int blockSizeKB, const char* fileId, bool hasSource);
    XdeltaError getHeaderRequest(int* size);
    XdeltaError copyHeaderData(char* ptr);

    XdeltaError process(XdeltaState* state);
    
    XdeltaError provideInputData(const char* ptr, int length, bool finalInput);

    XdeltaError getSourceRequest(int* block, int* blockSize);
    XdeltaError provideSourceData(const char* ptr, int length);

    XdeltaError getOutputRequest(int* size);
    XdeltaError copyOutputData(char* ptr);

public:

    XdeltaDecoder& operator = (const XdeltaDecoder&) = delete;
    XdeltaDecoder(const XdeltaDecoder&) = delete;

    XdeltaDecoder();
    ~XdeltaDecoder();

};

#endif
