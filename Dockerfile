ARG IMG_TAG=latest

# Compile the kyved binary
FROM golang:1.20-alpine AS kyved-builder

# Install make
RUN apk add --no-cache make

WORKDIR /go/src

# Install dependencies
COPY go.mod go.sum* ./
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/go/pkg/mod \
    go mod download
COPY . .

ENV ENV=mainnet
RUN make install

# Copy binary to a distroless container
FROM gcr.io/distroless/static-debian11:$IMG_TAG

COPY --from=kyved-builder "/go/bin/kyved" /usr/local/bin/
EXPOSE 26656 26657 1317 9090

ENTRYPOINT ["kyved"]