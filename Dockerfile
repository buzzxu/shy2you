FROM golang:bullseye as build

WORKDIR $GOPATH/src/shy2you
ADD . $GOPATH/src/shy2you
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io


RUN apt-get update && \
    apt-get install -y wget build-essential pkg-config --no-install-recommends


RUN cd && \
    cd $GOPATH/src/shy2you && \
    go build -o app; \
    mv app  /opt/app;


FROM debian:bullseye-slim

MAINTAINER buzzxu <downloadxu@163.com>

WORKDIR /app
COPY --from=build /opt/app /app

ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y wget build-essential pkg-config fontconfig --no-install-recommends && \
	ldconfig /usr/local/lib && \
	rm /etc/localtime && \
    ln -sv /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    apt-get remove --purge -y wget build-essential pkg-config && \
    apt-get clean && \
    apt-get autoremove -y && \
    apt-get autoclean && \
    rm -rf /var/lib/apt/lists/* && \
    rm -rf /tmp/*

ADD docker/app.yml /app/app.yml
ADD docker/run.sh /app/run.sh

ENV TZ Asia/Shanghai
ENV LANG C.UTF-8

EXPOSE 3000
#ENTRYPOINT ["/bin/bash","run.sh"]

