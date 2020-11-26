FROM ubuntu:18.04

RUN apt update && apt install ssh -y

# RUN groupadd -r deploy && useradd -m -d /home/deploy -g deploy deploy

# USER deploy

RUN mkdir ~/.ssh
RUN mkdir ~/keys

WORKDIR /root

RUN echo "Host * \n StrictHostKeyChecking no" > ~/.ssh/config
