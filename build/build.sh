#!bin/sh

APP_NAME=`awk -F '"' '/ProductName/{print $4;exit}' versioninfo.json`
APP_VER=`awk -F '"' '/ProductVersion/{print $4;exit}' versioninfo.json`
APP_DESC=`awk -F '"' '/FileDescription/{print $4;exit}' versioninfo.json`

docker build -t mdownloaderbuild:${APP_VER} -f ./Dockerfile .

docker run --rm -i \
    -e APP_NAME=$APP_NAME \
    -e APP_VER=$APP_VER \
    -e APP_DESC=$APP_DESC \
    -v $(dirname `pwd`):/workspace \
    mdownloaderbuild:${APP_VER}
