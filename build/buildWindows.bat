@echo off

set APP_NAME=multi-downloader
set APP_VER=0.0.1

set OUTPUT_DIR=.\output

set MINGW_32=C:\ProgramData\chocolatey
set MINGW_64=C:\ProgramData\chocolatey


::set MINGW=%MINGW_32
::set CC=i686-w64-mingw32-gcc
::call:Build windows 386

set MINGW=%MINGW_64
set CC=x86_64-w64-mingw32-gcc
call:Build windows amd64

goto:eof


rem Build Windows Executable
:Build
    set GOOS=%1
    set GOARCH=%2
    set CGO_ENABLED=1
    set PATH=%MINGW%\bin;%GOPATH%\bin;%PATH%

    set filename=%APP_NAME%-v%APP_VER%-%GOOS%-%GOARCH%.exe
    set filepath=%OUTPUT_DIR%\%filename%

    echo Building %filename%
    if not exist %OUTPUT_DIR% (
        md %OUTPUT_DIR%
    )

    echo pack icon...
    windres.exe -o ..\main.syso manifest.rc

    echo build exe...
    if exist %filepath% (
        del %filepath%
    )
    go build -o %filepath% -trimpath -ldflags="-s -w -H=windowsgui" ..

    echo rice pack...
    rice -i ../ append --exec %filepath%

    echo Done
goto:eof

