# Build
# Build frontend
FROM node:20-slim AS build-frontend
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
WORKDIR /build

COPY front /build/front

RUN corepack enable
RUN cd front && pnpm install && pnpm build

# Build backend
FROM golang:1.23 AS build-backend

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
COPY internal/ ./internal

# Enable static binary, target Linux, remove debug information and strip binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -s"

# Copy to our image
FROM gcr.io/distroless/static-debian12
LABEL org.opencontainers.image.authors="Ada <ada@gnous.eu>"
LABEL org.opencontainers.image.source="https://github.com/r4iponce/source"

WORKDIR /app

COPY --from=build-backend /build/status .
COPY --from=build-frontend /build/build /app/build

EXPOSE 3000/tcp

ENTRYPOINT ["/app/status", "/config.toml"]