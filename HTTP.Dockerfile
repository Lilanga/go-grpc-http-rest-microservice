FROM golang:stretch AS builder

RUN apt-get update
RUN apt-get install protobuf-compiler -y
# Get the source from GitHub
RUN go get google.golang.org/grpc
# Install protoc-gen-go
RUN go get github.com/golang/protobuf/protoc-gen-go
# Install grpc gateway
RUN go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
# Install protoc swagger generator
RUN go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
ENV GO111MODULE=on \
    CGO_ENABLED=1

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN cd ./third_party/ && ./protoc-gen.sh
RUN cd ./third_party/ && ./build-server.sh http
WORKDIR /dist

RUN cp -r /build/dist/ .

FROM debian:stretch-slim
WORKDIR /dist
COPY --from=builder /dist/dist/ .

ENTRYPOINT ["./server-http","-grpc-host=grpc-server" ,"-grpc-port=9090", "-http-port=8080", "-log-level=-1", "-log-time-format=2006-01-02T15:04:05.999999999Z07:00"]
EXPOSE 8080