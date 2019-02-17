# Xdelta for Go

This library provides a wrapper for the [Xdelta library](http://xdelta.org/) by Joshua MacDonald and others. 

**[Click here to open the GoDoc documentation.](https://godoc.org/github.com/konsorten/go-xdelta)**

## Example

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
