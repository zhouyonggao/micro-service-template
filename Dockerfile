FROM golang:1.20.1 AS builder

COPY . /src
WORKDIR /src

RUN GOPROXY=https://goproxy.cn make build

FROM debian:stable-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates  \
        netbase \
        && rm -rf /var/lib/apt/lists/ \
        && apt-get autoremove -y && apt-get autoclean -y

COPY --from=builder /src/bin /app
#正式环境不要打包 configs，本身也不会使用
COPY --from=builder /src/configs/ /app/configs/

WORKDIR /app

EXPOSE 9000
VOLUME /data/conf
ENV ROCKETMQ_GO_LOG_LEVEL=error TZ=Asia/Shanghai


CMD ["./server", "-conf", "./configs"]
