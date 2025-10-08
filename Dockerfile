# syntax=docker/dockerfile:1.4
FROM golang:1.22-alpine AS builder
WORKDIR /src

# Pre-fetch deps
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build static binary for target platform
ARG TARGETOS
ARG TARGETARCH
ENV CGO_ENABLED=0
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

FROM gcr.io/distroless/static:nonroot
WORKDIR /app

# App binary
COPY --from=builder /out/server /app/server

# Default lexicon files
COPY Vocabulary /app/Vocabulary

# Default envs
ENV PORT=8080
ENV LEXICON_DIR=Vocabulary

EXPOSE 8080
USER nonroot
ENTRYPOINT ["/app/server"]