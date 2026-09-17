# Build static files
FROM node:18-slim AS node-builder

WORKDIR /src

COPY package*.json .
RUN npm install

RUN mkdir cmd
RUN mkdir cmd/server
COPY cmd/server/static cmd/server/static
RUN mkdir web
COPY web/static web/static
RUN npm run build

# Build go binary
FROM golang:latest AS go-builder

WORKDIR /src

COPY --from=node-builder /src/cmd cmd

COPY go.mod go.sum .
RUN go mod download

RUN mkdir web/
COPY web/templates web/templates
COPY web/views web/views
RUN go tool templ generate

COPY internal internal
COPY config config
COPY cmd/server/main.go cmd/server/main.go

RUN CGO_ENABLED=0 go build -o /piggytron cmd/server/main.go

# Runtime — only the binary
FROM scratch

COPY --from=go-builder /piggytron /piggytron
EXPOSE 8080

ENTRYPOINT ["/piggytron"]
