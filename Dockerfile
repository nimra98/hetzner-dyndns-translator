# use golang image to copy ssl certs later
FROM golang:1.16

FROM scratch

# copy certs from golang:1.16 image
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENV GODEBUG=netdns=go
ADD ./hetzner-dyndns-translator /bin/

ENTRYPOINT ["/bin/hetzner-dyndns-translator"]
