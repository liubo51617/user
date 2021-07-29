package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/liubo51617/user/dao"
	"github.com/liubo51617/user/endpoint"
	"github.com/liubo51617/user/redis"
	"github.com/liubo51617/user/service"
	"github.com/liubo51617/user/transport"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main(){
	var (
		serverPoint = flag.Int("server.point", 10086, "server point")
	)

	flag.Parse()
	time.Sleep(10 * time.Second) // 延时启动，等待 MySQL 和 Redis 准备好
	ctx := context.Background()
	errchan := make(chan error)

	err := dao.InitMysql("127.0.0.1", "3306", "root", "123456", "user")
	if err != nil {
		log.Fatal(err)
	}

	err = redis.InitRedis("127.0.0.1", "6379", "")
	if err != nil {
		log.Fatal(err)
	}

	userService := service.MakeUserServiceImpl(&dao.UserDAOImpl{})

	userEndpoints := &endpoint.UserEndpoints{
		endpoint.MakeRegisterEndpoint(userService),
		endpoint.MakeLoginEndpoint(userService),
	}

	r := transport.MakeHttpHandler(ctx, userEndpoints)

	go func() {
		errchan <- http.ListenAndServe(":"+strconv.Itoa(*serverPoint), r)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errchan <- fmt.Errorf("%s", <-c)
	}()

	error := <-errchan
	log.Println(error)
	
}