FROM golang:1.17.13-buster as builder

WORKDIR /build

COPY . .

RUN export GO111MODULE=on&& \
export GOPROXY=https://goproxy.cn&& \
go mod download
RUN mkdir -p ./release
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix cgo -o ./release/http-server ./internal/http
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix cgo -o ./release/ws-server ./internal/websocket
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix cgo -o ./release/cmd ./internal/cmd

FROM alpine

WORKDIR /app
# RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY --from=builder /build/release/* .
COPY --from=builder /build/start.sh .
COPY --from=builder /build/config.prod.yaml ./config.yaml

RUN chmod +x /app/start.sh

EXPOSE 9503 9504

CMD ["./start.sh"]