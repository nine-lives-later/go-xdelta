package xdelta

import (
	"context"
	"fmt"
	"io"
	"runtime"
	"time"
	"unsafe"

	"github.com/nine-lives-later/go-xdelta/xdelta-lib"
)

type Encoder struct {
	io.Closer

	handle     unsafe.Pointer
	inputFile  io.Reader
	sourceFile io.ReadSeeker
	patchFile  io.Writer
	stats      *Stats

	inputBuffer  []byte
	sourceBuffer []byte
	outputBuffer []byte
}

type EncoderOptions struct {
	BlockSizeKB int
	FileID      string

	FromFile  io.ReadSeeker
	ToFile    io.Reader
	PatchFile io.Writer

	Header []byte

	EnableStats bool
}

func NewEncoder(options EncoderOptions) (*Encoder, error) {
	// create the new encoder
	handle, err := lib.NewEncoder()
	if err != nil {
		return nil, err
	}

	// initialize
	if options.BlockSizeKB <= 0 {
		options.BlockSizeKB = (8 * 1024) // 8 MB
	}

	err = lib.EncoderInit(handle, options.BlockSizeKB, options.FileID, options.FromFile != nil)
	if err != nil {
		lib.FreeEncoder(handle)

		return nil, err
	}

	// set header
	if options.Header != nil {
		err = lib.EncoderSetHeader(handle, unsafe.Pointer(&options.Header[0]), len(options.Header))
		if err != nil {
			lib.FreeEncoder(handle)

			return nil, err
		}
	}

	// setup encoder object
	ret := &Encoder{
		handle:       handle,
		inputFile:    options.ToFile,
		sourceFile:   options.FromFile,
		patchFile:    options.PatchFile,
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
	runtime.SetFinalizer(ret, freeEncoder)

	return ret, nil
}

func (enc *Encoder) GetStreamError() error {
	return lib.EncoderGetStreamError(enc.handle)
}

func (enc *Encoder) DumpStatsToStdout() {
	if enc.stats != nil {
		enc.stats.DumpToStdout()
	}
}

func (enc *Encoder) Process(ctx context.Context) error {
	var isFinal bool
	var perfStart time.Time

	for {
		// prepare gathering stats
		if enc.stats != nil {
			perfStart = time.Now()
		}

		// retrieve the current state
		state, err := lib.EncoderProcess(enc.handle)
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
			if enc.stats != nil {
				enc.stats.addDataBytes(state, n)
			}

			err = lib.EncoderProvideInputData(enc.handle, unsafe.Pointer(&enc.inputBuffer[0]), n, isFinal)
			if err != nil {
				return fmt.Errorf("Failed to provide data from TO/input file: %v", err)
			}
			break

		case lib.XdeltaState_OUTPUT:
			length, err := lib.EncoderGetOutputRequest(enc.handle)
			if err != nil {
				return fmt.Errorf("Failed to request data for PATCH/output file: %v", err)
			}
			if length <= 0 { // nothing to write?
				break
			}
			if length > len(enc.outputBuffer) {
				return fmt.Errorf("Failed to consume data for PATCH/output file: output buffer overflow")
			}

			err = lib.EncoderCopyOutputData(enc.handle, unsafe.Pointer(&enc.outputBuffer[0]))
			if err != nil {
				return fmt.Errorf("Failed to consume data for PATCH/output file: %v", err)
			}

			written, err := enc.patchFile.Write(enc.outputBuffer[:length])
			if err != nil {
				return fmt.Errorf("Failed to write data to PATCH/output file: %v", err)
			}
			if written < length {
				return fmt.Errorf("Failed to write data to PATCH/output file: not enough data written (%v < %v)", written, length)
			}
			if enc.stats != nil {
				enc.stats.addDataBytes(state, written)
			}
			break

		case lib.XdeltaState_GETSRCBLK:
			if enc.sourceFile == nil || enc.sourceBuffer == nil {
				return fmt.Errorf("Failed to request data for FROM/source file: not available")
			}

			blockno, blocksize, err := lib.EncoderGetSourceRequest(enc.handle)
			if err != nil {
				return fmt.Errorf("Failed to request data for FROM/source file: %v", err)
			}
			if blocksize != len(enc.sourceBuffer) {
				return fmt.Errorf("Failed to request data for FROM/source file: source buffer does not match block size (%v != %v)", blocksize, len(enc.sourceBuffer))
			}

			_, err = enc.sourceFile.Seek(int64(blockno)*int64(blocksize), io.SeekStart)
			if err != nil {
				return fmt.Errorf("Failed to seek FROM/source file: %v", err)
			}

			n, err := enc.sourceFile.Read(enc.sourceBuffer)
			if err != nil {
				return fmt.Errorf("Failed to read from FROM/source file: %v", err)
			}
			if enc.stats != nil {
				enc.stats.addDataBytes(state, n)
			}

			err = lib.EncoderProvideSourceData(enc.handle, unsafe.Pointer(&enc.sourceBuffer[0]), n)
			if err != nil {
				return fmt.Errorf("Failed to provide data from FROM/source file: %v", err)
			}
			break

		case lib.XdeltaState_WINSTART, lib.XdeltaState_WINFINISH:
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

func (enc *Encoder) Close() error {
	return freeEncoder(enc)
}

func freeEncoder(enc *Encoder) (err error) {
	// nothing to do?
	if enc == nil || enc.handle == nil {
		return
	}

	err = lib.FreeEncoder(enc.handle)

	enc.handle = nil
	return
}
