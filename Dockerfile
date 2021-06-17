FROM alpine

RUN apk add curl bind-tools

COPY do-update.sh /opt/do-update.sh

CMD ["sh", "/opt/do-update.sh"]
