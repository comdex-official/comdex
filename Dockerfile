FROM golang:1.17-alpine3.14 as builder

WORKDIR /sources

# Add source files
COPY go.* .

# Install minimum necessary dependencies
RUN apk add --no-cache make gcc libc-dev

RUN make install

# ----------------------------

FROM alpine:3.14

COPY --from=builder /sources/build/ /usr/local/bin/

RUN addgroup ucmdx && adduser -S -G ucmdx ucmdx -h /data

USER ucmdx

WORKDIR /data

# p2p port
EXPOSE 26656/tcp

ENTRYPOINT ["comdex"]  
