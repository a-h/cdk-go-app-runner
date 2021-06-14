FROM ubuntu:latest
ADD app ./
RUN apt-get update && apt-get -y install ca-certificates
EXPOSE 8000
ENTRYPOINT ["/app"]
