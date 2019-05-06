FROM debian:stable-slim

RUN apt-get update
RUN apt-get install -y ca-certificates
RUN apt-get install -y curl
RUN curl https://s3.amazonaws.com/rds-downloads/rds-combined-ca-bundle.pem -o /usr/local/share/ca-certificates/rds.crt
RUN update-ca-certificates


COPY server .
ENTRYPOINT ["./server"]

