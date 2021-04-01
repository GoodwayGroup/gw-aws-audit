FROM alpine:3.13.4

COPY gw-aws-audit /usr/local/bin/gw-aws-audit
RUN chmod +x /usr/local/bin/gw-aws-audit

RUN mkdir /workdir
WORKDIR /workdir

ENTRYPOINT [ "/usr/local/bin/gw-aws-audit" ]