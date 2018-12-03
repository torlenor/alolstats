# STEP 1: Used to get SSL root certificates
FROM alpine:3.8

# Install SSL ca certificates
RUN apk update && apk add ca-certificates
# Install ALoLStats dependencies
RUN apk update && apk add g++ R R-dev R-doc libc-dev
RUN apk update && apk add ttf-liberation

COPY bin/alolstats /usr/bin/
COPY R/*.R /app/R/

RUN Rscript /app/R/install_packages.R

CMD ["/usr/bin/alolstats", "-c", "/app/config/config.toml"]

LABEL org.label-schema.vendor="Abyle.org" \
      org.label-schema.url="https://github.com/torlenor/alolstats" \
      org.label-schema.name="ALoLStats" \
      org.label-schema.description="A League of Legends Statistics aggregation and calculation server written in GO"

