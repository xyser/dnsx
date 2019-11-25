FROM golang:latest as builder
ENV GO111MODULE=on GOPROXY=https://goproxy.io,direct
WORKDIR /opt
COPY ./ /opt
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X 'dnsx/api/controller/v1.BuildTime=`date +"%Y-%m-%d %H:%M:%S"`' -X dnsx/api/controller/v1.BuildVersion=1.0.1" -tags=jsoniter -o dnsx .

FROM alpine:latest
LABEL maintainer="dingdayu <614422099@qq.com>"
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
    && apk --no-cache add ca-certificates tzdata curl
ENV DREAMENV=TEST DEPLOY_TYPE=DOCKER
WORKDIR /opt/dnsx
COPY --from=0 /opt/dnsx .
COPY --from=0 /opt/config/config.yaml /opt/dnsx/config/config.yaml
EXPOSE 8080 53 53/udp
HEALTHCHECK --interval=5s --timeout=3s \
  CMD curl -fs -X HEAD http://127.0.0.1:8080/health || exit 1
ENTRYPOINT ["/opt/dnsx/dnsx", "server"]
