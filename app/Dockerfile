FROM golang:1.14 as build

WORKDIR /go/src/app
COPY go.mod /go/src/app
COPY go.sum /go/src/app
RUN go mod download

COPY . /go/src/app
ENV CGO_ENABLED 0
RUN go build -o /go/bin/app

FROM golang:1.14 as dev
WORKDIR /go/src/app
COPY --from=build /go/src/app /go/src/app
RUN go get github.com/githubnemo/CompileDaemon
ENTRYPOINT CompileDaemon --build="go build -o /go/bin/app" --command=/go/bin/app

FROM gcr.io/distroless/static:nonroot as prod
COPY --from=build /go/bin/app /
CMD ["/app"]