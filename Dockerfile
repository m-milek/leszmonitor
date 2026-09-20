FROM node:24-alpine AS ui-builder

WORKDIR /app

COPY leszmonitor-ui/package*.json ./
RUN npm ci

COPY leszmonitor-ui/ .

RUN npm run build

FROM golang:1.27.1-alpine AS server-builder

ARG VERSION
ARG GIT_COMMIT
ARG CI_BUILD_NUMBER
ARG IMAGE_TAG

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

COPY leszmonitor-server/src/go.mod leszmonitor-server/src/go.sum* ./

RUN go mod download

COPY leszmonitor-server/src/ ./src

COPY --from=ui-builder /app/dist ./src/static

WORKDIR /app/src

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w \
    -X github.com/m-milek/leszmonitor/platform/meta.CIBuildNumber=$CI_BUILD_NUMBER \
    -X github.com/m-milek/leszmonitor/platform/meta.GitCommit=$GIT_COMMIT \
    -X github.com/m-milek/leszmonitor/platform/meta.ImageTag=$IMAGE_TAG \
    -X github.com/m-milek/leszmonitor/platform/meta.Version=$VERSION" \
    -o main .
RUN mkdir -p /var/log/leszmonitor

FROM scratch

COPY --from=server-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=server-builder --chown=65532:65532 /var/log/leszmonitor /var/log/leszmonitor
COPY --from=server-builder --chown=65532:65532 /app/src/main /app/main

WORKDIR /app
USER 65532:65532

EXPOSE 7001

CMD ["./main"]
