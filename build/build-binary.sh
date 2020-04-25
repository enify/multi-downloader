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
        (cd ./build && goversioninfo -64 -o ../main.syso ./versioninfo.json)
    fi

    out=${APP_NAME}-v${APP_VER}-${GOOS}-${GOARCH}
    outpath=${TMP_OUTPUT}/${out}
    filepath=${outpath}/${APP_NAME}${suffix}

    echo "Building $out"
    mkdir -p $TMP_OUTPUT

    echo "build executable..."
    go build -o $filepath -trimpath -ldflags="-X main.AppName=${APP_NAME} -X main.AppVer=${APP_VER} -X main.AppDesc=${APP_DESC} ${ldflags}" .

    echo "rice pack..."
    rice -i ./ append --exec $filepath

    echo "copy dynamic library..."
    if [ $GOOS = "windows" -a $GOARCH = "amd64" ]; then
        cp ./build/resource/sciter.dll $outpath
    elif [ $GOOS = "linux" -a $GOARCH = "amd64" ]; then
        cp ./build/resource/libsciter-gtk.so $outpath
    fi

    echo "zip..."
    mkdir -p $OUTPUT
    (cd $TMP_OUTPUT && zip -q -r "${OUTPUT}/${out}.zip" ./$out)

    echo "Done"
}


go mod download

export CC=x86_64-w64-mingw32-gcc
Build windows amd64

export CC=gcc
Build linux amd64

echo "All Done"
