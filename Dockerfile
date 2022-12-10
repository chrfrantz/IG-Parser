FROM golang:1.19 as builder

LABEL maintainer="cf@christopherfrantz.org"
LABEL stage=builder

WORKDIR /go/src/IG-Parser/web

# Copy all relevant folders from repository
COPY ./app /go/src/IG-Parser/app
COPY ./config /go/src/IG-Parser/config
COPY ./exporter /go/src/IG-Parser/exporter
COPY ./parser /go/src/IG-Parser/parser
COPY ./shared /go/src/IG-Parser/shared
COPY ./tree /go/src/IG-Parser/tree
COPY ./web /go/src/IG-Parser/web
COPY ./go.mod /go/src/IG-Parser/go.mod
COPY ./go.sum /go/src/IG-Parser/go.sum

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o main


# Target container
FROM scratch

LABEL maintainer="cf@christopherfrantz.org"

WORKDIR /

# Retrieve binary from builder container
COPY --from=builder /go/src/IG-Parser/web .

CMD ["/main"]
