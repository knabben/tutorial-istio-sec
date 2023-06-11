FROM golang:latest as build

WORKDIR /go/src/app
COPY . .

RUN go mod tidy
RUN CGO_ENABLED=0 go build -o /go/bin/server

FROM gcr.io/distroless/static-debian11
COPY --from=build /go/bin/server /
CMD ["/server"]
