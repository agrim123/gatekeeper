FROM ubuntu:18.04

RUN apt update && apt install ssh -y

RUN groupadd -r deploy && useradd -m -d /home/deploy -g deploy deploy

RUN chown -R deploy:deploy /home/deploy/keys

USER deploy

RUN mkdir ~/.ssh

RUN echo "Host * \n StrictHostKeyChecking no" > ~/.ssh/config
