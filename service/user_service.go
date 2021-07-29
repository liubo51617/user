package service

import (
	"context"
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/liubo51617/user/dao"
	"github.com/liubo51617/user/redis"
	"log"
	"time"
)

var (
	ErrUserExisted = errors.New("user is existed")
	ErrPassword = errors.New("email and password are not match")
	ErrRegistering = errors.New("email is registering")
)

type UserInfoDTO struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type RegisterUserVO struct {
	Username string
	Email string
	Password string
}

type UserService interface {
	Login(ctx context.Context, email, password string) (*UserInfoDTO, error)
	Register(ctx context.Context, vo *RegisterUserVO) (*UserInfoDTO, error)
}

type UserServiceImpl struct {
	userDAO dao.UserDAO
}

func MakeUserServiceImpl(userDAO dao.UserDAO) UserService {
	return &UserServiceImpl{
		userDAO: userDAO,
	}
}

func (userService *UserServiceImpl) Login(ctx context.Context, email, password string) (*UserInfoDTO, error)  {
	user, err := userService.userDAO.SelectByEmail(email)
	if err == nil {
		if user.Password == password {
			return &UserInfoDTO{
				ID:       user.ID,
				Username: user.Username,
				Email:    user.Email,
			}, nil
		} else {
			return nil, ErrPassword
		}
	}else {
		log.Printf("err : %s",err)
	}
	return nil ,err
}

func (userService *UserServiceImpl) Register(ctx context.Context, vo *RegisterUserVO) (*UserInfoDTO, error) {
	lock := redis.GetRedisLock(vo.Email, time.Duration(5)*time.Second)
	err := lock.Lock()
	if err != nil {
		log.Printf("err : %s",err)
		return nil, ErrRegistering
	}
	defer lock.Unlock()

	existUser, err := userService.userDAO.SelectByEmail(vo.Email)
	if (err == nil && existUser == nil) || err == gorm.ErrRecordNotFound {
		newsUser := &dao.UserEntity{
			Username:  vo.Username,
			Password:  vo.Password,
			Email:     vo.Email,
		}
		err = userService.userDAO.Save(newsUser)
		if err == nil {
			return &UserInfoDTO{
				ID:       newsUser.ID,
				Username: newsUser.Username,
				Email:    newsUser.Email,
			} , nil
		} else {
			return nil, err
		}
	}
	if err != nil {
		err = ErrUserExisted
	}
	return nil, err
}




