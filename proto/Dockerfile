FROM bufbuild/buf:latest as BUILDER
FROM golang:1.23-alpine

# Install packages
RUN apk add --no-cache \
  nodejs \
  npm \
  git \
  make

# Install protoc plugins
RUN go install github.com/cosmos/gogoproto/protoc-gen-gocosmos@latest
RUN go install github.com/cosmos/cosmos-proto/cmd/protoc-gen-go-pulsar@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@latest
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

RUN npm install -g swagger-combine

COPY --from=BUILDER /usr/local/bin /usr/local/bin

# Inject user id and group id to avoid permission issues when running as a root user
ARG USER_ID=1000
ARG USER_GID=1000

# Create group if it doesn't exist
RUN if ! getent group $USER_GID; then \
  addgroup -g $USER_GID mygroup; \
fi

# Create user if it doesn't exist
RUN if ! getent passwd $USER_ID; then \
  adduser -D -u $USER_ID -G $(getent group $USER_GID | cut -d: -f1) myuser; \
fi

# Change to user
USER $USER_ID
