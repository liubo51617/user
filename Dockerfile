#jenkins build start
FROM golang:latest
WORKDIR /root/github.com/liubo51617/user
COPY / /root/github.com/liubo51617/user
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go build -o user
EXPOSE 10086
ENTRYPOINT ./user -mysql.addr $mysqlAddr -redis.addr $redisAddr

#jenkins build end





#以下构建user的镜像
#FROM golang:latest
#WORKDIR /var/www/go/src/github.com/liubo51617/user
#COPY / /var/www/go/src/github.com/liubo51617/user
#RUN go env -w GOPROXY=https://goproxy.cn,direct
#RUN go build -o user
#ENTRYPOINT ["./user"]
#构建user的镜像结束

#下面是构建mysql镜像
#FROM mysql:5.7
#
#WORKDIR /docker-entrypoint-initdb.d
#
#ENV LANG=C.UTF-8
#
#COPY init.sql .
#构建mysql镜像结束
