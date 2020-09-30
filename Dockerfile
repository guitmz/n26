FROM golang AS build-env

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

COPY . /go/src/github.com/guitmz/n26
WORKDIR /go/src/github.com/guitmz/n26

RUN go mod download
RUN go install -v -a -ldflags '-s -w -extldflags "-static"' github.com/guitmz/n26/cmd/n26

# final stage
FROM alpine

COPY --from=build-env /go/bin/n26 /

RUN apk add --no-cache ca-certificates

ENTRYPOINT [ "/n26" ]
