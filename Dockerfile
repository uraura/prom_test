FROM golang:1.17 AS build

WORKDIR /go/src/app
COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN go build -o app

FROM clearlinux

COPY --from=build /go/src/app/app /app

ENTRYPOINT ["/app"]
