ARG PHP_VERSION

FROM ghcr.io/shopwell-shop/shopwell-cli-base:${PHP_VERSION}

ARG TARGETPLATFORM

COPY $TARGETPLATFORM/shopwell-cli /usr/local/bin/

ENTRYPOINT ["/entrypoint", "/usr/local/bin/entrypoint.sh", "/usr/local/bin/shopwell-cli"]
CMD ["--help"]
