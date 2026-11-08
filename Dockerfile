# Build static files
FROM node:18-slim AS node-builder

WORKDIR /src

COPY package*.json .
RUN npm install

RUN mkdir -p internal/interface/http/assets/static
COPY internal/interface/http/assets/static internal/interface/http/assets/static
RUN mkdir web
COPY web/static web/static
RUN npm run build

RUN apt-get update \
    && apt-get install -y --no-install-recommends brotli \
    && rm -rf /var/lib/apt/lists/*
RUN gzip -9 -k internal/interface/http/assets/static/*.css
RUN gzip -9 -k internal/interface/http/assets/static/*.svg
RUN brotli -q 11 -k internal/interface/http/assets/static/*.js
RUN brotli -q 11 -k internal/interface/http/assets/static/*.css
RUN brotli -q 11 -k internal/interface/http/assets/static/*.svg

# Build go binary
FROM golang:latest AS go-builder

WORKDIR /src

COPY go.mod go.sum .
RUN go mod download

RUN mkdir web/
COPY web/templates web/templates
COPY web/views web/views
RUN go tool templ generate

COPY internal internal
COPY --from=node-builder /src/internal/interface/http/assets/static internal/interface/http/assets/static
COPY config config
COPY cmd/server/main.go cmd/server/main.go

RUN CGO_ENABLED=0 go build -o /piggytron cmd/server/main.go

# Runtime — only the binary
FROM scratch

COPY --from=go-builder /piggytron /piggytron
EXPOSE 8080

ENTRYPOINT ["/piggytron"]
