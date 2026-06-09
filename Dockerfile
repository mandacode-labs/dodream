#syntax=docker/dockerfile:1.4

ARG VERSION=0.0.1
ARG GIT_COMMIT=unknown

############################
# 1. Build Stage
############################
FROM golang:1.26.4-alpine AS builder

RUN apk add --no-cache \
      git \
      gcc \
      musl-dev

WORKDIR /build

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download && go mod verify

COPY . .

ARG TARGETOS
ARG TARGETARCH
ARG VERSION
ARG GIT_COMMIT

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} \
    go build -a \
    -ldflags="-w -s \
              -X main.version=${VERSION} \
              -X main.gitCommit=${GIT_COMMIT}" \
    -o dodream ./cmd/dodream

############################
# 2. Runtime Stage
############################
FROM alpine:3.21 AS runtime

ARG VERSION

LABEL org.opencontainers.image.title="Dodream Engine" \
      org.opencontainers.image.description="Flashcard SRS engine" \
      org.opencontainers.image.version="${VERSION}" \
      org.opencontainers.image.source="https://github.com/mandacode-labs/dodream"

RUN apk add --no-cache \
      ca-certificates \
      tzdata \
      curl && \
    adduser -D -u 1001 dodream

WORKDIR /app

COPY --from=builder /build/dodream /app/dodream

RUN mkdir -p /app/config && chown -R dodream:dodream /app

USER dodream

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

ENTRYPOINT ["/app/dodream"]
CMD ["api-server", "run"]
