FROM golang:1.19

WORKDIR /testserver
COPY / ./
RUN go get
RUN CGO_ENABLED=0 go build -o /usr/bin/testserver
EXPOSE 9003
EXPOSE 1601
RUN chmod +x ./entrypoint.sh
RUN chmod +x /usr/bin/testserver
RUN apt-get update && apt-get install netcat -y
ENTRYPOINT ["./entrypoint.sh"]