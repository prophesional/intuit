FROM debian:stable-slim


ARG DB_DATABASE_TYPE
ARG DB_DATABASE_SECRET_KEY
ARG AWS_DEFAULT_REGION


RUN apt-get update
RUN apt-get install -y ca-certificates
RUN apt-get install -y curl
RUN curl https://s3.amazonaws.com/rds-downloads/rds-combined-ca-bundle.pem -o /usr/local/share/ca-certificates/rds.crt
RUN update-ca-certificates

ENV DB_DATABASE_TYPE=$DB_DATABASE_TYPE
ENV DB_DATABASE_SECRET_KEY=$DB_DATABASE_SECRET_KEY
ENV AWS_DEFAULT_REGION=$AWS_DEFAULT_REGION



COPY server .
ENTRYPOINT ["./server"]

