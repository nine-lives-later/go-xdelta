package xdelta

import (
	"context"
	"fmt"
	"io"
	"runtime"
	"unsafe"

	"github.com/konsorten/go-xdelta/xdelta-lib"
)

type Decoder struct {
	io.Closer

	handle     uintptr
	inputFile  io.Reader
	sourceFile io.ReadSeeker
	outputFile io.Writer

	inputBuffer  []byte
	sourceBuffer []byte
	outputBuffer []byte

	Header chan<- []byte
}

type DecoderOptions struct {
	BlockSizeKB int
	FileID      string

	FromFile  io.ReadSeeker
	ToFile    io.Writer
	PatchFile io.Reader
}

func NewDecoder(options DecoderOptions) (*Decoder, error) {
	// create the new decoder
	var handle uintptr
	err := lib.CallToError(lib.NewDecoder.Call(uintptr(unsafe.Pointer(&handle))))
	if err != nil {
		return nil, err
	}

	// initialize
	if options.BlockSizeKB <= 0 {
		options.BlockSizeKB = (8 * 1024) // 8 MB
	}

	var hasSource uintptr
	if options.FromFile != nil {
		hasSource = 1
	}

	err = lib.CallToError(lib.DecoderInit.Call(handle, uintptr(options.BlockSizeKB), lib.FromString(options.FileID), hasSource))
	if err != nil {
		lib.FreeDecoder.Call(uintptr(unsafe.Pointer(&handle)))

		return nil, err
	}

	// setup decoder object
	ret := &Decoder{
		handle:       handle,
		inputFile:    options.PatchFile,
		sourceFile:   options.FromFile,
		outputFile:   options.ToFile,
		inputBuffer:  make([]byte, options.BlockSizeKB*1024),
		outputBuffer: make([]byte, options.BlockSizeKB*1024),
	}

	if options.FromFile != nil {
		ret.sourceBuffer = make([]byte, options.BlockSizeKB*1024)
	}

	// ensure shutdown
	runtime.SetFinalizer(ret, freeDecoder)

	return ret, nil
}

func (enc *Decoder) GetStreamError() error {
	var s uintptr
	err := lib.CallToError(lib.DecoderGetStreamError.Call(enc.handle, uintptr(unsafe.Pointer(&s))))
	if err != nil {
		return err
	}
	if s == 0 {
		return nil
	}

	return fmt.Errorf("%v", lib.ToString(s, true))
}

func (enc *Decoder) Process(ctx context.Context) error {
	var isFinal uintptr

	for {
		// retrieve the current state
		var state lib.XdeltaState
		err := lib.CallToError(lib.DecoderProcess.Call(enc.handle, uintptr(unsafe.Pointer(&state))))
		if err != nil {
			return err
		}

		switch state {
		case lib.XdeltaState_INPUT:
			// done?
			if isFinal != 0 {
				return nil
			}

			// try read input data
			n, err := enc.inputFile.Read(enc.inputBuffer)
			if err != nil && err != io.EOF {
				return fmt.Errorf("Failed to read from TO/input file: %v", err)
			}

			if n <= 0 { // no more data?
				isFinal = 1
			}

			err = lib.CallToError(lib.DecoderProvideInputData.Call(enc.handle, uintptr(unsafe.Pointer(&enc.inputBuffer[0])), uintptr(n), isFinal))
			if err != nil {
				return fmt.Errorf("Failed to provide data from TO/input file: %v", err)
			}
			break

		case lib.XdeltaState_OUTPUT:
			var length int
			err := lib.CallToError(lib.DecoderGetOutputRequest.Call(enc.handle, uintptr(unsafe.Pointer(&length))))
			if err != nil {
				return fmt.Errorf("Failed to request data for PATCH/output file: %v", err)
			}
			if length <= 0 { // nothing to write?
				break
			}
			if length > len(enc.outputBuffer) {
				return fmt.Errorf("Failed to consume data for PATCH/output file: output buffer overflow")
			}

			err = lib.CallToError(lib.DecoderCopyOutputData.Call(enc.handle, uintptr(unsafe.Pointer(&enc.outputBuffer[0]))))
			if err != nil {
				return fmt.Errorf("Failed to consume data for PATCH/output file: %v", err)
			}

			written, err := enc.outputFile.Write(enc.outputBuffer[:length])
			if err != nil {
				return fmt.Errorf("Failed to write data to PATCH/output file: %v", err)
			}
			if written < length {
				return fmt.Errorf("Failed to write data to PATCH/output file: not enough data written (%v < %v)", written, length)
			}
			break

		case lib.XdeltaState_GOTHEADER:
			if enc.Header == nil { // nobody there to receive
				break
			}

			var length int
			err := lib.CallToError(lib.DecoderGetHeaderRequest.Call(enc.handle, uintptr(unsafe.Pointer(&length))))
			if err != nil {
				return fmt.Errorf("Failed to request header from PATCH/output file: %v", err)
			}
			if length <= 0 { // nothing to write?
				break
			}

			headerData := make([]byte, length)

			err = lib.CallToError(lib.DecoderCopyHeaderData.Call(enc.handle, uintptr(unsafe.Pointer(&headerData[0]))))
			if err != nil {
				return fmt.Errorf("Failed to consume header from PATCH/output file: %v", err)
			}

			enc.Header <- headerData
			break

		case lib.XdeltaState_GETSRCBLK:
			if enc.sourceFile == nil || enc.sourceBuffer == nil {
				return fmt.Errorf("Failed to request data for FROM/source file: not available")
			}

			var blockno, blocksize int
			err := lib.CallToError(lib.DecoderGetSourceRequest.Call(enc.handle, uintptr(unsafe.Pointer(&blockno)), uintptr(unsafe.Pointer(&blocksize))))
			if err != nil {
				return fmt.Errorf("Failed to request data for FROM/source file: %v", err)
			}
			if blocksize > len(enc.sourceBuffer) {
				return fmt.Errorf("Failed to request data for FROM/source file: source buffer overflow (%v > %v)", blocksize, len(enc.sourceBuffer))
			}

			_, err = enc.sourceFile.Seek(int64(blockno)*int64(blocksize), io.SeekStart)
			if err != nil {
				return fmt.Errorf("Failed to seek FROM/source file: %v", err)
			}

			n, err := enc.sourceFile.Read(enc.sourceBuffer)
			if err != nil {
				return fmt.Errorf("Failed to read from FROM/source file: %v", err)
			}

			err = lib.CallToError(lib.DecoderProvideSourceData.Call(enc.handle, uintptr(unsafe.Pointer(&enc.sourceBuffer[0])), uintptr(n)))
			if err != nil {
				return fmt.Errorf("Failed to provide data from FROM/source file: %v", err)
			}
			break

		case lib.XdeltaState_WINSTART:
		case lib.XdeltaState_WINFINISH:
			break

		default:
			return fmt.Errorf("Unknown state: %v", state)
		}

		// check if cancelled
		err = ctx.Err()
		if err != nil {
			return err
		}
	}
}

func (enc *Decoder) Close() error {
	return freeDecoder(enc)
}

func freeDecoder(enc *Decoder) error {
	// nothing to do?
	if enc == nil || enc.handle == 0 {
		return nil
	}

	// create the new decoder
	return lib.CallToError(lib.FreeDecoder.Call(uintptr(unsafe.Pointer(&enc.handle))))
}
