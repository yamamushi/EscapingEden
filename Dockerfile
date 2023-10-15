# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/yamamushi/EscapingEden

# Create our shared volume
RUN mkdir /data

# Build the EQB command inside the container.
RUN cd /go/src/github.com/yamamushi/EscapingEden && go get -v ./... && go build -v ./... && go install

# Run the EQB command by default when the container starts.
WORKDIR /data
ENTRYPOINT /go/bin/EscapingEden

# Set the working directory to /data/
VOLUME /data
WORKDIR /data

EXPOSE 3000