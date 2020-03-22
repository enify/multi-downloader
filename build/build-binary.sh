#!/bin/sh

OUTPUT=$(cd `dirname $0`; pwd)/output
TMP_OUTPUT=/output


Build() {
    export GOOS=$1
    export GOARCH=$2
    export CGO_ENABLED=1
    export PATH="$GOPATH/bin:$PATH"

    suffix=""
    ldflags="-s -w"
    if [ $GOOS = "windows" ]; then
        suffix=".exe"
        ldflags="-s -w -H=windowsgui"
    fi

    filename=${APP_NAME}-v${APP_VER}-${GOOS}-${GOARCH}${suffix}
    filepath=${TMP_OUTPUT}/${filename}

    echo Building $filename
    mkdir -p $TMP_OUTPUT

    echo build executable...
    go build -o $filepath -trimpath -ldflags="${ldflags}" .

    echo rice pack...
    rice -i ./ append --exec $filepath

    echo Done
}


go mod download

export CC=x86_64-w64-mingw32-gcc
Build windows amd64

export CC=gcc
Build linux amd64

mkdir -p $OUTPUT
mv -v $TMP_OUTPUT/* $OUTPUT
echo All Done
