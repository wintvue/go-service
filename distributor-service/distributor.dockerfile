FROM golang:1.19-alpine as builder

WORKDIR /app

# COPY package*.json ./

COPY . /app

RUN CGO_ENABLED=0 go build -o distributorApp ./cmd/api 

RUN chmod +x /app/distributorApp

FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/distributorApp /app

CMD [ "/app/distributorApp" ]