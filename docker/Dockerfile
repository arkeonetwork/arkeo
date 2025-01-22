FROM alpine:latest

COPY ./arkeod /usr/bin/arkeod

WORKDIR /root/.arkeo

RUN apk add --no-cache \
    bash \
    coreutils \
    jq \
    ca-certificates \
    curl

EXPOSE 1317
EXPOSE 26656
EXPOSE 26657

ENTRYPOINT [ "/bin/sh" ]

CMD [ "/usr/bin/arkeod", "help" ]