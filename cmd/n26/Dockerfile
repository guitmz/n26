FROM golang AS build-env

WORKDIR /go/src/app
COPY * ./

RUN go get -v -d
RUN go install -v
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-s -w -extldflags "-static"' -o n26 .

# final stage
FROM scratch

COPY --from=build-env /go/src/app/n26 /

ENTRYPOINT [ "/n26" ]
