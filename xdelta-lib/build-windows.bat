@echo off

pushd "%~dp0"

setlocal

if exist "C:\Program Files (x86)\Microsoft Visual Studio 14.0\VC" (
    call "C:\Program Files (x86)\Microsoft Visual Studio 14.0\VC\vcvarsall.bat" amd64
) else if exist "C:\Program Files (x86)\Microsoft Visual Studio\2019\Enterprise\VC\Auxiliary\Build\vcvarsall.bat" (
    call "C:\Program Files (x86)\Microsoft Visual Studio\2019\Enterprise\VC\Auxiliary\Build\vcvarsall.bat" amd64
)

if exist "C:\BuildTools\VC" (
    call "C:\BuildTools\VC\vcvarsall.bat" amd64
)

cl.exe /nologo /I src /MT /LD /GL /Fe:go-xdelta-lib.dll xdelta.cpp xdelta-encoder.cpp xdelta-decoder.cpp xdelta-go-helpers.cpp xdelta-go-encoder.cpp xdelta-go-decoder.cpp /link /RELEASE /LTCG /NOLOGO /VERSION:3.1

del /q go-xdelta-lib.lib go-xdelta-lib.exp *.obj

REM dumpbin.exe /nologo /exports go-xdelta-lib.dll

endlocal

popd

if not exist "%~dp0go-xdelta-lib.dll" (
    exit /b 1
)

exit /b 0
