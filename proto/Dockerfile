FROM bufbuild/buf:latest as BUILDER
FROM golang:1.20-alpine

RUN apk add --no-cache \
  nodejs \
  npm \
  git \
  make

RUN go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@latest
RUN go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@latest

RUN git clone https://github.com/cosmos/gogoproto.git; \
  cd gogoproto; \
  go mod download; \
  make install

RUN npm install -g swagger-combine

COPY --from=BUILDER /usr/local/bin /usr/local/bin

# Inject user id and group id to avoid permission issues when running as a root user
ARG USER_ID=1000
ARG USER_GID=1000

# Create group and user
RUN addgroup -g $USER_GID mygroup && \
    adduser -D -u $USER_ID -G mygroup myuser

# Change to user
USER $USER_ID
