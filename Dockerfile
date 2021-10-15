FROM golang:alpine

LABEL maintainer="Ayaan Zaidi <hi@obviy.us>"

# Install gcc
RUN apk add build-base

WORKDIR "/go/src/github.com/downcount/api"
COPY . .

# Build Downcount server
RUN go build -a -v cmd/tldw/tldw.go

# Run server
CMD ["./tldw", "start"]