FROM kong/httpbin:latest

WORKDIR /httpbin

RUN apt-get update
RUN DEBIAN_FRONTEND=noninteractive apt-get install -y curl \
    tshark \
    tcpdump \
    lsof \
    vim \
    zsh \
    strace \
    net-tools

RUN curl https://hey-release.s3.us-east-2.amazonaws.com/hey_linux_amd64 -o /httpbin/hey 
RUN chmod +x /httpbin/hey

COPY ./curl.sh .
