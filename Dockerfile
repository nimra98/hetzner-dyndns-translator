# syntax=docker/dockerfile:1

FROM golang:1.23.0 AS builder

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY *.go ./

# Build
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /hetzner-dyndns-translator

FROM scratch

WORKDIR /app

COPY --from=builder hetzner-dyndns-translator hetzner-dyndns-translator
# copy certs from golang:1.23 image (from=builder specifies the first layer)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Run
CMD ["/app/hetzner-dyndns-translator"]