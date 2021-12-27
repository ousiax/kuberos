FROM golang:1.17-bullseye
ENV GOPROXY="https://goproxy.cn,direct"
WORKDIR /build

COPY go.* ./
RUN go mod download

COPY . .
RUN go build -o kuberos

FROM debian:bullseye-slim
COPY --from=0 /build/kuberos /kuberos
ENTRYPOINT ["/kuberos"]
