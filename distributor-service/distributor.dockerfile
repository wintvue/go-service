FROM alpine:latest

RUN mkdir /app

COPY distributorApp /app

CMD [ "/app/distributorApp" ]