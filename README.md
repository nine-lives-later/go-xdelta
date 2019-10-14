# Xdelta for Go

This library provides a wrapper for the [Xdelta library](http://xdelta.org/) by Joshua MacDonald and others. 

**[Click here to open the GoDoc documentation.](https://godoc.org/github.com/konsorten/go-xdelta)**

## Getting Started

Patches are being created using the *encoder*, while applying the resulting patches is done by the *decoder*. The following workflows do exist:

| Title | Data Flow | Description |
| --- | --- | --- |
| Encoding (changed file) | <nobr>`FROM -> TO => PATCH`</nobr> | The encoding reads the new TO file and compares the data to the original FROM file and outputs the resulting PATCH file. |
| Encoding (new file) | <nobr>`TO => PATCH`</nobr> | The encoding reads the new TO file outputs the resulting PATCH file. |
| Decoding (changed file) | <nobr>`PATCH -> FROM => TO`</nobr> | The encoding reads the PATCH file and applies its operations to the original FROM file and outputs the resulting new TO file. |
| Decoding (new file) | <nobr>`PATCH => TO`</nobr> | The encoding reads the PATCH file outputs the resulting new TO file. |

There is no process for deleting files (see *Best Practices* below).

The following example is more or less pseudo-code. (It should be easy enough to understand.)

```go
import "github.com/konsorten/go-xdelta"

options := xdelta.EncoderOptions{
    FileID:    "myfile.ext",
    FromFile:  fromFileReaderSeeker,
    ToFile:    toFileReader,
    PatchFile: patchFileWriter,
}

enc, err := xdelta.NewEncoder(options)
if err != nil {
    return err
}
defer enc.Close()

// create the patch
err = enc.Process(context.TODO())
if err != nil {
    return err
}
```

The decoder works the same way.

## Tracking Progress

The easiest way to track the progress is for encoding/creating to determine how much data has been read from the FROM file. And for the decoding/patching to take the PATCH file's read progress.

## Best Practices

1. Pre-allocate the TO/patched file when applying patches. This will reduce the fragmentation on the file system as it can reserve a spot (on the disk drive) that is large enough for the new file. For this to work, make sure to store the TO/patched file size so it can be read upfront.

1. Check FROM file hash, before decoding/patching. Ensure that the FROM and PATCH files are correct, before starting the decoding/patching process.

1. Handle deletion of files! The patching mechanism does not handle the deletion of existing files. Handle this yourself based on your meta-information. Be aware of the difference between an empty file (filesize of 0) and a deleted one.

1. Do not use the patch file header. It is convenient place to store information, but usually you need security related information like the patch file hash upfront, anyway. Store other meta-information like the TO/patched file size and FROM file hashes, too.

1. Sign the meta-information (like file sizes and hashes) with an asymmetric encryption key! Do this by calculating a hash and signing that one (never encrypt the file content itself). Sign using a private key and check the signature with the public one. Make sure to have the public key be embedded into your client to prevent man-in-the-middle attacks.

## Building

> On Windows, the project requires the *xdelta-lib* native C++ library to be built into a DLL file, before it will work. See *Native Library* below for details. (This is not needed for macOS and Linux.)

To build this project, simply run the following command:

```
go build
```

To run all the tests (including a patch roundtrip test), run the following command:

```
go test -v
```

## Native Library (Windows only)

On Windows, no CGO is used. Instead the native library has to be built separately and provided in your project's directory.

To build it, run the build script. The native library is saved in the *xdelta-lib* sub-directory.

`./xdelta-lib/build-windows.bat`

You can also obtain a [pre-compiled version](https://github.com/konsorten/go-xdelta/releases).

## Build Status

The current status of the *master* branch:

| OS | Status |
| --- | --- |
| Windows AMD64 | [![Build status](https://ci.appveyor.com/api/projects/status/w8cqnh1a2a23kkov/branch/master?svg=true)](https://ci.appveyor.com/project/fkollmann/go-xdelta/branch/master) |
| Ubuntu 18 AMD64 | [![Build status](https://ci.appveyor.com/api/projects/status/907d9w99x0pvlf0k/branch/master?svg=true)](https://ci.appveyor.com/project/fkollmann/go-xdelta-d0wsq/branch/master) |
| macOS AMD64 | supported, but no CI |

## Authors

The library is sponsored by the [marvin + konsorten GmbH](http://www.konsorten.de).

We thank all the authors who provided code to this library:

* Felix Kollmann

It is also based on the work by [Joshua MacDonald](https://github.com/jmacd) and others.

## License

(The Apache 2 License)

Copyright 2019 marvin + konsorten GmbH (open-source@konsorten.de)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
