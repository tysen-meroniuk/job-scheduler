# syntax=docker/dockerfile:1

# ─── Stage 1: build the dashboard ──────────────────────────────
FROM node:20-alpine AS dashboard-build
WORKDIR /app
COPY dashboard/package*.json ./
RUN npm ci
COPY dashboard/ ./
RUN npm run build

# ─── Stage 2: build the Go binary (with embedded dashboard) ────
FROM golang:1.22-alpine AS go-build
WORKDIR /src
COPY go.mod ./
COPY go.sum* ./
RUN go mod download
COPY . .
# Copy built dashboard into the embed directory before go build.
RUN rm -rf internal/web/static && mkdir -p internal/web/static
COPY --from=dashboard-build /app/dist/ ./internal/web/static/
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" \
    -o /out/jobqueue ./cmd/jobqueue

# ─── Stage 3: runtime ──────────────────────────────────────────
FROM alpine:3.19
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=go-build /out/jobqueue /usr/local/bin/jobqueue
COPY --from=go-build /src/migrations /app/migrations
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/jobqueue"]
CMD ["api"]
