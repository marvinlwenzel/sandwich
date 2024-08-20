# https://stackoverflow.com/a/45397221

FROM alpine as alp
RUN apk update && apk add ca-certificates
RUN mkdir /sandwich
COPY sandwich.stripped /sandwich/sandwich

FROM scratch
LABEL org.opencontainers.image.authors="Marvin Lukas Wenzel<mlw@mlw.wtf>"
LABEL org.opencontainers.image.url="https://github.com/marvinlwenzel/sandwich"
LABEL org.opencontainers.image.documentation="https://github.com/marvinlwenzel/sandwich"
LABEL org.opencontainers.image.source="https://github.com/marvinlwenzel/sandwich"
LABEL org.opencontainers.image.title="S.A.N.D.W.I.C.H."
LABEL org.opencontainers.image.description="Basic monitoring tool for checking remote webserver uptime and reporting them to discord hooks."

COPY --from=alp /etc/ssl/certs /etc/ssl/certs
COPY --from=alp /sandwich/sandwich /sandwich/sandwich

CMD ["/sandwich/sandwich"]
