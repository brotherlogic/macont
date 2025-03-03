# syntax=docker/dockerfile:1

FROM golang:1.24 AS build

WORKDIR $GOPATH/src/github.com/brotherlogic/macont

COPY go.mod ./
COPY go.sum ./

RUN mkdir proto
COPY proto/*.go ./proto/

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 go build -o /macont

##
## Deploy
##
FROM ubuntu:22.04
USER root:root

WORKDIR /
COPY --from=build /macont /macont

EXPOSE 8080
EXPOSE 8081

ENTRYPOINT ["/macont"]