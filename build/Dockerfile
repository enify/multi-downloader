FROM golang:1.13-stretch

RUN sed -i 's/deb.debian.org/mirrors.ustc.edu.cn/g' /etc/apt/sources.list
RUN apt-get update && apt-get install -y pkg-config libgtk-3-dev gcc-mingw-w64 golang-rice zip

ENV GOPROXY=https://goproxy.io

WORKDIR /workspace
ENTRYPOINT ["/workspace/build/build-binary.sh"]
