package xdelta

import (
	"context"
	"fmt"
	"io"
	"runtime"
	"time"
	"unsafe"

	"github.com/konsorten/go-xdelta/xdelta-lib"
)

type Decoder struct {
	io.Closer

	handle     unsafe.Pointer
	inputFile  io.Reader
	sourceFile io.ReadSeeker
	outputFile io.Writer
	stats      *Stats

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

	EnableStats bool
}

func NewDecoder(options DecoderOptions) (*Decoder, error) {
	// create the new decoder
	handle, err := lib.NewDecoder()
	if err != nil {
		return nil, err
	}

	// initialize
	if options.BlockSizeKB <= 0 {
		options.BlockSizeKB = (8 * 1024) // 8 MB
	}

	err = lib.DecoderInit(handle, options.BlockSizeKB, options.FileID, options.FromFile != nil)
	if err != nil {
		lib.FreeDecoder(handle)

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

	if options.EnableStats {
		ret.stats = newStats()
	}

	// ensure shutdown
	runtime.SetFinalizer(ret, freeDecoder)

	return ret, nil
}

func (enc *Decoder) GetStreamError() error {
	return lib.DecoderGetStreamError(enc.handle)
}

func (enc *Decoder) DumpStatsToStdout() {
	if enc.stats != nil {
		enc.stats.DumpToStdout()
	}
}

func (enc *Decoder) Process(ctx context.Context) error {
	var isFinal bool
	var perfStart time.Time

	for {
		// prepare gathering stats
		if enc.stats != nil {
			perfStart = time.Now()
		}

		// retrieve the current state
		state, err := lib.DecoderProcess(enc.handle)
		if err != nil {
			return err
		}

		switch state {
		case lib.XdeltaState_INPUT:
			// done?
			if isFinal {
				return nil
			}

			// try read input data
			n, err := enc.inputFile.Read(enc.inputBuffer)
			if err != nil && err != io.EOF {
				return fmt.Errorf("Failed to read from TO/input file: %v", err)
			}

			if n <= 0 { // no more data?
				isFinal = true
			}

			err = lib.DecoderProvideInputData(enc.handle, unsafe.Pointer(&enc.inputBuffer[0]), n, isFinal)
			if err != nil {
				return fmt.Errorf("Failed to provide data from TO/input file: %v", err)
			}
			break

		case lib.XdeltaState_OUTPUT:
			length, err := lib.DecoderGetOutputRequest(enc.handle)
			if err != nil {
				return fmt.Errorf("Failed to request data for PATCH/output file: %v", err)
			}
			if length <= 0 { // nothing to write?
				break
			}
			if length > len(enc.outputBuffer) {
				return fmt.Errorf("Failed to consume data for PATCH/output file: output buffer overflow")
			}

			err = lib.DecoderCopyOutputData(enc.handle, unsafe.Pointer(&enc.outputBuffer[0]))
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

			length, err := lib.DecoderGetHeaderRequest(enc.handle)
			if err != nil {
				return fmt.Errorf("Failed to request header from PATCH/output file: %v", err)
			}
			if length <= 0 { // nothing to write?
				break
			}

			headerData := make([]byte, length)

			err = lib.DecoderCopyHeaderData(enc.handle, unsafe.Pointer(&headerData[0]))
			if err != nil {
				return fmt.Errorf("Failed to consume header from PATCH/output file: %v", err)
			}

			enc.Header <- headerData
			break

		case lib.XdeltaState_GETSRCBLK:
			if enc.sourceFile == nil || enc.sourceBuffer == nil {
				return fmt.Errorf("Failed to request data for FROM/source file: not available")
			}

			blockno, blocksize, err := lib.DecoderGetSourceRequest(enc.handle)
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

			err = lib.DecoderProvideSourceData(enc.handle, unsafe.Pointer(&enc.sourceBuffer[0]), n)
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

		// measure time
		if enc.stats != nil {
			enc.stats.addStateTime(state, time.Since(perfStart))
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

func freeDecoder(enc *Decoder) (err error) {
	// nothing to do?
	if enc == nil || enc.handle == nil {
		return
	}

	err = lib.FreeDecoder(enc.handle)
	enc.handle = nil
	return
}
