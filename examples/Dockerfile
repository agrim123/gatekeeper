FROM alpine:3.12	

RUN apk add --update openssh-client \	
  ca-certificates	

RUN addgroup -S deploy && adduser --disabled-password --home /home/deploy --ingroup deploy deploy	

RUN echo "Host *" >> /etc/ssh/ssh_config	
RUN echo "  StrictHostKeyChecking no" >> /etc/ssh/ssh_config	

USER deploy	
WORKDIR /home/deploy	

RUN mkdir ~/.ssh	
RUN mkdir ~/keys
