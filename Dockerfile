# Compile stage
FROM golang:1.14.7-alpine3.12 AS build-env

# Install some dependencies needed to build the project
RUN apk update && apk add bash ca-certificates git gcc g++ libc-dev openssh-client

# Build Delve
RUN go get github.com/go-delve/delve/cmd/dlv

ADD . /src
WORKDIR /src
RUN go build -o /server
# Compile the application with the optimizations turned off
# This is important for the debugger to correctly work with the binary
RUN go build -gcflags "all=-N -l" -o /server

EXPOSE 3000 40000

VOLUME /src
ADD . /src
WORKDIR /src
CMD ['sh', '-c', 'go build -gcflags "all=-N -l" -o /server && /go/bin/dlv --headless --listen=0.0.0.0:40000 --api-version=2 exec /server']
