#!bin/sh

APP_NAME=multi-downloader
APP_VER=0.0.1

docker build -t mdownloaderbuild:${APP_VER} -f ./Dockerfile .

docker run --rm -i \
    -e APP_NAME=$APP_NAME \
    -e APP_VER=$APP_VER \
    -v $(dirname `pwd`):/workspace \
    mdownloaderbuild:${APP_VER}
