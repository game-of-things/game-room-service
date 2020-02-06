FROM golang:alpine AS build

WORKDIR /go/source/app
COPY . .

RUN go get -d -v ./... && \
    go install -v ./...

FROM alpine:latest

EXPOSE 8080

COPY --from=build /go/bin/game-room-service /usr/local/bin

CMD ["game-room-service"]