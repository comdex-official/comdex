FROM golang:1.17-alpine3.14 as builder

WORKDIR /sources

# Add source files
COPY go.* .

# Install minimum necessary dependencies
RUN apk add --no-cache make gcc libc-dev curl

RUN make install && curl -O https://raw.githubusercontent.com/comdex-official/networks/main/mainnet/genesis.json

# ----------------------------

FROM alpine:3.14

COPY --from=builder /sources/build/ /usr/local/bin/
COPY --from=builder /sources/genesis.json /data/.comdex/config/genesis.json

RUN addgroup ucmdx && adduser -S -G ucmdx ucmdx -h /data

USER ucmdx

WORKDIR /data

# p2p port
EXPOSE 26656/tcp

ENTRYPOINT ["comdex"]  
