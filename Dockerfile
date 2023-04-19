FROM golang AS build-env

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux

COPY . /go/src/github.com/guitmz/n26
WORKDIR /go/src/github.com/guitmz/n26

RUN go mod download \
    && go install -v -a -ldflags '-s -w -extldflags "-static"' github.com/guitmz/n26/cmd/n26

# final stage
FROM gcr.io/distroless/static-debian11
COPY --from=build-env /go/bin/n26 /
CMD [ "/n26" ]
