FROM golang:1.20.1-buster as build-env

ADD . /app
WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

# -----------------------------------------------------------------------------

FROM scratch

COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build-env /app/beckart /bin/beckart

USER 1664

ENTRYPOINT ["/bin/beckart"]
CMD ["--config", "/beckart.yaml"]
