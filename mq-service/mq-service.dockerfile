FROM alpine:latest

RUN mkdir /app

COPY mqApp /app

CMD [ "/app/mqApp" ]