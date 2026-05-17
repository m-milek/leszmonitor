FROM node:24-alpine AS frontend-builder

WORKDIR /app

COPY frontend/package*.json ./
RUN npm ci

COPY frontend/ .

RUN npm run build

FROM golang:1.26.3-alpine AS backend-builder

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

COPY backend/src/go.mod backend/src/go.sum* ./

RUN go mod download

COPY backend/src/ ./src

COPY --from=frontend-builder /app/dist ./src/static

WORKDIR /app/src

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main .
RUN mkdir -p /var/log/leszmonitor

FROM scratch

COPY --from=backend-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=backend-builder /var/log/leszmonitor /var/log/leszmonitor

WORKDIR /app

COPY --from=backend-builder /app/src/main .

EXPOSE 7001

CMD ["./main"]