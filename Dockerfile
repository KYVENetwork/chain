ARG IMG_TAG=latest

# Compile the kyved binary
FROM golang:1.19-alpine AS kyved-builder

WORKDIR /go/src
COPY go.mod go.sum* ./
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/go/pkg/mod \
    go mod download
COPY . .
ENV PACKAGES make
RUN apk add --no-cache $PACKAGES
RUN CGO_ENABLED=0 make install


# Copy binary to a distroless container
FROM gcr.io/distroless/static-debian11:$IMG_TAG

COPY --from=kyved-builder "/go/bin/kyved" /usr/local/bin/
EXPOSE 26656 26657 1317 9090

ENTRYPOINT ["kyved"]
