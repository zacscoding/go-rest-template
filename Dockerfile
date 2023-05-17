FROM golang:1.20-alpine AS build

RUN mkdir -p /go/src/github.com/zacscoding/go-rest-template ~/.ssh && \
    apk add --no-cache git openssh-client make gcc libc-dev
WORKDIR /go/src/github.com/zacscoding/go-rest-template
COPY . .
RUN make build

FROM alpine:latest
COPY --from=build /go/src/github.com/zacscoding/go-rest-template/build/bin/apiserver /usr/bin/apiserver
COPY --from=build /go/src/github.com/zacscoding/go-rest-template/migrations /etc/apiserver/migrations
COPY --from=build /go/src/github.com/zacscoding/go-rest-template/docs/docs.html /etc/apiserver/docs/docs.html

EXPOSE 8080
CMD ["/usr/bin/apiserver"]