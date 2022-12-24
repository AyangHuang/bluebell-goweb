package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/utils/jwt"
	"bluebell/utils/snowflake"
	"database/sql"
	"errors"
)

func SignUp(p *models.ParamSignUp) (err error) {
	// 1.判断用户存不存在
	if exit, err := mysql.CheckUserExist(p.Username); err == nil && !exit {
		// 雪花算法生成唯一 id
		userID := snowflake.GenID()
		user := &models.MysqlUser{
			UserID:   userID,
			Username: p.Username,
			Password: p.Password,
		}
		// 2.数据库持久化
		if err = mysql.InsertUser(user); err != nil {
			return err
		}
	}
	return mysql.ErrorUserExist
}

// Login 登录
func Login(p *models.ParamLogin) (acToken, reToken string, err error) {
	// 1. 验证用户存在与否和密码
	user := new(models.MysqlUser)
	// 传递的是指针，就能拿到user.UserID
	if err = mysql.FindUserByUsername(p.Username, user); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", mysql.ErrorUserNotExist
		}
		return
	}
	if p.Password != user.Password || p.Username != user.Password {
		return "", "", mysql.ErrorInvalidPassword
	}
	// 2. 生成JWT
	acToken, reToken, err = jwt.GetTowToken(user.UserID, user.Username)
	return
}

func GetAccessToken(userID int64) (string, error) {
	username, err := mysql.FindUserNameByUserID(userID)
	if err != nil {
		return "", err
	}
	return jwt.GetToken(userID, username, jwt.RefreshToken)
}
