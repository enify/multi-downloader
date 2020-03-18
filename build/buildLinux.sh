#!/bin/sh

APP_NAME=multi-downloader
APP_VER=0.0.1

OUTPUT_DIR=./output


Build() {
    export GOOS=$1
    export GOARCH=$2
    export CGO_ENABLED=1
    export PATH="$GOPATH/bin:$PATH"

    filename=${APP_NAME}-v${APP_VER}-${GOOS}-${GOARCH}
    filepath=${OUTPUT_DIR}/${filename}

    echo Building $filename
    if [ ! -d $OUTPUT_DIR ]; then
        mkdir $OUTPUT_DIR
    fi

    echo build executable...
    go build -o $filepath -trimpath -ldflags="-s -w" ..

    echo rice pack...
    rice -i ../ append --exec $filepath

    echo Done
}


Build linux amd64
