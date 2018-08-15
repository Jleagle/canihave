# Build image
FROM golang:1.10-alpine AS build-env
WORKDIR /go/src/github.com/Jleagle/canihave/
COPY . /go/src/github.com/Jleagle/canihave/
RUN apk update && apk add curl git openssh
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo

# Runtime image
FROM alpine:3.8
WORKDIR /root/
COPY --from=build-env /go/src/github.com/Jleagle/canihave/canihave ./
COPY templates ./templates
COPY assets ./assets
COPY location/GeoLite2-Country.mmdb ./location/GeoLite2-Country.mmdb
RUN apk update && apk add ca-certificates bash nano
EXPOSE 8080
CMD ["./canihave"]
