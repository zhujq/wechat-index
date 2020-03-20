FROM alpine:latest

ADD entrypoint.sh /entrypoint.sh
ADD wechat-index /wechat-index

RUN  chmod +x /wechat-index && chmod 777 /entrypoint.sh

ENTRYPOINT  /entrypoint.sh 

EXPOSE 80
