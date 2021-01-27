FROM golang:1.15.7-alpine
WORKDIR /home
# 代理加速
ENV GOPROXY=https://goproxy.cn
# 设置时区
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN apk add yarn
RUN yarn global add pm2
COPY . .
RUN go build -o server
CMD ["/home/server"]

