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

COPY proto proto

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

RUN chown -R $USER_ID:$USER_GID proto

# Change to user
USER $USER_ID

RUN cd proto && buf mod update