# Using the heighliner image as a base, to fetch the comdex binaries
FROM ghcr.io/strangelove-ventures/heighliner/comdex:v11.5.1 as source

FROM alpine:3.16

COPY --from=source /bin /usr/bin
COPY --from=source /lib /usr/lib

# Set up dependencies for Starship
ENV PACKAGES curl make bash jq sed

# Install minimum necessary dependencies
RUN apk add --no-cache $PACKAGES

WORKDIR /opt
