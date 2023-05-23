FROM bufbuild/buf:latest as BUILDER
FROM golang:1.20-alpine

RUN apk add --no-cache \
  nodejs \
  npm \
  git \
  make

RUN go install github.com/cosmos/cosmos-proto/cmd/protoc-gen-go-pulsar@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@latest
RUN go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@latest

RUN git clone https://github.com/cosmos/gogoproto.git; \
  cd gogoproto; \
  go mod download; \
  make install

RUN npm install -g swagger-combine

COPY --from=BUILDER /usr/local/bin /usr/local/bin
