# https://stackoverflow.com/a/45397221

FROM alpine as certs
RUN apk update && apk add ca-certificates


FROM busybox
LABEL org.opencontainers.image.authors="Marvin Lukas Wenzel@mlw@mlw.wtf"

COPY --from=certs /etc/ssl/certs /etc/ssl/certs

# https://stackoverflow.com/questions/75696690/how-to-resolve-tls-failed-to-verify-certificate-x509-certificate-signed-by-un
# check certer.sh as a little help
# COPY /etc/pki/tls/certs/ca-bundle.crt /etc/ssl/certs/
# COPY /etc/pki/tls/certs/ca-bundle.trust.crt /etc/ssl/certs/

RUN mkdir /sandwich
WORKDIR /sandwich

COPY sandwich ./

CMD ["./sandwich"]
