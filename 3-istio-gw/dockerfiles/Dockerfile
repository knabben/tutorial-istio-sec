FROM kong/httpbin:latest

WORKDIR /httpbin

RUN apt-get update
RUN DEBIAN_FRONTEND=noninteractive apt-get install -y curl tshark

COPY ./curl.sh .
