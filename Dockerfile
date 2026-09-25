# ---- build ----
FROM golang:1.25-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

ARG APP_VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/say-hi ./cmd

# ---- runtime ----
# Distroless: no shell or package manager, CA certificates included, and it
# runs as an unprivileged user (uid 65532).
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app

COPY --from=build /out/say-hi ./say-hi
COPY views ./views
COPY public ./public

ARG APP_VERSION=dev
ENV APP_VERSION=${APP_VERSION} \
    PORT=8080

EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/say-hi"]
