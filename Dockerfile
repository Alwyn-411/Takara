FROM node:22-alpine AS frontend
WORKDIR /app

RUN corepack enable

COPY package.json pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

COPY . .
RUN pnpm run build

FROM golang:1.26-alpine AS backend
WORKDIR /src

COPY services/go.mod services/go.sum ./
RUN go mod download

COPY services/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /bin/server .

FROM alpine:3.20
WORKDIR /app

RUN apk add --no-cache ca-certificates && \
    mkdir -p /data

COPY --from=backend /bin/server /app/server
COPY --from=frontend /app/dist /app/dist

ENV GIN_MODE=release \
    PORT=8080 \
    DB_PATH=/data/takara.db \
    STATIC_DIR=/app/dist 

EXPOSE 8080

VOLUME ["/data"]

ENTRYPOINT ["/app/server"]