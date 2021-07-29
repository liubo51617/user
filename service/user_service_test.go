package service

import (
	"context"
	"github.com/liubo51617/user/dao"
	"github.com/liubo51617/user/redis"
	"testing"
)

func TestUserServiceImpl_Login(t *testing.T) {
	err := dao.InitMysql("127.0.0.1", "3306", "root", "root", "user")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = redis.InitRedis("127.0.0.1","6379", "" )
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	userService := &UserServiceImpl{userDAO:&dao.UserDAOImpl{}}
	user, err := userService.Login(context.Background(), "alexliu@qq.com","alexliu")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("user is %s", user.Username)
}

func TestUserServiceImpl_Register(t *testing.T) {
	err := dao.InitMysql("127.0.0.1", "3306", "root", "root", "user")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = redis.InitRedis("127.0.0.1","6379", "" )
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	userService := &UserServiceImpl{userDAO:&dao.UserDAOImpl{}}
	user, err := userService.Register(context.Background(),&RegisterUserVO{
		Username: "alexliu",
		Email:    "alexliu",
		Password: "alexliu@qq.com",
	})
	if err != nil{
		t.Error(err)
		t.FailNow()
	}

	t.Logf("user id is %d", user.ID)
}
