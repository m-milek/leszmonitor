FROM node:24-alpine AS frontend-builder

WORKDIR /app

ARG VITE_BACKEND_URL
ENV VITE_BACKEND_URL=${VITE_BACKEND_URL}

COPY frontend/package*.json ./
RUN npm ci

COPY frontend/ .

RUN npm run build

FROM golang:1.26.1-alpine AS backend-builder

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

COPY backend/src/go.mod backend/src/go.sum* ./

RUN go mod download

COPY backend/src/ ./src

COPY --from=frontend-builder /app/dist ./src/static

WORKDIR /app/src

RUN CGO_ENABLED=1 GOOS=linux go build -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates sqlite-libs

WORKDIR /app

COPY --from=backend-builder /app/src/main .

RUN mkdir -p /var/log/leszmonitor

EXPOSE 7001

CMD ["./main"]