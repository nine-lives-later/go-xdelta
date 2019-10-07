// +build !windows

#include "xdelta.h"
#include "xdelta-encoder.h"

XdeltaError XdeltaEncoder::init(int blockSizeKB, const char* fileId, bool hasSource) {
    if (blockSizeKB < 0)
        return XdeltaError_ArgumentOutOfRange;
    if (fileId == nullptr)
        return XdeltaError_ArgumentNull;
        
    // configure stream
    if (blockSizeKB <= 0)
        blockSizeKB = (8 * 1024); // 8 MB
    
    _config.winsize = blockSizeKB * 1024;
    _config.flags = XD3_NOCOMPRESS;

	auto r = xd3_config_stream(&_stream, &_config);

    if (r != XdeltaError_OK)
        return r;

    // configure source
    _source.ioh = this;  // pass this pointer
    _source.blksize = _config.winsize;
    _source.name = strdup(fileId == nullptr ? "" : fileId);

    if (hasSource) {
        _source.curblkno = -1;
        _source.curblk = nullptr;

        r = xd3_set_source(&_stream, &_source);
    
        if (r != XdeltaError_OK)
            return r;
    }

    return XdeltaError_OK;
}

XdeltaError XdeltaEncoder::setHeader(const char* ptr, int length) {
    if (length < 0)
        return XdeltaError_ArgumentOutOfRange;

    // free old data
    if (_headerData != nullptr) {
        free(_headerData);
        _headerData = nullptr;
    }

    if (ptr == nullptr)
        return XdeltaError_OK;

    // copy new data
    _headerData = malloc(length);
    memcpy(_headerData, ptr, length);

    xd3_set_appheader(&_stream, (const uint8_t*)_headerData, length);

    return XdeltaError_OK;
}

char* XdeltaEncoder::getStreamError() const {
    if (_stream.msg == nullptr)
        return nullptr;

    const auto len = 1024;
    auto s = (char*)malloc(len);

    snprintf(s, len, "STREAM_ERROR: %s", xd3_errstring(&_stream));

    return s;
}

XdeltaError XdeltaEncoder::process(XdeltaState* state) {
    if (state == nullptr)
        return XdeltaError_ArgumentNull;

    auto r = xd3_encode_input(&_stream);

    *state = r;

    if (isXdeltaStateError(r))
        return r;

    return XdeltaError_OK;
}

XdeltaError XdeltaEncoder::provideInputData(const char* ptr, int length, bool finalInput) {
    if (ptr == nullptr)
        return XdeltaError_ArgumentNull;
    if (length < 0)
        return XdeltaError_ArgumentOutOfRange;

    if (finalInput)
        _stream.flags |= XD3_FLUSH;

    xd3_avail_input(&_stream, (uint8_t*)ptr, length);

    return XdeltaError_OK;
}

XdeltaError XdeltaEncoder::getSourceRequest(int* block, int* blockSize) {
    if (block == nullptr)
        return XdeltaError_ArgumentNull;
    if (blockSize == nullptr)
        return XdeltaError_ArgumentNull;

    *block = _source.getblkno;
    *blockSize = _source.blksize;

    return XdeltaError_OK;
}

XdeltaError XdeltaEncoder::provideSourceData(const char* ptr, int length) {
    if (ptr == nullptr)
        return XdeltaError_ArgumentNull;
    if (length < 0)
        return XdeltaError_ArgumentOutOfRange;

    _source.curblkno = _source.getblkno;
    _source.curblk = (uint8_t*)ptr;
    _source.onblk = length;

    return XdeltaError_OK;
}

XdeltaError XdeltaEncoder::getOutputRequest(int* size) {
    if (size == nullptr)
        return XdeltaError_ArgumentNull;

    *size = _stream.avail_out;

    return XdeltaError_OK;
}

XdeltaError XdeltaEncoder::copyOutputData(char* ptr) {
    if (ptr == nullptr)
        return XdeltaError_ArgumentNull;

    memcpy(ptr, _stream.next_out, _stream.avail_out);

    xd3_consume_output(&_stream);

    return XdeltaError_OK;
}

XdeltaEncoder::XdeltaEncoder() {
    memset(&_config, 0, sizeof(_config));
    memset(&_stream, 0, sizeof(_stream));
    memset(&_source, 0, sizeof(_source));
}

XdeltaEncoder::~XdeltaEncoder() {
    xd3_close_stream(&_stream);
    xd3_free_stream(&_stream);

    if (_source.name != nullptr) {
        free((void*)_source.name);
        _source.name = nullptr;
    }

    if (_headerData != nullptr)
        free(_headerData);
}
