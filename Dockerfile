FROM golang:1.20 as builder

LABEL maintainer="cf@christopherfrantz.org"
LABEL stage=builder

WORKDIR /go/src/IG-Parser/web

# Copy all relevant folders from repository
COPY ./core /go/src/IG-Parser/core
COPY ./web /go/src/IG-Parser/web
COPY ./go.mod /go/src/IG-Parser/go.mod
COPY ./go.sum /go/src/IG-Parser/go.sum

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o main


# Target container
FROM scratch

LABEL maintainer="cf@christopherfrantz.org"
LABEL stage=runner

WORKDIR /

# Retrieve binary from builder container
COPY --from=builder /go/src/IG-Parser/web .

CMD ["/main"]
