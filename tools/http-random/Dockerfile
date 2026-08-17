FROM golang:1.26.1-alpine AS build

WORKDIR /go/src/app

COPY main.go go.mod ./

RUN CGO_ENABLED=0 go build -o /go/bin/app

FROM gcr.io/distroless/static-debian12

EXPOSE 6767

COPY --from=build /go/bin/app /
CMD ["/app"]