FROM docker.io/redhat/ubi9-micro:9.2
# 设置时区为上海，ubi9-micro内置了tzdata
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo "Asia/Shanghai" > /etc/timezone
COPY go_web_server /
CMD ["/go_web_server"]