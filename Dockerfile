FROM golang:1.14-alpine AS build

WORKDIR /go/source/app
COPY . .

RUN go get -d -v && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /bin/app

FROM scratch

EXPOSE 8080

COPY --from=build /bin/app /

CMD ["/app"]