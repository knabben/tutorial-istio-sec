# dynamic config
ARG             BUILD_DATE
ARG             VCS_REF
ARG             VERSION

# build
FROM            golang:1.18.0-alpine as builder
RUN             apk add --no-cache git gcc musl-dev make curl
ENV             GO111MODULE=on
WORKDIR         /go/src/moul.io/grpcbin
COPY            go.* ./
RUN             go mod download
COPY            . ./
RUN             curl -L -s https://github.com/fullstorydev/grpcurl/releases/download/v1.8.7/grpcurl_1.8.7_linux_x86_64.tar.gz -o grpcurl.tar.gz
RUN             tar xzvf grpcurl.tar.gz
RUN             mv grpcurl /go/bin/
RUN             go build -o /go/bin/grpcbin -ldflags "-extldflags \"-static\"" -v

# minimalist runtime
FROM alpine:3.15.1

LABEL           org.label-schema.build-date=$BUILD_DATE \
                org.label-schema.name="grpcbin" \
                org.label-schema.description="" \
                org.label-schema.url="https://moul.io/grpcbin/" \
                org.label-schema.vcs-ref=$VCS_REF \
                org.label-schema.vcs-url="https://github.com/moul/grpcbin" \
                org.label-schema.vendor="Manfred Touron" \
                org.label-schema.version=$VERSION \
                org.label-schema.schema-version="1.0" \
                org.label-schema.cmd="docker run -i -t --rm moul/grpcbin" \
                org.label-schema.help="docker exec -it $CONTAINER grpcbin --help"

RUN             apk update && apk add ca-certificates bash curl tcpdump net-tools && rm -rf /var/cache/apk/*
COPY            grpcurl.sh /bin/grpcurl.sh
COPY            --from=builder /go/bin/grpcbin /bin/grpcbin
COPY            --from=builder /go/bin/grpcurl /bin/grpcurl
WORKDIR         /root
EXPOSE          9000
CMD             ["/bin/grpcbin"]