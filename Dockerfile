FROM node:24-alpine AS frontend-builder

WORKDIR /app

COPY frontend/package*.json ./
RUN npm ci

COPY frontend/ .

RUN npm run build

FROM golang:1.26.3-alpine AS server-builder

ARG VERSION
ARG GIT_COMMIT
ARG CI_BUILD_NUMBER
ARG IMAGE_TAG

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

COPY leszmonitor-server/src/go.mod leszmonitor-server/src/go.sum* ./

RUN go mod download

COPY leszmonitor-server/src/ ./src

COPY --from=frontend-builder /app/dist ./src/static

WORKDIR /app/src

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w \
    -X github.com/m-milek/leszmonitor/meta.CIBuildNumber=$CI_BUILD_NUMBER \
    -X github.com/m-milek/leszmonitor/meta.GitCommit=$GIT_COMMIT \
    -X github.com/m-milek/leszmonitor/meta.ImageTag=$IMAGE_TAG \
    -X github.com/m-milek/leszmonitor/meta.Version=$VERSION" \
    -o main .
RUN mkdir -p /var/log/leszmonitor

FROM scratch

COPY --from=server-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=server-builder /var/log/leszmonitor /var/log/leszmonitor

WORKDIR /app

COPY --from=server-builder /app/src/main .

EXPOSE 7001

CMD ["./main"]