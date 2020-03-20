FROM alpine:latest

ADD entrypoint.sh /entrypoint.sh
ADD wechat-index /wechat-index
ADD wechat-db /wechat-db
RUN  chmod +x /wechat-index /wechat-db && chmod 777 /entrypoint.sh
ENTRYPOINT  /entrypoint.sh 

EXPOSE 80
