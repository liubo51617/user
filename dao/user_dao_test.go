package dao

import (
	"testing"
)

func TestUserDAOImpl_Save(t *testing.T) {
	UserDAO := &UserDAOImpl{}
	err := InitMysql("127.0.0.1", "3306", "root", "root", "user")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	user := &UserEntity{
		Username:  "alexliu",
		Password:  "alexliu",
		Email:     "alexliu@qq.com",
	}

	err = UserDAO.Save(user)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("new user id is %d", user.ID)
}

func TestUserDAOImpl_SelectByEmail(t *testing.T) {
	userDAO := &UserDAOImpl{}
	err := InitMysql("127.0.0.1", "3306", "root", "root", "user")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	var user *UserEntity
	user,err = userDAO.SelectByEmail("alexliu@qq.com")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("username is %s",user.Username)

}
