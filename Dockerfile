# STEP 1: Used to get SSL root certificates
FROM alpine:3.8 as builder

# Install SSL ca certificates
RUN apk update && apk add git && apk add ca-certificates

COPY bin/alolstats /usr/bin/

CMD ["/usr/bin/alolstats", "-c", "/app/config/config.toml"]

LABEL org.label-schema.vendor="Abyle.org" \
      org.label-schema.url="https://github.com/torlenor/alolstats" \
      org.label-schema.name="ALoLStats" \
      org.label-schema.description="A League of Legends Statistics aggregation and calculation server written in GO"

