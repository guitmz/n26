FROM golang AS build-env

WORKDIR /go/src/app
COPY . /go/src/github.com/guitmz/n26

RUN go get -v -d github.com/guitmz/n26/cmd/n26
RUN CGO_ENABLED=0 GOOS=linux go install -v -a -ldflags '-s -w -extldflags "-static"' github.com/guitmz/n26/cmd/n26

# final stage
FROM alpine

COPY --from=build-env /go/bin/n26 /

RUN apk add --no-cache ca-certificates

ENTRYPOINT [ "/n26" ]
